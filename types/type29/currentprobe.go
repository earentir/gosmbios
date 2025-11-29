// Package type29 implements SMBIOS Type 29 - Electrical Current Probe
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type29

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Electrical Current Probe
const StructureType uint8 = 29

// CurrentProbe represents Type 29 - Electrical Current Probe
type CurrentProbe struct {
	Header            gosmbios.Header
	Description       string
	LocationAndStatus LocationAndStatus
	MaximumValue      uint16 // In milliamps
	MinimumValue      uint16 // In milliamps
	Resolution        uint16 // In tenths of milliamps
	Tolerance         uint16 // In +/- milliamps
	Accuracy          uint16 // In +/- 1/100 percent
	OEMDefined        uint32
	NominalValue      uint16 // In milliamps (SMBIOS 2.2+)
}

// LocationAndStatus represents the location and status byte
type LocationAndStatus uint8

// Location returns the probe location
func (l LocationAndStatus) Location() ProbeLocation {
	return ProbeLocation(l & 0x1F)
}

// Status returns the probe status
func (l LocationAndStatus) Status() ProbeStatus {
	return ProbeStatus((l >> 5) & 0x07)
}

func (l LocationAndStatus) String() string {
	return fmt.Sprintf("Location: %s, Status: %s", l.Location().String(), l.Status().String())
}

// ProbeLocation identifies the probe location
type ProbeLocation uint8

const (
	LocationOther         ProbeLocation = 0x01
	LocationUnknown       ProbeLocation = 0x02
	LocationProcessor     ProbeLocation = 0x03
	LocationDisk          ProbeLocation = 0x04
	LocationPeripheralBay ProbeLocation = 0x05
	LocationSMM           ProbeLocation = 0x06
	LocationMotherboard   ProbeLocation = 0x07
	LocationMemoryModule  ProbeLocation = 0x08
	LocationProcessorModule ProbeLocation = 0x09
	LocationPowerUnit     ProbeLocation = 0x0A
	LocationAddInCard     ProbeLocation = 0x0B
)

func (p ProbeLocation) String() string {
	switch p {
	case LocationOther:
		return "Other"
	case LocationUnknown:
		return "Unknown"
	case LocationProcessor:
		return "Processor"
	case LocationDisk:
		return "Disk"
	case LocationPeripheralBay:
		return "Peripheral Bay"
	case LocationSMM:
		return "System Management Module"
	case LocationMotherboard:
		return "Motherboard"
	case LocationMemoryModule:
		return "Memory Module"
	case LocationProcessorModule:
		return "Processor Module"
	case LocationPowerUnit:
		return "Power Unit"
	case LocationAddInCard:
		return "Add-in Card"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(p))
	}
}

// ProbeStatus identifies the probe status
type ProbeStatus uint8

const (
	StatusOther         ProbeStatus = 0x01
	StatusUnknown       ProbeStatus = 0x02
	StatusOK            ProbeStatus = 0x03
	StatusNonCritical   ProbeStatus = 0x04
	StatusCritical      ProbeStatus = 0x05
	StatusNonRecoverable ProbeStatus = 0x06
)

func (p ProbeStatus) String() string {
	switch p {
	case StatusOther:
		return "Other"
	case StatusUnknown:
		return "Unknown"
	case StatusOK:
		return "OK"
	case StatusNonCritical:
		return "Non-critical"
	case StatusCritical:
		return "Critical"
	case StatusNonRecoverable:
		return "Non-recoverable"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(p))
	}
}

// Parse parses an Electrical Current Probe structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*CurrentProbe, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 20 bytes
	if len(s.Data) < 20 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &CurrentProbe{
		Header:            s.Header,
		Description:       s.GetString(s.GetByte(0x04)),
		LocationAndStatus: LocationAndStatus(s.GetByte(0x05)),
		MaximumValue:      s.GetWord(0x06),
		MinimumValue:      s.GetWord(0x08),
		Resolution:        s.GetWord(0x0A),
		Tolerance:         s.GetWord(0x0C),
		Accuracy:          s.GetWord(0x0E),
		OEMDefined:        s.GetDWord(0x10),
	}

	// SMBIOS 2.2+
	if len(s.Data) >= 22 {
		info.NominalValue = s.GetWord(0x14)
	}

	return info, nil
}

// Get retrieves the first Electrical Current Probe from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*CurrentProbe, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Electrical Current Probe structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*CurrentProbe, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var probes []*CurrentProbe
	for i := range structures {
		probe, err := Parse(&structures[i])
		if err == nil {
			probes = append(probes, probe)
		}
	}

	if len(probes) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return probes, nil
}

// MaximumValueString returns the maximum value as a string
func (c *CurrentProbe) MaximumValueString() string {
	if c.MaximumValue == 0x8000 {
		return "Unknown"
	}
	return fmt.Sprintf("%.3f A", float64(c.MaximumValue)/1000.0)
}

// MinimumValueString returns the minimum value as a string
func (c *CurrentProbe) MinimumValueString() string {
	if c.MinimumValue == 0x8000 {
		return "Unknown"
	}
	return fmt.Sprintf("%.3f A", float64(c.MinimumValue)/1000.0)
}

// NominalValueString returns the nominal value as a string
func (c *CurrentProbe) NominalValueString() string {
	if c.NominalValue == 0x8000 {
		return "Unknown"
	}
	return fmt.Sprintf("%.3f A", float64(c.NominalValue)/1000.0)
}

// ResolutionString returns the resolution as a string
func (c *CurrentProbe) ResolutionString() string {
	if c.Resolution == 0x8000 {
		return "Unknown"
	}
	return fmt.Sprintf("%.4f A", float64(c.Resolution)/10000.0)
}

// AccuracyString returns the accuracy as a string
func (c *CurrentProbe) AccuracyString() string {
	if c.Accuracy == 0x8000 {
		return "Unknown"
	}
	return fmt.Sprintf("+/- %.2f%%", float64(c.Accuracy)/100.0)
}
