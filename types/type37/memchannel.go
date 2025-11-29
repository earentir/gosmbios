// Package type37 implements SMBIOS Type 37 - Memory Channel
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type37

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Memory Channel
const StructureType uint8 = 37

// MemoryChannel represents Type 37 - Memory Channel
type MemoryChannel struct {
	Header              gosmbios.Header
	ChannelType         ChannelType
	MaximumChannelLoad  uint8
	MemoryDeviceCount   uint8
	MemoryDevices       []MemoryDeviceInfo
}

// MemoryDeviceInfo represents information about a memory device in the channel
type MemoryDeviceInfo struct {
	MemoryDeviceLoad   uint8
	MemoryDeviceHandle uint16
}

// ChannelType identifies the memory channel type
type ChannelType uint8

const (
	ChannelTypeOther         ChannelType = 0x01
	ChannelTypeUnknown       ChannelType = 0x02
	ChannelTypeRamBus        ChannelType = 0x03
	ChannelTypeSyncLink      ChannelType = 0x04
)

func (c ChannelType) String() string {
	switch c {
	case ChannelTypeOther:
		return "Other"
	case ChannelTypeUnknown:
		return "Unknown"
	case ChannelTypeRamBus:
		return "RamBus"
	case ChannelTypeSyncLink:
		return "SyncLink"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(c))
	}
}

// Parse parses a Memory Channel structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*MemoryChannel, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 7 bytes
	if len(s.Data) < 7 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &MemoryChannel{
		Header:             s.Header,
		ChannelType:        ChannelType(s.GetByte(0x04)),
		MaximumChannelLoad: s.GetByte(0x05),
		MemoryDeviceCount:  s.GetByte(0x06),
	}

	// Read memory device entries (3 bytes each: load + handle)
	offset := 0x07
	for i := uint8(0); i < info.MemoryDeviceCount; i++ {
		if offset+2 >= len(s.Data) {
			break
		}
		device := MemoryDeviceInfo{
			MemoryDeviceLoad:   s.GetByte(offset),
			MemoryDeviceHandle: s.GetWord(offset + 1),
		}
		info.MemoryDevices = append(info.MemoryDevices, device)
		offset += 3
	}

	return info, nil
}

// Get retrieves the first Memory Channel from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*MemoryChannel, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Memory Channel structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*MemoryChannel, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var channels []*MemoryChannel
	for i := range structures {
		ch, err := Parse(&structures[i])
		if err == nil {
			channels = append(channels, ch)
		}
	}

	if len(channels) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return channels, nil
}
