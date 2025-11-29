// Package type35 implements SMBIOS Type 35 - Management Device Component
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type35

import (
	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Management Device Component
const StructureType uint8 = 35

// ManagementDeviceComponent represents Type 35 - Management Device Component
type ManagementDeviceComponent struct {
	Header                       gosmbios.Header
	Description                  string
	ManagementDeviceHandle       uint16
	ComponentHandle              uint16
	ThresholdHandle              uint16
}

// Parse parses a Management Device Component structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*ManagementDeviceComponent, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 11 bytes
	if len(s.Data) < 11 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &ManagementDeviceComponent{
		Header:                 s.Header,
		Description:            s.GetString(s.GetByte(0x04)),
		ManagementDeviceHandle: s.GetWord(0x05),
		ComponentHandle:        s.GetWord(0x07),
		ThresholdHandle:        s.GetWord(0x09),
	}

	return info, nil
}

// Get retrieves the first Management Device Component from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*ManagementDeviceComponent, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Management Device Component structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*ManagementDeviceComponent, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var components []*ManagementDeviceComponent
	for i := range structures {
		comp, err := Parse(&structures[i])
		if err == nil {
			components = append(components, comp)
		}
	}

	if len(components) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return components, nil
}

// HasThreshold returns true if a threshold handle is associated
func (m *ManagementDeviceComponent) HasThreshold() bool {
	return m.ThresholdHandle != 0xFFFF
}
