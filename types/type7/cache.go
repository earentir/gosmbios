// Package type7 implements SMBIOS Type 7 - Cache Information
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type7

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Cache Information
const StructureType uint8 = 7

// CacheInfo represents Type 7 - Cache Information
type CacheInfo struct {
	Header              gosmbios.Header
	SocketDesignation   string
	Configuration       CacheConfiguration
	MaximumSize         uint32 // In KB
	InstalledSize       uint32 // In KB
	SupportedSRAMType   SRAMType
	CurrentSRAMType     SRAMType
	CacheSpeed          uint8               // In nanoseconds (SMBIOS 2.1+)
	ErrorCorrectionType ErrorCorrectionType // SMBIOS 2.1+
	SystemCacheType     CacheType           // SMBIOS 2.1+
	Associativity       CacheAssociativity  // SMBIOS 2.1+
	MaximumSize2        uint32              // In KB (SMBIOS 3.1+, used when MaximumSize is 0xFFFF)
	InstalledSize2      uint32              // In KB (SMBIOS 3.1+, used when InstalledSize is 0xFFFF)
}

// CacheConfiguration represents the cache configuration word
type CacheConfiguration uint16

// Level returns the cache level (1, 2, 3, etc.)
func (cc CacheConfiguration) Level() int {
	return int((cc & 0x07) + 1)
}

// IsSocketed returns true if the cache is socketed
func (cc CacheConfiguration) IsSocketed() bool {
	return (cc & 0x08) != 0
}

// Location returns the cache location
func (cc CacheConfiguration) Location() CacheLocation {
	return CacheLocation((cc >> 5) & 0x03)
}

// IsEnabled returns true if the cache is enabled
func (cc CacheConfiguration) IsEnabled() bool {
	return (cc & 0x80) != 0
}

// OperationalMode returns the cache operational mode
func (cc CacheConfiguration) OperationalMode() CacheMode {
	return CacheMode((cc >> 8) & 0x03)
}

// CacheLocation represents the cache location
type CacheLocation uint8

// Cache location values
const (
	CacheLocationInternal CacheLocation = 0x00
	CacheLocationExternal CacheLocation = 0x01
	CacheLocationReserved CacheLocation = 0x02
	CacheLocationUnknown  CacheLocation = 0x03
)

// String returns a human-readable cache location description
func (cl CacheLocation) String() string {
	switch cl {
	case CacheLocationInternal:
		return "Internal"
	case CacheLocationExternal:
		return "External"
	case CacheLocationReserved:
		return "Reserved"
	case CacheLocationUnknown:
		return "Unknown"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(cl))
	}
}

// CacheMode represents the cache operational mode
type CacheMode uint8

// Cache mode values
const (
	CacheModeWriteThrough CacheMode = 0x00
	CacheModeWriteBack    CacheMode = 0x01
	CacheModeVaries       CacheMode = 0x02
	CacheModeUnknown      CacheMode = 0x03
)

// String returns a human-readable cache mode description
func (cm CacheMode) String() string {
	switch cm {
	case CacheModeWriteThrough:
		return "Write Through"
	case CacheModeWriteBack:
		return "Write Back"
	case CacheModeVaries:
		return "Varies with Memory Address"
	case CacheModeUnknown:
		return "Unknown"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(cm))
	}
}

// SRAMType represents supported SRAM types
type SRAMType uint16

// SRAM type bit definitions
const (
	SRAMTypeOther         SRAMType = 1 << 0
	SRAMTypeUnknown       SRAMType = 1 << 1
	SRAMTypeNonBurst      SRAMType = 1 << 2
	SRAMTypeBurst         SRAMType = 1 << 3
	SRAMTypePipelineBurst SRAMType = 1 << 4
	SRAMTypeSynchronous   SRAMType = 1 << 5
	SRAMTypeAsynchronous  SRAMType = 1 << 6
)

// Has checks if an SRAM type is set
func (st SRAMType) Has(flag SRAMType) bool {
	return st&flag != 0
}

// String returns a human-readable SRAM type description
func (st SRAMType) String() string {
	var types []string
	if st.Has(SRAMTypeOther) {
		types = append(types, "Other")
	}
	if st.Has(SRAMTypeUnknown) {
		types = append(types, "Unknown")
	}
	if st.Has(SRAMTypeNonBurst) {
		types = append(types, "Non-Burst")
	}
	if st.Has(SRAMTypeBurst) {
		types = append(types, "Burst")
	}
	if st.Has(SRAMTypePipelineBurst) {
		types = append(types, "Pipeline Burst")
	}
	if st.Has(SRAMTypeSynchronous) {
		types = append(types, "Synchronous")
	}
	if st.Has(SRAMTypeAsynchronous) {
		types = append(types, "Asynchronous")
	}

	if len(types) == 0 {
		return "None"
	}

	result := types[0]
	for i := 1; i < len(types); i++ {
		result += ", " + types[i]
	}
	return result
}

// ErrorCorrectionType represents the cache error correction type
type ErrorCorrectionType uint8

// Error correction type values
const (
	ErrorCorrectionOther        ErrorCorrectionType = 0x01
	ErrorCorrectionUnknown      ErrorCorrectionType = 0x02
	ErrorCorrectionNone         ErrorCorrectionType = 0x03
	ErrorCorrectionParity       ErrorCorrectionType = 0x04
	ErrorCorrectionSingleBitECC ErrorCorrectionType = 0x05
	ErrorCorrectionMultiBitECC  ErrorCorrectionType = 0x06
)

// String returns a human-readable error correction type description
func (ect ErrorCorrectionType) String() string {
	switch ect {
	case ErrorCorrectionOther:
		return "Other"
	case ErrorCorrectionUnknown:
		return "Unknown"
	case ErrorCorrectionNone:
		return "None"
	case ErrorCorrectionParity:
		return "Parity"
	case ErrorCorrectionSingleBitECC:
		return "Single-bit ECC"
	case ErrorCorrectionMultiBitECC:
		return "Multi-bit ECC"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(ect))
	}
}

// CacheType represents the system cache type
type CacheType uint8

// Cache type values
const (
	CacheTypeOther       CacheType = 0x01
	CacheTypeUnknown     CacheType = 0x02
	CacheTypeInstruction CacheType = 0x03
	CacheTypeData        CacheType = 0x04
	CacheTypeUnified     CacheType = 0x05
)

// String returns a human-readable cache type description
func (ct CacheType) String() string {
	switch ct {
	case CacheTypeOther:
		return "Other"
	case CacheTypeUnknown:
		return "Unknown"
	case CacheTypeInstruction:
		return "Instruction"
	case CacheTypeData:
		return "Data"
	case CacheTypeUnified:
		return "Unified"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(ct))
	}
}

// CacheAssociativity represents cache associativity
type CacheAssociativity uint8

// Cache associativity values
const (
	AssociativityOther            CacheAssociativity = 0x01
	AssociativityUnknown          CacheAssociativity = 0x02
	AssociativityDirectMapped     CacheAssociativity = 0x03
	Associativity2Way             CacheAssociativity = 0x04
	Associativity4Way             CacheAssociativity = 0x05
	AssociativityFullyAssociative CacheAssociativity = 0x06
	Associativity8Way             CacheAssociativity = 0x07
	Associativity16Way            CacheAssociativity = 0x08
	Associativity12Way            CacheAssociativity = 0x09
	Associativity24Way            CacheAssociativity = 0x0A
	Associativity32Way            CacheAssociativity = 0x0B
	Associativity48Way            CacheAssociativity = 0x0C
	Associativity64Way            CacheAssociativity = 0x0D
	Associativity20Way            CacheAssociativity = 0x0E
)

// String returns a human-readable associativity description
func (ca CacheAssociativity) String() string {
	switch ca {
	case AssociativityOther:
		return "Other"
	case AssociativityUnknown:
		return "Unknown"
	case AssociativityDirectMapped:
		return "Direct Mapped"
	case Associativity2Way:
		return "2-way Set-Associative"
	case Associativity4Way:
		return "4-way Set-Associative"
	case AssociativityFullyAssociative:
		return "Fully Associative"
	case Associativity8Way:
		return "8-way Set-Associative"
	case Associativity16Way:
		return "16-way Set-Associative"
	case Associativity12Way:
		return "12-way Set-Associative"
	case Associativity24Way:
		return "24-way Set-Associative"
	case Associativity32Way:
		return "32-way Set-Associative"
	case Associativity48Way:
		return "48-way Set-Associative"
	case Associativity64Way:
		return "64-way Set-Associative"
	case Associativity20Way:
		return "20-way Set-Associative"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(ca))
	}
}

// Parse parses a Cache Information structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*CacheInfo, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 15 bytes (SMBIOS 2.0)
	if len(s.Data) < 15 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &CacheInfo{
		Header:            s.Header,
		SocketDesignation: s.GetString(s.GetByte(0x04)),
		Configuration:     CacheConfiguration(s.GetWord(0x05)),
		SupportedSRAMType: SRAMType(s.GetWord(0x0B)),
		CurrentSRAMType:   SRAMType(s.GetWord(0x0D)),
	}

	// Parse maximum size (16-bit field)
	maxSizeRaw := s.GetWord(0x07)
	if maxSizeRaw&0x8000 != 0 {
		// Bit 15 set: granularity is 64KB
		info.MaximumSize = uint32(maxSizeRaw&0x7FFF) * 64
	} else {
		// Bit 15 clear: granularity is 1KB
		info.MaximumSize = uint32(maxSizeRaw)
	}

	// Parse installed size (16-bit field)
	instSizeRaw := s.GetWord(0x09)
	if instSizeRaw&0x8000 != 0 {
		// Bit 15 set: granularity is 64KB
		info.InstalledSize = uint32(instSizeRaw&0x7FFF) * 64
	} else {
		// Bit 15 clear: granularity is 1KB
		info.InstalledSize = uint32(instSizeRaw)
	}

	// SMBIOS 2.1+
	if len(s.Data) >= 19 {
		info.CacheSpeed = s.GetByte(0x0F)
		info.ErrorCorrectionType = ErrorCorrectionType(s.GetByte(0x10))
		info.SystemCacheType = CacheType(s.GetByte(0x11))
		info.Associativity = CacheAssociativity(s.GetByte(0x12))
	}

	// SMBIOS 3.1+
	if len(s.Data) >= 27 {
		// Maximum Size 2
		maxSize2Raw := s.GetDWord(0x13)
		if maxSizeRaw == 0xFFFF {
			if maxSize2Raw&0x80000000 != 0 {
				// Bit 31 set: granularity is 64KB
				info.MaximumSize = (maxSize2Raw & 0x7FFFFFFF) * 64
			} else {
				// Bit 31 clear: granularity is 1KB
				info.MaximumSize = maxSize2Raw
			}
		}
		info.MaximumSize2 = maxSize2Raw

		// Installed Size 2
		instSize2Raw := s.GetDWord(0x17)
		if instSizeRaw == 0xFFFF {
			if instSize2Raw&0x80000000 != 0 {
				// Bit 31 set: granularity is 64KB
				info.InstalledSize = (instSize2Raw & 0x7FFFFFFF) * 64
			} else {
				// Bit 31 clear: granularity is 1KB
				info.InstalledSize = instSize2Raw
			}
		}
		info.InstalledSize2 = instSize2Raw
	}

	return info, nil
}

// Get retrieves the first Cache Information from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*CacheInfo, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Cache Information structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*CacheInfo, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var caches []*CacheInfo
	for i := range structures {
		cache, err := Parse(&structures[i])
		if err == nil {
			caches = append(caches, cache)
		}
	}

	if len(caches) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return caches, nil
}

// GetByHandle retrieves a Cache Information structure by its handle
func GetByHandle(sm *gosmbios.SMBIOS, handle uint16) (*CacheInfo, error) {
	structures := sm.GetStructures(StructureType)
	for i := range structures {
		if structures[i].Header.Handle == handle {
			return Parse(&structures[i])
		}
	}
	return nil, gosmbios.ErrNotFound
}

// MaximumSizeString returns a human-readable maximum size
func (c *CacheInfo) MaximumSizeString() string {
	return formatSize(c.MaximumSize)
}

// InstalledSizeString returns a human-readable installed size
func (c *CacheInfo) InstalledSizeString() string {
	return formatSize(c.InstalledSize)
}

// Level returns the cache level (1, 2, 3, etc.)
func (c *CacheInfo) Level() int {
	return c.Configuration.Level()
}

// formatSize formats a size in KB to a human-readable string
func formatSize(sizeKB uint32) string {
	if sizeKB == 0 {
		return "Unknown"
	}
	if sizeKB >= 1024*1024 {
		return fmt.Sprintf("%d GB", sizeKB/(1024*1024))
	}
	if sizeKB >= 1024 {
		return fmt.Sprintf("%d MB", sizeKB/1024)
	}
	return fmt.Sprintf("%d KB", sizeKB)
}
