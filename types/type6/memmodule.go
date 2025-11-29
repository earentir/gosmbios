// Package type6 implements SMBIOS Type 6 - Memory Module Information (Obsolete)
// Per DSP0134 SMBIOS Reference Specification 3.9.0
// This structure is obsolete starting with SMBIOS 2.1 specification
package type6

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Memory Module Information
const StructureType uint8 = 6

// MemoryModule represents Type 6 - Memory Module Information (Obsolete)
type MemoryModule struct {
	Header            gosmbios.Header
	SocketDesignation string
	BankConnections   uint8
	CurrentSpeed      uint8 // In nanoseconds
	CurrentMemoryType MemoryType
	InstalledSize     MemorySize
	EnabledSize       MemorySize
	ErrorStatus       ErrorStatus
}

// MemoryType represents memory module types
type MemoryType uint16

const (
	MemTypeOther         MemoryType = 1 << 0
	MemTypeUnknown       MemoryType = 1 << 1
	MemTypeStandard      MemoryType = 1 << 2
	MemTypeFastPageMode  MemoryType = 1 << 3
	MemTypeEDO           MemoryType = 1 << 4
	MemTypeParity        MemoryType = 1 << 5
	MemTypeECC           MemoryType = 1 << 6
	MemTypeSIMM          MemoryType = 1 << 7
	MemTypeDIMM          MemoryType = 1 << 8
	MemTypeBurstEDO      MemoryType = 1 << 9
	MemTypeSDRAM         MemoryType = 1 << 10
)

func (m MemoryType) String() string {
	var types []string
	if m&MemTypeOther != 0 {
		types = append(types, "Other")
	}
	if m&MemTypeUnknown != 0 {
		types = append(types, "Unknown")
	}
	if m&MemTypeStandard != 0 {
		types = append(types, "Standard")
	}
	if m&MemTypeFastPageMode != 0 {
		types = append(types, "Fast Page Mode")
	}
	if m&MemTypeEDO != 0 {
		types = append(types, "EDO")
	}
	if m&MemTypeParity != 0 {
		types = append(types, "Parity")
	}
	if m&MemTypeECC != 0 {
		types = append(types, "ECC")
	}
	if m&MemTypeSIMM != 0 {
		types = append(types, "SIMM")
	}
	if m&MemTypeDIMM != 0 {
		types = append(types, "DIMM")
	}
	if m&MemTypeBurstEDO != 0 {
		types = append(types, "Burst EDO")
	}
	if m&MemTypeSDRAM != 0 {
		types = append(types, "SDRAM")
	}
	if len(types) == 0 {
		return "None"
	}
	return fmt.Sprintf("%v", types)
}

// MemorySize represents memory module size
type MemorySize uint8

const (
	MemSizeNotInstalled MemorySize = 0x7D
	MemSizeNotEnabled   MemorySize = 0x7E
	MemSizeNotDetermined MemorySize = 0x7F
)

func (m MemorySize) String() string {
	switch m {
	case MemSizeNotInstalled:
		return "Not Installed"
	case MemSizeNotEnabled:
		return "Not Enabled"
	case MemSizeNotDetermined:
		return "Not Determinable"
	default:
		// Size = 2^n MB, where n is bits 0-6
		// Bit 7 indicates double-bank connection
		size := m & 0x7F
		doubleBank := m&0x80 != 0
		mb := uint64(1) << size
		if doubleBank {
			return fmt.Sprintf("%d MB (Double-bank)", mb*2)
		}
		return fmt.Sprintf("%d MB", mb)
	}
}

// SizeMB returns the size in MB, 0 if not installed/determinable
func (m MemorySize) SizeMB() uint64 {
	if m == MemSizeNotInstalled || m == MemSizeNotEnabled || m == MemSizeNotDetermined {
		return 0
	}
	size := m & 0x7F
	doubleBank := m&0x80 != 0
	mb := uint64(1) << size
	if doubleBank {
		return mb * 2
	}
	return mb
}

// ErrorStatus represents memory module error status
type ErrorStatus uint8

const (
	ErrStatusUncorrectable  ErrorStatus = 1 << 0
	ErrStatusCorrectable    ErrorStatus = 1 << 1
	ErrStatusEventLogged    ErrorStatus = 1 << 2
)

func (e ErrorStatus) String() string {
	if e == 0 {
		return "OK"
	}
	var status []string
	if e&ErrStatusUncorrectable != 0 {
		status = append(status, "Uncorrectable errors")
	}
	if e&ErrStatusCorrectable != 0 {
		status = append(status, "Correctable errors")
	}
	if e&ErrStatusEventLogged != 0 {
		status = append(status, "Event logged")
	}
	return fmt.Sprintf("%v", status)
}

// Parse parses a Memory Module Information structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*MemoryModule, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 12 bytes
	if len(s.Data) < 12 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &MemoryModule{
		Header:            s.Header,
		SocketDesignation: s.GetString(s.GetByte(0x04)),
		BankConnections:   s.GetByte(0x05),
		CurrentSpeed:      s.GetByte(0x06),
		CurrentMemoryType: MemoryType(s.GetWord(0x07)),
		InstalledSize:     MemorySize(s.GetByte(0x09)),
		EnabledSize:       MemorySize(s.GetByte(0x0A)),
		ErrorStatus:       ErrorStatus(s.GetByte(0x0B)),
	}

	return info, nil
}

// Get retrieves the first Memory Module Information from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*MemoryModule, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Memory Module structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*MemoryModule, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var modules []*MemoryModule
	for i := range structures {
		mod, err := Parse(&structures[i])
		if err == nil {
			modules = append(modules, mod)
		}
	}

	if len(modules) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return modules, nil
}

// BankConnectionString returns a human-readable bank connection string
func (m *MemoryModule) BankConnectionString() string {
	if m.BankConnections == 0xFF {
		return "No connection"
	}
	rAS0 := m.BankConnections & 0x0F
	rAS1 := (m.BankConnections >> 4) & 0x0F
	if rAS1 == 0x0F {
		return fmt.Sprintf("RAS %d", rAS0)
	}
	return fmt.Sprintf("RAS %d and %d", rAS0, rAS1)
}

// IsInstalled returns true if the memory module is installed
func (m *MemoryModule) IsInstalled() bool {
	return m.InstalledSize != MemSizeNotInstalled
}
