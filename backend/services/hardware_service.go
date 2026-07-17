package services

import (
	"runtime"
)
// HardwareInfo represents detected hardware capabilities
type HardwareInfo struct {
	RAMGB     int  `json:"ram_gb"`
	HasGPU    bool `json:"has_gpu"`
	CPUCores  int  `json:"cpu_cores"`
	CPUSpeed  int  `json:"cpu_speed_mhz"` // CPU speed in MHz
	DiskSpace int  `json:"disk_space_gb"` // Available disk space in GB
}

// HardwareService handles hardware detection and analysis
type HardwareService struct{}

// NewHardwareService creates a new hardware service
func NewHardwareService() *HardwareService {
	return &HardwareService{}
}

// DetectHardware attempts to detect system hardware capabilities using
// real OS queries (CPU count, installed RAM, free disk space). Fields that
// cannot be determined reliably without elevated APIs (e.g. GPU presence,
// exact CPU frequency) are reported conservatively.
func (s *HardwareService) DetectHardware() (*HardwareInfo, error) {
	info := &HardwareInfo{
		HasGPU:   detectGPU(),
		CPUCores: runtime.NumCPU(),
		CPUSpeed: detectCPUSpeedMHz(),
	}

	if ram, ok := detectTotalSystemRAMGB(); ok {
		info.RAMGB = ram
	} else {
		info.RAMGB = 8 // conservative fallback
	}

	if disk, ok := detectFreeDiskGB(); ok {
		info.DiskSpace = disk
	} else {
		info.DiskSpace = 50 // conservative fallback
	}

	return info, nil
}

// CheckCompatibility checks if hardware meets requirements for a model
func (s *HardwareService) CheckCompatibility(modelSizeGB float64, hardware *HardwareInfo) bool {
	// Estimate RAM requirements (rough approximation)
	estimatedRAM := modelSizeGB * 1.5 // Model size * 1.5 for inference overhead

	// Add buffer for system
	requiredRAM := estimatedRAM + 2.0 // 2GB system buffer

	return float64(hardware.RAMGB) >= requiredRAM
}

// GetHardwareTier determines the hardware tier based on capabilities
func (s *HardwareService) GetHardwareTier(hardware *HardwareInfo) string {
	if hardware.RAMGB >= 32 && hardware.HasGPU {
		return "high-end"
	} else if hardware.RAMGB >= 16 {
		return "standard"
	} else if hardware.RAMGB >= 8 {
		return "low-resource"
	} else {
		return "minimal"
	}
}
