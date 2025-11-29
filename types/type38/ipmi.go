// Package type38 implements SMBIOS Type 38 - IPMI Device Information
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type38

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for IPMI Device Information
const StructureType uint8 = 38

// IPMIDeviceInfo represents Type 38 - IPMI Device Information
type IPMIDeviceInfo struct {
	Header                    gosmbios.Header
	InterfaceType             InterfaceType
	IPMISpecificationRevision uint8
	I2CSlaveAddress           uint8
	NVStorageDeviceAddress    uint8
	BaseAddress               uint64
	BaseAddressModifier       BaseAddressModifier // SMBIOS 2.4+
	InterruptNumber           uint8               // SMBIOS 2.4+
}

// InterfaceType identifies the IPMI interface type
type InterfaceType uint8

const (
	InterfaceTypeUnknown       InterfaceType = 0x00
	InterfaceTypeKCS           InterfaceType = 0x01
	InterfaceTypeSMIC          InterfaceType = 0x02
	InterfaceTypeBT            InterfaceType = 0x03
	InterfaceTypeSSIF          InterfaceType = 0x04
)

func (i InterfaceType) String() string {
	switch i {
	case InterfaceTypeUnknown:
		return "Unknown"
	case InterfaceTypeKCS:
		return "KCS (Keyboard Controller Style)"
	case InterfaceTypeSMIC:
		return "SMIC (Server Management Interface Chip)"
	case InterfaceTypeBT:
		return "BT (Block Transfer)"
	case InterfaceTypeSSIF:
		return "SSIF (SMBus System Interface)"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(i))
	}
}

// BaseAddressModifier represents the base address modifier/interrupt info
type BaseAddressModifier uint8

// RegisterSpacing returns the register spacing
func (b BaseAddressModifier) RegisterSpacing() RegisterSpacing {
	return RegisterSpacing((b >> 6) & 0x03)
}

// IsLSBit returns true if the LSB of base address is zero
func (b BaseAddressModifier) IsLSBit() bool {
	return b&0x10 != 0
}

// IsIOSpace returns true if the address is in I/O space (vs memory space)
func (b BaseAddressModifier) IsIOSpace() bool {
	return b&0x01 != 0
}

// InterruptPolarity returns true if active high
func (b BaseAddressModifier) InterruptPolarity() bool {
	return b&0x02 != 0
}

// InterruptTriggerMode returns true if level triggered
func (b BaseAddressModifier) InterruptTriggerMode() bool {
	return b&0x04 != 0
}

// InterruptEnabled returns true if interrupt generation is enabled
func (b BaseAddressModifier) InterruptEnabled() bool {
	return b&0x08 != 0
}

func (b BaseAddressModifier) String() string {
	space := "Memory"
	if b.IsIOSpace() {
		space = "I/O"
	}
	return fmt.Sprintf("Address Space: %s, Register Spacing: %s",
		space, b.RegisterSpacing().String())
}

// RegisterSpacing identifies register spacing
type RegisterSpacing uint8

const (
	RegisterSpacingByte      RegisterSpacing = 0x00
	RegisterSpacing4Byte     RegisterSpacing = 0x01
	RegisterSpacing16Byte    RegisterSpacing = 0x02
	RegisterSpacingReserved  RegisterSpacing = 0x03
)

func (r RegisterSpacing) String() string {
	switch r {
	case RegisterSpacingByte:
		return "Successive Byte Boundaries"
	case RegisterSpacing4Byte:
		return "32-bit Boundaries"
	case RegisterSpacing16Byte:
		return "16-byte Boundaries"
	default:
		return "Reserved"
	}
}

// Parse parses an IPMI Device Information structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*IPMIDeviceInfo, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 16 bytes
	if len(s.Data) < 16 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &IPMIDeviceInfo{
		Header:                    s.Header,
		InterfaceType:             InterfaceType(s.GetByte(0x04)),
		IPMISpecificationRevision: s.GetByte(0x05),
		I2CSlaveAddress:           s.GetByte(0x06),
		NVStorageDeviceAddress:    s.GetByte(0x07),
		BaseAddress:               s.GetQWord(0x08),
	}

	// SMBIOS 2.4+
	if len(s.Data) >= 18 {
		info.BaseAddressModifier = BaseAddressModifier(s.GetByte(0x10))
		info.InterruptNumber = s.GetByte(0x11)
	}

	return info, nil
}

// Get retrieves the IPMI Device Information from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*IPMIDeviceInfo, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// SpecificationRevisionString returns the IPMI spec revision as a string
func (i *IPMIDeviceInfo) SpecificationRevisionString() string {
	major := (i.IPMISpecificationRevision >> 4) & 0x0F
	minor := i.IPMISpecificationRevision & 0x0F
	return fmt.Sprintf("%d.%d", major, minor)
}

// I2CAddressString returns the I2C slave address as a string
func (i *IPMIDeviceInfo) I2CAddressString() string {
	// Address is in bits 7:1, bit 0 is unused
	return fmt.Sprintf("0x%02X", i.I2CSlaveAddress>>1)
}

// BaseAddressString returns the base address as a string
func (i *IPMIDeviceInfo) BaseAddressString() string {
	if i.BaseAddressModifier.IsIOSpace() {
		return fmt.Sprintf("I/O 0x%04X", i.BaseAddress)
	}
	return fmt.Sprintf("Memory 0x%016X", i.BaseAddress)
}

// InterruptNumberString returns the interrupt number as a string
func (i *IPMIDeviceInfo) InterruptNumberString() string {
	if i.InterruptNumber == 0 {
		return "Not Specified"
	}
	return fmt.Sprintf("IRQ %d", i.InterruptNumber)
}
