// Package type44 implements SMBIOS Type 44 - Processor Additional Information
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type44

import (
	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Processor Additional Information
const StructureType uint8 = 44

// ProcessorAdditionalInfo represents Type 44 - Processor Additional Information
type ProcessorAdditionalInfo struct {
	Header                    gosmbios.Header
	ReferencedHandle          uint16
	ProcessorSpecificBlock    ProcessorSpecificBlock
}

// ProcessorSpecificBlock contains processor-specific information
type ProcessorSpecificBlock struct {
	Length             uint8
	ProcessorType      ProcessorType
	Data               []byte
}

// ProcessorType identifies the processor type for the specific block
type ProcessorType uint8

const (
	ProcessorTypeIA32   ProcessorType = 0x01
	ProcessorTypeX64    ProcessorType = 0x02
	ProcessorTypeIA64   ProcessorType = 0x03
	ProcessorTypeARM32  ProcessorType = 0x04
	ProcessorTypeARM64  ProcessorType = 0x05
	ProcessorTypeRISCV32 ProcessorType = 0x06
	ProcessorTypeRISCV64 ProcessorType = 0x07
	ProcessorTypeRISCV128 ProcessorType = 0x08
	ProcessorTypeLoongArch32 ProcessorType = 0x09
	ProcessorTypeLoongArch64 ProcessorType = 0x0A
)

func (p ProcessorType) String() string {
	switch p {
	case ProcessorTypeIA32:
		return "IA32 (x86)"
	case ProcessorTypeX64:
		return "x64 (x86-64, AMD64)"
	case ProcessorTypeIA64:
		return "Intel Itanium"
	case ProcessorTypeARM32:
		return "32-bit ARM (Aarch32)"
	case ProcessorTypeARM64:
		return "64-bit ARM (Aarch64)"
	case ProcessorTypeRISCV32:
		return "32-bit RISC-V (RV32)"
	case ProcessorTypeRISCV64:
		return "64-bit RISC-V (RV64)"
	case ProcessorTypeRISCV128:
		return "128-bit RISC-V (RV128)"
	case ProcessorTypeLoongArch32:
		return "32-bit LoongArch"
	case ProcessorTypeLoongArch64:
		return "64-bit LoongArch"
	default:
		return "Unknown"
	}
}

// Parse parses a Processor Additional Information structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*ProcessorAdditionalInfo, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 6 bytes
	if len(s.Data) < 6 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &ProcessorAdditionalInfo{
		Header:           s.Header,
		ReferencedHandle: s.GetWord(0x04),
	}

	// Parse processor specific block
	offset := 0x06
	if offset < len(s.Data) {
		blockLen := s.GetByte(offset)
		if blockLen >= 2 && offset+int(blockLen) <= len(s.Data) {
			info.ProcessorSpecificBlock = ProcessorSpecificBlock{
				Length:        blockLen,
				ProcessorType: ProcessorType(s.GetByte(offset + 1)),
			}

			// Copy the remaining data
			if blockLen > 2 {
				info.ProcessorSpecificBlock.Data = make([]byte, blockLen-2)
				copy(info.ProcessorSpecificBlock.Data, s.Data[offset+2:offset+int(blockLen)])
			}
		}
	}

	return info, nil
}

// Get retrieves the first Processor Additional Information from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*ProcessorAdditionalInfo, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Processor Additional Information structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*ProcessorAdditionalInfo, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var infos []*ProcessorAdditionalInfo
	for i := range structures {
		info, err := Parse(&structures[i])
		if err == nil {
			infos = append(infos, info)
		}
	}

	if len(infos) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return infos, nil
}
