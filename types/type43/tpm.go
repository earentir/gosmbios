// Package type43 implements SMBIOS Type 43 - TPM Device
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type43

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for TPM Device
const StructureType uint8 = 43

// TPMDevice represents Type 43 - TPM Device
type TPMDevice struct {
	Header              gosmbios.Header
	VendorID            [4]byte
	MajorSpecVersion    uint8
	MinorSpecVersion    uint8
	FirmwareVersion1    uint32
	FirmwareVersion2    uint32
	Description         string
	Characteristics     Characteristics
	OEMDefined          uint32
}

// Characteristics represents TPM device characteristics
type Characteristics uint64

const (
	CharTPMDeviceNotSupported         Characteristics = 1 << 2
	CharTPMDeviceFamilyConfigurable   Characteristics = 1 << 3
	CharTPMDeviceFamilyIsTPM2_0       Characteristics = 1 << 4
	CharTPMDeviceFamilyIsTPM1_2       Characteristics = 1 << 5
)

// IsSupported returns true if TPM device is supported
func (c Characteristics) IsSupported() bool {
	return c&CharTPMDeviceNotSupported == 0
}

// IsFamilyConfigurable returns true if TPM family is configurable by platform software
func (c Characteristics) IsFamilyConfigurable() bool {
	return c&CharTPMDeviceFamilyConfigurable != 0
}

// IsTPM2_0 returns true if TPM family is TPM 2.0
func (c Characteristics) IsTPM2_0() bool {
	return c&CharTPMDeviceFamilyIsTPM2_0 != 0
}

// IsTPM1_2 returns true if TPM family is TPM 1.2
func (c Characteristics) IsTPM1_2() bool {
	return c&CharTPMDeviceFamilyIsTPM1_2 != 0
}

func (c Characteristics) String() string {
	if c&CharTPMDeviceNotSupported != 0 {
		return "TPM Device Not Supported"
	}

	family := "Unknown"
	if c.IsTPM2_0() {
		family = "TPM 2.0"
	} else if c.IsTPM1_2() {
		family = "TPM 1.2"
	}

	configurable := "No"
	if c.IsFamilyConfigurable() {
		configurable = "Yes"
	}

	return fmt.Sprintf("Family: %s, Configurable: %s", family, configurable)
}

// Parse parses a TPM Device structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*TPMDevice, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 31 bytes
	if len(s.Data) < 31 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &TPMDevice{
		Header:           s.Header,
		MajorSpecVersion: s.GetByte(0x08),
		MinorSpecVersion: s.GetByte(0x09),
		FirmwareVersion1: s.GetDWord(0x0A),
		FirmwareVersion2: s.GetDWord(0x0E),
		Description:      s.GetString(s.GetByte(0x12)),
		Characteristics:  Characteristics(s.GetQWord(0x13)),
		OEMDefined:       s.GetDWord(0x1B),
	}

	// Copy vendor ID
	copy(info.VendorID[:], s.Data[0x04:0x08])

	return info, nil
}

// Get retrieves the TPM Device from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*TPMDevice, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// VendorIDString returns the vendor ID as a string
func (t *TPMDevice) VendorIDString() string {
	// Vendor ID is typically 4 ASCII characters
	// Check if all printable
	printable := true
	for _, b := range t.VendorID {
		if b < 32 || b > 126 {
			printable = false
			break
		}
	}

	if printable {
		return string(t.VendorID[:])
	}
	return fmt.Sprintf("%02X%02X%02X%02X", t.VendorID[0], t.VendorID[1], t.VendorID[2], t.VendorID[3])
}

// SpecVersionString returns the specification version as a string
func (t *TPMDevice) SpecVersionString() string {
	return fmt.Sprintf("%d.%d", t.MajorSpecVersion, t.MinorSpecVersion)
}

// FirmwareVersionString returns the firmware version as a string
func (t *TPMDevice) FirmwareVersionString() string {
	if t.FirmwareVersion1 == 0 && t.FirmwareVersion2 == 0 {
		return "Not Reported"
	}
	// Format depends on TPM family
	if t.Characteristics.IsTPM2_0() {
		return fmt.Sprintf("%d.%d", t.FirmwareVersion1, t.FirmwareVersion2)
	}
	// TPM 1.2 uses BCD format
	return fmt.Sprintf("%08X.%08X", t.FirmwareVersion1, t.FirmwareVersion2)
}

// IsSupported returns true if the TPM is supported
func (t *TPMDevice) IsSupported() bool {
	return t.Characteristics.IsSupported()
}

// Family returns the TPM family as a string
func (t *TPMDevice) Family() string {
	if t.Characteristics.IsTPM2_0() {
		return "TPM 2.0"
	}
	if t.Characteristics.IsTPM1_2() {
		return "TPM 1.2"
	}
	return "Unknown"
}
