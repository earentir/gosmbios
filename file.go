package gosmbios

import (
	"bytes"
	"encoding/binary"
	"os"
)

// File format magic bytes
const (
	fileMagic   = "SMBIOSRAW"
	fileVersion = 1
)

// RawFileHeader represents the header of a raw SMBIOS dump file
// This is a simple header followed by the raw SMBIOS table bytes
type RawFileHeader struct {
	Magic          [9]byte // "SMBIOSRAW"
	Version        uint8   // File format version (1)
	EntryPointType uint8   // 0 = 32-bit, 1 = 64-bit
	MajorVersion   uint8   // SMBIOS major version
	MinorVersion   uint8   // SMBIOS minor version
	Revision       uint8   // SMBIOS revision (3.x only)
	Reserved       uint8   // Padding
	TableLength    uint32  // Length of raw table data that follows
	TableAddress   uint64  // Original table address (for reference)
}

// readSMBIOSFromFile reads SMBIOS data from a raw dump file
func readSMBIOSFromFile(filename string) (*SMBIOS, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Check minimum size for header (9 + 1 + 1 + 1 + 1 + 1 + 1 + 4 + 8 = 28 bytes)
	headerSize := 28
	if len(data) < headerSize {
		return nil, ErrInvalidStructure
	}

	// Check magic
	if string(data[0:9]) != fileMagic {
		return nil, ErrInvalidStructure
	}

	// Parse header manually for cross-platform compatibility
	version := data[9]
	if version != fileVersion {
		return nil, ErrInvalidStructure
	}

	entryPointType := data[10]
	majorVersion := data[11]
	minorVersion := data[12]
	revision := data[13]
	// reserved at 14
	tableLength := binary.LittleEndian.Uint32(data[15:19])
	tableAddress := binary.LittleEndian.Uint64(data[19:27])

	// Verify we have the complete table
	if len(data) < headerSize+int(tableLength) {
		return nil, ErrInvalidStructure
	}

	// Extract raw table data
	tableData := data[headerSize : headerSize+int(tableLength)]

	// Parse entry point
	var epType EntryPointType
	if entryPointType == 1 {
		epType = EntryPoint64Bit
	} else {
		epType = EntryPoint32Bit
	}

	ep := EntryPoint{
		Type:         epType,
		MajorVersion: majorVersion,
		MinorVersion: minorVersion,
		Revision:     revision,
		TableAddress: tableAddress,
		TableLength:  tableLength,
	}

	// Parse structures from raw table data (same as reading from system)
	structures, err := ParseStructures(tableData, 0)
	if err != nil {
		return nil, err
	}

	return &SMBIOS{
		EntryPoint: ep,
		Structures: structures,
	}, nil
}

// writeSMBIOSToFile writes SMBIOS data to a raw dump file
// The file contains a small header followed by the reconstructed raw SMBIOS table
func writeSMBIOSToFile(sm *SMBIOS, filename string) error {
	// First, reconstruct the raw table data exactly as it appears in memory
	var tableData bytes.Buffer

	for _, s := range sm.Structures {
		// Write the raw formatted section (includes header)
		tableData.Write(s.Data)

		// Write string table
		if len(s.Strings) == 0 {
			// Empty string table: two null bytes
			tableData.WriteByte(0)
			tableData.WriteByte(0)
		} else {
			// Write each string followed by null terminator
			for _, str := range s.Strings {
				tableData.WriteString(str)
				tableData.WriteByte(0)
			}
			// End of string table: additional null byte
			tableData.WriteByte(0)
		}
	}

	rawTable := tableData.Bytes()

	// Create output file
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write header
	var header [28]byte
	copy(header[0:9], fileMagic)
	header[9] = fileVersion
	if sm.EntryPoint.Type == EntryPoint64Bit {
		header[10] = 1
	} else {
		header[10] = 0
	}
	header[11] = sm.EntryPoint.MajorVersion
	header[12] = sm.EntryPoint.MinorVersion
	header[13] = sm.EntryPoint.Revision
	header[14] = 0 // reserved
	binary.LittleEndian.PutUint32(header[15:19], uint32(len(rawTable)))
	binary.LittleEndian.PutUint64(header[19:27], sm.EntryPoint.TableAddress)
	header[27] = 0 // padding

	if _, err := f.Write(header[:]); err != nil {
		return err
	}

	// Write raw table data
	if _, err := f.Write(rawTable); err != nil {
		return err
	}

	return nil
}
