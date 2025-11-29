// Package type0 implements SMBIOS Type 0 - BIOS Information
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type0

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for BIOS Information
const StructureType uint8 = 0

// BIOSInfo represents Type 0 - BIOS Information
type BIOSInfo struct {
	Header                         gosmbios.Header
	Vendor                         string
	Version                        string
	StartingAddressSegment         uint16
	ReleaseDate                    string
	ROMSize                        uint8  // In 64K blocks, minus 1
	ROMSizeBytes                   uint64 // Calculated actual ROM size
	Characteristics                Characteristics
	CharacteristicsExt1            CharacteristicsExt1
	CharacteristicsExt2            CharacteristicsExt2
	SystemBIOSMajorRelease         uint8
	SystemBIOSMinorRelease         uint8
	EmbeddedControllerMajorRelease uint8
	EmbeddedControllerMinorRelease uint8
	ExtendedROMSize                uint16      // SMBIOS 3.1+
	ExtendedROMSizeUnit            ROMSizeUnit // Calculated from ExtendedROMSize
}

// ROMSizeUnit indicates the unit for extended ROM size
type ROMSizeUnit int

const (
	ROMSizeUnitMB ROMSizeUnit = iota // Megabytes
	ROMSizeUnitGB                    // Gigabytes
)

// Characteristics represents BIOS characteristics (64-bit field)
type Characteristics uint64

// BIOS Characteristics bit definitions
const (
	CharReserved0               Characteristics = 1 << 0  // Reserved
	CharReserved1               Characteristics = 1 << 1  // Reserved
	CharUnknown                 Characteristics = 1 << 2  // Unknown
	CharNotSupported            Characteristics = 1 << 3  // BIOS Characteristics are not supported
	CharISASupported            Characteristics = 1 << 4  // ISA is supported
	CharMCASupported            Characteristics = 1 << 5  // MCA is supported
	CharEISASupported           Characteristics = 1 << 6  // EISA is supported
	CharPCISupported            Characteristics = 1 << 7  // PCI is supported
	CharPCMCIASupported         Characteristics = 1 << 8  // PC Card (PCMCIA) is supported
	CharPlugAndPlaySupported    Characteristics = 1 << 9  // Plug and Play is supported
	CharAPMSupported            Characteristics = 1 << 10 // APM is supported
	CharBIOSUpgradeable         Characteristics = 1 << 11 // BIOS is Upgradeable (Flash)
	CharBIOSShadowingAllowed    Characteristics = 1 << 12 // BIOS Shadowing is allowed
	CharVLVESASupported         Characteristics = 1 << 13 // VL-VESA is supported
	CharESCDSupported           Characteristics = 1 << 14 // ESCD support is available
	CharBootFromCDSupported     Characteristics = 1 << 15 // Boot from CD is supported
	CharSelectableBootSupported Characteristics = 1 << 16 // Selectable Boot is supported
	CharBIOSROMSocketed         Characteristics = 1 << 17 // BIOS ROM is socketed
	CharBootFromPCMCIASupported Characteristics = 1 << 18 // Boot from PC Card (PCMCIA) is supported
	CharEDDSupported            Characteristics = 1 << 19 // EDD specification is supported
	CharFloppyNEC9800           Characteristics = 1 << 20 // Int 13h — Japanese floppy for NEC 9800 1.2 MB
	CharFloppyToshiba           Characteristics = 1 << 21 // Int 13h — Japanese floppy for Toshiba 1.2 MB
	CharFloppy360KB             Characteristics = 1 << 22 // Int 13h — 5.25" / 360 KB floppy services
	CharFloppy1200KB            Characteristics = 1 << 23 // Int 13h — 5.25" /1.2 MB floppy services
	CharFloppy720KB             Characteristics = 1 << 24 // Int 13h — 3.5" / 720 KB floppy services
	CharFloppy2880KB            Characteristics = 1 << 25 // Int 13h — 3.5" / 2.88 MB floppy services
	CharPrintScreenSupported    Characteristics = 1 << 26 // Int 5h, print screen Service is supported
	CharKeyboard8042Supported   Characteristics = 1 << 27 // Int 9h, 8042 keyboard services are supported
	CharSerialSupported         Characteristics = 1 << 28 // Int 14h, serial services are supported
	CharPrinterSupported        Characteristics = 1 << 29 // Int 17h, printer services are supported
	CharCGAMonoSupported        Characteristics = 1 << 30 // Int 10h, CGA/Mono Video Services are supported
	CharNECPC98                 Characteristics = 1 << 31 // NEC PC-98
	// Bits 32-63 are reserved for BIOS vendor
)

// Has checks if a characteristic is set
func (c Characteristics) Has(flag Characteristics) bool {
	return c&flag != 0
}

// CharacteristicsExt1 represents BIOS characteristics extension byte 1
type CharacteristicsExt1 uint8

// Extension byte 1 bit definitions
const (
	CharExt1ACPISupported         CharacteristicsExt1 = 1 << 0 // ACPI is supported
	CharExt1USBLegacySupported    CharacteristicsExt1 = 1 << 1 // USB Legacy is supported
	CharExt1AGPSupported          CharacteristicsExt1 = 1 << 2 // AGP is supported
	CharExt1I2OBootSupported      CharacteristicsExt1 = 1 << 3 // I2O boot is supported
	CharExt1LS120BootSupported    CharacteristicsExt1 = 1 << 4 // LS-120 SuperDisk boot is supported
	CharExt1ATAPIZIPBootSupported CharacteristicsExt1 = 1 << 5 // ATAPI ZIP drive boot is supported
	CharExt11394BootSupported     CharacteristicsExt1 = 1 << 6 // 1394 boot is supported
	CharExt1SmartBatterySupported CharacteristicsExt1 = 1 << 7 // Smart battery is supported
)

// Has checks if an extension characteristic is set
func (c CharacteristicsExt1) Has(flag CharacteristicsExt1) bool {
	return c&flag != 0
}

// CharacteristicsExt2 represents BIOS characteristics extension byte 2
type CharacteristicsExt2 uint8

// Extension byte 2 bit definitions
const (
	CharExt2BIOSBootSpecSupported       CharacteristicsExt2 = 1 << 0 // BIOS Boot Specification is supported
	CharExt2FnKeyNetworkBootSupported   CharacteristicsExt2 = 1 << 1 // Function key-initiated network boot is supported
	CharExt2TargetedContentDistribution CharacteristicsExt2 = 1 << 2 // Enable targeted content distribution
	CharExt2UEFISupported               CharacteristicsExt2 = 1 << 3 // UEFI Specification is supported
	CharExt2VirtualMachine              CharacteristicsExt2 = 1 << 4 // SMBIOS table describes a virtual machine
	CharExt2ManufacturingModeSupported  CharacteristicsExt2 = 1 << 5 // Manufacturing mode is supported
	CharExt2ManufacturingModeEnabled    CharacteristicsExt2 = 1 << 6 // Manufacturing mode is enabled
	// Bit 7 reserved
)

// Has checks if an extension characteristic is set
func (c CharacteristicsExt2) Has(flag CharacteristicsExt2) bool {
	return c&flag != 0
}

// Parse parses a BIOS Information structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*BIOSInfo, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	if len(s.Data) < 18 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &BIOSInfo{
		Header:                 s.Header,
		Vendor:                 s.GetString(s.GetByte(0x04)),
		Version:                s.GetString(s.GetByte(0x05)),
		StartingAddressSegment: s.GetWord(0x06),
		ReleaseDate:            s.GetString(s.GetByte(0x08)),
		ROMSize:                s.GetByte(0x09),
		Characteristics:        Characteristics(s.GetQWord(0x0A)),
	}

	// Calculate ROM size in bytes
	if info.ROMSize != 0xFF {
		info.ROMSizeBytes = uint64(info.ROMSize+1) * 64 * 1024
	}

	// Extension bytes (SMBIOS 2.4+)
	if len(s.Data) >= 19 {
		info.CharacteristicsExt1 = CharacteristicsExt1(s.GetByte(0x12))
	}
	if len(s.Data) >= 20 {
		info.CharacteristicsExt2 = CharacteristicsExt2(s.GetByte(0x13))
	}

	// BIOS release info (SMBIOS 2.4+)
	if len(s.Data) >= 22 {
		info.SystemBIOSMajorRelease = s.GetByte(0x14)
		info.SystemBIOSMinorRelease = s.GetByte(0x15)
	}

	// Embedded controller release info (SMBIOS 2.4+)
	if len(s.Data) >= 24 {
		info.EmbeddedControllerMajorRelease = s.GetByte(0x16)
		info.EmbeddedControllerMinorRelease = s.GetByte(0x17)
	}

	// Extended ROM size (SMBIOS 3.1+)
	if len(s.Data) >= 26 && info.ROMSize == 0xFF {
		extSize := s.GetWord(0x18)
		info.ExtendedROMSize = extSize & 0x3FFF
		if extSize&0xC000 == 0 {
			info.ExtendedROMSizeUnit = ROMSizeUnitMB
			info.ROMSizeBytes = uint64(info.ExtendedROMSize) * 1024 * 1024
		} else {
			info.ExtendedROMSizeUnit = ROMSizeUnitGB
			info.ROMSizeBytes = uint64(info.ExtendedROMSize) * 1024 * 1024 * 1024
		}
	}

	return info, nil
}

// Get retrieves the BIOS Information from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*BIOSInfo, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// BIOSVersionString returns a formatted BIOS version string
func (b *BIOSInfo) BIOSVersionString() string {
	if b.SystemBIOSMajorRelease != 0xFF {
		return fmt.Sprintf("%d.%d", b.SystemBIOSMajorRelease, b.SystemBIOSMinorRelease)
	}
	return b.Version
}

// ECVersionString returns a formatted Embedded Controller version string
func (b *BIOSInfo) ECVersionString() string {
	if b.EmbeddedControllerMajorRelease == 0xFF {
		return "Not Present"
	}
	return fmt.Sprintf("%d.%d", b.EmbeddedControllerMajorRelease, b.EmbeddedControllerMinorRelease)
}

// ROMSizeString returns a human-readable ROM size string
func (b *BIOSInfo) ROMSizeString() string {
	if b.ROMSizeBytes == 0 {
		return "Unknown"
	}
	if b.ROMSizeBytes >= 1024*1024*1024 {
		return fmt.Sprintf("%d GB", b.ROMSizeBytes/(1024*1024*1024))
	}
	if b.ROMSizeBytes >= 1024*1024 {
		return fmt.Sprintf("%d MB", b.ROMSizeBytes/(1024*1024))
	}
	return fmt.Sprintf("%d KB", b.ROMSizeBytes/1024)
}

// IsUEFI returns true if UEFI is supported
func (b *BIOSInfo) IsUEFI() bool {
	return b.CharacteristicsExt2.Has(CharExt2UEFISupported)
}

// IsVirtualMachine returns true if running in a virtual machine
func (b *BIOSInfo) IsVirtualMachine() bool {
	return b.CharacteristicsExt2.Has(CharExt2VirtualMachine)
}
