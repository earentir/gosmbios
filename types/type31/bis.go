// Package type31 implements SMBIOS Type 31 - Boot Integrity Services Entry Point
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type31

import (
	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Boot Integrity Services Entry Point
const StructureType uint8 = 31

// BISEntryPoint represents Type 31 - Boot Integrity Services Entry Point
type BISEntryPoint struct {
	Header        gosmbios.Header
	Checksum      uint8
	Reserved1     uint8
	Reserved2     uint16
	BISEntryPoint uint32 // Physical address of BIS entry point
}

// Parse parses a Boot Integrity Services Entry Point structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*BISEntryPoint, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Structure length is variable based on the BIS Entry Point structure
	// Minimum is 28 bytes per spec, but the SMBIOS structure itself may be smaller
	if len(s.Data) < 4 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &BISEntryPoint{
		Header: s.Header,
	}

	// The BIS Entry Point structure is 28 bytes starting at offset 4
	if len(s.Data) >= 8 {
		info.Checksum = s.GetByte(0x04)
		info.Reserved1 = s.GetByte(0x05)
		info.Reserved2 = s.GetWord(0x06)
	}

	if len(s.Data) >= 12 {
		info.BISEntryPoint = s.GetDWord(0x08)
	}

	return info, nil
}

// Get retrieves the Boot Integrity Services Entry Point from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*BISEntryPoint, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}
