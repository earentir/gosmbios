// Package type12 implements SMBIOS Type 12 - System Configuration Options
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type12

import (
	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for System Configuration Options
const StructureType uint8 = 12

// SystemConfigOptions represents Type 12 - System Configuration Options
type SystemConfigOptions struct {
	Header  gosmbios.Header
	Count   uint8
	Options []string
}

// Parse parses a System Configuration Options structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*SystemConfigOptions, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 5 bytes
	if len(s.Data) < 5 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &SystemConfigOptions{
		Header:  s.Header,
		Count:   s.GetByte(0x04),
		Options: s.Strings,
	}

	return info, nil
}

// Get retrieves the System Configuration Options from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*SystemConfigOptions, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all System Configuration Options structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*SystemConfigOptions, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var configs []*SystemConfigOptions
	for i := range structures {
		cfg, err := Parse(&structures[i])
		if err == nil {
			configs = append(configs, cfg)
		}
	}

	if len(configs) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return configs, nil
}
