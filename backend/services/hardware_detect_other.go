//go:build !windows
// +build !windows

package services

import (
	"os"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

// detectTotalSystemRAMGB returns the total installed system RAM in gigabytes.
func detectTotalSystemRAMGB() (int, bool) {
	var info unix.Sysinfo_t
	if err := unix.Sysinfo(&info); err != nil {
		return 0, false
	}
	gb := int(info.Totalram / (1024 * 1024 * 1024))
	if gb < 1 {
		gb = int((info.Totalram / (1024 * 1024))) / 1024
	}
	if gb < 1 {
		return 0, false
	}
	return gb, true
}

// detectFreeDiskGB returns free disk space (in GB) for the current working
// directory's volume.
func detectFreeDiskGB() (int, bool) {
	var stat unix.Statfs_t
	wd, err := os.Getwd()
	if err != nil {
		return 0, false
	}
	if err := unix.Statfs(wd, &stat); err != nil {
		return 0, false
	}
	bavail := stat.Bavail * uint64(stat.Bsize)
	gb := int(bavail / (1024 * 1024 * 1024))
	return gb, true
}

// detectGPUPlatform is the non-Windows variant; a reliable GPU probe needs
// vendor SDKs, so it conservatively reports false.
func detectGPUPlatform() bool {
	return false
}

// detectCPUSpeedMHz returns the detected CPU base frequency in MHz, or 0 if
// it cannot be determined.
func detectCPUSpeedMHz() int {
	// Linux: read from /proc/cpuinfo
	if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "cpu MHz") || strings.HasPrefix(line, "clock") {
				if idx := strings.Index(line, ":"); idx >= 0 {
					s := strings.TrimSpace(strings.TrimSuffix(line[idx+1:], "MHz"))
					if mhz, err := strconv.Atoi(s); err == nil {
						return mhz
					}
				}
			}
		}
	}
	// macOS: sysctl
	if out, err := runCmd("sysctl", "-n", "hw.cpufrequency"); err == nil {
		if hz, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64); err == nil && hz > 0 {
			return int(hz / 1_000_000)
		}
	}
	return 0
}
