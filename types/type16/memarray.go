// Package type16 implements SMBIOS Type 16 - Physical Memory Array
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type16

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Physical Memory Array
const StructureType uint8 = 16

// MemoryArray represents Type 16 - Physical Memory Array
type MemoryArray struct {
	Header                  gosmbios.Header
	Location                MemoryArrayLocation
	Use                     MemoryArrayUse
	ErrorCorrection         MemoryArrayErrorCorrection
	MaximumCapacity         uint64 // In KB
	ErrorInformationHandle  uint16
	NumberOfMemoryDevices   uint16
	ExtendedMaximumCapacity uint64 // SMBIOS 2.7+ (in bytes)
}

// MemoryArrayLocation identifies where the memory array is located
type MemoryArrayLocation uint8

// Memory array location values
const (
	LocationOther                 MemoryArrayLocation = 0x01
	LocationUnknown               MemoryArrayLocation = 0x02
	LocationSystemBoard           MemoryArrayLocation = 0x03
	LocationISAAddonCard          MemoryArrayLocation = 0x04
	LocationEISAAddonCard         MemoryArrayLocation = 0x05
	LocationPCIAddonCard          MemoryArrayLocation = 0x06
	LocationMCAAddonCard          MemoryArrayLocation = 0x07
	LocationPCMCIAAddonCard       MemoryArrayLocation = 0x08
	LocationProprietaryAddonCard  MemoryArrayLocation = 0x09
	LocationNuBus                 MemoryArrayLocation = 0x0A
	LocationPC98C20AddonCard      MemoryArrayLocation = 0xA0
	LocationPC98C24AddonCard      MemoryArrayLocation = 0xA1
	LocationPC98EAddonCard        MemoryArrayLocation = 0xA2
	LocationPC98LocalBusAddonCard MemoryArrayLocation = 0xA3
	LocationCXLAddonCard          MemoryArrayLocation = 0xA4
)

// String returns a human-readable location description
func (l MemoryArrayLocation) String() string {
	locations := map[MemoryArrayLocation]string{
		LocationOther:                 "Other",
		LocationUnknown:               "Unknown",
		LocationSystemBoard:           "System Board or Motherboard",
		LocationISAAddonCard:          "ISA Add-on Card",
		LocationEISAAddonCard:         "EISA Add-on Card",
		LocationPCIAddonCard:          "PCI Add-on Card",
		LocationMCAAddonCard:          "MCA Add-on Card",
		LocationPCMCIAAddonCard:       "PCMCIA Add-on Card",
		LocationProprietaryAddonCard:  "Proprietary Add-on Card",
		LocationNuBus:                 "NuBus",
		LocationPC98C20AddonCard:      "PC-98/C20 Add-on Card",
		LocationPC98C24AddonCard:      "PC-98/C24 Add-on Card",
		LocationPC98EAddonCard:        "PC-98/E Add-on Card",
		LocationPC98LocalBusAddonCard: "PC-98/Local Bus Add-on Card",
		LocationCXLAddonCard:          "CXL Add-on Card",
	}

	if name, ok := locations[l]; ok {
		return name
	}
	return fmt.Sprintf("Unknown (0x%02X)", uint8(l))
}

// MemoryArrayUse identifies the function for which the array is used
type MemoryArrayUse uint8

// Memory array use values
const (
	UseOther          MemoryArrayUse = 0x01
	UseUnknown        MemoryArrayUse = 0x02
	UseSystemMemory   MemoryArrayUse = 0x03
	UseVideoMemory    MemoryArrayUse = 0x04
	UseFlashMemory    MemoryArrayUse = 0x05
	UseNonVolatileRAM MemoryArrayUse = 0x06
	UseCacheMemory    MemoryArrayUse = 0x07
)

// String returns a human-readable use description
func (u MemoryArrayUse) String() string {
	switch u {
	case UseOther:
		return "Other"
	case UseUnknown:
		return "Unknown"
	case UseSystemMemory:
		return "System Memory"
	case UseVideoMemory:
		return "Video Memory"
	case UseFlashMemory:
		return "Flash Memory"
	case UseNonVolatileRAM:
		return "Non-volatile RAM"
	case UseCacheMemory:
		return "Cache Memory"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(u))
	}
}

// MemoryArrayErrorCorrection identifies the error correction type
type MemoryArrayErrorCorrection uint8

// Error correction values
const (
	ErrorCorrectionOther        MemoryArrayErrorCorrection = 0x01
	ErrorCorrectionUnknown      MemoryArrayErrorCorrection = 0x02
	ErrorCorrectionNone         MemoryArrayErrorCorrection = 0x03
	ErrorCorrectionParity       MemoryArrayErrorCorrection = 0x04
	ErrorCorrectionSingleBitECC MemoryArrayErrorCorrection = 0x05
	ErrorCorrectionMultiBitECC  MemoryArrayErrorCorrection = 0x06
	ErrorCorrectionCRC          MemoryArrayErrorCorrection = 0x07
)

// String returns a human-readable error correction description
func (ec MemoryArrayErrorCorrection) String() string {
	switch ec {
	case ErrorCorrectionOther:
		return "Other"
	case ErrorCorrectionUnknown:
		return "Unknown"
	case ErrorCorrectionNone:
		return "None"
	case ErrorCorrectionParity:
		return "Parity"
	case ErrorCorrectionSingleBitECC:
		return "Single-bit ECC"
	case ErrorCorrectionMultiBitECC:
		return "Multi-bit ECC"
	case ErrorCorrectionCRC:
		return "CRC"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(ec))
	}
}

// Parse parses a Physical Memory Array structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*MemoryArray, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 15 bytes (SMBIOS 2.1)
	if len(s.Data) < 15 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &MemoryArray{
		Header:                 s.Header,
		Location:               MemoryArrayLocation(s.GetByte(0x04)),
		Use:                    MemoryArrayUse(s.GetByte(0x05)),
		ErrorCorrection:        MemoryArrayErrorCorrection(s.GetByte(0x06)),
		MaximumCapacity:        uint64(s.GetDWord(0x07)),
		ErrorInformationHandle: s.GetWord(0x0B),
		NumberOfMemoryDevices:  s.GetWord(0x0D),
	}

	// SMBIOS 2.7+ Extended Maximum Capacity
	if len(s.Data) >= 23 && info.MaximumCapacity == 0x80000000 {
		info.ExtendedMaximumCapacity = s.GetQWord(0x0F)
		// Convert to KB for consistency
		info.MaximumCapacity = info.ExtendedMaximumCapacity / 1024
	}

	return info, nil
}

// Get retrieves the first Physical Memory Array from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*MemoryArray, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Physical Memory Array structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*MemoryArray, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var arrays []*MemoryArray
	for i := range structures {
		arr, err := Parse(&structures[i])
		if err == nil {
			arrays = append(arrays, arr)
		}
	}

	if len(arrays) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return arrays, nil
}

// MaximumCapacityString returns a human-readable maximum capacity
func (m *MemoryArray) MaximumCapacityString() string {
	if m.MaximumCapacity == 0 {
		return "Unknown"
	}

	kb := m.MaximumCapacity
	if kb >= 1024*1024*1024 {
		return fmt.Sprintf("%d TB", kb/(1024*1024*1024))
	}
	if kb >= 1024*1024 {
		return fmt.Sprintf("%d GB", kb/(1024*1024))
	}
	if kb >= 1024 {
		return fmt.Sprintf("%d MB", kb/1024)
	}
	return fmt.Sprintf("%d KB", kb)
}

// IsSystemMemory returns true if this is the main system memory array
func (m *MemoryArray) IsSystemMemory() bool {
	return m.Use == UseSystemMemory && m.Location == LocationSystemBoard
}
