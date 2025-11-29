// Package type45 implements SMBIOS Type 45 - Firmware Inventory Information
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type45

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Firmware Inventory Information
const StructureType uint8 = 45

// FirmwareInventory represents Type 45 - Firmware Inventory Information
type FirmwareInventory struct {
	Header                   gosmbios.Header
	FirmwareComponentName    string
	FirmwareVersion          string
	VersionFormat            VersionFormat
	FirmwareID               string
	FirmwareIDFormat         FirmwareIDFormat
	ReleaseDate              string
	Manufacturer             string
	LowestSupportedVersion   string
	ImageSize                uint64
	Characteristics          Characteristics
	State                    FirmwareState
	AssociatedComponentCount uint8
	AssociatedComponentHandles []uint16
}

// VersionFormat identifies the firmware version format
type VersionFormat uint8

const (
	VersionFormatFreeForm      VersionFormat = 0x00
	VersionFormatMajorMinor    VersionFormat = 0x01
	VersionFormat32BitHex      VersionFormat = 0x02
	VersionFormat64BitHex      VersionFormat = 0x03
)

func (v VersionFormat) String() string {
	switch v {
	case VersionFormatFreeForm:
		return "Free-form String"
	case VersionFormatMajorMinor:
		return "Major.Minor"
	case VersionFormat32BitHex:
		return "32-bit Hex"
	case VersionFormat64BitHex:
		return "64-bit Hex"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(v))
	}
}

// FirmwareIDFormat identifies the firmware ID format
type FirmwareIDFormat uint8

const (
	FirmwareIDFormatFreeForm FirmwareIDFormat = 0x00
	FirmwareIDFormatUEFI     FirmwareIDFormat = 0x01
)

func (f FirmwareIDFormat) String() string {
	switch f {
	case FirmwareIDFormatFreeForm:
		return "Free-form String"
	case FirmwareIDFormatUEFI:
		return "UEFI ESRT FwClass GUID"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(f))
	}
}

// Characteristics represents firmware inventory characteristics
type Characteristics uint16

const (
	CharUpdatable              Characteristics = 1 << 0
	CharWriteProtected         Characteristics = 1 << 1
)

// IsUpdatable returns true if the firmware is updatable
func (c Characteristics) IsUpdatable() bool {
	return c&CharUpdatable != 0
}

// IsWriteProtected returns true if the firmware is write-protected
func (c Characteristics) IsWriteProtected() bool {
	return c&CharWriteProtected != 0
}

func (c Characteristics) String() string {
	var chars []string
	if c.IsUpdatable() {
		chars = append(chars, "Updatable")
	}
	if c.IsWriteProtected() {
		chars = append(chars, "Write-Protected")
	}
	if len(chars) == 0 {
		return "None"
	}
	return fmt.Sprintf("%v", chars)
}

// FirmwareState represents the firmware state
type FirmwareState uint8

const (
	StateOther           FirmwareState = 0x01
	StateUnknown         FirmwareState = 0x02
	StateDisabled        FirmwareState = 0x03
	StateEnabled         FirmwareState = 0x04
	StateAbsent          FirmwareState = 0x05
	StateStandbyOffline  FirmwareState = 0x06
	StateStandbySpare    FirmwareState = 0x07
	StateUnavailableOffline FirmwareState = 0x08
)

func (s FirmwareState) String() string {
	switch s {
	case StateOther:
		return "Other"
	case StateUnknown:
		return "Unknown"
	case StateDisabled:
		return "Disabled"
	case StateEnabled:
		return "Enabled"
	case StateAbsent:
		return "Absent"
	case StateStandbyOffline:
		return "Standby Offline"
	case StateStandbySpare:
		return "Standby Spare"
	case StateUnavailableOffline:
		return "Unavailable Offline"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(s))
	}
}

// Parse parses a Firmware Inventory Information structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*FirmwareInventory, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 24 bytes
	if len(s.Data) < 24 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &FirmwareInventory{
		Header:                 s.Header,
		FirmwareComponentName:  s.GetString(s.GetByte(0x04)),
		FirmwareVersion:        s.GetString(s.GetByte(0x05)),
		VersionFormat:          VersionFormat(s.GetByte(0x06)),
		FirmwareID:             s.GetString(s.GetByte(0x07)),
		FirmwareIDFormat:       FirmwareIDFormat(s.GetByte(0x08)),
		ReleaseDate:            s.GetString(s.GetByte(0x09)),
		Manufacturer:           s.GetString(s.GetByte(0x0A)),
		LowestSupportedVersion: s.GetString(s.GetByte(0x0B)),
		ImageSize:              s.GetQWord(0x0C),
		Characteristics:        Characteristics(s.GetWord(0x14)),
		State:                  FirmwareState(s.GetByte(0x16)),
		AssociatedComponentCount: s.GetByte(0x17),
	}

	// Read associated component handles
	offset := 0x18
	for i := uint8(0); i < info.AssociatedComponentCount; i++ {
		if offset+1 >= len(s.Data) {
			break
		}
		handle := s.GetWord(offset)
		info.AssociatedComponentHandles = append(info.AssociatedComponentHandles, handle)
		offset += 2
	}

	return info, nil
}

// Get retrieves the first Firmware Inventory Information from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*FirmwareInventory, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Firmware Inventory Information structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*FirmwareInventory, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var inventories []*FirmwareInventory
	for i := range structures {
		inv, err := Parse(&structures[i])
		if err == nil {
			inventories = append(inventories, inv)
		}
	}

	if len(inventories) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return inventories, nil
}

// ImageSizeString returns the image size as a human-readable string
func (f *FirmwareInventory) ImageSizeString() string {
	if f.ImageSize == 0xFFFFFFFFFFFFFFFF {
		return "Unknown"
	}
	if f.ImageSize >= 1024*1024*1024 {
		return fmt.Sprintf("%.2f GB", float64(f.ImageSize)/(1024*1024*1024))
	}
	if f.ImageSize >= 1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(f.ImageSize)/(1024*1024))
	}
	if f.ImageSize >= 1024 {
		return fmt.Sprintf("%.2f KB", float64(f.ImageSize)/1024)
	}
	return fmt.Sprintf("%d bytes", f.ImageSize)
}
