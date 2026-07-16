// backend/services/speech_local_stt.go
package services

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/uuid"
)

// LocalSTTService performs speech-to-text using platform-native or open-source
// recognition engines (whisper.cpp, Kaldi, PocketSphinx, Julius).
type LocalSTTService struct{}

// NewLocalSTTService creates a new LocalSTTService.
func NewLocalSTTService() *LocalSTTService {
	return &LocalSTTService{}
}

// Transcribe converts audio bytes to text using the best available local engine
// for the current operating system. whisper.cpp is checked first on every
// platform since it's a single portable binary that runs on Windows, macOS,
// and Linux alike; platform/engine-specific fallbacks follow.
func (s *LocalSTTService) Transcribe(audioData []byte, options SpeechOptions) (string, error) {
	if len(audioData) == 0 {
		return "", fmt.Errorf("no audio data provided for transcription")
	}

	if bin := s.whisperCppBinary(); bin != "" {
		return s.transcribeWithWhisperCpp(bin, audioData, options)
	}

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

// whisperCppBinary returns the configured or discovered whisper.cpp executable
// name, or "" if none is available. The binary name/path can be overridden via
// WHISPER_CPP_BINARY (e.g. "whisper-cli", "main", or a full path).
func (s *LocalSTTService) whisperCppBinary() string {
	if custom := os.Getenv("WHISPER_CPP_BINARY"); custom != "" {
		if s.commandExists(custom) {
			return custom
		}
		return ""
	}
	for _, candidate := range []string{"whisper-cli", "whisper-cpp", "whisper"} {
		if s.commandExists(candidate) {
			return candidate
		}
	}
	return ""
}

// whisperCppModel returns the GGML model path for whisper.cpp, configurable via
// WHISPER_CPP_MODEL. Defaults to "models/ggml-base.en.bin" relative to the
// working directory, which is where whisper.cpp's own download script places it.
func (s *LocalSTTService) whisperCppModel() string {
	if model := os.Getenv("WHISPER_CPP_MODEL"); model != "" {
		return model
	}
	return filepath.Join("models", "ggml-base.en.bin")
}

// transcribeWithWhisperCpp shells out to a locally installed whisper.cpp binary
// and parses its JSON output for the transcript text.
func (s *LocalSTTService) transcribeWithWhisperCpp(binary string, audioData []byte, options SpeechOptions) (string, error) {
	model := s.whisperCppModel()
	if _, err := os.Stat(model); err != nil {
		return "", fmt.Errorf("whisper.cpp model not found at %q (set WHISPER_CPP_MODEL): %v", model, err)
	}

	tempAudioFile := filepath.Join(os.TempDir(), fmt.Sprintf("whisper_input_%s.wav", uuid.NewString()))
	if err := os.WriteFile(tempAudioFile, audioData, 0644); err != nil {
		return "", fmt.Errorf("failed to save audio file: %v", err)
	}
	defer os.Remove(tempAudioFile)

	outputPrefix := strings.TrimSuffix(tempAudioFile, filepath.Ext(tempAudioFile))
	args := []string{
		"-m", model,
		"-f", tempAudioFile,
		"-oj", "-of", outputPrefix,
		"-np",
	}
	if options.Language != "" {
		args = append(args, "-l", options.Language)
	}

	cmd := exec.Command(binary, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("whisper.cpp failed: %v, output: %s", err, string(output))
	}

	jsonPath := outputPrefix + ".json"
	defer os.Remove(jsonPath)
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return "", fmt.Errorf("whisper.cpp did not produce output: %v", err)
	}

	var result struct {
		Transcription []struct {
			Text string `json:"text"`
		} `json:"transcription"`
	}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return "", fmt.Errorf("failed to parse whisper.cpp output: %v", err)
	}

	var sb strings.Builder
	for _, segment := range result.Transcription {
		sb.WriteString(segment.Text)
	}
	text := strings.TrimSpace(sb.String())
	if text == "" {
		return "", fmt.Errorf("whisper.cpp produced no transcript")
	}
	return text, nil
}

// transcribeWindows — Windows built-in recognition is real-time only, and no
// portable file-based engine is installed unless whisper.cpp is available.
func (s *LocalSTTService) transcribeWindows(audioData []byte, options SpeechOptions) (string, error) {
	return "", fmt.Errorf("Windows file-based STT requires whisper.cpp (set WHISPER_CPP_BINARY/WHISPER_CPP_MODEL) or a cloud service like OpenAI Whisper or Azure Speech Service")
}

// transcribeMacOS — macOS lacks a file-based STT API without native bindings,
// and no portable file-based engine is installed unless whisper.cpp is available.
func (s *LocalSTTService) transcribeMacOS(audioData []byte, options SpeechOptions) (string, error) {
	return "", fmt.Errorf("macOS file-based STT requires whisper.cpp (set WHISPER_CPP_BINARY/WHISPER_CPP_MODEL) or a cloud service like OpenAI Whisper or Azure Speech Service")
}

// transcribeLinux attempts Kaldi, then PocketSphinx, then Julius.
func (s *LocalSTTService) transcribeLinux(audioData []byte, options SpeechOptions) (string, error) {
	if s.commandExists("kaldi-gstreamer-server") || s.commandExists("online2-wav-nnet3-latgen-faster") {
		return s.transcribeWithKaldi(audioData, options)
	}
	if s.commandExists("pocketsphinx_continuous") {
		return s.transcribeWithPocketSphinx(audioData, options)
	}
	if s.commandExists("julius") {
		return s.transcribeWithJulius(audioData, options)
	}

	return "", fmt.Errorf("no local STT engine found. Install whisper.cpp, kaldi, pocketsphinx, or julius, or use cloud services like OpenAI Whisper")
}

// transcribeWithKaldi uses Kaldi for high-accuracy transcription.
// Real Kaldi deployments require acoustic and language model files, which vary
// too much per-install to drive generically here — use whisper.cpp instead.
func (s *LocalSTTService) transcribeWithKaldi(audioData []byte, options SpeechOptions) (string, error) {
	tempAudioFile := filepath.Join(os.TempDir(), fmt.Sprintf("kaldi_input_%s.wav", uuid.NewString()))
	if err := os.WriteFile(tempAudioFile, audioData, 0644); err != nil {
		return "", fmt.Errorf("failed to save audio file: %v", err)
	}
	defer os.Remove(tempAudioFile)

	return "", fmt.Errorf("Kaldi STT requires a custom decoding graph and models specific to your install; use whisper.cpp (set WHISPER_CPP_BINARY) or a cloud service like OpenAI Whisper instead")
}

// transcribeWithPocketSphinx uses PocketSphinx for simpler local transcription.
// Model paths default to the standard Debian/Ubuntu package locations but can
// be overridden via POCKETSPHINX_HMM, POCKETSPHINX_LM, and POCKETSPHINX_DICT
// for other distros or custom-trained models.
func (s *LocalSTTService) transcribeWithPocketSphinx(audioData []byte, options SpeechOptions) (string, error) {
	tempAudioFile := filepath.Join(os.TempDir(), fmt.Sprintf("pocketsphinx_input_%s.wav", uuid.NewString()))
	if err := os.WriteFile(tempAudioFile, audioData, 0644); err != nil {
		return "", fmt.Errorf("failed to save audio file: %v", err)
	}
	defer os.Remove(tempAudioFile)

	hmm := envOrDefault("POCKETSPHINX_HMM", "/usr/share/pocketsphinx/model/en-us/en-us")
	lm := envOrDefault("POCKETSPHINX_LM", "/usr/share/pocketsphinx/model/en-us/en-us.lm.bin")
	dict := envOrDefault("POCKETSPHINX_DICT", "/usr/share/pocketsphinx/model/en-us/cmudict-en-us.dict")

	for name, path := range map[string]string{"hmm": hmm, "lm": lm, "dict": dict} {
		if _, err := os.Stat(path); err != nil {
			return "", fmt.Errorf("PocketSphinx %s model not found at %q (set POCKETSPHINX_%s): %v", name, path, strings.ToUpper(name), err)
		}
	}

	args := []string{
		"-inmic", "no",
		"-infile", tempAudioFile,
		"-bestpath", "1",
		"-hmm", hmm,
		"-lm", lm,
		"-dict", dict,
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

// transcribeWithJulius requires model files specific to the install and isn't
// generically drivable — surface that honestly rather than faking output.
func (s *LocalSTTService) transcribeWithJulius(audioData []byte, options SpeechOptions) (string, error) {
	tempAudioFile := filepath.Join(os.TempDir(), fmt.Sprintf("julius_input_%s.wav", uuid.NewString()))
	if err := os.WriteFile(tempAudioFile, audioData, 0644); err != nil {
		return "", fmt.Errorf("failed to save audio file: %v", err)
	}
	defer os.Remove(tempAudioFile)

	return "", fmt.Errorf("Julius STT requires specific acoustic and language model files; use whisper.cpp (set WHISPER_CPP_BINARY) or a cloud service like OpenAI Whisper or Google Speech-to-Text instead")
}

// commandExists returns true if the named executable is on PATH.
func (s *LocalSTTService) commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// envOrDefault returns the named environment variable, falling back to def if unset.
func envOrDefault(name, def string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return def
}
