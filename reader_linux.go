//go:build linux

package gosmbios

import (
	"os"
)

const (
	// Linux sysfs paths for SMBIOS data
	sysfsEntryPoint = "/sys/firmware/dmi/tables/smbios_entry_point"
	sysfsDMITable   = "/sys/firmware/dmi/tables/DMI"

	// Legacy paths (fallback)
	legacyEntryPoint = "/sys/firmware/dmi/tables/smbios_entry_point"
	legacyDMITable   = "/sys/firmware/dmi/tables/DMI"

	// EFI systab path for SMBIOS address discovery
	efiSystab = "/sys/firmware/efi/systab"
)

// readSMBIOS reads SMBIOS data on Linux systems
func readSMBIOS() (*SMBIOS, error) {
	// Try reading from sysfs first (preferred method, works without root on most systems)
	entryPointData, err := os.ReadFile(sysfsEntryPoint)
	if err != nil {
		return nil, ErrNotFound
	}

	tableData, err := os.ReadFile(sysfsDMITable)
	if err != nil {
		return nil, ErrNotFound
	}

	// Try parsing as 64-bit entry point first (SMBIOS 3.x)
	entryPoint, err := ParseEntryPoint64(entryPointData)
	if err != nil {
		// Fall back to 32-bit entry point (SMBIOS 2.x)
		entryPoint, err = ParseEntryPoint32(entryPointData)
		if err != nil {
			return nil, err
		}
	}

	// Parse structures
	maxStructures := 0
	if entryPoint.Type == EntryPoint32Bit {
		maxStructures = int(entryPoint.StructureCount)
	}

	structures, err := ParseStructures(tableData, maxStructures)
	if err != nil {
		return nil, err
	}

	// Update table length from actual data if not set
	if entryPoint.TableLength == 0 {
		entryPoint.TableLength = uint32(len(tableData))
	}

	return &SMBIOS{
		EntryPoint: *entryPoint,
		Structures: structures,
	}, nil
}
