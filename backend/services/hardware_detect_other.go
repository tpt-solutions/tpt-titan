//go:build !windows
// +build !windows

package services

import (
	"os"

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
