// backend/services/speech_piper.go
package services

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// PiperService synthesises speech using the Piper TTS engine.
// Piper must be installed locally (https://github.com/rhasspy/piper).
// It does not support speech recognition.
type PiperService struct{}

// NewPiperService creates a new PiperService.
func NewPiperService() *PiperService {
	return &PiperService{}
}

// Synthesize converts text to WAV audio bytes using a locally installed Piper binary.
func (s *PiperService) Synthesize(text string, options SpeechOptions) ([]byte, error) {
	if !s.commandExists("piper") {
		return nil, fmt.Errorf("Piper TTS requires local installation. Install from https://github.com/rhasspy/piper or use cloud services")
	}
	return s.synthesizeLocalPiper(text, options)
}

// synthesizeLocalPiper runs the piper binary and returns the generated WAV bytes.
func (s *PiperService) synthesizeLocalPiper(text string, options SpeechOptions) ([]byte, error) {
	tempOutputFile := "/tmp/piper_output.wav"

	args := []string{
		"--model", "en_US-lessac-medium", // Default English model — override via voice option.
		"--output_file", tempOutputFile,
	}

	if options.Voice != "" {
		// Replace the default model with the caller-specified voice/model path.
		args[1] = options.Voice
	}

	cmd := exec.Command("piper", args...)
	cmd.Stdin = strings.NewReader(text)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Piper TTS failed: %v, output: %s", err, string(output))
	}

	audioData, err := os.ReadFile(tempOutputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read Piper output: %v", err)
	}

	os.Remove(tempOutputFile)
	return audioData, nil
}

// Transcribe is not supported by Piper (TTS engine only).
func (s *PiperService) Transcribe(_ []byte, _ SpeechOptions) (string, error) {
	return "", fmt.Errorf("Piper does not support speech-to-text. Use cloud services like OpenAI Whisper")
}

// commandExists returns true if the named executable is on PATH.
func (s *PiperService) commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
