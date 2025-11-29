// Package type2 implements SMBIOS Type 2 - Baseboard (Module) Information
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type2

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Baseboard Information
const StructureType uint8 = 2

// BaseboardInfo represents Type 2 - Baseboard (Module) Information
type BaseboardInfo struct {
	Header                   gosmbios.Header
	Manufacturer             string
	Product                  string
	Version                  string
	SerialNumber             string
	AssetTag                 string
	FeatureFlags             FeatureFlags
	LocationInChassis        string
	ChassisHandle            uint16
	BoardType                BoardType
	NumberOfContainedHandles uint8
	ContainedObjectHandles   []uint16
}

// FeatureFlags represents baseboard feature flags
type FeatureFlags uint8

// Feature flag bit definitions
const (
	FeatureHostingBoard     FeatureFlags = 1 << 0 // Board is a hosting board
	FeatureRequiresDaughter FeatureFlags = 1 << 1 // Board requires at least one daughter board
	FeatureRemovable        FeatureFlags = 1 << 2 // Board is removable
	FeatureReplaceable      FeatureFlags = 1 << 3 // Board is replaceable
	FeatureHotSwappable     FeatureFlags = 1 << 4 // Board is hot swappable
)

// Has checks if a feature flag is set
func (f FeatureFlags) Has(flag FeatureFlags) bool {
	return f&flag != 0
}

// IsHostingBoard returns true if this is the main hosting board (motherboard)
func (f FeatureFlags) IsHostingBoard() bool {
	return f.Has(FeatureHostingBoard)
}

// BoardType identifies the type of baseboard
type BoardType uint8

// Board type values
const (
	BoardTypeUnknown            BoardType = 0x01
	BoardTypeOther              BoardType = 0x02
	BoardTypeServerBlade        BoardType = 0x03
	BoardTypeConnectivitySwitch BoardType = 0x04
	BoardTypeSystemManagement   BoardType = 0x05
	BoardTypeProcessorModule    BoardType = 0x06
	BoardTypeIOModule           BoardType = 0x07
	BoardTypeMemoryModule       BoardType = 0x08
	BoardTypeDaughterBoard      BoardType = 0x09
	BoardTypeMotherboard        BoardType = 0x0A
	BoardTypeProcessorMemModule BoardType = 0x0B
	BoardTypeProcessorIOModule  BoardType = 0x0C
	BoardTypeInterconnectBoard  BoardType = 0x0D
)

// String returns a human-readable board type description
func (bt BoardType) String() string {
	switch bt {
	case BoardTypeUnknown:
		return "Unknown"
	case BoardTypeOther:
		return "Other"
	case BoardTypeServerBlade:
		return "Server Blade"
	case BoardTypeConnectivitySwitch:
		return "Connectivity Switch"
	case BoardTypeSystemManagement:
		return "System Management Module"
	case BoardTypeProcessorModule:
		return "Processor Module"
	case BoardTypeIOModule:
		return "I/O Module"
	case BoardTypeMemoryModule:
		return "Memory Module"
	case BoardTypeDaughterBoard:
		return "Daughter Board"
	case BoardTypeMotherboard:
		return "Motherboard"
	case BoardTypeProcessorMemModule:
		return "Processor/Memory Module"
	case BoardTypeProcessorIOModule:
		return "Processor/I/O Module"
	case BoardTypeInterconnectBoard:
		return "Interconnect Board"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(bt))
	}
}

// Parse parses a Baseboard Information structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*BaseboardInfo, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 8 bytes
	if len(s.Data) < 8 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &BaseboardInfo{
		Header:       s.Header,
		Manufacturer: s.GetString(s.GetByte(0x04)),
		Product:      s.GetString(s.GetByte(0x05)),
		Version:      s.GetString(s.GetByte(0x06)),
		SerialNumber: s.GetString(s.GetByte(0x07)),
	}

	// Asset Tag (optional)
	if len(s.Data) >= 9 {
		info.AssetTag = s.GetString(s.GetByte(0x08))
	}

	// Feature Flags
	if len(s.Data) >= 10 {
		info.FeatureFlags = FeatureFlags(s.GetByte(0x09))
	}

	// Location in Chassis
	if len(s.Data) >= 11 {
		info.LocationInChassis = s.GetString(s.GetByte(0x0A))
	}

	// Chassis Handle
	if len(s.Data) >= 13 {
		info.ChassisHandle = s.GetWord(0x0B)
	}

	// Board Type
	if len(s.Data) >= 14 {
		info.BoardType = BoardType(s.GetByte(0x0D))
	}

	// Number of Contained Object Handles
	if len(s.Data) >= 15 {
		info.NumberOfContainedHandles = s.GetByte(0x0E)

		// Parse contained handles
		handleOffset := 0x0F
		for i := uint8(0); i < info.NumberOfContainedHandles; i++ {
			if handleOffset+2 <= len(s.Data) {
				info.ContainedObjectHandles = append(info.ContainedObjectHandles, s.GetWord(handleOffset))
				handleOffset += 2
			}
		}
	}

	return info, nil
}

// Get retrieves the first Baseboard Information from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*BaseboardInfo, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Baseboard Information structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*BaseboardInfo, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var boards []*BaseboardInfo
	for i := range structures {
		board, err := Parse(&structures[i])
		if err == nil {
			boards = append(boards, board)
		}
	}

	if len(boards) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return boards, nil
}

// DisplayName returns a display-friendly baseboard name
func (b *BaseboardInfo) DisplayName() string {
	if b.Manufacturer != "" && b.Product != "" {
		return fmt.Sprintf("%s %s", b.Manufacturer, b.Product)
	}
	if b.Product != "" {
		return b.Product
	}
	return "Unknown Baseboard"
}
