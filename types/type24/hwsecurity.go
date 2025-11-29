// Package type24 implements SMBIOS Type 24 - Hardware Security
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type24

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Hardware Security
const StructureType uint8 = 24

// HardwareSecurity represents Type 24 - Hardware Security
type HardwareSecurity struct {
	Header           gosmbios.Header
	HardwareSettings HardwareSettings
}

// HardwareSettings represents hardware security settings
type HardwareSettings uint8

// PowerOnPasswordStatus returns the power-on password status
func (h HardwareSettings) PowerOnPasswordStatus() SecurityStatus {
	return SecurityStatus(h & 0x03)
}

// KeyboardPasswordStatus returns the keyboard password status
func (h HardwareSettings) KeyboardPasswordStatus() SecurityStatus {
	return SecurityStatus((h >> 2) & 0x03)
}

// AdministratorPasswordStatus returns the administrator password status
func (h HardwareSettings) AdministratorPasswordStatus() SecurityStatus {
	return SecurityStatus((h >> 4) & 0x03)
}

// FrontPanelResetStatus returns the front panel reset status
func (h HardwareSettings) FrontPanelResetStatus() SecurityStatus {
	return SecurityStatus((h >> 6) & 0x03)
}

func (h HardwareSettings) String() string {
	return fmt.Sprintf("Power-On: %s, Keyboard: %s, Admin: %s, Front Panel: %s",
		h.PowerOnPasswordStatus().String(),
		h.KeyboardPasswordStatus().String(),
		h.AdministratorPasswordStatus().String(),
		h.FrontPanelResetStatus().String())
}

// SecurityStatus identifies the security status
type SecurityStatus uint8

const (
	SecurityStatusDisabled     SecurityStatus = 0x00
	SecurityStatusEnabled      SecurityStatus = 0x01
	SecurityStatusNotImplemented SecurityStatus = 0x02
	SecurityStatusUnknown      SecurityStatus = 0x03
)

func (s SecurityStatus) String() string {
	switch s {
	case SecurityStatusDisabled:
		return "Disabled"
	case SecurityStatusEnabled:
		return "Enabled"
	case SecurityStatusNotImplemented:
		return "Not Implemented"
	case SecurityStatusUnknown:
		return "Unknown"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(s))
	}
}

// Parse parses a Hardware Security structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*HardwareSecurity, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 5 bytes
	if len(s.Data) < 5 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &HardwareSecurity{
		Header:           s.Header,
		HardwareSettings: HardwareSettings(s.GetByte(0x04)),
	}

	return info, nil
}

// Get retrieves the Hardware Security from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*HardwareSecurity, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}
