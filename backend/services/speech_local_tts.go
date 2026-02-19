// backend/services/speech_local_tts.go
package services

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// LocalTTSService performs text-to-speech using platform system voices.
// Windows uses PowerShell + System.Speech.Synthesis, macOS uses the 'say'
// command, and Linux uses espeak-ng / espeak / festival.
type LocalTTSService struct{}

// NewLocalTTSService creates a new LocalTTSService.
func NewLocalTTSService() *LocalTTSService {
	return &LocalTTSService{}
}

// Synthesize converts text to audio bytes using the host operating system's
// built-in TTS engine.
func (s *LocalTTSService) Synthesize(text string, options SpeechOptions) ([]byte, error) {
	switch runtime.GOOS {
	case "windows":
		return s.synthesizeWindows(text, options)
	case "darwin":
		return s.synthesizeMacOS(text, options)
	case "linux":
		return s.synthesizeLinux(text, options)
	default:
		return nil, fmt.Errorf("local TTS not supported on platform: %s", runtime.GOOS)
	}
}

// synthesizeWindows uses PowerShell and System.Speech.Synthesis.
func (s *LocalTTSService) synthesizeWindows(text string, options SpeechOptions) ([]byte, error) {
	escapedText := strings.ReplaceAll(text, "'", "''")

	voice := options.Voice
	if voice == "" {
		voice = "Microsoft Zira Desktop"
	}

	psScript := fmt.Sprintf(`
		Add-Type -AssemblyName System.Speech
		$synth = New-Object System.Speech.Synthesis.SpeechSynthesizer
		$synth.SelectVoice('%s')
		$synth.Rate = %d
		$synth.Volume = %d
		$synth.Speak('%s')
	`, voice, s.speedToInt(options.Speed), s.volumeToInt(options.Volume), escapedText)

	cmd := exec.Command("powershell", "-Command", psScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Windows TTS failed: %v, output: %s", err, string(output))
	}

	// PowerShell TTS speaks directly to speakers; return a success indicator.
	return []byte("tts_completed"), nil
}

// synthesizeMacOS uses the macOS 'say' command.
func (s *LocalTTSService) synthesizeMacOS(text string, options SpeechOptions) ([]byte, error) {
	escapedText := strings.ReplaceAll(text, "'", "\\'")

	args := []string{}

	if options.Voice != "" {
		args = append(args, "-v", options.Voice)
	} else {
		args = append(args, "-v", "Samantha")
	}

	rate := int(200 * options.Speed)
	if rate < 90 {
		rate = 90
	}
	if rate > 600 {
		rate = 600
	}
	args = append(args, "-r", strconv.Itoa(rate))
	args = append(args, escapedText)

	cmd := exec.Command("say", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("macOS TTS failed: %v, output: %s", err, string(output))
	}

	// 'say' speaks directly to speakers; return a success indicator.
	return []byte("tts_completed"), nil
}

// synthesizeLinux attempts espeak-ng, espeak, or festival in order.
func (s *LocalTTSService) synthesizeLinux(text string, options SpeechOptions) ([]byte, error) {
	var cmd *exec.Cmd

	if s.commandExists("espeak-ng") {
		args := []string{"-v", "en-us", "-s", s.speedToEspeakRate(options.Speed)}
		if options.Voice != "" {
			args = append(args, "-v", options.Voice)
		}
		pitch := int(50 + (options.Pitch * 25))
		if pitch < 0 {
			pitch = 0
		}
		if pitch > 100 {
			pitch = 100
		}
		args = append(args, "-p", strconv.Itoa(pitch))
		args = append(args, text)
		cmd = exec.Command("espeak-ng", args...)
	} else if s.commandExists("espeak") {
		args := []string{"-v", "en-us", "-s", s.speedToEspeakRate(options.Speed)}
		if options.Voice != "" {
			args = append(args, "-v", options.Voice)
		}
		args = append(args, text)
		cmd = exec.Command("espeak", args...)
	} else if s.commandExists("festival") {
		return s.synthesizeFestival(text, options)
	} else {
		return nil, fmt.Errorf("no TTS engine found. Please install espeak-ng, espeak, or festival")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Linux TTS failed: %v, output: %s", err, string(output))
	}

	return []byte("tts_completed"), nil
}

// synthesizeFestival uses the Festival TTS engine as a Linux fallback.
func (s *LocalTTSService) synthesizeFestival(text string, options SpeechOptions) ([]byte, error) {
	festivalScript := fmt.Sprintf(`
		(Parameter.set 'Audio_Method 'Audio_Command)
		(Parameter.set 'Audio_Command "aplay -q -c 1 -t raw -f s16 -r $SRATE $FILE")
		(SayText "%s")
	`, strings.ReplaceAll(text, `"`, `\"`))

	cmd := exec.Command("festival", "--tts")
	cmd.Stdin = strings.NewReader(festivalScript)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Festival TTS failed: %v, output: %s", err, string(output))
	}

	return []byte("tts_completed"), nil
}

// speedToInt converts a speed multiplier to a PowerShell rate value (-10 to 10).
func (s *LocalTTSService) speedToInt(speed float64) int {
	rate := int((speed - 1.0) * 5)
	if rate < -10 {
		rate = -10
	}
	if rate > 10 {
		rate = 10
	}
	return rate
}

// volumeToInt converts a 0.0–1.0 volume to the 0–100 range PowerShell expects.
func (s *LocalTTSService) volumeToInt(volume float64) int {
	return int(volume * 100)
}

// speedToEspeakRate converts a speed multiplier to espeak words-per-minute (80–450).
func (s *LocalTTSService) speedToEspeakRate(speed float64) string {
	rate := int(175 * speed)
	if rate < 80 {
		rate = 80
	}
	if rate > 450 {
		rate = 450
	}
	return strconv.Itoa(rate)
}

// commandExists returns true if the named executable is on PATH.
func (s *LocalTTSService) commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
