// Package type19 implements SMBIOS Type 19 - Memory Array Mapped Address
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type19

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Memory Array Mapped Address
const StructureType uint8 = 19

// MemoryArrayMappedAddress represents Type 19 - Memory Array Mapped Address
type MemoryArrayMappedAddress struct {
	Header                    gosmbios.Header
	StartingAddress           uint32 // In KB
	EndingAddress             uint32 // In KB
	MemoryArrayHandle         uint16
	PartitionWidth            uint8
	ExtendedStartingAddress   uint64 // In bytes (SMBIOS 2.7+)
	ExtendedEndingAddress     uint64 // In bytes (SMBIOS 2.7+)
}

// Parse parses a Memory Array Mapped Address structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*MemoryArrayMappedAddress, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 15 bytes
	if len(s.Data) < 15 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &MemoryArrayMappedAddress{
		Header:            s.Header,
		StartingAddress:   s.GetDWord(0x04),
		EndingAddress:     s.GetDWord(0x08),
		MemoryArrayHandle: s.GetWord(0x0C),
		PartitionWidth:    s.GetByte(0x0E),
	}

	// Extended addresses (SMBIOS 2.7+)
	if len(s.Data) >= 31 && info.StartingAddress == 0xFFFFFFFF {
		info.ExtendedStartingAddress = s.GetQWord(0x0F)
		info.ExtendedEndingAddress = s.GetQWord(0x17)
	}

	return info, nil
}

// Get retrieves the first Memory Array Mapped Address from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*MemoryArrayMappedAddress, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Memory Array Mapped Address structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*MemoryArrayMappedAddress, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var addrs []*MemoryArrayMappedAddress
	for i := range structures {
		addr, err := Parse(&structures[i])
		if err == nil {
			addrs = append(addrs, addr)
		}
	}

	if len(addrs) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return addrs, nil
}

// GetStartingAddressBytes returns the starting address in bytes
func (m *MemoryArrayMappedAddress) GetStartingAddressBytes() uint64 {
	if m.StartingAddress == 0xFFFFFFFF {
		return m.ExtendedStartingAddress
	}
	return uint64(m.StartingAddress) * 1024
}

// GetEndingAddressBytes returns the ending address in bytes
func (m *MemoryArrayMappedAddress) GetEndingAddressBytes() uint64 {
	if m.EndingAddress == 0xFFFFFFFF {
		return m.ExtendedEndingAddress
	}
	return uint64(m.EndingAddress) * 1024
}

// GetSizeBytes returns the size of the mapped region in bytes
func (m *MemoryArrayMappedAddress) GetSizeBytes() uint64 {
	return m.GetEndingAddressBytes() - m.GetStartingAddressBytes() + 1
}

// GetSizeString returns a human-readable size string
func (m *MemoryArrayMappedAddress) GetSizeString() string {
	size := m.GetSizeBytes()
	if size >= 1024*1024*1024*1024 {
		return fmt.Sprintf("%d TB", size/(1024*1024*1024*1024))
	}
	if size >= 1024*1024*1024 {
		return fmt.Sprintf("%d GB", size/(1024*1024*1024))
	}
	if size >= 1024*1024 {
		return fmt.Sprintf("%d MB", size/(1024*1024))
	}
	if size >= 1024 {
		return fmt.Sprintf("%d KB", size/1024)
	}
	return fmt.Sprintf("%d bytes", size)
}
