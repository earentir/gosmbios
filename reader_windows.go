//go:build windows

package gosmbios

import (
	"syscall"
	"unsafe"
)

var (
	kernel32                   = syscall.NewLazyDLL("kernel32.dll")
	procGetSystemFirmwareTable = kernel32.NewProc("GetSystemFirmwareTable")
)

const (
	// FirmwareTableProviderSignature for SMBIOS
	// 'RSMB' in little-endian
	firmwareTableIDRSMB = 0x52534D42
)

// RawSMBIOSData represents the Windows raw SMBIOS data structure
type rawSMBIOSData struct {
	Used20CallingMethod uint8
	MajorVersion        uint8
	MinorVersion        uint8
	DMIRevision         uint8
	Length              uint32
	// SMBIOSTableData follows immediately after
}

// readSMBIOS reads SMBIOS data on Windows systems
func readSMBIOS() (*SMBIOS, error) {
	// First call to get the required buffer size
	size, _, _ := procGetSystemFirmwareTable.Call(
		uintptr(firmwareTableIDRSMB),
		0,
		0,
		0,
	)

	if size == 0 {
		return nil, ErrNotFound
	}

	// Allocate buffer and get the data
	buffer := make([]byte, size)
	ret, _, err := procGetSystemFirmwareTable.Call(
		uintptr(firmwareTableIDRSMB),
		0,
		uintptr(unsafe.Pointer(&buffer[0])),
		size,
	)

	if ret == 0 {
		if err != nil && err != syscall.Errno(0) {
			return nil, ErrAccessDenied
		}
		return nil, ErrNotFound
	}

	// Parse the raw SMBIOS data header
	if len(buffer) < 8 {
		return nil, ErrInvalidStructure
	}

	rawHeader := (*rawSMBIOSData)(unsafe.Pointer(&buffer[0]))

	// Create entry point from Windows data
	entryPoint := &EntryPoint{
		MajorVersion: rawHeader.MajorVersion,
		MinorVersion: rawHeader.MinorVersion,
		Revision:     rawHeader.DMIRevision,
		TableLength:  rawHeader.Length,
	}

	// Determine entry point type based on version
	if rawHeader.MajorVersion >= 3 {
		entryPoint.Type = EntryPoint64Bit
	} else {
		entryPoint.Type = EntryPoint32Bit
	}

	// Extract table data (starts after the 8-byte header)
	headerSize := 8
	if len(buffer) < headerSize+int(rawHeader.Length) {
		// Use available data if Length is larger than buffer
		rawHeader.Length = uint32(len(buffer) - headerSize)
	}

	tableData := buffer[headerSize : headerSize+int(rawHeader.Length)]

	// Parse structures
	structures, err := ParseStructures(tableData, 0)
	if err != nil {
		return nil, err
	}

	return &SMBIOS{
		EntryPoint: *entryPoint,
		Structures: structures,
	}, nil
}
