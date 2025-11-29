package gosmbios

import (
	"encoding/binary"
)

// ParseStructures parses raw SMBIOS table data into individual structures
func ParseStructures(tableData []byte, maxStructures int) ([]Structure, error) {
	var structures []Structure
	offset := 0

	for offset < len(tableData) {
		// Check if we have enough data for the header
		if offset+4 > len(tableData) {
			break
		}

		// Parse header
		header := Header{
			Type:   tableData[offset],
			Length: tableData[offset+1],
			Handle: binary.LittleEndian.Uint16(tableData[offset+2:]),
		}

		// End-of-Table structure (Type 127)
		if header.Type == 127 {
			structures = append(structures, Structure{
				Header:  header,
				Data:    tableData[offset : offset+int(header.Length)],
				Strings: nil,
			})
			break
		}

		// Validate length
		if header.Length < 4 {
			break
		}

		// Check if we have the full formatted section
		if offset+int(header.Length) > len(tableData) {
			break
		}

		// Extract formatted section data
		formattedSection := make([]byte, header.Length)
		copy(formattedSection, tableData[offset:offset+int(header.Length)])

		// Parse string table
		stringStart := offset + int(header.Length)
		strings, stringEnd := parseStringTable(tableData, stringStart)

		structures = append(structures, Structure{
			Header:  header,
			Data:    formattedSection,
			Strings: strings,
		})

		offset = stringEnd

		// Safety check for maxStructures (0 means no limit)
		if maxStructures > 0 && len(structures) >= maxStructures {
			break
		}
	}

	return structures, nil
}

// parseStringTable parses the null-terminated string table following a structure
// Returns the strings and the offset after the string table (after double-null terminator)
// Per SMBIOS spec: strings are null-terminated, table ends with additional null (double-null)
// Empty string table is just \0\0 (two consecutive nulls)
func parseStringTable(data []byte, start int) ([]string, int) {
	var strings []string
	current := start

	for current < len(data) {
		// Check for null byte - indicates either empty table or end of strings
		if data[current] == 0 {
			// Found null - this could be:
			// 1. First null of empty string table (\0\0) - skip to current+2
			// 2. Second null after last string's terminator - skip to current+1
			// Case 2 only happens after we've parsed at least one string
			if len(strings) == 0 {
				// Empty string table: \0\0 - skip both nulls
				return strings, current + 2
			}
			// End of string table after last string
			return strings, current + 1
		}

		// Find end of current string (look for null terminator)
		end := current
		for end < len(data) && data[end] != 0 {
			end++
		}

		if end > current {
			strings = append(strings, string(data[current:end]))
		}

		// Move past the null terminator of this string
		current = end + 1
	}

	return strings, current
}

// ParseEntryPoint32 parses a 32-bit SMBIOS entry point (_SM_)
func ParseEntryPoint32(data []byte) (*EntryPoint, error) {
	if len(data) < 31 {
		return nil, ErrInvalidStructure
	}

	// Verify anchor string "_SM_"
	if string(data[0:4]) != "_SM_" {
		return nil, ErrNotFound
	}

	// Verify checksum
	epLength := data[5]
	if epLength > byte(len(data)) {
		return nil, ErrInvalidStructure
	}

	var checksum uint8
	for i := uint8(0); i < epLength; i++ {
		checksum += data[i]
	}
	if checksum != 0 {
		return nil, ErrInvalidChecksum
	}

	// Verify intermediate anchor "_DMI_"
	if string(data[16:21]) != "_DMI_" {
		return nil, ErrNotFound
	}

	// Verify intermediate checksum
	var intermediateChecksum uint8
	for i := 16; i < 31; i++ {
		intermediateChecksum += data[i]
	}
	if intermediateChecksum != 0 {
		return nil, ErrInvalidChecksum
	}

	return &EntryPoint{
		Type:             EntryPoint32Bit,
		MajorVersion:     data[6],
		MinorVersion:     data[7],
		TableLength:      binary.LittleEndian.Uint32(data[22:26]) & 0xFFFF, // 16-bit value
		TableAddress:     uint64(binary.LittleEndian.Uint32(data[24:28])),
		StructureCount:   binary.LittleEndian.Uint16(data[28:30]),
		BCDRevision:      data[30],
		EntryPointLength: epLength,
	}, nil
}

// ParseEntryPoint64 parses a 64-bit SMBIOS entry point (_SM3_)
func ParseEntryPoint64(data []byte) (*EntryPoint, error) {
	if len(data) < 24 {
		return nil, ErrInvalidStructure
	}

	// Verify anchor string "_SM3_"
	if string(data[0:5]) != "_SM3_" {
		return nil, ErrNotFound
	}

	// Verify checksum
	epLength := data[6]
	if epLength > byte(len(data)) {
		return nil, ErrInvalidStructure
	}

	var checksum uint8
	for i := uint8(0); i < epLength; i++ {
		checksum += data[i]
	}
	if checksum != 0 {
		return nil, ErrInvalidChecksum
	}

	return &EntryPoint{
		Type:             EntryPoint64Bit,
		MajorVersion:     data[7],
		MinorVersion:     data[8],
		Revision:         data[9],
		EntryPointLength: epLength,
		TableMaxSize:     binary.LittleEndian.Uint32(data[12:16]),
		TableAddress:     binary.LittleEndian.Uint64(data[16:24]),
	}, nil
}
