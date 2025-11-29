// Package type41 implements SMBIOS Type 41 - Onboard Devices Extended Information
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type41

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Onboard Devices Extended Information
const StructureType uint8 = 41

// OnboardDeviceExtended represents Type 41 - Onboard Devices Extended Information
type OnboardDeviceExtended struct {
	Header             gosmbios.Header
	ReferenceDesignation string
	DeviceType         DeviceType
	DeviceTypeInstance uint8
	SegmentGroupNumber uint16
	BusNumber          uint8
	DeviceFunctionNumber uint8
}

// DeviceType identifies the type of onboard device
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
	// Remove enabled bit for type lookup
	t := d & 0x7F
	switch t {
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
		return fmt.Sprintf("Unknown (0x%02X)", uint8(t))
	}
}

// IsEnabled returns true if the device is enabled
func (d DeviceType) IsEnabled() bool {
	return d&0x80 != 0
}

// Type returns the device type (without enabled bit)
func (d DeviceType) Type() DeviceType {
	return d & 0x7F
}

// Parse parses an Onboard Devices Extended Information structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*OnboardDeviceExtended, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 11 bytes
	if len(s.Data) < 11 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &OnboardDeviceExtended{
		Header:               s.Header,
		ReferenceDesignation: s.GetString(s.GetByte(0x04)),
		DeviceType:           DeviceType(s.GetByte(0x05)),
		DeviceTypeInstance:   s.GetByte(0x06),
		SegmentGroupNumber:   s.GetWord(0x07),
		BusNumber:            s.GetByte(0x09),
		DeviceFunctionNumber: s.GetByte(0x0A),
	}

	return info, nil
}

// Get retrieves the first Onboard Devices Extended from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*OnboardDeviceExtended, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Onboard Devices Extended structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*OnboardDeviceExtended, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var devices []*OnboardDeviceExtended
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

// IsEnabled returns true if the device is enabled
func (o *OnboardDeviceExtended) IsEnabled() bool {
	return o.DeviceType.IsEnabled()
}

// TypeString returns the device type as a string
func (o *OnboardDeviceExtended) TypeString() string {
	return o.DeviceType.Type().String()
}

// PCIAddress returns the PCI address as a formatted string
func (o *OnboardDeviceExtended) PCIAddress() string {
	device := (o.DeviceFunctionNumber >> 3) & 0x1F
	function := o.DeviceFunctionNumber & 0x07
	return fmt.Sprintf("%04X:%02X:%02X.%X", o.SegmentGroupNumber, o.BusNumber, device, function)
}

// StatusString returns the enabled/disabled status as a string
func (o *OnboardDeviceExtended) StatusString() string {
	if o.IsEnabled() {
		return "Enabled"
	}
	return "Disabled"
}
