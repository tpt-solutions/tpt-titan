package services

import (
	"testing"
)

func TestDetectHardwareRealValues(t *testing.T) {
	hs := NewHardwareService()
	info, err := hs.DetectHardware()
	if err != nil {
		t.Fatalf("DetectHardware: %v", err)
	}
	if info.CPUCores < 1 {
		t.Errorf("CPUCores should be >= 1, got %d", info.CPUCores)
	}
	if info.RAMGB < 1 {
		t.Errorf("RAMGB should be detected >= 1, got %d", info.RAMGB)
	}
	if info.DiskSpace < 1 {
		t.Errorf("DiskSpace should be detected >= 1, got %d", info.DiskSpace)
	}
}

func TestCheckCompatibilityReflectsSize(t *testing.T) {
	hs := NewHardwareService()
	small := hs.CheckCompatibility(1.0, &HardwareInfo{RAMGB: 8})
	if !small {
		t.Errorf("1GB model should be compatible with 8GB RAM")
	}
	big := hs.CheckCompatibility(64.0, &HardwareInfo{RAMGB: 8})
	if big {
		t.Errorf("64GB model should NOT be compatible with 8GB RAM")
	}
}
