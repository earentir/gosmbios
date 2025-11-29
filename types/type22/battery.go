// Package type22 implements SMBIOS Type 22 - Portable Battery
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type22

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Portable Battery
const StructureType uint8 = 22

// PortableBattery represents Type 22 - Portable Battery
type PortableBattery struct {
	Header                       gosmbios.Header
	Location                     string
	Manufacturer                 string
	ManufactureDate              string
	SerialNumber                 string
	DeviceName                   string
	DeviceChemistry              DeviceChemistry
	DesignCapacity               uint16 // In mWh
	DesignVoltage                uint16 // In mV
	SBDSVersionNumber            string
	MaximumErrorInBatteryData    uint8  // As percentage
	SBDSSerialNumber             uint16 // SMBIOS 2.2+
	SBDSManufactureDate          uint16 // SMBIOS 2.2+
	SBDSDeviceChemistry          string // SMBIOS 2.2+
	DesignCapacityMultiplier     uint8  // SMBIOS 2.2+
	OEMSpecific                  uint32 // SMBIOS 2.2+
}

// DeviceChemistry identifies the battery chemistry
type DeviceChemistry uint8

const (
	ChemistryOther             DeviceChemistry = 0x01
	ChemistryUnknown           DeviceChemistry = 0x02
	ChemistryLeadAcid          DeviceChemistry = 0x03
	ChemistryNickelCadmium     DeviceChemistry = 0x04
	ChemistryNickelMetalHydride DeviceChemistry = 0x05
	ChemistryLithiumIon        DeviceChemistry = 0x06
	ChemistryZincAir           DeviceChemistry = 0x07
	ChemistryLithiumPolymer    DeviceChemistry = 0x08
)

func (d DeviceChemistry) String() string {
	switch d {
	case ChemistryOther:
		return "Other"
	case ChemistryUnknown:
		return "Unknown"
	case ChemistryLeadAcid:
		return "Lead Acid"
	case ChemistryNickelCadmium:
		return "Nickel Cadmium"
	case ChemistryNickelMetalHydride:
		return "Nickel Metal Hydride"
	case ChemistryLithiumIon:
		return "Lithium-ion"
	case ChemistryZincAir:
		return "Zinc Air"
	case ChemistryLithiumPolymer:
		return "Lithium Polymer"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(d))
	}
}

// Parse parses a Portable Battery structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*PortableBattery, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 16 bytes (SMBIOS 2.1)
	if len(s.Data) < 16 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &PortableBattery{
		Header:                    s.Header,
		Location:                  s.GetString(s.GetByte(0x04)),
		Manufacturer:              s.GetString(s.GetByte(0x05)),
		ManufactureDate:           s.GetString(s.GetByte(0x06)),
		SerialNumber:              s.GetString(s.GetByte(0x07)),
		DeviceName:                s.GetString(s.GetByte(0x08)),
		DeviceChemistry:           DeviceChemistry(s.GetByte(0x09)),
		DesignCapacity:            s.GetWord(0x0A),
		DesignVoltage:             s.GetWord(0x0C),
		SBDSVersionNumber:         s.GetString(s.GetByte(0x0E)),
		MaximumErrorInBatteryData: s.GetByte(0x0F),
	}

	// SMBIOS 2.2+
	if len(s.Data) >= 26 {
		info.SBDSSerialNumber = s.GetWord(0x10)
		info.SBDSManufactureDate = s.GetWord(0x12)
		info.SBDSDeviceChemistry = s.GetString(s.GetByte(0x14))
		info.DesignCapacityMultiplier = s.GetByte(0x15)
		info.OEMSpecific = s.GetDWord(0x16)
	}

	return info, nil
}

// Get retrieves the first Portable Battery from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*PortableBattery, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Portable Battery structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*PortableBattery, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var batteries []*PortableBattery
	for i := range structures {
		bat, err := Parse(&structures[i])
		if err == nil {
			batteries = append(batteries, bat)
		}
	}

	if len(batteries) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return batteries, nil
}

// DesignCapacityString returns the design capacity as a string
func (p *PortableBattery) DesignCapacityString() string {
	if p.DesignCapacity == 0 {
		return "Unknown"
	}
	capacity := uint32(p.DesignCapacity)
	if p.DesignCapacityMultiplier > 0 {
		capacity *= uint32(p.DesignCapacityMultiplier)
	}
	if capacity >= 1000 {
		return fmt.Sprintf("%.2f Wh", float64(capacity)/1000.0)
	}
	return fmt.Sprintf("%d mWh", capacity)
}

// DesignVoltageString returns the design voltage as a string
func (p *PortableBattery) DesignVoltageString() string {
	if p.DesignVoltage == 0 {
		return "Unknown"
	}
	return fmt.Sprintf("%.2f V", float64(p.DesignVoltage)/1000.0)
}

// SBDSManufactureDateString returns the SBDS manufacture date as a formatted string
func (p *PortableBattery) SBDSManufactureDateString() string {
	if p.SBDSManufactureDate == 0 {
		return "Unknown"
	}
	// Date format: bits 0-4 = day (1-31), bits 5-8 = month (1-12), bits 9-15 = year (0=1980)
	day := p.SBDSManufactureDate & 0x1F
	month := (p.SBDSManufactureDate >> 5) & 0x0F
	year := ((p.SBDSManufactureDate >> 9) & 0x7F) + 1980
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

// MaximumErrorString returns the maximum error as a string
func (p *PortableBattery) MaximumErrorString() string {
	if p.MaximumErrorInBatteryData == 0xFF {
		return "Unknown"
	}
	return fmt.Sprintf("%d%%", p.MaximumErrorInBatteryData)
}
