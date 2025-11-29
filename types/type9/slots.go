// Package type9 implements SMBIOS Type 9 - System Slots
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type9

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for System Slots
const StructureType uint8 = 9

// SlotInfo represents Type 9 - System Slots
type SlotInfo struct {
	Header               gosmbios.Header
	Designation          string
	SlotType             SlotType
	SlotDataBusWidth     SlotDataBusWidth
	CurrentUsage         SlotUsage
	SlotLength           SlotLength
	SlotID               uint16
	Characteristics1     SlotCharacteristics1
	Characteristics2     SlotCharacteristics2 // SMBIOS 2.1+
	SegmentGroupNumber   uint16               // SMBIOS 2.6+
	BusNumber            uint8                // SMBIOS 2.6+
	DeviceFunctionNumber uint8                // SMBIOS 2.6+
	DataBusWidth         uint8                // SMBIOS 3.2+
	PeerGroups           []SlotPeerGroup      // SMBIOS 3.2+
	SlotInformation      uint8                // SMBIOS 3.4+
	SlotPhysicalWidth    uint8                // SMBIOS 3.4+
	SlotPitch            uint16               // SMBIOS 3.4+ (in 1/100 mm)
	SlotHeight           SlotHeight           // SMBIOS 3.5+
}

// SlotType identifies the slot type
type SlotType uint8

// Slot type values
const (
	SlotTypeOther                   SlotType = 0x01
	SlotTypeUnknown                 SlotType = 0x02
	SlotTypeISA                     SlotType = 0x03
	SlotTypeMCA                     SlotType = 0x04
	SlotTypeEISA                    SlotType = 0x05
	SlotTypePCI                     SlotType = 0x06
	SlotTypePCMCIA                  SlotType = 0x07
	SlotTypeVLVESA                  SlotType = 0x08
	SlotTypeProprietary             SlotType = 0x09
	SlotTypeProcessorCardSlot       SlotType = 0x0A
	SlotTypeProprietaryMemory       SlotType = 0x0B
	SlotTypeIORiserCard             SlotType = 0x0C
	SlotTypeNuBus                   SlotType = 0x0D
	SlotTypePCI66MHz                SlotType = 0x0E
	SlotTypeAGP                     SlotType = 0x0F
	SlotTypeAGP2X                   SlotType = 0x10
	SlotTypeAGP4X                   SlotType = 0x11
	SlotTypePCIX                    SlotType = 0x12
	SlotTypeAGP8X                   SlotType = 0x13
	SlotTypeM2Socket1DP             SlotType = 0x14
	SlotTypeM2Socket1SD             SlotType = 0x15
	SlotTypeM2Socket2               SlotType = 0x16
	SlotTypeM2Socket3               SlotType = 0x17
	SlotTypeMXMTypeI                SlotType = 0x18
	SlotTypeMXMTypeII               SlotType = 0x19
	SlotTypeMXMTypeIIIStandard      SlotType = 0x1A
	SlotTypeMXMTypeIIIHE            SlotType = 0x1B
	SlotTypeMXMTypeIV               SlotType = 0x1C
	SlotTypeMXM30TypeA              SlotType = 0x1D
	SlotTypeMXM30TypeB              SlotType = 0x1E
	SlotTypePCIExpressGen2SFF8639   SlotType = 0x1F
	SlotTypePCIExpressGen3SFF8639   SlotType = 0x20
	SlotTypePCIExpressMini52pin     SlotType = 0x21
	SlotTypePCIExpressMini52pinFull SlotType = 0x22
	SlotTypePCIExpressMini76pin     SlotType = 0x23
	SlotTypePCIExpressGen4SFF8639   SlotType = 0x24
	SlotTypePCIExpressGen5SFF8639   SlotType = 0x25
	SlotTypeOCPNIC30SmallFormFactor SlotType = 0x26
	SlotTypeOCPNIC30LargeFormFactor SlotType = 0x27
	SlotTypeOCPNICPrior30           SlotType = 0x28
	SlotTypeCXLFlexbus10            SlotType = 0x30
	SlotTypePC98C20                 SlotType = 0xA0
	SlotTypePC98C24                 SlotType = 0xA1
	SlotTypePC98E                   SlotType = 0xA2
	SlotTypePC98LocalBus            SlotType = 0xA3
	SlotTypePC98Card                SlotType = 0xA4
	SlotTypePCIExpress              SlotType = 0xA5
	SlotTypePCIExpressX1            SlotType = 0xA6
	SlotTypePCIExpressX2            SlotType = 0xA7
	SlotTypePCIExpressX4            SlotType = 0xA8
	SlotTypePCIExpressX8            SlotType = 0xA9
	SlotTypePCIExpressX16           SlotType = 0xAA
	SlotTypePCIExpressGen2          SlotType = 0xAB
	SlotTypePCIExpressGen2X1        SlotType = 0xAC
	SlotTypePCIExpressGen2X2        SlotType = 0xAD
	SlotTypePCIExpressGen2X4        SlotType = 0xAE
	SlotTypePCIExpressGen2X8        SlotType = 0xAF
	SlotTypePCIExpressGen2X16       SlotType = 0xB0
	SlotTypePCIExpressGen3          SlotType = 0xB1
	SlotTypePCIExpressGen3X1        SlotType = 0xB2
	SlotTypePCIExpressGen3X2        SlotType = 0xB3
	SlotTypePCIExpressGen3X4        SlotType = 0xB4
	SlotTypePCIExpressGen3X8        SlotType = 0xB5
	SlotTypePCIExpressGen3X16       SlotType = 0xB6
	SlotTypePCIExpressGen4          SlotType = 0xB8
	SlotTypePCIExpressGen4X1        SlotType = 0xB9
	SlotTypePCIExpressGen4X2        SlotType = 0xBA
	SlotTypePCIExpressGen4X4        SlotType = 0xBB
	SlotTypePCIExpressGen4X8        SlotType = 0xBC
	SlotTypePCIExpressGen4X16       SlotType = 0xBD
	SlotTypePCIExpressGen5          SlotType = 0xBE
	SlotTypePCIExpressGen5X1        SlotType = 0xBF
	SlotTypePCIExpressGen5X2        SlotType = 0xC0
	SlotTypePCIExpressGen5X4        SlotType = 0xC1
	SlotTypePCIExpressGen5X8        SlotType = 0xC2
	SlotTypePCIExpressGen5X16       SlotType = 0xC3
	SlotTypePCIExpressGen6          SlotType = 0xC4
	SlotTypeEDSFF_E1                SlotType = 0xC5
	SlotTypeEDSFF_E3                SlotType = 0xC6
)

// String returns a human-readable slot type description
func (st SlotType) String() string {
	types := map[SlotType]string{
		SlotTypeOther:             "Other",
		SlotTypeUnknown:           "Unknown",
		SlotTypeISA:               "ISA",
		SlotTypeMCA:               "MCA",
		SlotTypeEISA:              "EISA",
		SlotTypePCI:               "PCI",
		SlotTypePCMCIA:            "PCMCIA",
		SlotTypeVLVESA:            "VL-VESA",
		SlotTypeProprietary:       "Proprietary",
		SlotTypeAGP:               "AGP",
		SlotTypeAGP2X:             "AGP 2X",
		SlotTypeAGP4X:             "AGP 4X",
		SlotTypeAGP8X:             "AGP 8X",
		SlotTypePCIX:              "PCI-X",
		SlotTypeM2Socket1DP:       "M.2 Socket 1-DP",
		SlotTypeM2Socket1SD:       "M.2 Socket 1-SD",
		SlotTypeM2Socket2:         "M.2 Socket 2",
		SlotTypeM2Socket3:         "M.2 Socket 3",
		SlotTypePCIExpress:        "PCI Express",
		SlotTypePCIExpressX1:      "PCI Express x1",
		SlotTypePCIExpressX2:      "PCI Express x2",
		SlotTypePCIExpressX4:      "PCI Express x4",
		SlotTypePCIExpressX8:      "PCI Express x8",
		SlotTypePCIExpressX16:     "PCI Express x16",
		SlotTypePCIExpressGen2:    "PCI Express Gen 2",
		SlotTypePCIExpressGen2X1:  "PCI Express Gen 2 x1",
		SlotTypePCIExpressGen2X4:  "PCI Express Gen 2 x4",
		SlotTypePCIExpressGen2X8:  "PCI Express Gen 2 x8",
		SlotTypePCIExpressGen2X16: "PCI Express Gen 2 x16",
		SlotTypePCIExpressGen3:    "PCI Express Gen 3",
		SlotTypePCIExpressGen3X1:  "PCI Express Gen 3 x1",
		SlotTypePCIExpressGen3X4:  "PCI Express Gen 3 x4",
		SlotTypePCIExpressGen3X8:  "PCI Express Gen 3 x8",
		SlotTypePCIExpressGen3X16: "PCI Express Gen 3 x16",
		SlotTypePCIExpressGen4:    "PCI Express Gen 4",
		SlotTypePCIExpressGen4X1:  "PCI Express Gen 4 x1",
		SlotTypePCIExpressGen4X4:  "PCI Express Gen 4 x4",
		SlotTypePCIExpressGen4X8:  "PCI Express Gen 4 x8",
		SlotTypePCIExpressGen4X16: "PCI Express Gen 4 x16",
		SlotTypePCIExpressGen5:    "PCI Express Gen 5",
		SlotTypePCIExpressGen5X1:  "PCI Express Gen 5 x1",
		SlotTypePCIExpressGen5X4:  "PCI Express Gen 5 x4",
		SlotTypePCIExpressGen5X8:  "PCI Express Gen 5 x8",
		SlotTypePCIExpressGen5X16: "PCI Express Gen 5 x16",
		SlotTypePCIExpressGen6:    "PCI Express Gen 6",
		SlotTypeCXLFlexbus10:      "CXL Flexbus 1.0",
	}

	if name, ok := types[st]; ok {
		return name
	}
	return fmt.Sprintf("Unknown (0x%02X)", uint8(st))
}

// SlotDataBusWidth identifies the slot data bus width
type SlotDataBusWidth uint8

// Slot data bus width values
const (
	SlotDataBusWidthOther   SlotDataBusWidth = 0x01
	SlotDataBusWidthUnknown SlotDataBusWidth = 0x02
	SlotDataBusWidth8Bit    SlotDataBusWidth = 0x03
	SlotDataBusWidth16Bit   SlotDataBusWidth = 0x04
	SlotDataBusWidth32Bit   SlotDataBusWidth = 0x05
	SlotDataBusWidth64Bit   SlotDataBusWidth = 0x06
	SlotDataBusWidth128Bit  SlotDataBusWidth = 0x07
	SlotDataBusWidth1X      SlotDataBusWidth = 0x08
	SlotDataBusWidth2X      SlotDataBusWidth = 0x09
	SlotDataBusWidth4X      SlotDataBusWidth = 0x0A
	SlotDataBusWidth8X      SlotDataBusWidth = 0x0B
	SlotDataBusWidth12X     SlotDataBusWidth = 0x0C
	SlotDataBusWidth16X     SlotDataBusWidth = 0x0D
	SlotDataBusWidth32X     SlotDataBusWidth = 0x0E
)

// String returns a human-readable data bus width description
func (dbw SlotDataBusWidth) String() string {
	widths := map[SlotDataBusWidth]string{
		SlotDataBusWidthOther:   "Other",
		SlotDataBusWidthUnknown: "Unknown",
		SlotDataBusWidth8Bit:    "8 bit",
		SlotDataBusWidth16Bit:   "16 bit",
		SlotDataBusWidth32Bit:   "32 bit",
		SlotDataBusWidth64Bit:   "64 bit",
		SlotDataBusWidth128Bit:  "128 bit",
		SlotDataBusWidth1X:      "1x",
		SlotDataBusWidth2X:      "2x",
		SlotDataBusWidth4X:      "4x",
		SlotDataBusWidth8X:      "8x",
		SlotDataBusWidth12X:     "12x",
		SlotDataBusWidth16X:     "16x",
		SlotDataBusWidth32X:     "32x",
	}

	if name, ok := widths[dbw]; ok {
		return name
	}
	return fmt.Sprintf("Unknown (0x%02X)", uint8(dbw))
}

// SlotUsage identifies the current slot usage
type SlotUsage uint8

// Slot usage values
const (
	SlotUsageOther       SlotUsage = 0x01
	SlotUsageUnknown     SlotUsage = 0x02
	SlotUsageAvailable   SlotUsage = 0x03
	SlotUsageInUse       SlotUsage = 0x04
	SlotUsageUnavailable SlotUsage = 0x05
)

// String returns a human-readable usage description
func (su SlotUsage) String() string {
	switch su {
	case SlotUsageOther:
		return "Other"
	case SlotUsageUnknown:
		return "Unknown"
	case SlotUsageAvailable:
		return "Available"
	case SlotUsageInUse:
		return "In Use"
	case SlotUsageUnavailable:
		return "Unavailable"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(su))
	}
}

// SlotLength identifies the slot length
type SlotLength uint8

// Slot length values
const (
	SlotLengthOther   SlotLength = 0x01
	SlotLengthUnknown SlotLength = 0x02
	SlotLengthShort   SlotLength = 0x03
	SlotLengthLong    SlotLength = 0x04
	SlotLength2_5     SlotLength = 0x05
	SlotLength3_5     SlotLength = 0x06
)

// String returns a human-readable length description
func (sl SlotLength) String() string {
	switch sl {
	case SlotLengthOther:
		return "Other"
	case SlotLengthUnknown:
		return "Unknown"
	case SlotLengthShort:
		return "Short"
	case SlotLengthLong:
		return "Long"
	case SlotLength2_5:
		return "2.5\" drive form factor"
	case SlotLength3_5:
		return "3.5\" drive form factor"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(sl))
	}
}

// SlotCharacteristics1 represents slot characteristics byte 1
type SlotCharacteristics1 uint8

// Characteristics1 bit definitions
const (
	SlotChar1Unknown         SlotCharacteristics1 = 1 << 0 // Characteristics unknown
	SlotChar1Provides5V      SlotCharacteristics1 = 1 << 1 // Provides 5.0 volts
	SlotChar1Provides3_3V    SlotCharacteristics1 = 1 << 2 // Provides 3.3 volts
	SlotChar1Shared          SlotCharacteristics1 = 1 << 3 // Slot's opening is shared
	SlotChar1PCCard16        SlotCharacteristics1 = 1 << 4 // PC Card slot supports PC Card-16
	SlotChar1CardBus         SlotCharacteristics1 = 1 << 5 // PC Card slot supports CardBus
	SlotChar1ZoomVideo       SlotCharacteristics1 = 1 << 6 // PC Card slot supports Zoom Video
	SlotChar1ModemRingResume SlotCharacteristics1 = 1 << 7 // PC Card slot supports Modem Ring Resume
)

// Has checks if a characteristic is set
func (sc SlotCharacteristics1) Has(flag SlotCharacteristics1) bool {
	return sc&flag != 0
}

// SlotCharacteristics2 represents slot characteristics byte 2
type SlotCharacteristics2 uint8

// Characteristics2 bit definitions
const (
	SlotChar2PMESignal        SlotCharacteristics2 = 1 << 0 // PCI slot supports PME# signal
	SlotChar2HotPlugDevices   SlotCharacteristics2 = 1 << 1 // Slot supports hot-plug devices
	SlotChar2SMBusSignal      SlotCharacteristics2 = 1 << 2 // PCI slot supports SMBus signal
	SlotChar2Bifurcation      SlotCharacteristics2 = 1 << 3 // PCIe slot supports bifurcation
	SlotChar2SurpriseRemoval  SlotCharacteristics2 = 1 << 4 // Slot supports async/surprise removal
	SlotChar2FlexbusSlotCXL10 SlotCharacteristics2 = 1 << 5 // Flexbus slot, CXL 1.0 capable
	SlotChar2FlexbusSlotCXL20 SlotCharacteristics2 = 1 << 6 // Flexbus slot, CXL 2.0 capable
	SlotChar2FlexbusSlotCXL30 SlotCharacteristics2 = 1 << 7 // Flexbus slot, CXL 3.0 capable
)

// Has checks if a characteristic is set
func (sc SlotCharacteristics2) Has(flag SlotCharacteristics2) bool {
	return sc&flag != 0
}

// SlotPeerGroup represents a slot peer group (SMBIOS 3.2+)
type SlotPeerGroup struct {
	SegmentGroupNumber   uint16
	BusNumber            uint8
	DeviceFunctionNumber uint8
	DataBusWidth         uint8
}

// SlotHeight identifies the slot height (SMBIOS 3.5+)
type SlotHeight uint8

// Slot height values
const (
	SlotHeightNotApplicable SlotHeight = 0x00
	SlotHeightOther         SlotHeight = 0x01
	SlotHeightUnknown       SlotHeight = 0x02
	SlotHeightFullHeight    SlotHeight = 0x03
	SlotHeightLowProfile    SlotHeight = 0x04
)

// String returns a human-readable height description
func (sh SlotHeight) String() string {
	switch sh {
	case SlotHeightNotApplicable:
		return "Not Applicable"
	case SlotHeightOther:
		return "Other"
	case SlotHeightUnknown:
		return "Unknown"
	case SlotHeightFullHeight:
		return "Full Height"
	case SlotHeightLowProfile:
		return "Low Profile"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(sh))
	}
}

// Parse parses a System Slots structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*SlotInfo, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 12 bytes (SMBIOS 2.0)
	if len(s.Data) < 12 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &SlotInfo{
		Header:           s.Header,
		Designation:      s.GetString(s.GetByte(0x04)),
		SlotType:         SlotType(s.GetByte(0x05)),
		SlotDataBusWidth: SlotDataBusWidth(s.GetByte(0x06)),
		CurrentUsage:     SlotUsage(s.GetByte(0x07)),
		SlotLength:       SlotLength(s.GetByte(0x08)),
		SlotID:           s.GetWord(0x09),
		Characteristics1: SlotCharacteristics1(s.GetByte(0x0B)),
	}

	// SMBIOS 2.1+
	if len(s.Data) >= 13 {
		info.Characteristics2 = SlotCharacteristics2(s.GetByte(0x0C))
	}

	// SMBIOS 2.6+
	if len(s.Data) >= 17 {
		info.SegmentGroupNumber = s.GetWord(0x0D)
		info.BusNumber = s.GetByte(0x0F)
		info.DeviceFunctionNumber = s.GetByte(0x10)
	}

	// SMBIOS 3.2+
	if len(s.Data) >= 18 {
		info.DataBusWidth = s.GetByte(0x11)
	}

	// Peer groups (SMBIOS 3.2+)
	if len(s.Data) >= 19 {
		peerGroupCount := s.GetByte(0x12)
		if peerGroupCount > 0 && len(s.Data) >= 19+int(peerGroupCount)*5 {
			offset := 0x13
			for i := uint8(0); i < peerGroupCount; i++ {
				pg := SlotPeerGroup{
					SegmentGroupNumber:   s.GetWord(offset),
					BusNumber:            s.GetByte(offset + 2),
					DeviceFunctionNumber: s.GetByte(offset + 3),
					DataBusWidth:         s.GetByte(offset + 4),
				}
				info.PeerGroups = append(info.PeerGroups, pg)
				offset += 5
			}

			// SMBIOS 3.4+
			if len(s.Data) > offset+3 {
				info.SlotInformation = s.GetByte(offset)
				info.SlotPhysicalWidth = s.GetByte(offset + 1)
				info.SlotPitch = s.GetWord(offset + 2)
			}

			// SMBIOS 3.5+
			if len(s.Data) > offset+4 {
				info.SlotHeight = SlotHeight(s.GetByte(offset + 4))
			}
		}
	}

	return info, nil
}

// Get retrieves the first System Slot from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*SlotInfo, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all System Slots structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*SlotInfo, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var slots []*SlotInfo
	for i := range structures {
		slot, err := Parse(&structures[i])
		if err == nil {
			slots = append(slots, slot)
		}
	}

	if len(slots) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return slots, nil
}

// PCIAddress returns the PCI address as a formatted string
func (s *SlotInfo) PCIAddress() string {
	device := (s.DeviceFunctionNumber >> 3) & 0x1F
	function := s.DeviceFunctionNumber & 0x07
	return fmt.Sprintf("%04X:%02X:%02X.%X", s.SegmentGroupNumber, s.BusNumber, device, function)
}

// IsInUse returns true if the slot is currently in use
func (s *SlotInfo) IsInUse() bool {
	return s.CurrentUsage == SlotUsageInUse
}

// SupportsHotPlug returns true if the slot supports hot-plug
func (s *SlotInfo) SupportsHotPlug() bool {
	return s.Characteristics2.Has(SlotChar2HotPlugDevices)
}
