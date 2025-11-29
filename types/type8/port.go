// Package type8 implements SMBIOS Type 8 - Port Connector Information
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type8

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Port Connector Information
const StructureType uint8 = 8

// PortConnector represents Type 8 - Port Connector Information
type PortConnector struct {
	Header                      gosmbios.Header
	InternalReferenceDesignator string
	InternalConnectorType       ConnectorType
	ExternalReferenceDesignator string
	ExternalConnectorType       ConnectorType
	PortType                    PortType
}

// ConnectorType identifies the physical connector type
type ConnectorType uint8

// Connector type values
const (
	ConnectorNone                  ConnectorType = 0x00
	ConnectorCentronics            ConnectorType = 0x01
	ConnectorMiniCentronics        ConnectorType = 0x02
	ConnectorProprietary           ConnectorType = 0x03
	ConnectorDB25PinMale           ConnectorType = 0x04
	ConnectorDB25PinFemale         ConnectorType = 0x05
	ConnectorDB15PinMale           ConnectorType = 0x06
	ConnectorDB15PinFemale         ConnectorType = 0x07
	ConnectorDB9PinMale            ConnectorType = 0x08
	ConnectorDB9PinFemale          ConnectorType = 0x09
	ConnectorRJ11                  ConnectorType = 0x0A
	ConnectorRJ45                  ConnectorType = 0x0B
	Connector50PinMiniSCSI         ConnectorType = 0x0C
	ConnectorMiniDIN               ConnectorType = 0x0D
	ConnectorMicroDIN              ConnectorType = 0x0E
	ConnectorPS2                   ConnectorType = 0x0F
	ConnectorInfrared              ConnectorType = 0x10
	ConnectorHPHIL                 ConnectorType = 0x11
	ConnectorAccessBusUSB          ConnectorType = 0x12
	ConnectorSSASCSI               ConnectorType = 0x13
	ConnectorCircularDIN8Male      ConnectorType = 0x14
	ConnectorCircularDIN8Female    ConnectorType = 0x15
	ConnectorOnBoardIDE            ConnectorType = 0x16
	ConnectorOnBoardFloppy         ConnectorType = 0x17
	Connector9PinDualInline        ConnectorType = 0x18
	Connector25PinDualInline       ConnectorType = 0x19
	Connector50PinDualInline       ConnectorType = 0x1A
	Connector68PinDualInline       ConnectorType = 0x1B
	ConnectorOnBoardSoundFromCD    ConnectorType = 0x1C
	ConnectorMiniCentronicsType14  ConnectorType = 0x1D
	ConnectorMiniCentronicsType26  ConnectorType = 0x1E
	ConnectorMiniJack              ConnectorType = 0x1F
	ConnectorBNC                   ConnectorType = 0x20
	Connector1394                  ConnectorType = 0x21
	ConnectorSASSATAPlugReceptacle ConnectorType = 0x22
	ConnectorUSBTypeC              ConnectorType = 0x23
	ConnectorPC98                  ConnectorType = 0xA0
	ConnectorPC98Hireso            ConnectorType = 0xA1
	ConnectorPCH98                 ConnectorType = 0xA2
	ConnectorPC98Note              ConnectorType = 0xA3
	ConnectorPC98Full              ConnectorType = 0xA4
	ConnectorOther                 ConnectorType = 0xFF
)

// String returns a human-readable connector type description
func (ct ConnectorType) String() string {
	types := map[ConnectorType]string{
		ConnectorNone:                  "None",
		ConnectorCentronics:            "Centronics",
		ConnectorMiniCentronics:        "Mini Centronics",
		ConnectorProprietary:           "Proprietary",
		ConnectorDB25PinMale:           "DB-25 pin male",
		ConnectorDB25PinFemale:         "DB-25 pin female",
		ConnectorDB15PinMale:           "DB-15 pin male",
		ConnectorDB15PinFemale:         "DB-15 pin female",
		ConnectorDB9PinMale:            "DB-9 pin male",
		ConnectorDB9PinFemale:          "DB-9 pin female",
		ConnectorRJ11:                  "RJ-11",
		ConnectorRJ45:                  "RJ-45",
		Connector50PinMiniSCSI:         "50-pin MiniSCSI",
		ConnectorMiniDIN:               "Mini-DIN",
		ConnectorMicroDIN:              "Micro-DIN",
		ConnectorPS2:                   "PS/2",
		ConnectorInfrared:              "Infrared",
		ConnectorHPHIL:                 "HP-HIL",
		ConnectorAccessBusUSB:          "Access Bus (USB)",
		ConnectorSSASCSI:               "SSA SCSI",
		ConnectorCircularDIN8Male:      "Circular DIN-8 male",
		ConnectorCircularDIN8Female:    "Circular DIN-8 female",
		ConnectorOnBoardIDE:            "On Board IDE",
		ConnectorOnBoardFloppy:         "On Board Floppy",
		Connector9PinDualInline:        "9-pin Dual Inline",
		Connector25PinDualInline:       "25-pin Dual Inline",
		Connector50PinDualInline:       "50-pin Dual Inline",
		Connector68PinDualInline:       "68-pin Dual Inline",
		ConnectorOnBoardSoundFromCD:    "On Board Sound Input from CD-ROM",
		ConnectorMiniCentronicsType14:  "Mini-Centronics Type-14",
		ConnectorMiniCentronicsType26:  "Mini-Centronics Type-26",
		ConnectorMiniJack:              "Mini-jack (headphones)",
		ConnectorBNC:                   "BNC",
		Connector1394:                  "1394",
		ConnectorSASSATAPlugReceptacle: "SAS/SATA Plug Receptacle",
		ConnectorUSBTypeC:              "USB Type-C Receptacle",
		ConnectorOther:                 "Other",
	}

	if name, ok := types[ct]; ok {
		return name
	}
	return fmt.Sprintf("Unknown (0x%02X)", uint8(ct))
}

// PortType identifies the function of the port
type PortType uint8

// Port type values
const (
	PortTypeNone               PortType = 0x00
	PortTypeParallelXTAT       PortType = 0x01
	PortTypeParallelPS2        PortType = 0x02
	PortTypeParallelECP        PortType = 0x03
	PortTypeParallelEPP        PortType = 0x04
	PortTypeParallelECPEPP     PortType = 0x05
	PortTypeSerialXTAT         PortType = 0x06
	PortTypeSerial16450        PortType = 0x07
	PortTypeSerial16550        PortType = 0x08
	PortTypeSerial16550A       PortType = 0x09
	PortTypeSCSI               PortType = 0x0A
	PortTypeMIDI               PortType = 0x0B
	PortTypeJoystick           PortType = 0x0C
	PortTypeKeyboard           PortType = 0x0D
	PortTypeMouse              PortType = 0x0E
	PortTypeSSASCSI            PortType = 0x0F
	PortTypeUSB                PortType = 0x10
	PortTypeFireWire           PortType = 0x11
	PortTypePCMCIATypeI        PortType = 0x12
	PortTypePCMCIATypeII       PortType = 0x13
	PortTypePCMCIATypeIII      PortType = 0x14
	PortTypeCardbus            PortType = 0x15
	PortTypeAccessBus          PortType = 0x16
	PortTypeSCSIII             PortType = 0x17
	PortTypeSCSIWide           PortType = 0x18
	PortTypePC98               PortType = 0x19
	PortTypePC98Hireso         PortType = 0x1A
	PortTypePCH98              PortType = 0x1B
	PortTypeVideoPort          PortType = 0x1C
	PortTypeAudioPort          PortType = 0x1D
	PortTypeModemPort          PortType = 0x1E
	PortTypeNetworkPort        PortType = 0x1F
	PortTypeSATA               PortType = 0x20
	PortTypeSAS                PortType = 0x21
	PortTypeMFDP               PortType = 0x22 // Multi-Function Display Port
	PortTypeThunderbolt        PortType = 0x23
	PortType8251Compatible     PortType = 0xA0
	PortType8251FIFOCompatible PortType = 0xA1
	PortTypeOther              PortType = 0xFF
)

// String returns a human-readable port type description
func (pt PortType) String() string {
	types := map[PortType]string{
		PortTypeNone:               "None",
		PortTypeParallelXTAT:       "Parallel Port XT/AT Compatible",
		PortTypeParallelPS2:        "Parallel Port PS/2",
		PortTypeParallelECP:        "Parallel Port ECP",
		PortTypeParallelEPP:        "Parallel Port EPP",
		PortTypeParallelECPEPP:     "Parallel Port ECP/EPP",
		PortTypeSerialXTAT:         "Serial Port XT/AT Compatible",
		PortTypeSerial16450:        "Serial Port 16450 Compatible",
		PortTypeSerial16550:        "Serial Port 16550 Compatible",
		PortTypeSerial16550A:       "Serial Port 16550A Compatible",
		PortTypeSCSI:               "SCSI Port",
		PortTypeMIDI:               "MIDI Port",
		PortTypeJoystick:           "Joy Stick Port",
		PortTypeKeyboard:           "Keyboard Port",
		PortTypeMouse:              "Mouse Port",
		PortTypeSSASCSI:            "SSA SCSI",
		PortTypeUSB:                "USB",
		PortTypeFireWire:           "FireWire (IEEE P1394)",
		PortTypePCMCIATypeI:        "PCMCIA Type I",
		PortTypePCMCIATypeII:       "PCMCIA Type II",
		PortTypePCMCIATypeIII:      "PCMCIA Type III",
		PortTypeCardbus:            "Cardbus",
		PortTypeAccessBus:          "Access Bus Port",
		PortTypeSCSIII:             "SCSI II",
		PortTypeSCSIWide:           "SCSI Wide",
		PortTypeVideoPort:          "Video Port",
		PortTypeAudioPort:          "Audio Port",
		PortTypeModemPort:          "Modem Port",
		PortTypeNetworkPort:        "Network Port",
		PortTypeSATA:               "SATA",
		PortTypeSAS:                "SAS",
		PortTypeMFDP:               "Multi-Function Display Port (MFDP)",
		PortTypeThunderbolt:        "Thunderbolt",
		PortType8251Compatible:     "8251 Compatible",
		PortType8251FIFOCompatible: "8251 FIFO Compatible",
		PortTypeOther:              "Other",
	}

	if name, ok := types[pt]; ok {
		return name
	}
	return fmt.Sprintf("Unknown (0x%02X)", uint8(pt))
}

// Parse parses a Port Connector Information structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*PortConnector, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 9 bytes
	if len(s.Data) < 9 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &PortConnector{
		Header:                      s.Header,
		InternalReferenceDesignator: s.GetString(s.GetByte(0x04)),
		InternalConnectorType:       ConnectorType(s.GetByte(0x05)),
		ExternalReferenceDesignator: s.GetString(s.GetByte(0x06)),
		ExternalConnectorType:       ConnectorType(s.GetByte(0x07)),
		PortType:                    PortType(s.GetByte(0x08)),
	}

	return info, nil
}

// Get retrieves the first Port Connector from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*PortConnector, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Port Connector structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*PortConnector, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var ports []*PortConnector
	for i := range structures {
		port, err := Parse(&structures[i])
		if err == nil {
			ports = append(ports, port)
		}
	}

	if len(ports) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return ports, nil
}

// DisplayName returns a display-friendly port description
func (p *PortConnector) DisplayName() string {
	if p.ExternalReferenceDesignator != "" {
		return p.ExternalReferenceDesignator
	}
	if p.InternalReferenceDesignator != "" {
		return p.InternalReferenceDesignator
	}
	return p.PortType.String()
}
