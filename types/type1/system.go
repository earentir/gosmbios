// Package type1 implements SMBIOS Type 1 - System Information
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type1

import (
	"encoding/hex"
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for System Information
const StructureType uint8 = 1

// SystemInfo represents Type 1 - System Information
type SystemInfo struct {
	Header       gosmbios.Header
	Manufacturer string
	ProductName  string
	Version      string
	SerialNumber string
	UUID         UUID
	WakeUpType   WakeUpType
	SKUNumber    string // SMBIOS 2.4+
	Family       string // SMBIOS 2.4+
}

// UUID represents a 128-bit Universal Unique Identifier
type UUID [16]byte

// WakeUpType identifies the event that caused the system to power up
type WakeUpType uint8

// Wake-up type values
const (
	WakeUpReserved        WakeUpType = 0x00
	WakeUpOther           WakeUpType = 0x01
	WakeUpUnknown         WakeUpType = 0x02
	WakeUpAPMTimer        WakeUpType = 0x03
	WakeUpModemRing       WakeUpType = 0x04
	WakeUpLANRemote       WakeUpType = 0x05
	WakeUpPowerSwitch     WakeUpType = 0x06
	WakeUpPCIPME          WakeUpType = 0x07
	WakeUpACPowerRestored WakeUpType = 0x08
)

// String returns a human-readable wake-up type description
func (w WakeUpType) String() string {
	switch w {
	case WakeUpReserved:
		return "Reserved"
	case WakeUpOther:
		return "Other"
	case WakeUpUnknown:
		return "Unknown"
	case WakeUpAPMTimer:
		return "APM Timer"
	case WakeUpModemRing:
		return "Modem Ring"
	case WakeUpLANRemote:
		return "LAN Remote"
	case WakeUpPowerSwitch:
		return "Power Switch"
	case WakeUpPCIPME:
		return "PCI PME#"
	case WakeUpACPowerRestored:
		return "AC Power Restored"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(w))
	}
}

// String returns the UUID in standard format (8-4-4-4-12)
func (u UUID) String() string {
	// SMBIOS uses mixed-endian UUID format
	// First 3 fields are little-endian, last 2 are big-endian
	return fmt.Sprintf("%02X%02X%02X%02X-%02X%02X-%02X%02X-%02X%02X-%02X%02X%02X%02X%02X%02X",
		u[3], u[2], u[1], u[0], // Time-Low (LE)
		u[5], u[4], // Time-Mid (LE)
		u[7], u[6], // Time-High-And-Version (LE)
		u[8], u[9], // Clock-Seq (BE)
		u[10], u[11], u[12], u[13], u[14], u[15]) // Node (BE)
}

// Bytes returns the raw UUID bytes
func (u UUID) Bytes() []byte {
	return u[:]
}

// Hex returns the UUID as a hexadecimal string without dashes
func (u UUID) Hex() string {
	return hex.EncodeToString(u[:])
}

// IsZero returns true if UUID is all zeros (not set)
func (u UUID) IsZero() bool {
	for _, b := range u {
		if b != 0 {
			return false
		}
	}
	return true
}

// IsInvalid returns true if UUID is all 0xFF (not settable)
func (u UUID) IsInvalid() bool {
	for _, b := range u {
		if b != 0xFF {
			return false
		}
	}
	return true
}

// Parse parses a System Information structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*SystemInfo, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length for SMBIOS 2.0 is 8 bytes
	if len(s.Data) < 8 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &SystemInfo{
		Header:       s.Header,
		Manufacturer: s.GetString(s.GetByte(0x04)),
		ProductName:  s.GetString(s.GetByte(0x05)),
		Version:      s.GetString(s.GetByte(0x06)),
		SerialNumber: s.GetString(s.GetByte(0x07)),
	}

	// UUID (SMBIOS 2.1+)
	if len(s.Data) >= 25 {
		copy(info.UUID[:], s.Data[0x08:0x18])
	}

	// Wake-up Type (SMBIOS 2.1+)
	if len(s.Data) >= 25 {
		info.WakeUpType = WakeUpType(s.GetByte(0x18))
	}

	// SKU Number and Family (SMBIOS 2.4+)
	if len(s.Data) >= 27 {
		info.SKUNumber = s.GetString(s.GetByte(0x19))
		info.Family = s.GetString(s.GetByte(0x1A))
	}

	return info, nil
}

// Get retrieves the System Information from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*SystemInfo, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// DisplayName returns a display-friendly system name
func (si *SystemInfo) DisplayName() string {
	if si.Manufacturer != "" && si.ProductName != "" {
		return fmt.Sprintf("%s %s", si.Manufacturer, si.ProductName)
	}
	if si.ProductName != "" {
		return si.ProductName
	}
	return "Unknown System"
}
