// Package type3 implements SMBIOS Type 3 - System Enclosure or Chassis
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type3

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for System Enclosure
const StructureType uint8 = 3

// ChassisInfo represents Type 3 - System Enclosure or Chassis
type ChassisInfo struct {
	Header             gosmbios.Header
	Manufacturer       string
	Type               ChassisType
	TypeLocked         bool // Chassis lock present
	Version            string
	SerialNumber       string
	AssetTag           string
	BootUpState        ChassisState
	PowerSupplyState   ChassisState
	ThermalState       ChassisState
	SecurityStatus     SecurityStatus
	OEMDefined         uint32             // SMBIOS 2.3+
	Height             uint8              // In U units (0 = unspecified)
	NumberOfPowerCords uint8              // 0 = unspecified
	ContainedElements  []ContainedElement // SMBIOS 2.3+
	SKUNumber          string             // SMBIOS 2.7+
}

// ChassisType identifies the enclosure type
type ChassisType uint8

// Chassis type values (bits 0-6, bit 7 indicates lock)
const (
	ChassisTypeOther               ChassisType = 0x01
	ChassisTypeUnknown             ChassisType = 0x02
	ChassisTypeDesktop             ChassisType = 0x03
	ChassisTypeLowProfileDesktop   ChassisType = 0x04
	ChassisTypePizzaBox            ChassisType = 0x05
	ChassisTypeMiniTower           ChassisType = 0x06
	ChassisTypeTower               ChassisType = 0x07
	ChassisTypePortable            ChassisType = 0x08
	ChassisTypeLaptop              ChassisType = 0x09
	ChassisTypeNotebook            ChassisType = 0x0A
	ChassisTypeHandHeld            ChassisType = 0x0B
	ChassisTypeDockingStation      ChassisType = 0x0C
	ChassisTypeAllInOne            ChassisType = 0x0D
	ChassisTypeSubNotebook         ChassisType = 0x0E
	ChassisTypeSpaceSaving         ChassisType = 0x0F
	ChassisTypeLunchBox            ChassisType = 0x10
	ChassisTypeMainServerChassis   ChassisType = 0x11
	ChassisTypeExpansionChassis    ChassisType = 0x12
	ChassisTypeSubChassis          ChassisType = 0x13
	ChassisTypeBusExpansionChassis ChassisType = 0x14
	ChassisTypePeripheralChassis   ChassisType = 0x15
	ChassisTypeRAIDChassis         ChassisType = 0x16
	ChassisTypeRackMountChassis    ChassisType = 0x17
	ChassisTypeSealedCasePC        ChassisType = 0x18
	ChassisTypeMultiSystemChassis  ChassisType = 0x19
	ChassisTypeCompactPCI          ChassisType = 0x1A
	ChassisTypeAdvancedTCA         ChassisType = 0x1B
	ChassisTypeBlade               ChassisType = 0x1C
	ChassisTypeBladeEnclosure      ChassisType = 0x1D
	ChassisTypeTablet              ChassisType = 0x1E
	ChassisTypeConvertible         ChassisType = 0x1F
	ChassisTypeDetachable          ChassisType = 0x20
	ChassisTypeIoTGateway          ChassisType = 0x21
	ChassisTypeEmbeddedPC          ChassisType = 0x22
	ChassisTypeMiniPC              ChassisType = 0x23
	ChassisTypeStickPC             ChassisType = 0x24
)

// String returns a human-readable chassis type description
func (ct ChassisType) String() string {
	switch ct {
	case ChassisTypeOther:
		return "Other"
	case ChassisTypeUnknown:
		return "Unknown"
	case ChassisTypeDesktop:
		return "Desktop"
	case ChassisTypeLowProfileDesktop:
		return "Low Profile Desktop"
	case ChassisTypePizzaBox:
		return "Pizza Box"
	case ChassisTypeMiniTower:
		return "Mini Tower"
	case ChassisTypeTower:
		return "Tower"
	case ChassisTypePortable:
		return "Portable"
	case ChassisTypeLaptop:
		return "Laptop"
	case ChassisTypeNotebook:
		return "Notebook"
	case ChassisTypeHandHeld:
		return "Hand Held"
	case ChassisTypeDockingStation:
		return "Docking Station"
	case ChassisTypeAllInOne:
		return "All in One"
	case ChassisTypeSubNotebook:
		return "Sub Notebook"
	case ChassisTypeSpaceSaving:
		return "Space-saving"
	case ChassisTypeLunchBox:
		return "Lunch Box"
	case ChassisTypeMainServerChassis:
		return "Main Server Chassis"
	case ChassisTypeExpansionChassis:
		return "Expansion Chassis"
	case ChassisTypeSubChassis:
		return "SubChassis"
	case ChassisTypeBusExpansionChassis:
		return "Bus Expansion Chassis"
	case ChassisTypePeripheralChassis:
		return "Peripheral Chassis"
	case ChassisTypeRAIDChassis:
		return "RAID Chassis"
	case ChassisTypeRackMountChassis:
		return "Rack Mount Chassis"
	case ChassisTypeSealedCasePC:
		return "Sealed-case PC"
	case ChassisTypeMultiSystemChassis:
		return "Multi-system Chassis"
	case ChassisTypeCompactPCI:
		return "Compact PCI"
	case ChassisTypeAdvancedTCA:
		return "Advanced TCA"
	case ChassisTypeBlade:
		return "Blade"
	case ChassisTypeBladeEnclosure:
		return "Blade Enclosure"
	case ChassisTypeTablet:
		return "Tablet"
	case ChassisTypeConvertible:
		return "Convertible"
	case ChassisTypeDetachable:
		return "Detachable"
	case ChassisTypeIoTGateway:
		return "IoT Gateway"
	case ChassisTypeEmbeddedPC:
		return "Embedded PC"
	case ChassisTypeMiniPC:
		return "Mini PC"
	case ChassisTypeStickPC:
		return "Stick PC"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(ct))
	}
}

// IsPortable returns true if the chassis type indicates a portable device
func (ct ChassisType) IsPortable() bool {
	switch ct {
	case ChassisTypePortable, ChassisTypeLaptop, ChassisTypeNotebook,
		ChassisTypeHandHeld, ChassisTypeSubNotebook, ChassisTypeTablet,
		ChassisTypeConvertible, ChassisTypeDetachable:
		return true
	}
	return false
}

// ChassisState represents the state of the chassis
type ChassisState uint8

// Chassis state values
const (
	ChassisStateOther          ChassisState = 0x01
	ChassisStateUnknown        ChassisState = 0x02
	ChassisStateSafe           ChassisState = 0x03
	ChassisStateWarning        ChassisState = 0x04
	ChassisStateCritical       ChassisState = 0x05
	ChassisStateNonRecoverable ChassisState = 0x06
)

// String returns a human-readable state description
func (cs ChassisState) String() string {
	switch cs {
	case ChassisStateOther:
		return "Other"
	case ChassisStateUnknown:
		return "Unknown"
	case ChassisStateSafe:
		return "Safe"
	case ChassisStateWarning:
		return "Warning"
	case ChassisStateCritical:
		return "Critical"
	case ChassisStateNonRecoverable:
		return "Non-recoverable"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(cs))
	}
}

// SecurityStatus represents the chassis security status
type SecurityStatus uint8

// Security status values
const (
	SecurityOther                      SecurityStatus = 0x01
	SecurityUnknown                    SecurityStatus = 0x02
	SecurityNone                       SecurityStatus = 0x03
	SecurityExternalInterfaceLockedOut SecurityStatus = 0x04
	SecurityExternalInterfaceEnabled   SecurityStatus = 0x05
)

// String returns a human-readable security status description
func (ss SecurityStatus) String() string {
	switch ss {
	case SecurityOther:
		return "Other"
	case SecurityUnknown:
		return "Unknown"
	case SecurityNone:
		return "None"
	case SecurityExternalInterfaceLockedOut:
		return "External Interface Locked Out"
	case SecurityExternalInterfaceEnabled:
		return "External Interface Enabled"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(ss))
	}
}

// ContainedElement represents an element contained in the chassis
type ContainedElement struct {
	Type    uint8
	Minimum uint8
	Maximum uint8
}

// Parse parses a System Enclosure structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*ChassisInfo, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 9 bytes (SMBIOS 2.0)
	if len(s.Data) < 9 {
		return nil, gosmbios.ErrInvalidStructure
	}

	typeField := s.GetByte(0x05)

	info := &ChassisInfo{
		Header:       s.Header,
		Manufacturer: s.GetString(s.GetByte(0x04)),
		Type:         ChassisType(typeField & 0x7F),
		TypeLocked:   (typeField & 0x80) != 0,
		Version:      s.GetString(s.GetByte(0x06)),
		SerialNumber: s.GetString(s.GetByte(0x07)),
		AssetTag:     s.GetString(s.GetByte(0x08)),
	}

	// SMBIOS 2.1+
	if len(s.Data) >= 13 {
		info.BootUpState = ChassisState(s.GetByte(0x09))
		info.PowerSupplyState = ChassisState(s.GetByte(0x0A))
		info.ThermalState = ChassisState(s.GetByte(0x0B))
		info.SecurityStatus = SecurityStatus(s.GetByte(0x0C))
	}

	// SMBIOS 2.3+
	if len(s.Data) >= 17 {
		info.OEMDefined = s.GetDWord(0x0D)
	}

	if len(s.Data) >= 19 {
		info.Height = s.GetByte(0x11)
		info.NumberOfPowerCords = s.GetByte(0x12)
	}

	// Contained elements (SMBIOS 2.3+)
	if len(s.Data) >= 21 {
		containedCount := s.GetByte(0x13)
		elementRecordLen := s.GetByte(0x14)

		if elementRecordLen >= 3 && containedCount > 0 {
			offset := 0x15
			for i := uint8(0); i < containedCount; i++ {
				if offset+int(elementRecordLen) <= len(s.Data) {
					elem := ContainedElement{
						Type:    s.GetByte(offset),
						Minimum: s.GetByte(offset + 1),
						Maximum: s.GetByte(offset + 2),
					}
					info.ContainedElements = append(info.ContainedElements, elem)
					offset += int(elementRecordLen)
				}
			}

			// SKU Number (SMBIOS 2.7+) - follows contained elements
			if offset < len(s.Data) {
				info.SKUNumber = s.GetString(s.GetByte(offset))
			}
		}
	}

	return info, nil
}

// Get retrieves the System Enclosure Information from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*ChassisInfo, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all System Enclosure structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*ChassisInfo, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var chassis []*ChassisInfo
	for i := range structures {
		ch, err := Parse(&structures[i])
		if err == nil {
			chassis = append(chassis, ch)
		}
	}

	if len(chassis) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return chassis, nil
}

// HeightString returns the chassis height in U units, or "Unspecified"
func (c *ChassisInfo) HeightString() string {
	if c.Height == 0 {
		return "Unspecified"
	}
	return fmt.Sprintf("%dU", c.Height)
}
