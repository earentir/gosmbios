// Package type4 implements SMBIOS Type 4 - Processor Information
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type4

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Processor Information
const StructureType uint8 = 4

// ProcessorInfo represents Type 4 - Processor Information
type ProcessorInfo struct {
	Header                   gosmbios.Header
	SocketDesignation        string
	ProcessorType            ProcessorType
	ProcessorFamily          ProcessorFamily
	ProcessorManufacturer    string
	ProcessorID              uint64
	ProcessorVersion         string
	Voltage                  Voltage
	ExternalClock            uint16 // MHz, 0 = unknown
	MaxSpeed                 uint16 // MHz, 0 = unknown
	CurrentSpeed             uint16 // MHz, 0 = unknown
	Status                   ProcessorStatus
	ProcessorUpgrade         ProcessorUpgrade
	L1CacheHandle            uint16                   // SMBIOS 2.1+
	L2CacheHandle            uint16                   // SMBIOS 2.1+
	L3CacheHandle            uint16                   // SMBIOS 2.1+
	SerialNumber             string                   // SMBIOS 2.3+
	AssetTag                 string                   // SMBIOS 2.3+
	PartNumber               string                   // SMBIOS 2.3+
	CoreCount                uint8                    // SMBIOS 2.5+
	CoreEnabled              uint8                    // SMBIOS 2.5+
	ThreadCount              uint8                    // SMBIOS 2.5+
	ProcessorCharacteristics ProcessorCharacteristics // SMBIOS 2.5+
	ProcessorFamily2         uint16                   // SMBIOS 2.6+
	CoreCount2               uint16                   // SMBIOS 3.0+ (used when CoreCount is 0xFF)
	CoreEnabled2             uint16                   // SMBIOS 3.0+ (used when CoreEnabled is 0xFF)
	ThreadCount2             uint16                   // SMBIOS 3.0+ (used when ThreadCount is 0xFF)
	ThreadEnabled            uint16                   // SMBIOS 3.6+
}

// ProcessorType identifies the processor type
type ProcessorType uint8

// Processor type values
const (
	ProcessorTypeOther            ProcessorType = 0x01
	ProcessorTypeUnknown          ProcessorType = 0x02
	ProcessorTypeCentralProcessor ProcessorType = 0x03
	ProcessorTypeMathProcessor    ProcessorType = 0x04
	ProcessorTypeDSP              ProcessorType = 0x05
	ProcessorTypeVideoProcessor   ProcessorType = 0x06
)

// String returns a human-readable processor type description
func (pt ProcessorType) String() string {
	switch pt {
	case ProcessorTypeOther:
		return "Other"
	case ProcessorTypeUnknown:
		return "Unknown"
	case ProcessorTypeCentralProcessor:
		return "Central Processor"
	case ProcessorTypeMathProcessor:
		return "Math Processor"
	case ProcessorTypeDSP:
		return "DSP Processor"
	case ProcessorTypeVideoProcessor:
		return "Video Processor"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(pt))
	}
}

// ProcessorFamily identifies the processor family
type ProcessorFamily uint16

// Common processor family values (subset of full list)
const (
	ProcessorFamilyOther             ProcessorFamily = 0x01
	ProcessorFamilyUnknown           ProcessorFamily = 0x02
	ProcessorFamily8086              ProcessorFamily = 0x03
	ProcessorFamily80486             ProcessorFamily = 0x06
	ProcessorFamilyPentium           ProcessorFamily = 0x0B
	ProcessorFamilyPentiumPro        ProcessorFamily = 0x0C
	ProcessorFamilyPentiumII         ProcessorFamily = 0x0D
	ProcessorFamilyPentiumIII        ProcessorFamily = 0x11
	ProcessorFamilyPentium4          ProcessorFamily = 0x15
	ProcessorFamilyXeon              ProcessorFamily = 0xB3
	ProcessorFamilyPentiumM          ProcessorFamily = 0xB5
	ProcessorFamilyCeleron           ProcessorFamily = 0x14
	ProcessorFamilyCeleronD          ProcessorFamily = 0xB8
	ProcessorFamilyPentiumD          ProcessorFamily = 0xB9
	ProcessorFamilyPentiumDualCore   ProcessorFamily = 0xBC
	ProcessorFamilyCore2             ProcessorFamily = 0xBF
	ProcessorFamilyCoreSolo          ProcessorFamily = 0xBD
	ProcessorFamilyCoreDuo           ProcessorFamily = 0xBE
	ProcessorFamilyIntelCore         ProcessorFamily = 0xC0
	ProcessorFamilyIntelCoreM        ProcessorFamily = 0xC1
	ProcessorFamilyIntelCorem3       ProcessorFamily = 0xC2
	ProcessorFamilyIntelCorem5       ProcessorFamily = 0xC3
	ProcessorFamilyIntelCorem7       ProcessorFamily = 0xC4
	ProcessorFamilyAMDK5             ProcessorFamily = 0x18
	ProcessorFamilyAMDK6             ProcessorFamily = 0x19
	ProcessorFamilyAMDAthlon         ProcessorFamily = 0x1D
	ProcessorFamilyAMDAthlon64       ProcessorFamily = 0x53
	ProcessorFamilyAMDOpteron        ProcessorFamily = 0x54
	ProcessorFamilyAMDSempron        ProcessorFamily = 0x55
	ProcessorFamilyAMDTurion         ProcessorFamily = 0x56
	ProcessorFamilyAMDPhenom         ProcessorFamily = 0x5C
	ProcessorFamilyAMDPhenomII       ProcessorFamily = 0x5D
	ProcessorFamilyAMDAthlonII       ProcessorFamily = 0x5E
	ProcessorFamilyAMDOpteronSixCore ProcessorFamily = 0x5F
	ProcessorFamilyAMDRyzen          ProcessorFamily = 0x6B
	ProcessorFamilyAMDRyzen3         ProcessorFamily = 0x6C
	ProcessorFamilyAMDRyzen5         ProcessorFamily = 0x6D
	ProcessorFamilyAMDRyzen7         ProcessorFamily = 0x6E
	ProcessorFamilyAMDRyzen9         ProcessorFamily = 0x6F
	ProcessorFamilyARM               ProcessorFamily = 0x100
	ProcessorFamilyARMv7             ProcessorFamily = 0x101
	ProcessorFamilyARMv8             ProcessorFamily = 0x102
	ProcessorFamilyARMv9             ProcessorFamily = 0x103
	ProcessorFamilyAppleSilicon      ProcessorFamily = 0x110 // Custom, not official
	ProcessorFamilyIndicatorFamily2  ProcessorFamily = 0xFE  // Use ProcessorFamily2 field
)

// String returns a human-readable processor family description
func (pf ProcessorFamily) String() string {
	families := map[ProcessorFamily]string{
		ProcessorFamilyOther:             "Other",
		ProcessorFamilyUnknown:           "Unknown",
		ProcessorFamily8086:              "8086",
		ProcessorFamily80486:             "80486",
		ProcessorFamilyPentium:           "Pentium",
		ProcessorFamilyPentiumPro:        "Pentium Pro",
		ProcessorFamilyPentiumII:         "Pentium II",
		ProcessorFamilyPentiumIII:        "Pentium III",
		ProcessorFamilyPentium4:          "Pentium 4",
		ProcessorFamilyXeon:              "Xeon",
		ProcessorFamilyPentiumM:          "Pentium M",
		ProcessorFamilyCeleron:           "Celeron",
		ProcessorFamilyCeleronD:          "Celeron D",
		ProcessorFamilyPentiumD:          "Pentium D",
		ProcessorFamilyPentiumDualCore:   "Pentium Dual-Core",
		ProcessorFamilyCore2:             "Core 2",
		ProcessorFamilyCoreSolo:          "Core Solo",
		ProcessorFamilyCoreDuo:           "Core Duo",
		ProcessorFamilyIntelCore:         "Intel Core",
		ProcessorFamilyIntelCoreM:        "Intel Core M",
		ProcessorFamilyIntelCorem3:       "Intel Core m3",
		ProcessorFamilyIntelCorem5:       "Intel Core m5",
		ProcessorFamilyIntelCorem7:       "Intel Core m7",
		ProcessorFamilyAMDK5:             "AMD K5",
		ProcessorFamilyAMDK6:             "AMD K6",
		ProcessorFamilyAMDAthlon:         "AMD Athlon",
		ProcessorFamilyAMDAthlon64:       "AMD Athlon 64",
		ProcessorFamilyAMDOpteron:        "AMD Opteron",
		ProcessorFamilyAMDSempron:        "AMD Sempron",
		ProcessorFamilyAMDTurion:         "AMD Turion",
		ProcessorFamilyAMDPhenom:         "AMD Phenom",
		ProcessorFamilyAMDPhenomII:       "AMD Phenom II",
		ProcessorFamilyAMDAthlonII:       "AMD Athlon II",
		ProcessorFamilyAMDOpteronSixCore: "AMD Opteron Six-Core",
		ProcessorFamilyAMDRyzen:          "AMD Ryzen",
		ProcessorFamilyAMDRyzen3:         "AMD Ryzen 3",
		ProcessorFamilyAMDRyzen5:         "AMD Ryzen 5",
		ProcessorFamilyAMDRyzen7:         "AMD Ryzen 7",
		ProcessorFamilyAMDRyzen9:         "AMD Ryzen 9",
		ProcessorFamilyARM:               "ARM",
		ProcessorFamilyARMv7:             "ARMv7",
		ProcessorFamilyARMv8:             "ARMv8",
		ProcessorFamilyARMv9:             "ARMv9",
	}

	if name, ok := families[pf]; ok {
		return name
	}
	return fmt.Sprintf("Unknown (0x%04X)", uint16(pf))
}

// Voltage represents processor voltage information
type Voltage uint8

// VoltageValue returns the voltage value in volts, or a list of supported voltages
func (v Voltage) VoltageValue() (float64, bool) {
	if v&0x80 != 0 {
		// Bit 7 set: bits 0-6 contain voltage * 10
		return float64(v&0x7F) / 10.0, true
	}
	return 0, false
}

// SupportedVoltages returns the list of supported voltages when bit 7 is not set
func (v Voltage) SupportedVoltages() []float64 {
	if v&0x80 != 0 {
		return nil
	}

	var voltages []float64
	if v&0x01 != 0 {
		voltages = append(voltages, 5.0)
	}
	if v&0x02 != 0 {
		voltages = append(voltages, 3.3)
	}
	if v&0x04 != 0 {
		voltages = append(voltages, 2.9)
	}
	return voltages
}

// String returns a human-readable voltage description
func (v Voltage) String() string {
	if val, ok := v.VoltageValue(); ok {
		return fmt.Sprintf("%.1f V", val)
	}

	voltages := v.SupportedVoltages()
	if len(voltages) == 0 {
		return "Unknown"
	}

	result := ""
	for i, volt := range voltages {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("%.1f V", volt)
	}
	return result
}

// ProcessorStatus represents the processor status field
type ProcessorStatus uint8

// CPUStatus returns the CPU status (bits 0-2)
func (ps ProcessorStatus) CPUStatus() uint8 {
	return uint8(ps) & 0x07
}

// IsPopulated returns true if a CPU is present in the socket
func (ps ProcessorStatus) IsPopulated() bool {
	return (ps & 0x40) != 0
}

// String returns a human-readable status description
func (ps ProcessorStatus) String() string {
	if !ps.IsPopulated() {
		return "Unpopulated"
	}

	status := ps.CPUStatus()
	switch status {
	case 0:
		return "Unknown"
	case 1:
		return "CPU Enabled"
	case 2:
		return "CPU Disabled by User"
	case 3:
		return "CPU Disabled by BIOS"
	case 4:
		return "CPU Idle"
	case 7:
		return "Other"
	default:
		return fmt.Sprintf("Reserved (0x%02X)", status)
	}
}

// ProcessorUpgrade identifies the processor upgrade method
type ProcessorUpgrade uint8

// Processor upgrade values
const (
	ProcessorUpgradeOther                ProcessorUpgrade = 0x01
	ProcessorUpgradeUnknown              ProcessorUpgrade = 0x02
	ProcessorUpgradeDaughterBoard        ProcessorUpgrade = 0x03
	ProcessorUpgradeZIFSocket            ProcessorUpgrade = 0x04
	ProcessorUpgradeReplaceablePiggyback ProcessorUpgrade = 0x05
	ProcessorUpgradeNone                 ProcessorUpgrade = 0x06
	ProcessorUpgradeLIFSocket            ProcessorUpgrade = 0x07
	ProcessorUpgradeSlot1                ProcessorUpgrade = 0x08
	ProcessorUpgradeSlot2                ProcessorUpgrade = 0x09
	ProcessorUpgrade370PinSocket         ProcessorUpgrade = 0x0A
	ProcessorUpgradeSlotA                ProcessorUpgrade = 0x0B
	ProcessorUpgradeSlotM                ProcessorUpgrade = 0x0C
	ProcessorUpgradeSocket423            ProcessorUpgrade = 0x0D
	ProcessorUpgradeSocketA              ProcessorUpgrade = 0x0E
	ProcessorUpgradeSocket478            ProcessorUpgrade = 0x0F
	ProcessorUpgradeSocket754            ProcessorUpgrade = 0x10
	ProcessorUpgradeSocket940            ProcessorUpgrade = 0x11
	ProcessorUpgradeSocket939            ProcessorUpgrade = 0x12
	ProcessorUpgradeSocketmPGA604        ProcessorUpgrade = 0x13
	ProcessorUpgradeSocketLGA771         ProcessorUpgrade = 0x14
	ProcessorUpgradeSocketLGA775         ProcessorUpgrade = 0x15
	ProcessorUpgradeSocketS1             ProcessorUpgrade = 0x16
	ProcessorUpgradeSocketAM2            ProcessorUpgrade = 0x17
	ProcessorUpgradeSocketF1207          ProcessorUpgrade = 0x18
	ProcessorUpgradeSocketLGA1366        ProcessorUpgrade = 0x19
	ProcessorUpgradeSocketG34            ProcessorUpgrade = 0x1A
	ProcessorUpgradeSocketAM3            ProcessorUpgrade = 0x1B
	ProcessorUpgradeSocketC32            ProcessorUpgrade = 0x1C
	ProcessorUpgradeSocketLGA1156        ProcessorUpgrade = 0x1D
	ProcessorUpgradeSocketLGA1567        ProcessorUpgrade = 0x1E
	ProcessorUpgradeSocketPGA988A        ProcessorUpgrade = 0x1F
	ProcessorUpgradeSocketBGA1288        ProcessorUpgrade = 0x20
	ProcessorUpgradeSocketrPGA988B       ProcessorUpgrade = 0x21
	ProcessorUpgradeSocketBGA1023        ProcessorUpgrade = 0x22
	ProcessorUpgradeSocketBGA1224        ProcessorUpgrade = 0x23
	ProcessorUpgradeSocketLGA1155        ProcessorUpgrade = 0x24
	ProcessorUpgradeSocketLGA1356        ProcessorUpgrade = 0x25
	ProcessorUpgradeSocketLGA2011        ProcessorUpgrade = 0x26
	ProcessorUpgradeSocketFS1            ProcessorUpgrade = 0x27
	ProcessorUpgradeSocketFS2            ProcessorUpgrade = 0x28
	ProcessorUpgradeSocketFM1            ProcessorUpgrade = 0x29
	ProcessorUpgradeSocketFM2            ProcessorUpgrade = 0x2A
	ProcessorUpgradeSocketLGA2011v3      ProcessorUpgrade = 0x2B
	ProcessorUpgradeSocketLGA1356v3      ProcessorUpgrade = 0x2C
	ProcessorUpgradeSocketLGA1150        ProcessorUpgrade = 0x2D
	ProcessorUpgradeSocketBGA1168        ProcessorUpgrade = 0x2E
	ProcessorUpgradeSocketBGA1234        ProcessorUpgrade = 0x2F
	ProcessorUpgradeSocketBGA1364        ProcessorUpgrade = 0x30
	ProcessorUpgradeSocketAM4            ProcessorUpgrade = 0x31
	ProcessorUpgradeSocketLGA1151        ProcessorUpgrade = 0x32
	ProcessorUpgradeSocketBGA1356        ProcessorUpgrade = 0x33
	ProcessorUpgradeSocketBGA1440        ProcessorUpgrade = 0x34
	ProcessorUpgradeSocketBGA1515        ProcessorUpgrade = 0x35
	ProcessorUpgradeSocketLGA3647v1      ProcessorUpgrade = 0x36
	ProcessorUpgradeSocketSP3            ProcessorUpgrade = 0x37
	ProcessorUpgradeSocketSP3r2          ProcessorUpgrade = 0x38
	ProcessorUpgradeSocketLGA2066        ProcessorUpgrade = 0x39
	ProcessorUpgradeSocketBGA1392        ProcessorUpgrade = 0x3A
	ProcessorUpgradeSocketBGA1510        ProcessorUpgrade = 0x3B
	ProcessorUpgradeSocketBGA1528        ProcessorUpgrade = 0x3C
	ProcessorUpgradeSocketLGA4189        ProcessorUpgrade = 0x3D
	ProcessorUpgradeSocketLGA1200        ProcessorUpgrade = 0x3E
	ProcessorUpgradeSocketLGA4677        ProcessorUpgrade = 0x3F
	ProcessorUpgradeSocketLGA1700        ProcessorUpgrade = 0x40
	ProcessorUpgradeSocketBGA1744        ProcessorUpgrade = 0x41
	ProcessorUpgradeSocketBGA1781        ProcessorUpgrade = 0x42
	ProcessorUpgradeSocketBGA1211        ProcessorUpgrade = 0x43
	ProcessorUpgradeSocketBGA2422        ProcessorUpgrade = 0x44
	ProcessorUpgradeSocketLGA1211        ProcessorUpgrade = 0x45
	ProcessorUpgradeSocketLGA2422        ProcessorUpgrade = 0x46
	ProcessorUpgradeSocketLGA5773        ProcessorUpgrade = 0x47
	ProcessorUpgradeSocketBGA5773        ProcessorUpgrade = 0x48
	ProcessorUpgradeSocketAM5            ProcessorUpgrade = 0x49
	ProcessorUpgradeSocketSP5            ProcessorUpgrade = 0x4A
	ProcessorUpgradeSocketSP6            ProcessorUpgrade = 0x4B
	ProcessorUpgradeSocketBGA883         ProcessorUpgrade = 0x4C
	ProcessorUpgradeSocketBGA1190        ProcessorUpgrade = 0x4D
	ProcessorUpgradeSocketBGA4129        ProcessorUpgrade = 0x4E
	ProcessorUpgradeSocketLGA4710        ProcessorUpgrade = 0x4F
	ProcessorUpgradeSocketLGA7529        ProcessorUpgrade = 0x50
)

// String returns a human-readable upgrade type description
func (pu ProcessorUpgrade) String() string {
	upgrades := map[ProcessorUpgrade]string{
		ProcessorUpgradeOther:         "Other",
		ProcessorUpgradeUnknown:       "Unknown",
		ProcessorUpgradeDaughterBoard: "Daughter Board",
		ProcessorUpgradeZIFSocket:     "ZIF Socket",
		ProcessorUpgradeNone:          "None",
		ProcessorUpgradeSocketAM4:     "Socket AM4",
		ProcessorUpgradeSocketAM5:     "Socket AM5",
		ProcessorUpgradeSocketLGA1151: "Socket LGA1151",
		ProcessorUpgradeSocketLGA1200: "Socket LGA1200",
		ProcessorUpgradeSocketLGA1700: "Socket LGA1700",
		ProcessorUpgradeSocketLGA2011: "Socket LGA2011-3",
		ProcessorUpgradeSocketLGA4189: "Socket LGA4189",
		ProcessorUpgradeSocketLGA4677: "Socket LGA4677",
		ProcessorUpgradeSocketSP3:     "Socket SP3",
		ProcessorUpgradeSocketSP5:     "Socket SP5",
	}

	if name, ok := upgrades[pu]; ok {
		return name
	}
	return fmt.Sprintf("Socket/Slot (0x%02X)", uint8(pu))
}

// ProcessorCharacteristics represents processor characteristics
type ProcessorCharacteristics uint16

// Processor characteristics bit definitions
const (
	CharReserved0               ProcessorCharacteristics = 1 << 0 // Reserved
	CharUnknown                 ProcessorCharacteristics = 1 << 1 // Unknown
	Char64BitCapable            ProcessorCharacteristics = 1 << 2 // 64-bit Capable
	CharMultiCore               ProcessorCharacteristics = 1 << 3 // Multi-Core
	CharHardwareThread          ProcessorCharacteristics = 1 << 4 // Hardware Thread
	CharExecuteProtection       ProcessorCharacteristics = 1 << 5 // Execute Protection
	CharEnhancedVirtualization  ProcessorCharacteristics = 1 << 6 // Enhanced Virtualization
	CharPowerPerformanceControl ProcessorCharacteristics = 1 << 7 // Power/Performance Control
	Char128BitCapable           ProcessorCharacteristics = 1 << 8 // 128-bit Capable
	CharArm64SocID              ProcessorCharacteristics = 1 << 9 // Arm64 SoC ID
)

// Has checks if a characteristic is set
func (pc ProcessorCharacteristics) Has(flag ProcessorCharacteristics) bool {
	return pc&flag != 0
}

// Is64Bit returns true if the processor is 64-bit capable
func (pc ProcessorCharacteristics) Is64Bit() bool {
	return pc.Has(Char64BitCapable)
}

// IsMultiCore returns true if the processor is multi-core
func (pc ProcessorCharacteristics) IsMultiCore() bool {
	return pc.Has(CharMultiCore)
}

// Parse parses a Processor Information structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*ProcessorInfo, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 26 bytes (SMBIOS 2.0)
	if len(s.Data) < 26 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &ProcessorInfo{
		Header:                s.Header,
		SocketDesignation:     s.GetString(s.GetByte(0x04)),
		ProcessorType:         ProcessorType(s.GetByte(0x05)),
		ProcessorFamily:       ProcessorFamily(s.GetByte(0x06)),
		ProcessorManufacturer: s.GetString(s.GetByte(0x07)),
		ProcessorID:           s.GetQWord(0x08),
		ProcessorVersion:      s.GetString(s.GetByte(0x10)),
		Voltage:               Voltage(s.GetByte(0x11)),
		ExternalClock:         s.GetWord(0x12),
		MaxSpeed:              s.GetWord(0x14),
		CurrentSpeed:          s.GetWord(0x16),
		Status:                ProcessorStatus(s.GetByte(0x18)),
		ProcessorUpgrade:      ProcessorUpgrade(s.GetByte(0x19)),
	}

	// SMBIOS 2.1+
	if len(s.Data) >= 32 {
		info.L1CacheHandle = s.GetWord(0x1A)
		info.L2CacheHandle = s.GetWord(0x1C)
		info.L3CacheHandle = s.GetWord(0x1E)
	}

	// SMBIOS 2.3+
	if len(s.Data) >= 35 {
		info.SerialNumber = s.GetString(s.GetByte(0x20))
		info.AssetTag = s.GetString(s.GetByte(0x21))
		info.PartNumber = s.GetString(s.GetByte(0x22))
	}

	// SMBIOS 2.5+
	if len(s.Data) >= 40 {
		info.CoreCount = s.GetByte(0x23)
		info.CoreEnabled = s.GetByte(0x24)
		info.ThreadCount = s.GetByte(0x25)
		info.ProcessorCharacteristics = ProcessorCharacteristics(s.GetWord(0x26))
	}

	// SMBIOS 2.6+
	if len(s.Data) >= 42 {
		info.ProcessorFamily2 = s.GetWord(0x28)
		if info.ProcessorFamily == ProcessorFamilyIndicatorFamily2 {
			info.ProcessorFamily = ProcessorFamily(info.ProcessorFamily2)
		}
	}

	// SMBIOS 3.0+
	if len(s.Data) >= 48 {
		info.CoreCount2 = s.GetWord(0x2A)
		info.CoreEnabled2 = s.GetWord(0x2C)
		info.ThreadCount2 = s.GetWord(0x2E)
	}

	// SMBIOS 3.6+
	if len(s.Data) >= 50 {
		info.ThreadEnabled = s.GetWord(0x30)
	}

	return info, nil
}

// Get retrieves the first Processor Information from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*ProcessorInfo, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Processor Information structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*ProcessorInfo, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var processors []*ProcessorInfo
	for i := range structures {
		proc, err := Parse(&structures[i])
		if err == nil {
			processors = append(processors, proc)
		}
	}

	if len(processors) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return processors, nil
}

// GetCoreCount returns the actual core count (handling the extension field)
func (p *ProcessorInfo) GetCoreCount() uint16 {
	if p.CoreCount == 0xFF && p.CoreCount2 != 0 {
		return p.CoreCount2
	}
	return uint16(p.CoreCount)
}

// GetCoreEnabled returns the actual enabled core count (handling the extension field)
func (p *ProcessorInfo) GetCoreEnabled() uint16 {
	if p.CoreEnabled == 0xFF && p.CoreEnabled2 != 0 {
		return p.CoreEnabled2
	}
	return uint16(p.CoreEnabled)
}

// GetThreadCount returns the actual thread count (handling the extension field)
func (p *ProcessorInfo) GetThreadCount() uint16 {
	if p.ThreadCount == 0xFF && p.ThreadCount2 != 0 {
		return p.ThreadCount2
	}
	return uint16(p.ThreadCount)
}

// DisplayName returns a display-friendly processor name
func (p *ProcessorInfo) DisplayName() string {
	if p.ProcessorVersion != "" {
		return p.ProcessorVersion
	}
	return p.ProcessorFamily.String()
}
