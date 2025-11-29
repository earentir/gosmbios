// Package type21 implements SMBIOS Type 21 - Built-in Pointing Device
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type21

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Built-in Pointing Device
const StructureType uint8 = 21

// PointingDevice represents Type 21 - Built-in Pointing Device
type PointingDevice struct {
	Header          gosmbios.Header
	DeviceType      DeviceType
	Interface       Interface
	NumberOfButtons uint8
}

// DeviceType identifies the type of pointing device
type DeviceType uint8

const (
	DeviceTypeOther           DeviceType = 0x01
	DeviceTypeUnknown         DeviceType = 0x02
	DeviceTypeMouse           DeviceType = 0x03
	DeviceTypeTrackBall       DeviceType = 0x04
	DeviceTypeTrackPoint      DeviceType = 0x05
	DeviceTypeGlidePoint      DeviceType = 0x06
	DeviceTypeTouchPad        DeviceType = 0x07
	DeviceTypeTouchScreen     DeviceType = 0x08
	DeviceTypeOpticalSensor   DeviceType = 0x09
)

func (d DeviceType) String() string {
	switch d {
	case DeviceTypeOther:
		return "Other"
	case DeviceTypeUnknown:
		return "Unknown"
	case DeviceTypeMouse:
		return "Mouse"
	case DeviceTypeTrackBall:
		return "Track Ball"
	case DeviceTypeTrackPoint:
		return "Track Point"
	case DeviceTypeGlidePoint:
		return "Glide Point"
	case DeviceTypeTouchPad:
		return "Touch Pad"
	case DeviceTypeTouchScreen:
		return "Touch Screen"
	case DeviceTypeOpticalSensor:
		return "Optical Sensor"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(d))
	}
}

// Interface identifies the interface type
type Interface uint8

const (
	InterfaceOther          Interface = 0x01
	InterfaceUnknown        Interface = 0x02
	InterfaceSerial         Interface = 0x03
	InterfacePS2            Interface = 0x04
	InterfaceInfrared       Interface = 0x05
	InterfaceHPHIL          Interface = 0x06
	InterfaceBusMouse       Interface = 0x07
	InterfaceADB            Interface = 0x08
	InterfaceBusMouseDB9    Interface = 0xA0
	InterfaceBusMouseMicroDIN Interface = 0xA1
	InterfaceUSB            Interface = 0xA2
	InterfaceI2C            Interface = 0xA3
	InterfaceSPI            Interface = 0xA4
)

func (i Interface) String() string {
	switch i {
	case InterfaceOther:
		return "Other"
	case InterfaceUnknown:
		return "Unknown"
	case InterfaceSerial:
		return "Serial"
	case InterfacePS2:
		return "PS/2"
	case InterfaceInfrared:
		return "Infrared"
	case InterfaceHPHIL:
		return "HP-HIL"
	case InterfaceBusMouse:
		return "Bus mouse"
	case InterfaceADB:
		return "ADB (Apple Desktop Bus)"
	case InterfaceBusMouseDB9:
		return "Bus mouse DB-9"
	case InterfaceBusMouseMicroDIN:
		return "Bus mouse Micro-DIN"
	case InterfaceUSB:
		return "USB"
	case InterfaceI2C:
		return "I2C"
	case InterfaceSPI:
		return "SPI"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(i))
	}
}

// Parse parses a Built-in Pointing Device structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*PointingDevice, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 7 bytes
	if len(s.Data) < 7 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &PointingDevice{
		Header:          s.Header,
		DeviceType:      DeviceType(s.GetByte(0x04)),
		Interface:       Interface(s.GetByte(0x05)),
		NumberOfButtons: s.GetByte(0x06),
	}

	return info, nil
}

// Get retrieves the first Built-in Pointing Device from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*PointingDevice, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Built-in Pointing Device structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*PointingDevice, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var devices []*PointingDevice
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
