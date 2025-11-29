// Package type13 implements SMBIOS Type 13 - BIOS Language Information
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type13

import (
	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for BIOS Language Information
const StructureType uint8 = 13

// BIOSLanguage represents Type 13 - BIOS Language Information
type BIOSLanguage struct {
	Header              gosmbios.Header
	InstallableLanguages uint8
	Flags               LanguageFlags
	Reserved            [15]byte
	CurrentLanguage     string
	Languages           []string
}

// LanguageFlags represents BIOS language format flags
type LanguageFlags uint8

// IsAbbreviatedFormat returns true if languages are in abbreviated format
// false means long format (e.g., "en|US|iso8859-1" vs "English")
func (f LanguageFlags) IsAbbreviatedFormat() bool {
	return f&0x01 != 0
}

// String returns the format description
func (f LanguageFlags) String() string {
	if f.IsAbbreviatedFormat() {
		return "Abbreviated"
	}
	return "Long"
}

// Parse parses a BIOS Language Information structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*BIOSLanguage, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 22 bytes
	if len(s.Data) < 22 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &BIOSLanguage{
		Header:              s.Header,
		InstallableLanguages: s.GetByte(0x04),
		Flags:               LanguageFlags(s.GetByte(0x05)),
		CurrentLanguage:     s.GetString(s.GetByte(0x15)),
		Languages:          s.Strings,
	}

	// Copy reserved bytes
	copy(info.Reserved[:], s.Data[0x06:0x15])

	return info, nil
}

// Get retrieves the BIOS Language Information from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*BIOSLanguage, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetCurrentLanguage returns the currently selected BIOS language
func (b *BIOSLanguage) GetCurrentLanguage() string {
	return b.CurrentLanguage
}

// GetAllLanguages returns all installable languages
func (b *BIOSLanguage) GetAllLanguages() []string {
	return b.Languages
}
