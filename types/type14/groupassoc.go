// Package type14 implements SMBIOS Type 14 - Group Associations
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type14

import (
	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Group Associations
const StructureType uint8 = 14

// GroupAssociations represents Type 14 - Group Associations
type GroupAssociations struct {
	Header    gosmbios.Header
	GroupName string
	Items     []GroupItem
}

// GroupItem represents a single item in the group
type GroupItem struct {
	ItemType   uint8  // SMBIOS structure type
	ItemHandle uint16 // Handle of the structure
}

// Parse parses a Group Associations structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*GroupAssociations, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 5 bytes
	if len(s.Data) < 5 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &GroupAssociations{
		Header:    s.Header,
		GroupName: s.GetString(s.GetByte(0x04)),
	}

	// Read group items (3 bytes each: type + handle)
	// Number of items = (length - 5) / 3
	numItems := (int(s.Header.Length) - 5) / 3
	offset := 0x05
	for i := 0; i < numItems && offset+2 < len(s.Data); i++ {
		item := GroupItem{
			ItemType:   s.GetByte(offset),
			ItemHandle: s.GetWord(offset + 1),
		}
		info.Items = append(info.Items, item)
		offset += 3
	}

	return info, nil
}

// Get retrieves the first Group Associations from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*GroupAssociations, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Group Associations structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*GroupAssociations, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var groups []*GroupAssociations
	for i := range structures {
		grp, err := Parse(&structures[i])
		if err == nil {
			groups = append(groups, grp)
		}
	}

	if len(groups) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return groups, nil
}

// GetItemCount returns the number of items in the group
func (g *GroupAssociations) GetItemCount() int {
	return len(g.Items)
}
