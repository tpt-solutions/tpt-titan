package services

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

// DetectHardware attempts to detect system hardware capabilities
func (s *HardwareService) DetectHardware() (*HardwareInfo, error) {
	// This is a basic implementation - in production, you'd use system calls
	// For now, we'll return conservative defaults
	// TODO: Implement actual hardware detection using runtime.NumCPU(), syscall, etc.

	info := &HardwareInfo{
		RAMGB:     8,  // Conservative default - most users have at least 8GB
		HasGPU:    false, // Conservative - assume no GPU unless detected
		CPUCores:  4,  // Most systems have at least 4 cores
		CPUSpeed:  2000, // 2GHz base speed
		DiskSpace: 50, // 50GB free space minimum
	}

	// Try to detect actual RAM (simplified)
	// In production, use: runtime.MemStats, syscall.Sysinfo, etc.
	// For now, return defaults that work for most systems

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
