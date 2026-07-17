//go:build windows
// +build windows

package services

import (
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

// memStatusEx mirrors the Win32 MEMORYSTATUSEX structure for GlobalMemoryStatusEx.
type memStatusEx struct {
	Length               uint32
	MemoryLoad           uint32
	TotalPhys            uint64
	AvailPhys            uint64
	TotalPageFile        uint64
	AvailPageFile        uint64
	TotalVirtual         uint64
	AvailVirtual         uint64
	AvailExtendedVirtual uint64
}

// detectTotalSystemRAMGB returns the total installed system RAM in gigabytes
// using the Win32 GlobalMemoryStatusEx API via the kernel32 lazy DLL.
func detectTotalSystemRAMGB() (int, bool) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("GlobalMemoryStatusEx")
	if proc.Find() != nil {
		return 0, false
	}

	var m memStatusEx
	m.Length = uint32(unsafe.Sizeof(m))
	r, _, _ := proc.Call(uintptr(unsafe.Pointer(&m)))
	if r == 0 {
		return 0, false
	}
	gb := int(m.TotalPhys / (1024 * 1024 * 1024))
	if gb < 1 {
		return 0, false
	}
	return gb, true
}

// detectFreeDiskGB returns free disk space (in GB) for the current working
// directory's volume using GetDiskFreeSpaceExW.
func detectFreeDiskGB() (int, bool) {
	wd, err := os.Getwd()
	if err != nil {
		return 0, false
	}
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("GetDiskFreeSpaceExW")
	if proc.Find() != nil {
		return 0, false
	}

	var freeBytesAvailable uint64
	pwd, _ := syscall.UTF16PtrFromString(wd)
	r, _, _ := proc.Call(
		uintptr(unsafe.Pointer(pwd)),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		0,
		0,
	)
	if r == 0 {
		return 0, false
	}
	gb := int(freeBytesAvailable / (1024 * 1024 * 1024))
	return gb, true
}

// detectGPUPlatform reports GPU availability on Windows via the display
// adapter registry key. Returns false if the query fails.
func detectGPUPlatform() bool {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Class\{4d36e968-e325-11ce-bfc1-08002be10318}`,
		registry.ENUMERATE_SUB_KEYS|registry.READ)
	if err != nil {
		return false
	}
	defer k.Close()
	subKeys, err := k.ReadSubKeyNames(-1)
	if err != nil {
		return false
	}
	return len(subKeys) > 0
}

// detectCPUSpeedWindowsMHz reads the CPU nominal frequency (MHz) from the
// registry. Returns ok=false if unavailable.
func detectCPUSpeedWindowsMHz() (int, bool) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`HARDWARE\DESCRIPTION\System\CentralProcessor\0`, registry.READ)
	if err != nil {
		return 0, false
	}
	defer k.Close()
	val, _, err := k.GetIntegerValue("~MHz")
	if err != nil {
		return 0, false
	}
	return int(val), true
}
