// Package type11 implements SMBIOS Type 11 - OEM Strings
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type11

import (
	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for OEM Strings
const StructureType uint8 = 11

// OEMStrings represents Type 11 - OEM Strings
type OEMStrings struct {
	Header  gosmbios.Header
	Count   uint8
	Strings []string
}

// Parse parses an OEM Strings structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*OEMStrings, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 5 bytes
	if len(s.Data) < 5 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &OEMStrings{
		Header:  s.Header,
		Count:   s.GetByte(0x04),
		Strings: s.Strings,
	}

	return info, nil
}

// Get retrieves the OEM Strings from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*OEMStrings, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all OEM Strings structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*OEMStrings, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var oemStrings []*OEMStrings
	for i := range structures {
		oem, err := Parse(&structures[i])
		if err == nil {
			oemStrings = append(oemStrings, oem)
		}
	}

	if len(oemStrings) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return oemStrings, nil
}
