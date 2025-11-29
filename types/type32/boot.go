// Package type32 implements SMBIOS Type 32 - System Boot Information
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type32

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for System Boot Information
const StructureType uint8 = 32

// BootInfo represents Type 32 - System Boot Information
type BootInfo struct {
	Header     gosmbios.Header
	Reserved   [6]byte
	BootStatus BootStatus
}

// BootStatus represents the system boot status
type BootStatus uint8

// Boot status values
const (
	BootStatusNoErrors                 BootStatus = 0
	BootStatusNoBootableMedia          BootStatus = 1
	BootStatusNormalOSLoadFailed       BootStatus = 2
	BootStatusFirmwareHardwareFailure  BootStatus = 3
	BootStatusOSHardwareFailure        BootStatus = 4
	BootStatusUserRequestedBoot        BootStatus = 5
	BootStatusSecurityViolation        BootStatus = 6
	BootStatusPreviouslyRequestedImage BootStatus = 7
	BootStatusWatchdogTimerExpired     BootStatus = 8
	// 9-127: Reserved for future use
	// 128-191: Vendor/OEM-specific implementations
	// 192-255: Product-specific implementations
)

// String returns a human-readable boot status description
func (bs BootStatus) String() string {
	switch {
	case bs == BootStatusNoErrors:
		return "No errors detected"
	case bs == BootStatusNoBootableMedia:
		return "No bootable media"
	case bs == BootStatusNormalOSLoadFailed:
		return "Normal operating system load failed"
	case bs == BootStatusFirmwareHardwareFailure:
		return "Firmware-detected hardware failure"
	case bs == BootStatusOSHardwareFailure:
		return "Operating system-detected hardware failure"
	case bs == BootStatusUserRequestedBoot:
		return "User-requested boot"
	case bs == BootStatusSecurityViolation:
		return "System security violation"
	case bs == BootStatusPreviouslyRequestedImage:
		return "Previously-requested image"
	case bs == BootStatusWatchdogTimerExpired:
		return "System watchdog timer expired"
	case bs >= 9 && bs <= 127:
		return fmt.Sprintf("Reserved (%d)", bs)
	case bs >= 128 && bs <= 191:
		return fmt.Sprintf("Vendor/OEM-specific (%d)", bs)
	case bs >= 192:
		return fmt.Sprintf("Product-specific (%d)", bs)
	default:
		return fmt.Sprintf("Unknown (%d)", bs)
	}
}

// IsSuccess returns true if the boot was successful
func (bs BootStatus) IsSuccess() bool {
	return bs == BootStatusNoErrors
}

// IsFailure returns true if the boot failed
func (bs BootStatus) IsFailure() bool {
	switch bs {
	case BootStatusNoBootableMedia, BootStatusNormalOSLoadFailed,
		BootStatusFirmwareHardwareFailure, BootStatusOSHardwareFailure,
		BootStatusSecurityViolation, BootStatusWatchdogTimerExpired:
		return true
	}
	return false
}

// Parse parses a System Boot Information structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*BootInfo, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 11 bytes
	if len(s.Data) < 11 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &BootInfo{
		Header:     s.Header,
		BootStatus: BootStatus(s.GetByte(0x0A)),
	}

	// Copy reserved bytes
	copy(info.Reserved[:], s.Data[0x04:0x0A])

	return info, nil
}

// Get retrieves the System Boot Information from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*BootInfo, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}
