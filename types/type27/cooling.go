// Package type27 implements SMBIOS Type 27 - Cooling Device
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type27

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Cooling Device
const StructureType uint8 = 27

// CoolingDevice represents Type 27 - Cooling Device
type CoolingDevice struct {
	Header                   gosmbios.Header
	TemperatureProbeHandle   uint16
	DeviceTypeAndStatus      DeviceTypeAndStatus
	CoolingUnitGroup         uint8
	OEMDefined               uint32
	NominalSpeed             uint16 // In rpm (SMBIOS 2.2+)
	Description              string // SMBIOS 2.7+
}

// DeviceTypeAndStatus represents the device type and status byte
type DeviceTypeAndStatus uint8

// DeviceType returns the cooling device type
func (d DeviceTypeAndStatus) DeviceType() CoolingType {
	return CoolingType(d & 0x1F)
}

// Status returns the device status
func (d DeviceTypeAndStatus) Status() DeviceStatus {
	return DeviceStatus((d >> 5) & 0x07)
}

func (d DeviceTypeAndStatus) String() string {
	return fmt.Sprintf("Type: %s, Status: %s", d.DeviceType().String(), d.Status().String())
}

// CoolingType identifies the cooling device type
type CoolingType uint8

const (
	CoolingTypeOther              CoolingType = 0x01
	CoolingTypeUnknown            CoolingType = 0x02
	CoolingTypeFan                CoolingType = 0x03
	CoolingTypeCentrifugalBlower  CoolingType = 0x04
	CoolingTypeChipFan            CoolingType = 0x05
	CoolingTypeCabinetFan         CoolingType = 0x06
	CoolingTypePowerSupplyFan     CoolingType = 0x07
	CoolingTypeHeatPipe           CoolingType = 0x08
	CoolingTypeIntegratedRefrigeration CoolingType = 0x09
	CoolingTypeActiveCooling      CoolingType = 0x10
	CoolingTypePassiveCooling     CoolingType = 0x11
)

func (c CoolingType) String() string {
	switch c {
	case CoolingTypeOther:
		return "Other"
	case CoolingTypeUnknown:
		return "Unknown"
	case CoolingTypeFan:
		return "Fan"
	case CoolingTypeCentrifugalBlower:
		return "Centrifugal Blower"
	case CoolingTypeChipFan:
		return "Chip Fan"
	case CoolingTypeCabinetFan:
		return "Cabinet Fan"
	case CoolingTypePowerSupplyFan:
		return "Power Supply Fan"
	case CoolingTypeHeatPipe:
		return "Heat Pipe"
	case CoolingTypeIntegratedRefrigeration:
		return "Integrated Refrigeration"
	case CoolingTypeActiveCooling:
		return "Active Cooling"
	case CoolingTypePassiveCooling:
		return "Passive Cooling"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(c))
	}
}

// DeviceStatus identifies the device status
type DeviceStatus uint8

const (
	DeviceStatusOther          DeviceStatus = 0x01
	DeviceStatusUnknown        DeviceStatus = 0x02
	DeviceStatusOK             DeviceStatus = 0x03
	DeviceStatusNonCritical    DeviceStatus = 0x04
	DeviceStatusCritical       DeviceStatus = 0x05
	DeviceStatusNonRecoverable DeviceStatus = 0x06
)

func (d DeviceStatus) String() string {
	switch d {
	case DeviceStatusOther:
		return "Other"
	case DeviceStatusUnknown:
		return "Unknown"
	case DeviceStatusOK:
		return "OK"
	case DeviceStatusNonCritical:
		return "Non-critical"
	case DeviceStatusCritical:
		return "Critical"
	case DeviceStatusNonRecoverable:
		return "Non-recoverable"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(d))
	}
}

// Parse parses a Cooling Device structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*CoolingDevice, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 12 bytes
	if len(s.Data) < 12 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &CoolingDevice{
		Header:                 s.Header,
		TemperatureProbeHandle: s.GetWord(0x04),
		DeviceTypeAndStatus:    DeviceTypeAndStatus(s.GetByte(0x06)),
		CoolingUnitGroup:       s.GetByte(0x07),
		OEMDefined:             s.GetDWord(0x08),
	}

	// SMBIOS 2.2+
	if len(s.Data) >= 14 {
		info.NominalSpeed = s.GetWord(0x0C)
	}

	// SMBIOS 2.7+
	if len(s.Data) >= 15 {
		info.Description = s.GetString(s.GetByte(0x0E))
	}

	return info, nil
}

// Get retrieves the first Cooling Device from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*CoolingDevice, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Cooling Device structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*CoolingDevice, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var devices []*CoolingDevice
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

// NominalSpeedString returns the nominal speed as a string
func (c *CoolingDevice) NominalSpeedString() string {
	if c.NominalSpeed == 0x8000 {
		return "Unknown"
	}
	if c.NominalSpeed == 0 {
		return "Non-rotating"
	}
	return fmt.Sprintf("%d rpm", c.NominalSpeed)
}

// CoolingUnitGroupString returns the cooling unit group as a string
func (c *CoolingDevice) CoolingUnitGroupString() string {
	if c.CoolingUnitGroup == 0 {
		return "None"
	}
	return fmt.Sprintf("Group %d", c.CoolingUnitGroup)
}
