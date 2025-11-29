// Package type46 implements SMBIOS Type 46 - String Property
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type46

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for String Property
const StructureType uint8 = 46

// StringProperty represents Type 46 - String Property
type StringProperty struct {
	Header              gosmbios.Header
	StringPropertyID    StringPropertyID
	StringPropertyValue string
	ParentHandle        uint16
}

// StringPropertyID identifies the string property type
type StringPropertyID uint16

const (
	StringPropertyReserved         StringPropertyID = 0x0000
	StringPropertyUEFIDevicePath   StringPropertyID = 0x0001
)

func (s StringPropertyID) String() string {
	switch s {
	case StringPropertyReserved:
		return "Reserved"
	case StringPropertyUEFIDevicePath:
		return "UEFI Device Path"
	default:
		if s >= 0x8000 {
			return fmt.Sprintf("BIOS Vendor Specific (0x%04X)", uint16(s))
		}
		return fmt.Sprintf("Unknown (0x%04X)", uint16(s))
	}
}

// Parse parses a String Property structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*StringProperty, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 9 bytes
	if len(s.Data) < 9 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &StringProperty{
		Header:              s.Header,
		StringPropertyID:    StringPropertyID(s.GetWord(0x04)),
		StringPropertyValue: s.GetString(s.GetByte(0x06)),
		ParentHandle:        s.GetWord(0x07),
	}

	return info, nil
}

// Get retrieves the first String Property from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*StringProperty, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all String Property structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*StringProperty, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var props []*StringProperty
	for i := range structures {
		prop, err := Parse(&structures[i])
		if err == nil {
			props = append(props, prop)
		}
	}

	if len(props) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return props, nil
}

// GetByID retrieves all String Properties with a specific ID from SMBIOS data
func GetByID(sm *gosmbios.SMBIOS, id StringPropertyID) ([]*StringProperty, error) {
	all, err := GetAll(sm)
	if err != nil {
		return nil, err
	}

	var filtered []*StringProperty
	for _, prop := range all {
		if prop.StringPropertyID == id {
			filtered = append(filtered, prop)
		}
	}

	if len(filtered) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return filtered, nil
}
