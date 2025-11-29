// Package gosmbios provides pure Go implementation for reading SMBIOS/DMI data
// across Windows, Linux, and macOS on both AMD64 and ARM64 architectures.
// Implements DSP0134 SMBIOS Reference Specification Version 3.9.0
package gosmbios

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// Version represents the SMBIOS specification version implemented
const Version = "3.9.0"

// Common errors
var (
	ErrNotFound         = errors.New("smbios: entry point not found")
	ErrInvalidChecksum  = errors.New("smbios: invalid checksum")
	ErrUnsupportedOS    = errors.New("smbios: unsupported operating system")
	ErrAccessDenied     = errors.New("smbios: access denied (try running as root/admin)")
	ErrInvalidStructure = errors.New("smbios: invalid structure format")
)

// EntryPointType indicates the SMBIOS entry point version
type EntryPointType int

const (
	EntryPoint32Bit EntryPointType = iota // SMBIOS 2.x 32-bit entry point
	EntryPoint64Bit                       // SMBIOS 3.x 64-bit entry point
)

// EntryPoint contains SMBIOS entry point information
type EntryPoint struct {
	Type             EntryPointType
	MajorVersion     uint8
	MinorVersion     uint8
	Revision         uint8 // Only for 3.x
	TableAddress     uint64
	TableLength      uint32
	TableMaxSize     uint32 // Only for 3.x
	StructureCount   uint16 // Only for 2.x (not reliable for 3.x)
	BCDRevision      uint8  // Only for 2.x
	EntryPointLength uint8
}

// String returns a human-readable version string
func (ep *EntryPoint) String() string {
	if ep.Type == EntryPoint64Bit {
		return fmt.Sprintf("SMBIOS %d.%d.%d", ep.MajorVersion, ep.MinorVersion, ep.Revision)
	}
	return fmt.Sprintf("SMBIOS %d.%d", ep.MajorVersion, ep.MinorVersion)
}

// Header represents the common SMBIOS structure header (4 bytes)
type Header struct {
	Type   uint8
	Length uint8
	Handle uint16
}

// Structure represents a single SMBIOS structure with its data and strings
type Structure struct {
	Header  Header
	Data    []byte   // Raw formatted section data (includes header)
	Strings []string // String table entries
}

// GetString returns a string from the string table (1-indexed as per SMBIOS spec)
// Returns empty string if index is 0 or out of bounds
func (s *Structure) GetString(index uint8) string {
	if index == 0 || int(index) > len(s.Strings) {
		return ""
	}
	return s.Strings[index-1]
}

// GetByte returns a byte at the given offset in the formatted section
func (s *Structure) GetByte(offset int) uint8 {
	if offset >= len(s.Data) {
		return 0
	}
	return s.Data[offset]
}

// GetWord returns a 16-bit little-endian value at the given offset
func (s *Structure) GetWord(offset int) uint16 {
	if offset+1 >= len(s.Data) {
		return 0
	}
	return binary.LittleEndian.Uint16(s.Data[offset:])
}

// GetDWord returns a 32-bit little-endian value at the given offset
func (s *Structure) GetDWord(offset int) uint32 {
	if offset+3 >= len(s.Data) {
		return 0
	}
	return binary.LittleEndian.Uint32(s.Data[offset:])
}

// GetQWord returns a 64-bit little-endian value at the given offset
func (s *Structure) GetQWord(offset int) uint64 {
	if offset+7 >= len(s.Data) {
		return 0
	}
	return binary.LittleEndian.Uint64(s.Data[offset:])
}

// SMBIOS holds all parsed SMBIOS data
type SMBIOS struct {
	EntryPoint EntryPoint
	Structures []Structure
}

// GetStructures returns all structures of the specified type
func (sm *SMBIOS) GetStructures(structType uint8) []Structure {
	var result []Structure
	for _, s := range sm.Structures {
		if s.Header.Type == structType {
			result = append(result, s)
		}
	}
	return result
}

// GetStructure returns the first structure of the specified type, or nil if not found
func (sm *SMBIOS) GetStructure(structType uint8) *Structure {
	for i := range sm.Structures {
		if sm.Structures[i].Header.Type == structType {
			return &sm.Structures[i]
		}
	}
	return nil
}

// Read reads and parses SMBIOS data from the system
// This is the main entry point for the library
func Read() (*SMBIOS, error) {
	return readSMBIOS()
}

// ReadFromFile reads SMBIOS data from a binary dump file
// The file format is a simple binary format:
// - 32 bytes: Entry point header
// - Remaining: Raw SMBIOS table data
func ReadFromFile(filename string) (*SMBIOS, error) {
	return readSMBIOSFromFile(filename)
}

// WriteToFile writes SMBIOS data to a binary dump file
func (sm *SMBIOS) WriteToFile(filename string) error {
	return writeSMBIOSToFile(sm, filename)
}
