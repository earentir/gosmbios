// Package type5 implements SMBIOS Type 5 - Memory Controller Information (Obsolete)
// Per DSP0134 SMBIOS Reference Specification 3.9.0
// This structure is obsolete starting with SMBIOS 2.1 specification
package type5

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Memory Controller Information
const StructureType uint8 = 5

// MemoryController represents Type 5 - Memory Controller Information (Obsolete)
type MemoryController struct {
	Header                        gosmbios.Header
	ErrorDetectingMethod          ErrorDetectingMethod
	ErrorCorrectingCapability     ErrorCorrectingCapability
	SupportedInterleave           Interleave
	CurrentInterleave             Interleave
	MaximumMemoryModuleSize       uint8    // Size = 2^n MB
	SupportedSpeeds               SpeedSet // Bit field
	SupportedMemoryTypes          uint16   // Bit field
	MemoryModuleVoltage           Voltage  // Bit field
	NumberOfAssociatedMemorySlots uint8
	MemoryModuleConfigHandles     []uint16 // Handles of Type 6 structures
	EnabledErrorCorrectingCaps    ErrorCorrectingCapability
}

// ErrorDetectingMethod represents memory error detecting methods
type ErrorDetectingMethod uint8

const (
	ErrorDetectOther        ErrorDetectingMethod = 0x01
	ErrorDetectUnknown      ErrorDetectingMethod = 0x02
	ErrorDetectNone         ErrorDetectingMethod = 0x03
	ErrorDetect8BitParity   ErrorDetectingMethod = 0x04
	ErrorDetect32BitECC     ErrorDetectingMethod = 0x05
	ErrorDetect64BitECC     ErrorDetectingMethod = 0x06
	ErrorDetect128BitECC    ErrorDetectingMethod = 0x07
	ErrorDetectCRC          ErrorDetectingMethod = 0x08
)

func (e ErrorDetectingMethod) String() string {
	switch e {
	case ErrorDetectOther:
		return "Other"
	case ErrorDetectUnknown:
		return "Unknown"
	case ErrorDetectNone:
		return "None"
	case ErrorDetect8BitParity:
		return "8-bit Parity"
	case ErrorDetect32BitECC:
		return "32-bit ECC"
	case ErrorDetect64BitECC:
		return "64-bit ECC"
	case ErrorDetect128BitECC:
		return "128-bit ECC"
	case ErrorDetectCRC:
		return "CRC"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(e))
	}
}

// ErrorCorrectingCapability represents error correcting capabilities
type ErrorCorrectingCapability uint8

const (
	ECCCapOther           ErrorCorrectingCapability = 1 << 0
	ECCCapUnknown         ErrorCorrectingCapability = 1 << 1
	ECCCapNone            ErrorCorrectingCapability = 1 << 2
	ECCCapSingleBitECC    ErrorCorrectingCapability = 1 << 3
	ECCCapDoubleBitECC    ErrorCorrectingCapability = 1 << 4
	ECCCapErrorScrubbing  ErrorCorrectingCapability = 1 << 5
)

func (e ErrorCorrectingCapability) String() string {
	if e == 0 {
		return "None"
	}
	var caps []string
	if e&ECCCapOther != 0 {
		caps = append(caps, "Other")
	}
	if e&ECCCapUnknown != 0 {
		caps = append(caps, "Unknown")
	}
	if e&ECCCapNone != 0 {
		caps = append(caps, "None")
	}
	if e&ECCCapSingleBitECC != 0 {
		caps = append(caps, "Single-bit ECC")
	}
	if e&ECCCapDoubleBitECC != 0 {
		caps = append(caps, "Double-bit ECC")
	}
	if e&ECCCapErrorScrubbing != 0 {
		caps = append(caps, "Error Scrubbing")
	}
	if len(caps) == 0 {
		return fmt.Sprintf("Unknown (0x%02X)", uint8(e))
	}
	return fmt.Sprintf("%v", caps)
}

// Interleave represents memory interleave support
type Interleave uint8

const (
	InterleaveOther     Interleave = 0x01
	InterleaveUnknown   Interleave = 0x02
	InterleaveOneWay    Interleave = 0x03
	InterleaveTwoWay    Interleave = 0x04
	InterleaveFourWay   Interleave = 0x05
	InterleaveEightWay  Interleave = 0x06
	InterleaveSixteenWay Interleave = 0x07
)

func (i Interleave) String() string {
	switch i {
	case InterleaveOther:
		return "Other"
	case InterleaveUnknown:
		return "Unknown"
	case InterleaveOneWay:
		return "One-Way Interleave"
	case InterleaveTwoWay:
		return "Two-Way Interleave"
	case InterleaveFourWay:
		return "Four-Way Interleave"
	case InterleaveEightWay:
		return "Eight-Way Interleave"
	case InterleaveSixteenWay:
		return "Sixteen-Way Interleave"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(i))
	}
}

// SpeedSet represents supported memory speeds
type SpeedSet uint16

const (
	SpeedOther SpeedSet = 1 << 0
	SpeedUnknown SpeedSet = 1 << 1
	Speed70ns SpeedSet = 1 << 2
	Speed60ns SpeedSet = 1 << 3
	Speed50ns SpeedSet = 1 << 4
)

func (s SpeedSet) String() string {
	var speeds []string
	if s&SpeedOther != 0 {
		speeds = append(speeds, "Other")
	}
	if s&SpeedUnknown != 0 {
		speeds = append(speeds, "Unknown")
	}
	if s&Speed70ns != 0 {
		speeds = append(speeds, "70ns")
	}
	if s&Speed60ns != 0 {
		speeds = append(speeds, "60ns")
	}
	if s&Speed50ns != 0 {
		speeds = append(speeds, "50ns")
	}
	if len(speeds) == 0 {
		return "None"
	}
	return fmt.Sprintf("%v", speeds)
}

// Voltage represents memory module voltage requirements
type Voltage uint8

const (
	Voltage5V  Voltage = 1 << 0
	Voltage3_3V Voltage = 1 << 1
	Voltage2_9V Voltage = 1 << 2
)

func (v Voltage) String() string {
	var volts []string
	if v&Voltage5V != 0 {
		volts = append(volts, "5V")
	}
	if v&Voltage3_3V != 0 {
		volts = append(volts, "3.3V")
	}
	if v&Voltage2_9V != 0 {
		volts = append(volts, "2.9V")
	}
	if len(volts) == 0 {
		return "None"
	}
	return fmt.Sprintf("%v", volts)
}

// Parse parses a Memory Controller Information structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*MemoryController, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 15 bytes
	if len(s.Data) < 15 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &MemoryController{
		Header:                        s.Header,
		ErrorDetectingMethod:          ErrorDetectingMethod(s.GetByte(0x04)),
		ErrorCorrectingCapability:     ErrorCorrectingCapability(s.GetByte(0x05)),
		SupportedInterleave:           Interleave(s.GetByte(0x06)),
		CurrentInterleave:             Interleave(s.GetByte(0x07)),
		MaximumMemoryModuleSize:       s.GetByte(0x08),
		SupportedSpeeds:               SpeedSet(s.GetWord(0x09)),
		SupportedMemoryTypes:          s.GetWord(0x0B),
		MemoryModuleVoltage:           Voltage(s.GetByte(0x0D)),
		NumberOfAssociatedMemorySlots: s.GetByte(0x0E),
	}

	// Read memory module configuration handles
	numSlots := int(info.NumberOfAssociatedMemorySlots)
	offset := 0x0F
	for i := 0; i < numSlots && offset+1 < len(s.Data); i++ {
		handle := s.GetWord(offset)
		info.MemoryModuleConfigHandles = append(info.MemoryModuleConfigHandles, handle)
		offset += 2
	}

	// Enabled error correcting capabilities (after handles)
	if offset < len(s.Data) {
		info.EnabledErrorCorrectingCaps = ErrorCorrectingCapability(s.GetByte(offset))
	}

	return info, nil
}

// Get retrieves the Memory Controller Information from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*MemoryController, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Memory Controller structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*MemoryController, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var controllers []*MemoryController
	for i := range structures {
		ctrl, err := Parse(&structures[i])
		if err == nil {
			controllers = append(controllers, ctrl)
		}
	}

	if len(controllers) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return controllers, nil
}

// MaxModuleSizeMB returns the maximum memory module size in MB
func (m *MemoryController) MaxModuleSizeMB() uint64 {
	return uint64(1) << m.MaximumMemoryModuleSize
}
