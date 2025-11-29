// Package type40 implements SMBIOS Type 40 - Additional Information
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type40

import (
	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Additional Information
const StructureType uint8 = 40

// AdditionalInformation represents Type 40 - Additional Information
type AdditionalInformation struct {
	Header               gosmbios.Header
	NumberOfEntries      uint8
	Entries              []AdditionalEntry
}

// AdditionalEntry represents an additional information entry
type AdditionalEntry struct {
	EntryLength          uint8
	ReferencedHandle     uint16
	ReferencedOffset     uint8
	String               string
	Value                []byte
}

// Parse parses an Additional Information structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*AdditionalInformation, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 5 bytes
	if len(s.Data) < 5 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &AdditionalInformation{
		Header:          s.Header,
		NumberOfEntries: s.GetByte(0x04),
	}

	// Parse entries
	offset := 0x05
	for i := uint8(0); i < info.NumberOfEntries; i++ {
		if offset >= len(s.Data) {
			break
		}

		entryLength := s.GetByte(offset)
		if entryLength < 6 || offset+int(entryLength) > len(s.Data) {
			break
		}

		entry := AdditionalEntry{
			EntryLength:      entryLength,
			ReferencedHandle: s.GetWord(offset + 1),
			ReferencedOffset: s.GetByte(offset + 3),
			String:           s.GetString(s.GetByte(offset + 4)),
		}

		// Copy value bytes (remaining bytes after the string index)
		valueLen := int(entryLength) - 5
		if valueLen > 0 {
			entry.Value = make([]byte, valueLen)
			copy(entry.Value, s.Data[offset+5:offset+int(entryLength)])
		}

		info.Entries = append(info.Entries, entry)
		offset += int(entryLength)
	}

	return info, nil
}

// Get retrieves the first Additional Information from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*AdditionalInformation, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Additional Information structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*AdditionalInformation, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var infos []*AdditionalInformation
	for i := range structures {
		info, err := Parse(&structures[i])
		if err == nil {
			infos = append(infos, info)
		}
	}

	if len(infos) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return infos, nil
}
