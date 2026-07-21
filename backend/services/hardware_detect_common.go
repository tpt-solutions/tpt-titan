package services

import (
	"os"
	"os/exec"
	"strings"
)

// detectGPU reports whether a discrete/compute GPU appears to be available.
// A fully reliable cross-platform GPU probe requires vendor SDKs (CUDA,
// Metal, D3D); this does a best-effort check of common environment hints and
// otherwise returns false so callers don't overcommit resources.
func detectGPU() bool {
	// Common environment hints that a GPU stack is configured.
	for _, env := range []string{"CUDA_VISIBLE_DEVICES", "NVIDIA_VISIBLE_DEVICES", "ROCR_VISIBLE_DEVICES"} {
		if v := os.Getenv(env); v != "" && v != "-1" {
			return true
		}
	}
	return detectGPUPlatform()
}

// detectCPUSpeedMHz returns the detected CPU base frequency in MHz, or 0 if it
// cannot be determined. Platform-specific implementations are in the
// hardware_detect_windows.go and hardware_detect_other.go files.

// runCmd runs a command and returns trimmed stdout, or an error.
func runCmd(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
