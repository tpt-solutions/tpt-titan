package services

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
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
// cannot be determined without elevated APIs.
func detectCPUSpeedMHz() int {
	switch runtime.GOOS {
	case "linux":
		if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				if strings.HasPrefix(line, "cpu MHz") || strings.HasPrefix(line, "clock") {
					if idx := strings.Index(line, ":"); idx >= 0 {
						if mhz, err := strconv.Atoi(strings.TrimSpace(strings.TrimSuffix(line[idx+1:], "MHz"))); err == nil {
							return mhz
						}
					}
				}
			}
		}
	case "darwin":
		if out, err := runCmd("sysctl", "-n", "hw.cpufrequency"); err == nil {
			if hz, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64); err == nil && hz > 0 {
				return int(hz / 1_000_000)
			}
		}
	case "windows":
		if mhz, ok := detectCPUSpeedWindowsMHz(); ok {
			return mhz
		}
	}
	return 0
}

// runCmd runs a command and returns trimmed stdout, or an error.
func runCmd(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
