// Package type17 implements SMBIOS Type 17 - Memory Device
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type17

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Memory Device
const StructureType uint8 = 17

// MemoryDevice represents Type 17 - Memory Device
type MemoryDevice struct {
	Header                                  gosmbios.Header
	PhysicalMemoryArrayHandle               uint16
	MemoryErrorInformationHandle            uint16
	TotalWidth                              uint16 // In bits, 0xFFFF = unknown
	DataWidth                               uint16 // In bits, 0xFFFF = unknown
	Size                                    uint64 // In MB (calculated)
	FormFactor                              MemoryFormFactor
	DeviceSet                               uint8
	DeviceLocator                           string
	BankLocator                             string
	MemoryType                              MemoryType
	TypeDetail                              MemoryTypeDetail
	Speed                                   uint16                  // In MT/s, 0 = unknown
	Manufacturer                            string                  // SMBIOS 2.3+
	SerialNumber                            string                  // SMBIOS 2.3+
	AssetTag                                string                  // SMBIOS 2.3+
	PartNumber                              string                  // SMBIOS 2.3+
	Attributes                              uint8                   // SMBIOS 2.6+
	ExtendedSize                            uint32                  // SMBIOS 2.7+ (in MB)
	ConfiguredMemorySpeed                   uint16                  // SMBIOS 2.7+ (in MT/s)
	MinimumVoltage                          uint16                  // SMBIOS 2.8+ (in mV)
	MaximumVoltage                          uint16                  // SMBIOS 2.8+ (in mV)
	ConfiguredVoltage                       uint16                  // SMBIOS 2.8+ (in mV)
	MemoryTechnology                        MemoryTechnology        // SMBIOS 3.2+
	MemoryOperatingModeCapability           OperatingModeCapability // SMBIOS 3.2+
	FirmwareVersion                         string                  // SMBIOS 3.2+
	ModuleManufacturerID                    uint16                  // SMBIOS 3.2+
	ModuleProductID                         uint16                  // SMBIOS 3.2+
	MemorySubsystemControllerManufacturerID uint16                  // SMBIOS 3.2+
	MemorySubsystemControllerProductID      uint16                  // SMBIOS 3.2+
	NonVolatileSize                         uint64                  // SMBIOS 3.2+ (in bytes)
	VolatileSize                            uint64                  // SMBIOS 3.2+ (in bytes)
	CacheSize                               uint64                  // SMBIOS 3.2+ (in bytes)
	LogicalSize                             uint64                  // SMBIOS 3.2+ (in bytes)
	ExtendedSpeed                           uint32                  // SMBIOS 3.3+ (in MT/s)
	ExtendedConfiguredMemorySpeed           uint32                  // SMBIOS 3.3+ (in MT/s)
	PMIC0ManufacturerID                     uint16                  // SMBIOS 3.7+
	PMIC0RevisionNumber                     uint16                  // SMBIOS 3.7+
	RCDManufacturerID                       uint16                  // SMBIOS 3.7+
	RCDRevisionNumber                       uint16                  // SMBIOS 3.7+
}

// MemoryFormFactor identifies the physical form factor of the memory device
type MemoryFormFactor uint8

// Memory form factor values
const (
	FormFactorOther           MemoryFormFactor = 0x01
	FormFactorUnknown         MemoryFormFactor = 0x02
	FormFactorSIMM            MemoryFormFactor = 0x03
	FormFactorSIP             MemoryFormFactor = 0x04
	FormFactorChip            MemoryFormFactor = 0x05
	FormFactorDIP             MemoryFormFactor = 0x06
	FormFactorZIP             MemoryFormFactor = 0x07
	FormFactorProprietaryCard MemoryFormFactor = 0x08
	FormFactorDIMM            MemoryFormFactor = 0x09
	FormFactorTSOP            MemoryFormFactor = 0x0A
	FormFactorRowOfChips      MemoryFormFactor = 0x0B
	FormFactorRIMM            MemoryFormFactor = 0x0C
	FormFactorSODIMM          MemoryFormFactor = 0x0D
	FormFactorSRIMM           MemoryFormFactor = 0x0E
	FormFactorFBDIMM          MemoryFormFactor = 0x0F
	FormFactorDie             MemoryFormFactor = 0x10
)

// String returns a human-readable form factor description
func (ff MemoryFormFactor) String() string {
	formFactors := map[MemoryFormFactor]string{
		FormFactorOther:           "Other",
		FormFactorUnknown:         "Unknown",
		FormFactorSIMM:            "SIMM",
		FormFactorSIP:             "SIP",
		FormFactorChip:            "Chip",
		FormFactorDIP:             "DIP",
		FormFactorZIP:             "ZIP",
		FormFactorProprietaryCard: "Proprietary Card",
		FormFactorDIMM:            "DIMM",
		FormFactorTSOP:            "TSOP",
		FormFactorRowOfChips:      "Row of chips",
		FormFactorRIMM:            "RIMM",
		FormFactorSODIMM:          "SO-DIMM",
		FormFactorSRIMM:           "SRIMM",
		FormFactorFBDIMM:          "FB-DIMM",
		FormFactorDie:             "Die",
	}

	if name, ok := formFactors[ff]; ok {
		return name
	}
	return fmt.Sprintf("Unknown (0x%02X)", uint8(ff))
}

// MemoryType identifies the type of memory
type MemoryType uint8

// Memory type values
const (
	MemTypeOther              MemoryType = 0x01
	MemTypeUnknown            MemoryType = 0x02
	MemTypeDRAM               MemoryType = 0x03
	MemTypeEDRAM              MemoryType = 0x04
	MemTypeVRAM               MemoryType = 0x05
	MemTypeSRAM               MemoryType = 0x06
	MemTypeRAM                MemoryType = 0x07
	MemTypeROM                MemoryType = 0x08
	MemTypeFlash              MemoryType = 0x09
	MemTypeEEPROM             MemoryType = 0x0A
	MemTypeFEPROM             MemoryType = 0x0B
	MemTypeEPROM              MemoryType = 0x0C
	MemTypeCDRAM              MemoryType = 0x0D
	MemType3DRAM              MemoryType = 0x0E
	MemTypeSDRAM              MemoryType = 0x0F
	MemTypeSGRAM              MemoryType = 0x10
	MemTypeRDRAM              MemoryType = 0x11
	MemTypeDDR                MemoryType = 0x12
	MemTypeDDR2               MemoryType = 0x13
	MemTypeDDR2FBDIMM         MemoryType = 0x14
	MemTypeDDR3               MemoryType = 0x18
	MemTypeFBD2               MemoryType = 0x19
	MemTypeDDR4               MemoryType = 0x1A
	MemTypeLPDDR              MemoryType = 0x1B
	MemTypeLPDDR2             MemoryType = 0x1C
	MemTypeLPDDR3             MemoryType = 0x1D
	MemTypeLPDDR4             MemoryType = 0x1E
	MemTypeLogicalNonVolatile MemoryType = 0x1F
	MemTypeHBM                MemoryType = 0x20
	MemTypeHBM2               MemoryType = 0x21
	MemTypeDDR5               MemoryType = 0x22
	MemTypeLPDDR5             MemoryType = 0x23
	MemTypeHBM3               MemoryType = 0x24
)

// String returns a human-readable memory type description
func (mt MemoryType) String() string {
	types := map[MemoryType]string{
		MemTypeOther:              "Other",
		MemTypeUnknown:            "Unknown",
		MemTypeDRAM:               "DRAM",
		MemTypeEDRAM:              "EDRAM",
		MemTypeVRAM:               "VRAM",
		MemTypeSRAM:               "SRAM",
		MemTypeRAM:                "RAM",
		MemTypeROM:                "ROM",
		MemTypeFlash:              "Flash",
		MemTypeEEPROM:             "EEPROM",
		MemTypeFEPROM:             "FEPROM",
		MemTypeEPROM:              "EPROM",
		MemTypeCDRAM:              "CDRAM",
		MemType3DRAM:              "3DRAM",
		MemTypeSDRAM:              "SDRAM",
		MemTypeSGRAM:              "SGRAM",
		MemTypeRDRAM:              "RDRAM",
		MemTypeDDR:                "DDR",
		MemTypeDDR2:               "DDR2",
		MemTypeDDR2FBDIMM:         "DDR2 FB-DIMM",
		MemTypeDDR3:               "DDR3",
		MemTypeFBD2:               "FBD2",
		MemTypeDDR4:               "DDR4",
		MemTypeLPDDR:              "LPDDR",
		MemTypeLPDDR2:             "LPDDR2",
		MemTypeLPDDR3:             "LPDDR3",
		MemTypeLPDDR4:             "LPDDR4",
		MemTypeLogicalNonVolatile: "Logical non-volatile device",
		MemTypeHBM:                "HBM",
		MemTypeHBM2:               "HBM2",
		MemTypeDDR5:               "DDR5",
		MemTypeLPDDR5:             "LPDDR5",
		MemTypeHBM3:               "HBM3",
	}

	if name, ok := types[mt]; ok {
		return name
	}
	return fmt.Sprintf("Unknown (0x%02X)", uint8(mt))
}

// IsDDR returns true if the memory type is DDR-based
func (mt MemoryType) IsDDR() bool {
	switch mt {
	case MemTypeDDR, MemTypeDDR2, MemTypeDDR2FBDIMM, MemTypeDDR3, MemTypeDDR4, MemTypeDDR5:
		return true
	}
	return false
}

// MemoryTypeDetail represents memory type detail flags
type MemoryTypeDetail uint16

// Memory type detail bit definitions
const (
	TypeDetailReserved0    MemoryTypeDetail = 1 << 0  // Reserved
	TypeDetailOther        MemoryTypeDetail = 1 << 1  // Other
	TypeDetailUnknown      MemoryTypeDetail = 1 << 2  // Unknown
	TypeDetailFastPaged    MemoryTypeDetail = 1 << 3  // Fast-paged
	TypeDetailStaticColumn MemoryTypeDetail = 1 << 4  // Static column
	TypeDetailPseudoStatic MemoryTypeDetail = 1 << 5  // Pseudo-static
	TypeDetailRAMBUS       MemoryTypeDetail = 1 << 6  // RAMBUS
	TypeDetailSynchronous  MemoryTypeDetail = 1 << 7  // Synchronous
	TypeDetailCMOS         MemoryTypeDetail = 1 << 8  // CMOS
	TypeDetailEDO          MemoryTypeDetail = 1 << 9  // EDO
	TypeDetailWindowDRAM   MemoryTypeDetail = 1 << 10 // Window DRAM
	TypeDetailCacheDRAM    MemoryTypeDetail = 1 << 11 // Cache DRAM
	TypeDetailNonVolatile  MemoryTypeDetail = 1 << 12 // Non-volatile
	TypeDetailRegistered   MemoryTypeDetail = 1 << 13 // Registered (Buffered)
	TypeDetailUnbuffered   MemoryTypeDetail = 1 << 14 // Unbuffered (Unregistered)
	TypeDetailLRDIMM       MemoryTypeDetail = 1 << 15 // LRDIMM
)

// Has checks if a type detail flag is set
func (td MemoryTypeDetail) Has(flag MemoryTypeDetail) bool {
	return td&flag != 0
}

// String returns a human-readable type detail description
func (td MemoryTypeDetail) String() string {
	var details []string

	if td.Has(TypeDetailOther) {
		details = append(details, "Other")
	}
	if td.Has(TypeDetailUnknown) {
		details = append(details, "Unknown")
	}
	if td.Has(TypeDetailFastPaged) {
		details = append(details, "Fast-paged")
	}
	if td.Has(TypeDetailStaticColumn) {
		details = append(details, "Static column")
	}
	if td.Has(TypeDetailPseudoStatic) {
		details = append(details, "Pseudo-static")
	}
	if td.Has(TypeDetailRAMBUS) {
		details = append(details, "RAMBUS")
	}
	if td.Has(TypeDetailSynchronous) {
		details = append(details, "Synchronous")
	}
	if td.Has(TypeDetailCMOS) {
		details = append(details, "CMOS")
	}
	if td.Has(TypeDetailEDO) {
		details = append(details, "EDO")
	}
	if td.Has(TypeDetailWindowDRAM) {
		details = append(details, "Window DRAM")
	}
	if td.Has(TypeDetailCacheDRAM) {
		details = append(details, "Cache DRAM")
	}
	if td.Has(TypeDetailNonVolatile) {
		details = append(details, "Non-volatile")
	}
	if td.Has(TypeDetailRegistered) {
		details = append(details, "Registered (Buffered)")
	}
	if td.Has(TypeDetailUnbuffered) {
		details = append(details, "Unbuffered (Unregistered)")
	}
	if td.Has(TypeDetailLRDIMM) {
		details = append(details, "LRDIMM")
	}

	if len(details) == 0 {
		return "None"
	}

	result := details[0]
	for i := 1; i < len(details); i++ {
		result += ", " + details[i]
	}
	return result
}

// MemoryTechnology identifies the memory technology (SMBIOS 3.2+)
type MemoryTechnology uint8

// Memory technology values
const (
	TechOther                 MemoryTechnology = 0x01
	TechUnknown               MemoryTechnology = 0x02
	TechDRAM                  MemoryTechnology = 0x03
	TechNVDIMMN               MemoryTechnology = 0x04
	TechNVDIMMF               MemoryTechnology = 0x05
	TechNVDIMMP               MemoryTechnology = 0x06
	TechIntelOptanePersistent MemoryTechnology = 0x07
)

// String returns a human-readable technology description
func (mt MemoryTechnology) String() string {
	switch mt {
	case TechOther:
		return "Other"
	case TechUnknown:
		return "Unknown"
	case TechDRAM:
		return "DRAM"
	case TechNVDIMMN:
		return "NVDIMM-N"
	case TechNVDIMMF:
		return "NVDIMM-F"
	case TechNVDIMMP:
		return "NVDIMM-P"
	case TechIntelOptanePersistent:
		return "Intel Optane Persistent Memory"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(mt))
	}
}

// OperatingModeCapability represents memory operating mode capability (SMBIOS 3.2+)
type OperatingModeCapability uint16

// Operating mode capability bit definitions
const (
	OpModeReserved0             OperatingModeCapability = 1 << 0 // Reserved
	OpModeOther                 OperatingModeCapability = 1 << 1 // Other
	OpModeUnknown               OperatingModeCapability = 1 << 2 // Unknown
	OpModeVolatile              OperatingModeCapability = 1 << 3 // Volatile memory
	OpModeByteAccessPersistent  OperatingModeCapability = 1 << 4 // Byte-accessible persistent memory
	OpModeBlockAccessPersistent OperatingModeCapability = 1 << 5 // Block-accessible persistent memory
)

// Has checks if an operating mode is supported
func (omc OperatingModeCapability) Has(flag OperatingModeCapability) bool {
	return omc&flag != 0
}

// Parse parses a Memory Device structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*MemoryDevice, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 21 bytes (SMBIOS 2.1)
	if len(s.Data) < 21 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &MemoryDevice{
		Header:                       s.Header,
		PhysicalMemoryArrayHandle:    s.GetWord(0x04),
		MemoryErrorInformationHandle: s.GetWord(0x06),
		TotalWidth:                   s.GetWord(0x08),
		DataWidth:                    s.GetWord(0x0A),
		FormFactor:                   MemoryFormFactor(s.GetByte(0x0E)),
		DeviceSet:                    s.GetByte(0x0F),
		DeviceLocator:                s.GetString(s.GetByte(0x10)),
		BankLocator:                  s.GetString(s.GetByte(0x11)),
		MemoryType:                   MemoryType(s.GetByte(0x12)),
		TypeDetail:                   MemoryTypeDetail(s.GetWord(0x13)),
	}

	// Parse size (16-bit field)
	sizeRaw := s.GetWord(0x0C)
	if sizeRaw == 0xFFFF {
		// Size unknown
		info.Size = 0
	} else if sizeRaw == 0x7FFF {
		// Extended size field is used
		info.Size = 0 // Will be set from ExtendedSize
	} else if sizeRaw&0x8000 != 0 {
		// Size in KB
		info.Size = uint64(sizeRaw&0x7FFF) / 1024
	} else {
		// Size in MB
		info.Size = uint64(sizeRaw)
	}

	// SMBIOS 2.3+
	if len(s.Data) >= 27 {
		info.Speed = s.GetWord(0x15)
		info.Manufacturer = s.GetString(s.GetByte(0x17))
		info.SerialNumber = s.GetString(s.GetByte(0x18))
		info.AssetTag = s.GetString(s.GetByte(0x19))
		info.PartNumber = s.GetString(s.GetByte(0x1A))
	}

	// SMBIOS 2.6+
	if len(s.Data) >= 28 {
		info.Attributes = s.GetByte(0x1B)
	}

	// SMBIOS 2.7+
	if len(s.Data) >= 34 {
		info.ExtendedSize = s.GetDWord(0x1C)
		info.ConfiguredMemorySpeed = s.GetWord(0x20)

		// Use extended size if primary size indicates it
		if sizeRaw == 0x7FFF {
			info.Size = uint64(info.ExtendedSize & 0x7FFFFFFF) // In MB
		}
	}

	// SMBIOS 2.8+
	if len(s.Data) >= 40 {
		info.MinimumVoltage = s.GetWord(0x22)
		info.MaximumVoltage = s.GetWord(0x24)
		info.ConfiguredVoltage = s.GetWord(0x26)
	}

	// SMBIOS 3.2+
	if len(s.Data) >= 84 {
		info.MemoryTechnology = MemoryTechnology(s.GetByte(0x28))
		info.MemoryOperatingModeCapability = OperatingModeCapability(s.GetWord(0x29))
		info.FirmwareVersion = s.GetString(s.GetByte(0x2B))
		info.ModuleManufacturerID = s.GetWord(0x2C)
		info.ModuleProductID = s.GetWord(0x2E)
		info.MemorySubsystemControllerManufacturerID = s.GetWord(0x30)
		info.MemorySubsystemControllerProductID = s.GetWord(0x32)
		info.NonVolatileSize = s.GetQWord(0x34)
		info.VolatileSize = s.GetQWord(0x3C)
		info.CacheSize = s.GetQWord(0x44)
		info.LogicalSize = s.GetQWord(0x4C)
	}

	// SMBIOS 3.3+
	if len(s.Data) >= 92 {
		info.ExtendedSpeed = s.GetDWord(0x54)
		info.ExtendedConfiguredMemorySpeed = s.GetDWord(0x58)

		// Use extended speed if primary indicates it
		if info.Speed == 0xFFFF {
			info.Speed = uint16(info.ExtendedSpeed & 0xFFFF)
		}
	}

	// SMBIOS 3.7+
	if len(s.Data) >= 100 {
		info.PMIC0ManufacturerID = s.GetWord(0x5C)
		info.PMIC0RevisionNumber = s.GetWord(0x5E)
		info.RCDManufacturerID = s.GetWord(0x60)
		info.RCDRevisionNumber = s.GetWord(0x62)
	}

	return info, nil
}

// Get retrieves the first Memory Device from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*MemoryDevice, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Memory Device structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*MemoryDevice, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var devices []*MemoryDevice
	for i := range structures {
		dev, err := Parse(&structures[i])
		if err == nil {
			devices = append(devices, dev)
		}
	}

	if len(devices) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return devices, nil
}

// GetPopulated retrieves only populated Memory Device structures (those with memory installed)
func GetPopulated(sm *gosmbios.SMBIOS) ([]*MemoryDevice, error) {
	all, err := GetAll(sm)
	if err != nil {
		return nil, err
	}

	var populated []*MemoryDevice
	for _, dev := range all {
		if dev.Size > 0 {
			populated = append(populated, dev)
		}
	}

	if len(populated) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return populated, nil
}

// SizeString returns a human-readable size string
func (m *MemoryDevice) SizeString() string {
	if m.Size == 0 {
		return "Unknown"
	}

	mb := m.Size
	if mb >= 1024*1024 {
		return fmt.Sprintf("%d TB", mb/(1024*1024))
	}
	if mb >= 1024 {
		return fmt.Sprintf("%d GB", mb/1024)
	}
	return fmt.Sprintf("%d MB", mb)
}

// SpeedString returns a human-readable speed string
func (m *MemoryDevice) SpeedString() string {
	speed := m.GetSpeed()
	if speed == 0 {
		return "Unknown"
	}
	return fmt.Sprintf("%d MT/s", speed)
}

// GetSpeed returns the effective speed in MT/s
func (m *MemoryDevice) GetSpeed() uint32 {
	if m.ExtendedSpeed > 0 && m.Speed == 0xFFFF {
		return m.ExtendedSpeed
	}
	return uint32(m.Speed)
}

// GetConfiguredSpeed returns the configured speed in MT/s
func (m *MemoryDevice) GetConfiguredSpeed() uint32 {
	if m.ExtendedConfiguredMemorySpeed > 0 && m.ConfiguredMemorySpeed == 0xFFFF {
		return m.ExtendedConfiguredMemorySpeed
	}
	return uint32(m.ConfiguredMemorySpeed)
}

// VoltageString returns a human-readable voltage string
func (m *MemoryDevice) VoltageString() string {
	if m.ConfiguredVoltage == 0 {
		return "Unknown"
	}
	return fmt.Sprintf("%.2f V", float64(m.ConfiguredVoltage)/1000.0)
}

// Ranks returns the number of ranks (from Attributes field)
func (m *MemoryDevice) Ranks() uint8 {
	return m.Attributes & 0x0F
}

// IsPopulated returns true if memory is installed in this slot
func (m *MemoryDevice) IsPopulated() bool {
	return m.Size > 0
}

// DisplayName returns a display-friendly memory description
func (m *MemoryDevice) DisplayName() string {
	if !m.IsPopulated() {
		return fmt.Sprintf("%s (Empty)", m.DeviceLocator)
	}

	return fmt.Sprintf("%s %s %s (%s)",
		m.Manufacturer,
		m.MemoryType.String(),
		m.SizeString(),
		m.SpeedString())
}
