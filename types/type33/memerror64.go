// Package type33 implements SMBIOS Type 33 - 64-Bit Memory Error Information
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type33

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for 64-Bit Memory Error Information
const StructureType uint8 = 33

// MemoryError64 represents Type 33 - 64-Bit Memory Error Information
type MemoryError64 struct {
	Header            gosmbios.Header
	ErrorType         ErrorType
	ErrorGranularity  ErrorGranularity
	ErrorOperation    ErrorOperation
	VendorSyndrome    uint32
	MemoryArrayErrorAddress   uint64
	DeviceErrorAddress        uint64
	ErrorResolution   uint32
}

// ErrorType identifies the type of memory error
type ErrorType uint8

const (
	ErrorTypeOther              ErrorType = 0x01
	ErrorTypeUnknown            ErrorType = 0x02
	ErrorTypeOK                 ErrorType = 0x03
	ErrorTypeBadRead            ErrorType = 0x04
	ErrorTypeParity             ErrorType = 0x05
	ErrorTypeSingleBit          ErrorType = 0x06
	ErrorTypeDoubleBit          ErrorType = 0x07
	ErrorTypeMultiBit           ErrorType = 0x08
	ErrorTypeNibble             ErrorType = 0x09
	ErrorTypeChecksum           ErrorType = 0x0A
	ErrorTypeCRC                ErrorType = 0x0B
	ErrorTypeCorrectedSingleBit ErrorType = 0x0C
	ErrorTypeCorrected          ErrorType = 0x0D
	ErrorTypeUncorrectable      ErrorType = 0x0E
)

func (e ErrorType) String() string {
	switch e {
	case ErrorTypeOther:
		return "Other"
	case ErrorTypeUnknown:
		return "Unknown"
	case ErrorTypeOK:
		return "OK"
	case ErrorTypeBadRead:
		return "Bad Read"
	case ErrorTypeParity:
		return "Parity Error"
	case ErrorTypeSingleBit:
		return "Single-Bit Error"
	case ErrorTypeDoubleBit:
		return "Double-Bit Error"
	case ErrorTypeMultiBit:
		return "Multi-Bit Error"
	case ErrorTypeNibble:
		return "Nibble Error"
	case ErrorTypeChecksum:
		return "Checksum Error"
	case ErrorTypeCRC:
		return "CRC Error"
	case ErrorTypeCorrectedSingleBit:
		return "Corrected Single-Bit Error"
	case ErrorTypeCorrected:
		return "Corrected Error"
	case ErrorTypeUncorrectable:
		return "Uncorrectable Error"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(e))
	}
}

// ErrorGranularity identifies the granularity of the error
type ErrorGranularity uint8

const (
	GranularityOther       ErrorGranularity = 0x01
	GranularityUnknown     ErrorGranularity = 0x02
	GranularityDeviceLevel ErrorGranularity = 0x03
	GranularityPartition   ErrorGranularity = 0x04
)

func (e ErrorGranularity) String() string {
	switch e {
	case GranularityOther:
		return "Other"
	case GranularityUnknown:
		return "Unknown"
	case GranularityDeviceLevel:
		return "Device Level"
	case GranularityPartition:
		return "Memory Partition Level"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(e))
	}
}

// ErrorOperation identifies the operation that caused the error
type ErrorOperation uint8

const (
	OperationOther   ErrorOperation = 0x01
	OperationUnknown ErrorOperation = 0x02
	OperationRead    ErrorOperation = 0x03
	OperationWrite   ErrorOperation = 0x04
	OperationPartial ErrorOperation = 0x05
)

func (e ErrorOperation) String() string {
	switch e {
	case OperationOther:
		return "Other"
	case OperationUnknown:
		return "Unknown"
	case OperationRead:
		return "Read"
	case OperationWrite:
		return "Write"
	case OperationPartial:
		return "Partial Write"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(e))
	}
}

// Parse parses a 64-Bit Memory Error Information structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*MemoryError64, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 31 bytes
	if len(s.Data) < 31 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &MemoryError64{
		Header:            s.Header,
		ErrorType:         ErrorType(s.GetByte(0x04)),
		ErrorGranularity:  ErrorGranularity(s.GetByte(0x05)),
		ErrorOperation:    ErrorOperation(s.GetByte(0x06)),
		VendorSyndrome:    s.GetDWord(0x07),
		MemoryArrayErrorAddress:   s.GetQWord(0x0B),
		DeviceErrorAddress:        s.GetQWord(0x13),
		ErrorResolution:   s.GetDWord(0x1B),
	}

	return info, nil
}

// Get retrieves the first 64-Bit Memory Error Information from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*MemoryError64, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all 64-Bit Memory Error Information structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*MemoryError64, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var errors []*MemoryError64
	for i := range structures {
		err, parseErr := Parse(&structures[i])
		if parseErr == nil {
			errors = append(errors, err)
		}
	}

	if len(errors) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return errors, nil
}

// IsAddressUnknown returns true if the memory array error address is unknown
func (m *MemoryError64) IsAddressUnknown() bool {
	return m.MemoryArrayErrorAddress == 0x8000000000000000
}

// IsDeviceAddressUnknown returns true if the device error address is unknown
func (m *MemoryError64) IsDeviceAddressUnknown() bool {
	return m.DeviceErrorAddress == 0x8000000000000000
}

// IsResolutionUnknown returns true if the error resolution is unknown
func (m *MemoryError64) IsResolutionUnknown() bool {
	return m.ErrorResolution == 0x80000000
}
