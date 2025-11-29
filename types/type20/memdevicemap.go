// Package type20 implements SMBIOS Type 20 - Memory Device Mapped Address
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type20

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Memory Device Mapped Address
const StructureType uint8 = 20

// MemoryDeviceMappedAddress represents Type 20 - Memory Device Mapped Address
type MemoryDeviceMappedAddress struct {
	Header                        gosmbios.Header
	StartingAddress               uint32 // In KB
	EndingAddress                 uint32 // In KB
	MemoryDeviceHandle            uint16
	MemoryArrayMappedAddressHandle uint16
	PartitionRowPosition          uint8
	InterleavePosition            uint8
	InterleavedDataDepth          uint8
	ExtendedStartingAddress       uint64 // In bytes (SMBIOS 2.7+)
	ExtendedEndingAddress         uint64 // In bytes (SMBIOS 2.7+)
}

// Parse parses a Memory Device Mapped Address structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*MemoryDeviceMappedAddress, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 19 bytes
	if len(s.Data) < 19 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &MemoryDeviceMappedAddress{
		Header:                        s.Header,
		StartingAddress:               s.GetDWord(0x04),
		EndingAddress:                 s.GetDWord(0x08),
		MemoryDeviceHandle:            s.GetWord(0x0C),
		MemoryArrayMappedAddressHandle: s.GetWord(0x0E),
		PartitionRowPosition:          s.GetByte(0x10),
		InterleavePosition:            s.GetByte(0x11),
		InterleavedDataDepth:          s.GetByte(0x12),
	}

	// Extended addresses (SMBIOS 2.7+)
	if len(s.Data) >= 35 && info.StartingAddress == 0xFFFFFFFF {
		info.ExtendedStartingAddress = s.GetQWord(0x13)
		info.ExtendedEndingAddress = s.GetQWord(0x1B)
	}

	return info, nil
}

// Get retrieves the first Memory Device Mapped Address from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*MemoryDeviceMappedAddress, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Memory Device Mapped Address structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*MemoryDeviceMappedAddress, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var addrs []*MemoryDeviceMappedAddress
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
func (m *MemoryDeviceMappedAddress) GetStartingAddressBytes() uint64 {
	if m.StartingAddress == 0xFFFFFFFF {
		return m.ExtendedStartingAddress
	}
	return uint64(m.StartingAddress) * 1024
}

// GetEndingAddressBytes returns the ending address in bytes
func (m *MemoryDeviceMappedAddress) GetEndingAddressBytes() uint64 {
	if m.EndingAddress == 0xFFFFFFFF {
		return m.ExtendedEndingAddress
	}
	return uint64(m.EndingAddress) * 1024
}

// GetSizeBytes returns the size of the mapped region in bytes
func (m *MemoryDeviceMappedAddress) GetSizeBytes() uint64 {
	return m.GetEndingAddressBytes() - m.GetStartingAddressBytes() + 1
}

// GetSizeString returns a human-readable size string
func (m *MemoryDeviceMappedAddress) GetSizeString() string {
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

// PartitionRowPositionString returns a human-readable partition row position
func (m *MemoryDeviceMappedAddress) PartitionRowPositionString() string {
	if m.PartitionRowPosition == 0 {
		return "Unknown"
	}
	if m.PartitionRowPosition == 0xFF {
		return "Not Applicable"
	}
	return fmt.Sprintf("Row %d", m.PartitionRowPosition)
}

// InterleavePositionString returns a human-readable interleave position
func (m *MemoryDeviceMappedAddress) InterleavePositionString() string {
	if m.InterleavePosition == 0 {
		return "Not Interleaved"
	}
	if m.InterleavePosition == 0xFF {
		return "Unknown"
	}
	return fmt.Sprintf("Position %d", m.InterleavePosition)
}
