// Package type127 implements SMBIOS Type 127 - End-of-Table
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type127

import (
	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for End-of-Table
const StructureType uint8 = 127

// EndOfTable represents Type 127 - End-of-Table
// This structure marks the end of the SMBIOS structure table
type EndOfTable struct {
	Header gosmbios.Header
}

// Parse parses an End-of-Table structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*EndOfTable, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	return &EndOfTable{
		Header: s.Header,
	}, nil
}

// Get retrieves the End-of-Table marker from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*EndOfTable, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}
