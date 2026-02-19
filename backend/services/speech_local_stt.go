// backend/services/speech_local_stt.go
package services

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// LocalSTTService performs speech-to-text using platform-native or open-source
// recognition engines (Kaldi, PocketSphinx, Julius).
type LocalSTTService struct{}

// NewLocalSTTService creates a new LocalSTTService.
func NewLocalSTTService() *LocalSTTService {
	return &LocalSTTService{}
}

// Transcribe converts audio bytes to text using the best available local engine
// for the current operating system.
func (s *LocalSTTService) Transcribe(audioData []byte, options SpeechOptions) (string, error) {
	switch runtime.GOOS {
	case "windows":
		return s.transcribeWindows(audioData, options)
	case "darwin":
		return s.transcribeMacOS(audioData, options)
	case "linux":
		return s.transcribeLinux(audioData, options)
	default:
		return "", fmt.Errorf("local STT not supported on platform: %s", runtime.GOOS)
	}
}

// transcribeWindows — Windows built-in recognition is real-time only.
func (s *LocalSTTService) transcribeWindows(audioData []byte, options SpeechOptions) (string, error) {
	if len(audioData) == 0 {
		return "", fmt.Errorf("no audio data provided for transcription")
	}
	return "", fmt.Errorf("Windows file-based STT not yet implemented. Use cloud services like OpenAI Whisper or Azure Speech Service")
}

// transcribeMacOS — macOS lacks a file-based STT API without native bindings.
func (s *LocalSTTService) transcribeMacOS(audioData []byte, options SpeechOptions) (string, error) {
	return "", fmt.Errorf("macOS built-in STT has limited file support. Consider using OpenAI Whisper or Azure Speech Service")
}

// transcribeLinux attempts Kaldi, then PocketSphinx, then Julius.
func (s *LocalSTTService) transcribeLinux(audioData []byte, options SpeechOptions) (string, error) {
	if len(audioData) == 0 {
		return "", fmt.Errorf("no audio data provided for transcription")
	}

	if s.commandExists("kaldi-gstreamer-server") || s.commandExists("online2-wav-nnet3-latgen-faster") {
		return s.transcribeWithKaldi(audioData, options)
	}
	if s.commandExists("pocketsphinx_continuous") {
		return s.transcribeWithPocketSphinx(audioData, options)
	}
	if s.commandExists("julius") {
		return s.transcribeWithJulius(audioData, options)
	}

	return "", fmt.Errorf("no local STT engine found. Please install kaldi, pocketsphinx, or julius, or use cloud services like OpenAI Whisper")
}

// transcribeWithKaldi uses Kaldi for high-accuracy transcription.
// Real Kaldi deployments require acoustic and language model files.
func (s *LocalSTTService) transcribeWithKaldi(audioData []byte, options SpeechOptions) (string, error) {
	tempAudioFile := "/tmp/kaldi_input.wav"
	if err := os.WriteFile(tempAudioFile, audioData, 0644); err != nil {
		return "", fmt.Errorf("failed to save audio file: %v", err)
	}
	defer os.Remove(tempAudioFile)

	return "", fmt.Errorf("Kaldi STT requires complex setup with acoustic and language models. Consider using cloud services like OpenAI Whisper or Google Speech-to-Text")
}

// transcribeWithPocketSphinx uses PocketSphinx for simpler local transcription.
func (s *LocalSTTService) transcribeWithPocketSphinx(audioData []byte, options SpeechOptions) (string, error) {
	tempAudioFile := "/tmp/pocketsphinx_input.wav"
	if err := os.WriteFile(tempAudioFile, audioData, 0644); err != nil {
		return "", fmt.Errorf("failed to save audio file: %v", err)
	}
	defer os.Remove(tempAudioFile)

	args := []string{
		"-inmic", "no",
		"-infile", tempAudioFile,
		"-bestpath", "1",
	}

	if options.Language == "en" || options.Language == "" {
		args = append(args, "-hmm", "/usr/share/pocketsphinx/model/en-us/en-us")
		args = append(args, "-lm", "/usr/share/pocketsphinx/model/en-us/en-us.lm.bin")
		args = append(args, "-dict", "/usr/share/pocketsphinx/model/en-us/cmudict-en-us.dict")
	}

	cmd := exec.Command("pocketsphinx_continuous", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("PocketSphinx failed: %v, output: %s", err, string(output))
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "INFO:") && strings.Contains(line, "READY") {
			continue
		}
		if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "INFO:") {
			return strings.TrimSpace(line), nil
		}
	}

	return "", fmt.Errorf("no transcript found in PocketSphinx output")
}

// transcribeWithJulius is a placeholder — Julius requires model files.
func (s *LocalSTTService) transcribeWithJulius(audioData []byte, options SpeechOptions) (string, error) {
	tempAudioFile := "/tmp/julius_input.wav"
	if err := os.WriteFile(tempAudioFile, audioData, 0644); err != nil {
		return "", fmt.Errorf("failed to save audio file: %v", err)
	}
	defer os.Remove(tempAudioFile)

	return "", fmt.Errorf("Julius STT requires specific acoustic and language model files. Consider using cloud services like OpenAI Whisper or Google Speech-to-Text")
}

// commandExists returns true if the named executable is on PATH.
func (s *LocalSTTService) commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
