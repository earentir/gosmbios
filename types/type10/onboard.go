// Package type10 implements SMBIOS Type 10 - On Board Devices Information (Obsolete)
// Per DSP0134 SMBIOS Reference Specification 3.9.0
// This structure is obsolete starting with SMBIOS 2.6 specification
// Use Type 41 (Onboard Devices Extended Information) instead
package type10

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for On Board Devices Information
const StructureType uint8 = 10

// OnBoardDevices represents Type 10 - On Board Devices Information (Obsolete)
type OnBoardDevices struct {
	Header  gosmbios.Header
	Devices []OnBoardDevice
}

// OnBoardDevice represents a single on-board device
type OnBoardDevice struct {
	DeviceType  DeviceType
	Description string
	Enabled     bool
}

// DeviceType identifies the type of on-board device
type DeviceType uint8

const (
	DeviceTypeOther          DeviceType = 0x01
	DeviceTypeUnknown        DeviceType = 0x02
	DeviceTypeVideo          DeviceType = 0x03
	DeviceTypeSCSIController DeviceType = 0x04
	DeviceTypeEthernet       DeviceType = 0x05
	DeviceTypeTokenRing      DeviceType = 0x06
	DeviceTypeSound          DeviceType = 0x07
	DeviceTypePATAController DeviceType = 0x08
	DeviceTypeSATAController DeviceType = 0x09
	DeviceTypeSASController  DeviceType = 0x0A
	DeviceTypeWirelessLAN    DeviceType = 0x0B
	DeviceTypeBluetooth      DeviceType = 0x0C
	DeviceTypeWWAN           DeviceType = 0x0D
	DeviceTypeeMMC           DeviceType = 0x0E
	DeviceTypeNVMe           DeviceType = 0x0F
	DeviceTypeUFS            DeviceType = 0x10
)

func (d DeviceType) String() string {
	switch d {
	case DeviceTypeOther:
		return "Other"
	case DeviceTypeUnknown:
		return "Unknown"
	case DeviceTypeVideo:
		return "Video"
	case DeviceTypeSCSIController:
		return "SCSI Controller"
	case DeviceTypeEthernet:
		return "Ethernet"
	case DeviceTypeTokenRing:
		return "Token Ring"
	case DeviceTypeSound:
		return "Sound"
	case DeviceTypePATAController:
		return "PATA Controller"
	case DeviceTypeSATAController:
		return "SATA Controller"
	case DeviceTypeSASController:
		return "SAS Controller"
	case DeviceTypeWirelessLAN:
		return "Wireless LAN"
	case DeviceTypeBluetooth:
		return "Bluetooth"
	case DeviceTypeWWAN:
		return "WWAN"
	case DeviceTypeeMMC:
		return "eMMC"
	case DeviceTypeNVMe:
		return "NVMe Controller"
	case DeviceTypeUFS:
		return "UFS Controller"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(d))
	}
}

// Parse parses an On Board Devices Information structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*OnBoardDevices, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 4 bytes (header only)
	if len(s.Data) < 4 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &OnBoardDevices{
		Header: s.Header,
	}

	// Number of devices = (structure length - 4) / 2
	numDevices := (int(s.Header.Length) - 4) / 2

	for i := 0; i < numDevices; i++ {
		offset := 0x04 + (i * 2)
		if offset+1 >= len(s.Data) {
			break
		}

		typeByte := s.GetByte(offset)
		enabled := typeByte&0x80 != 0
		deviceType := DeviceType(typeByte & 0x7F)
		descIndex := s.GetByte(offset + 1)

		device := OnBoardDevice{
			DeviceType:  deviceType,
			Description: s.GetString(descIndex),
			Enabled:     enabled,
		}
		info.Devices = append(info.Devices, device)
	}

	return info, nil
}

// Get retrieves the On Board Devices Information from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*OnBoardDevices, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all On Board Devices structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*OnBoardDevices, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var devices []*OnBoardDevices
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

// GetAllDevices returns all devices from all Type 10 structures as a flat list
func GetAllDevices(sm *gosmbios.SMBIOS) ([]OnBoardDevice, error) {
	all, err := GetAll(sm)
	if err != nil {
		return nil, err
	}

	var devices []OnBoardDevice
	for _, info := range all {
		devices = append(devices, info.Devices...)
	}
	return devices, nil
}
