// Package type39 implements SMBIOS Type 39 - System Power Supply
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type39

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for System Power Supply
const StructureType uint8 = 39

// SystemPowerSupply represents Type 39 - System Power Supply
type SystemPowerSupply struct {
	Header                        gosmbios.Header
	PowerUnitGroup                uint8
	Location                      string
	DeviceName                    string
	Manufacturer                  string
	SerialNumber                  string
	AssetTagNumber                string
	ModelPartNumber               string
	RevisionLevel                 string
	MaxPowerCapacity              uint16 // In watts
	Characteristics               Characteristics
	InputVoltageProbeHandle       uint16
	CoolingDeviceHandle           uint16
	InputCurrentProbeHandle       uint16
}

// Characteristics represents power supply characteristics
type Characteristics uint16

// IsHotReplaceable returns true if the power supply is hot-replaceable
func (c Characteristics) IsHotReplaceable() bool {
	return c&0x0001 != 0
}

// IsPresent returns true if the power supply is present
func (c Characteristics) IsPresent() bool {
	return c&0x0002 != 0
}

// IsUnplugged returns true if the power supply is unplugged from the wall
func (c Characteristics) IsUnplugged() bool {
	return c&0x0004 != 0
}

// InputVoltageRange returns the input voltage range
func (c Characteristics) InputVoltageRange() InputVoltageRange {
	return InputVoltageRange((c >> 3) & 0x07)
}

// Status returns the power supply status
func (c Characteristics) Status() PSUStatus {
	return PSUStatus((c >> 7) & 0x07)
}

// Type returns the power supply type
func (c Characteristics) Type() PSUType {
	return PSUType((c >> 10) & 0x0F)
}

func (c Characteristics) String() string {
	present := "Not Present"
	if c.IsPresent() {
		present = "Present"
	}
	plugged := "Plugged"
	if c.IsUnplugged() {
		plugged = "Unplugged"
	}
	hotReplace := "No"
	if c.IsHotReplaceable() {
		hotReplace = "Yes"
	}
	return fmt.Sprintf("Status: %s, %s, %s, Hot-Replaceable: %s, Type: %s",
		c.Status().String(), present, plugged, hotReplace, c.Type().String())
}

// InputVoltageRange identifies the input voltage range
type InputVoltageRange uint8

const (
	VoltageRangeOther       InputVoltageRange = 0x01
	VoltageRangeUnknown     InputVoltageRange = 0x02
	VoltageRangeManual      InputVoltageRange = 0x03
	VoltageRangeAutoSwitch  InputVoltageRange = 0x04
	VoltageRangeWideRange   InputVoltageRange = 0x05
	VoltageRangeNotApplicable InputVoltageRange = 0x06
)

func (v InputVoltageRange) String() string {
	switch v {
	case VoltageRangeOther:
		return "Other"
	case VoltageRangeUnknown:
		return "Unknown"
	case VoltageRangeManual:
		return "Manual"
	case VoltageRangeAutoSwitch:
		return "Auto-Switch"
	case VoltageRangeWideRange:
		return "Wide Range"
	case VoltageRangeNotApplicable:
		return "Not Applicable"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(v))
	}
}

// PSUStatus identifies the power supply status
type PSUStatus uint8

const (
	PSUStatusOther          PSUStatus = 0x01
	PSUStatusUnknown        PSUStatus = 0x02
	PSUStatusOK             PSUStatus = 0x03
	PSUStatusNonCritical    PSUStatus = 0x04
	PSUStatusCritical       PSUStatus = 0x05
	PSUStatusNonRecoverable PSUStatus = 0x06
)

func (p PSUStatus) String() string {
	switch p {
	case PSUStatusOther:
		return "Other"
	case PSUStatusUnknown:
		return "Unknown"
	case PSUStatusOK:
		return "OK"
	case PSUStatusNonCritical:
		return "Non-critical"
	case PSUStatusCritical:
		return "Critical"
	case PSUStatusNonRecoverable:
		return "Non-recoverable"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(p))
	}
}

// PSUType identifies the power supply type
type PSUType uint8

const (
	PSUTypeOther         PSUType = 0x01
	PSUTypeUnknown       PSUType = 0x02
	PSUTypeLinear        PSUType = 0x03
	PSUTypeSwitching     PSUType = 0x04
	PSUTypeBattery       PSUType = 0x05
	PSUTypeUPS           PSUType = 0x06
	PSUTypeConverter     PSUType = 0x07
	PSUTypeRegulator     PSUType = 0x08
)

func (p PSUType) String() string {
	switch p {
	case PSUTypeOther:
		return "Other"
	case PSUTypeUnknown:
		return "Unknown"
	case PSUTypeLinear:
		return "Linear"
	case PSUTypeSwitching:
		return "Switching"
	case PSUTypeBattery:
		return "Battery"
	case PSUTypeUPS:
		return "UPS"
	case PSUTypeConverter:
		return "Converter"
	case PSUTypeRegulator:
		return "Regulator"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(p))
	}
}

// Parse parses a System Power Supply structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*SystemPowerSupply, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 16 bytes
	if len(s.Data) < 16 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &SystemPowerSupply{
		Header:                  s.Header,
		PowerUnitGroup:          s.GetByte(0x04),
		Location:                s.GetString(s.GetByte(0x05)),
		DeviceName:              s.GetString(s.GetByte(0x06)),
		Manufacturer:            s.GetString(s.GetByte(0x07)),
		SerialNumber:            s.GetString(s.GetByte(0x08)),
		AssetTagNumber:          s.GetString(s.GetByte(0x09)),
		ModelPartNumber:         s.GetString(s.GetByte(0x0A)),
		RevisionLevel:           s.GetString(s.GetByte(0x0B)),
		MaxPowerCapacity:        s.GetWord(0x0C),
		Characteristics:         Characteristics(s.GetWord(0x0E)),
	}

	// Optional handles (SMBIOS 2.3.1+)
	if len(s.Data) >= 18 {
		info.InputVoltageProbeHandle = s.GetWord(0x10)
	}
	if len(s.Data) >= 20 {
		info.CoolingDeviceHandle = s.GetWord(0x12)
	}
	if len(s.Data) >= 22 {
		info.InputCurrentProbeHandle = s.GetWord(0x14)
	}

	return info, nil
}

// Get retrieves the first System Power Supply from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*SystemPowerSupply, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all System Power Supply structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*SystemPowerSupply, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var supplies []*SystemPowerSupply
	for i := range structures {
		psu, err := Parse(&structures[i])
		if err == nil {
			supplies = append(supplies, psu)
		}
	}

	if len(supplies) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return supplies, nil
}

// MaxPowerCapacityString returns the max power capacity as a string
func (s *SystemPowerSupply) MaxPowerCapacityString() string {
	if s.MaxPowerCapacity == 0x8000 {
		return "Unknown"
	}
	return fmt.Sprintf("%d W", s.MaxPowerCapacity)
}

// PowerUnitGroupString returns the power unit group as a string
func (s *SystemPowerSupply) PowerUnitGroupString() string {
	if s.PowerUnitGroup == 0 {
		return "None"
	}
	return fmt.Sprintf("Group %d", s.PowerUnitGroup)
}

// HasInputVoltageProbe returns true if an input voltage probe is associated
func (s *SystemPowerSupply) HasInputVoltageProbe() bool {
	return s.InputVoltageProbeHandle != 0xFFFF
}

// HasCoolingDevice returns true if a cooling device is associated
func (s *SystemPowerSupply) HasCoolingDevice() bool {
	return s.CoolingDeviceHandle != 0xFFFF
}

// HasInputCurrentProbe returns true if an input current probe is associated
func (s *SystemPowerSupply) HasInputCurrentProbe() bool {
	return s.InputCurrentProbeHandle != 0xFFFF
}
