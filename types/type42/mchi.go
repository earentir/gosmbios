// Package type42 implements SMBIOS Type 42 - Management Controller Host Interface
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type42

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Management Controller Host Interface
const StructureType uint8 = 42

// ManagementControllerHostInterface represents Type 42 - Management Controller Host Interface
type ManagementControllerHostInterface struct {
	Header                      gosmbios.Header
	InterfaceType               InterfaceType
	InterfaceTypeSpecificData   []byte
	ProtocolRecords             []ProtocolRecord
}

// InterfaceType identifies the management controller interface type
type InterfaceType uint8

const (
	InterfaceTypeKCS               InterfaceType = 0x02
	InterfaceType8250UART          InterfaceType = 0x03
	InterfaceType16450UART         InterfaceType = 0x04
	InterfaceType16550UART         InterfaceType = 0x05
	InterfaceType16650UART         InterfaceType = 0x06
	InterfaceType16750UART         InterfaceType = 0x07
	InterfaceType16850UART         InterfaceType = 0x08
	InterfaceTypeNetworkHostIF     InterfaceType = 0x40
	InterfaceTypeOEMDefined        InterfaceType = 0xF0
)

func (i InterfaceType) String() string {
	switch i {
	case InterfaceTypeKCS:
		return "KCS (Keyboard Controller Style)"
	case InterfaceType8250UART:
		return "8250 UART Register Compatible"
	case InterfaceType16450UART:
		return "16450 UART Register Compatible"
	case InterfaceType16550UART:
		return "16550/16550A UART Register Compatible"
	case InterfaceType16650UART:
		return "16650/16650A UART Register Compatible"
	case InterfaceType16750UART:
		return "16750/16750A UART Register Compatible"
	case InterfaceType16850UART:
		return "16850/16850A UART Register Compatible"
	case InterfaceTypeNetworkHostIF:
		return "Network Host Interface"
	default:
		if i >= 0xF0 {
			return fmt.Sprintf("OEM Defined (0x%02X)", uint8(i))
		}
		return fmt.Sprintf("Unknown (0x%02X)", uint8(i))
	}
}

// ProtocolRecord represents a protocol record
type ProtocolRecord struct {
	ProtocolType         ProtocolType
	ProtocolTypeSpecific []byte
}

// ProtocolType identifies the protocol type
type ProtocolType uint8

const (
	ProtocolTypeIPMI           ProtocolType = 0x02
	ProtocolTypeMCTP           ProtocolType = 0x03
	ProtocolTypeRedfishOverIP  ProtocolType = 0x04
	ProtocolTypeOEMDefined     ProtocolType = 0xF0
)

func (p ProtocolType) String() string {
	switch p {
	case ProtocolTypeIPMI:
		return "IPMI"
	case ProtocolTypeMCTP:
		return "MCTP"
	case ProtocolTypeRedfishOverIP:
		return "Redfish over IP"
	default:
		if p >= 0xF0 {
			return fmt.Sprintf("OEM Defined (0x%02X)", uint8(p))
		}
		return fmt.Sprintf("Unknown (0x%02X)", uint8(p))
	}
}

// Parse parses a Management Controller Host Interface structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*ManagementControllerHostInterface, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 9 bytes (SMBIOS 3.0)
	if len(s.Data) < 9 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &ManagementControllerHostInterface{
		Header:        s.Header,
		InterfaceType: InterfaceType(s.GetByte(0x04)),
	}

	// Interface type specific data length (SMBIOS 3.2+)
	if len(s.Data) >= 6 {
		ifDataLen := s.GetByte(0x05)
		if ifDataLen > 0 && len(s.Data) >= 6+int(ifDataLen) {
			info.InterfaceTypeSpecificData = make([]byte, ifDataLen)
			copy(info.InterfaceTypeSpecificData, s.Data[0x06:0x06+ifDataLen])
		}

		// Parse protocol records
		offset := 0x06 + int(ifDataLen)
		if offset < len(s.Data) {
			protocolCount := s.GetByte(offset)
			offset++

			for i := uint8(0); i < protocolCount && offset < len(s.Data); i++ {
				if offset >= len(s.Data) {
					break
				}
				protocolType := ProtocolType(s.GetByte(offset))
				offset++

				if offset >= len(s.Data) {
					break
				}
				protocolDataLen := s.GetByte(offset)
				offset++

				var protocolData []byte
				if protocolDataLen > 0 && offset+int(protocolDataLen) <= len(s.Data) {
					protocolData = make([]byte, protocolDataLen)
					copy(protocolData, s.Data[offset:offset+int(protocolDataLen)])
					offset += int(protocolDataLen)
				}

				record := ProtocolRecord{
					ProtocolType:         protocolType,
					ProtocolTypeSpecific: protocolData,
				}
				info.ProtocolRecords = append(info.ProtocolRecords, record)
			}
		}
	}

	return info, nil
}

// Get retrieves the first Management Controller Host Interface from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*ManagementControllerHostInterface, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Management Controller Host Interface structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*ManagementControllerHostInterface, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var interfaces []*ManagementControllerHostInterface
	for i := range structures {
		iface, err := Parse(&structures[i])
		if err == nil {
			interfaces = append(interfaces, iface)
		}
	}

	if len(interfaces) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return interfaces, nil
}
