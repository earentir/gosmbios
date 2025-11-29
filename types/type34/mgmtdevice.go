// Package type34 implements SMBIOS Type 34 - Management Device
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type34

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Management Device
const StructureType uint8 = 34

// ManagementDevice represents Type 34 - Management Device
type ManagementDevice struct {
	Header      gosmbios.Header
	Description string
	DeviceType  DeviceType
	Address     uint32
	AddressType AddressType
}

// DeviceType identifies the type of management device
type DeviceType uint8

const (
	DeviceTypeOther          DeviceType = 0x01
	DeviceTypeUnknown        DeviceType = 0x02
	DeviceTypeLM75           DeviceType = 0x03
	DeviceTypeLM78           DeviceType = 0x04
	DeviceTypeLM79           DeviceType = 0x05
	DeviceTypeLM80           DeviceType = 0x06
	DeviceTypeLM81           DeviceType = 0x07
	DeviceTypeADM9240        DeviceType = 0x08
	DeviceTypeDS1780         DeviceType = 0x09
	DeviceTypeMaxim1617      DeviceType = 0x0A
	DeviceTypeGL518SM        DeviceType = 0x0B
	DeviceTypeW83781D        DeviceType = 0x0C
	DeviceTypeHT82H791       DeviceType = 0x0D
)

func (d DeviceType) String() string {
	switch d {
	case DeviceTypeOther:
		return "Other"
	case DeviceTypeUnknown:
		return "Unknown"
	case DeviceTypeLM75:
		return "National Semiconductor LM75"
	case DeviceTypeLM78:
		return "National Semiconductor LM78"
	case DeviceTypeLM79:
		return "National Semiconductor LM79"
	case DeviceTypeLM80:
		return "National Semiconductor LM80"
	case DeviceTypeLM81:
		return "National Semiconductor LM81"
	case DeviceTypeADM9240:
		return "Analog Devices ADM9240"
	case DeviceTypeDS1780:
		return "Dallas Semiconductor DS1780"
	case DeviceTypeMaxim1617:
		return "Maxim 1617"
	case DeviceTypeGL518SM:
		return "Genesys GL518SM"
	case DeviceTypeW83781D:
		return "Winbond W83781D"
	case DeviceTypeHT82H791:
		return "Holtek HT82H791"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(d))
	}
}

// AddressType identifies the address type
type AddressType uint8

const (
	AddressTypeOther   AddressType = 0x01
	AddressTypeUnknown AddressType = 0x02
	AddressTypeIO      AddressType = 0x03
	AddressTypeMemory  AddressType = 0x04
	AddressTypeSMBus   AddressType = 0x05
)

func (a AddressType) String() string {
	switch a {
	case AddressTypeOther:
		return "Other"
	case AddressTypeUnknown:
		return "Unknown"
	case AddressTypeIO:
		return "I/O Port"
	case AddressTypeMemory:
		return "Memory"
	case AddressTypeSMBus:
		return "SMBus"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(a))
	}
}

// Parse parses a Management Device structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*ManagementDevice, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 11 bytes
	if len(s.Data) < 11 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &ManagementDevice{
		Header:      s.Header,
		Description: s.GetString(s.GetByte(0x04)),
		DeviceType:  DeviceType(s.GetByte(0x05)),
		Address:     s.GetDWord(0x06),
		AddressType: AddressType(s.GetByte(0x0A)),
	}

	return info, nil
}

// Get retrieves the first Management Device from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*ManagementDevice, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Management Device structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*ManagementDevice, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var devices []*ManagementDevice
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
