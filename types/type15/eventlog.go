// Package type15 implements SMBIOS Type 15 - System Event Log
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type15

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for System Event Log
const StructureType uint8 = 15

// SystemEventLog represents Type 15 - System Event Log
type SystemEventLog struct {
	Header                      gosmbios.Header
	LogAreaLength               uint16
	LogHeaderStartOffset        uint16
	LogDataStartOffset          uint16
	AccessMethod                AccessMethod
	LogStatus                   LogStatus
	LogChangeToken              uint32
	AccessMethodAddress         uint32
	LogHeaderFormat             LogHeaderFormat
	NumberOfSupportedLogTypes   uint8
	LengthOfEachLogTypeDesc     uint8
	SupportedEventLogTypes      []LogTypeDescriptor
}

// AccessMethod identifies the method to access the log
type AccessMethod uint8

const (
	AccessIndexedIO8Bit        AccessMethod = 0x00
	AccessIndexedIO2x8Bit      AccessMethod = 0x01
	AccessIndexedIO16Bit       AccessMethod = 0x02
	AccessMemoryMapped32Bit    AccessMethod = 0x03
	AccessGPNV                 AccessMethod = 0x04
)

func (a AccessMethod) String() string {
	switch a {
	case AccessIndexedIO8Bit:
		return "Indexed I/O, one 8-bit index port, one 8-bit data port"
	case AccessIndexedIO2x8Bit:
		return "Indexed I/O, two 8-bit index ports, one 8-bit data port"
	case AccessIndexedIO16Bit:
		return "Indexed I/O, one 16-bit index port, one 8-bit data port"
	case AccessMemoryMapped32Bit:
		return "Memory-mapped physical 32-bit address"
	case AccessGPNV:
		return "General Purpose Non-Volatile Data functions"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(a))
	}
}

// LogStatus represents the current status of the event log
type LogStatus uint8

const (
	LogStatusValid    LogStatus = 1 << 0 // Log area valid
	LogStatusFull     LogStatus = 1 << 1 // Log area full
)

func (l LogStatus) String() string {
	var status []string
	if l&LogStatusValid != 0 {
		status = append(status, "Valid")
	} else {
		status = append(status, "Invalid")
	}
	if l&LogStatusFull != 0 {
		status = append(status, "Full")
	}
	if len(status) == 0 {
		return "Unknown"
	}
	return fmt.Sprintf("%v", status)
}

// IsValid returns true if the log area is valid
func (l LogStatus) IsValid() bool {
	return l&LogStatusValid != 0
}

// IsFull returns true if the log area is full
func (l LogStatus) IsFull() bool {
	return l&LogStatusFull != 0
}

// LogHeaderFormat identifies the format of the log header
type LogHeaderFormat uint8

const (
	LogHeaderNone   LogHeaderFormat = 0x00
	LogHeaderType1  LogHeaderFormat = 0x01
)

func (l LogHeaderFormat) String() string {
	switch l {
	case LogHeaderNone:
		return "No header"
	case LogHeaderType1:
		return "Type 1 log header"
	default:
		if l >= 0x80 {
			return fmt.Sprintf("OEM-specific (0x%02X)", uint8(l))
		}
		return fmt.Sprintf("Unknown (0x%02X)", uint8(l))
	}
}

// LogTypeDescriptor describes a supported event log type
type LogTypeDescriptor struct {
	LogType            EventLogType
	VariableDataFormat VariableDataFormat
}

// EventLogType identifies the type of event
type EventLogType uint8

const (
	EventLogReserved                  EventLogType = 0x00
	EventLogSingleBitECC              EventLogType = 0x01
	EventLogMultiBitECC               EventLogType = 0x02
	EventLogParityMemory              EventLogType = 0x03
	EventLogBusTimeout                EventLogType = 0x04
	EventLogIOChannelCheck            EventLogType = 0x05
	EventLogSoftwareNMI               EventLogType = 0x06
	EventLogPOSTMemoryResize          EventLogType = 0x07
	EventLogPOSTError                 EventLogType = 0x08
	EventLogPCIParityError            EventLogType = 0x09
	EventLogPCISystemError            EventLogType = 0x0A
	EventLogCPUFailure                EventLogType = 0x0B
	EventLogEISAFailsafe              EventLogType = 0x0C
	EventLogLogCleared                EventLogType = 0x0D
	EventLogBootEvent                 EventLogType = 0x0E
	EventLogLogDisabled               EventLogType = 0x0F
	EventLogLogLimit                  EventLogType = 0x10
	EventLogHardwareWatchdog          EventLogType = 0x11
	EventLogSystemStart               EventLogType = 0x12
	EventLogSystemHardwareSecViol     EventLogType = 0x13
	EventLogAuxLog                    EventLogType = 0x14
	EventLogPCIResource               EventLogType = 0x15
	EventLogSBMCMessage               EventLogType = 0x16
	EventLogManagementHW              EventLogType = 0x17
	EventLogEndOfLog                  EventLogType = 0xFF
)

func (e EventLogType) String() string {
	switch e {
	case EventLogReserved:
		return "Reserved"
	case EventLogSingleBitECC:
		return "Single-bit ECC memory error"
	case EventLogMultiBitECC:
		return "Multi-bit ECC memory error"
	case EventLogParityMemory:
		return "Parity memory error"
	case EventLogBusTimeout:
		return "Bus timeout"
	case EventLogIOChannelCheck:
		return "I/O channel check"
	case EventLogSoftwareNMI:
		return "Software NMI"
	case EventLogPOSTMemoryResize:
		return "POST memory resize"
	case EventLogPOSTError:
		return "POST error"
	case EventLogPCIParityError:
		return "PCI parity error"
	case EventLogPCISystemError:
		return "PCI system error"
	case EventLogCPUFailure:
		return "CPU failure"
	case EventLogEISAFailsafe:
		return "EISA FailSafe timer timeout"
	case EventLogLogCleared:
		return "Log cleared"
	case EventLogBootEvent:
		return "Boot event"
	case EventLogLogDisabled:
		return "Logging disabled"
	case EventLogLogLimit:
		return "Log limit reached"
	case EventLogHardwareWatchdog:
		return "Hardware watchdog reset"
	case EventLogSystemStart:
		return "System start"
	case EventLogSystemHardwareSecViol:
		return "System hardware security violation"
	case EventLogAuxLog:
		return "Auxiliary log entry support"
	case EventLogPCIResource:
		return "PCI resource configuration"
	case EventLogSBMCMessage:
		return "System board MC message"
	case EventLogManagementHW:
		return "Management hardware"
	case EventLogEndOfLog:
		return "End of log"
	default:
		if e >= 0x80 && e <= 0xFE {
			return fmt.Sprintf("OEM-specific (0x%02X)", uint8(e))
		}
		return fmt.Sprintf("Unknown (0x%02X)", uint8(e))
	}
}

// VariableDataFormat identifies the format of variable data
type VariableDataFormat uint8

const (
	VarDataNone      VariableDataFormat = 0x00
	VarDataHandle    VariableDataFormat = 0x01
	VarDataMultiple  VariableDataFormat = 0x02
	VarDataPOSTCodes VariableDataFormat = 0x03
	VarDataTimeStamp VariableDataFormat = 0x04
	VarDataTime      VariableDataFormat = 0x05
	VarDataSymbol    VariableDataFormat = 0x06
)

func (v VariableDataFormat) String() string {
	switch v {
	case VarDataNone:
		return "None"
	case VarDataHandle:
		return "Handle"
	case VarDataMultiple:
		return "Multiple-Event"
	case VarDataPOSTCodes:
		return "Multiple-Event Handle"
	case VarDataTimeStamp:
		return "POST Results Bitmap"
	case VarDataTime:
		return "System Management Type"
	case VarDataSymbol:
		return "Multiple-Event System Management Type"
	default:
		if v >= 0x80 {
			return fmt.Sprintf("OEM-specific (0x%02X)", uint8(v))
		}
		return fmt.Sprintf("Unknown (0x%02X)", uint8(v))
	}
}

// Parse parses a System Event Log structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*SystemEventLog, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 20 bytes
	if len(s.Data) < 20 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &SystemEventLog{
		Header:                      s.Header,
		LogAreaLength:               s.GetWord(0x04),
		LogHeaderStartOffset:        s.GetWord(0x06),
		LogDataStartOffset:          s.GetWord(0x08),
		AccessMethod:                AccessMethod(s.GetByte(0x0A)),
		LogStatus:                   LogStatus(s.GetByte(0x0B)),
		LogChangeToken:              s.GetDWord(0x0C),
		AccessMethodAddress:         s.GetDWord(0x10),
		LogHeaderFormat:             LogHeaderFormat(s.GetByte(0x14)),
		NumberOfSupportedLogTypes:   s.GetByte(0x15),
		LengthOfEachLogTypeDesc:     s.GetByte(0x16),
	}

	// Read supported event log type descriptors
	if len(s.Data) >= 23 && info.LengthOfEachLogTypeDesc >= 2 {
		offset := 0x17
		for i := uint8(0); i < info.NumberOfSupportedLogTypes; i++ {
			if offset+1 >= len(s.Data) {
				break
			}
			desc := LogTypeDescriptor{
				LogType:            EventLogType(s.GetByte(offset)),
				VariableDataFormat: VariableDataFormat(s.GetByte(offset + 1)),
			}
			info.SupportedEventLogTypes = append(info.SupportedEventLogTypes, desc)
			offset += int(info.LengthOfEachLogTypeDesc)
		}
	}

	return info, nil
}

// Get retrieves the System Event Log from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*SystemEventLog, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all System Event Log structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*SystemEventLog, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var logs []*SystemEventLog
	for i := range structures {
		log, err := Parse(&structures[i])
		if err == nil {
			logs = append(logs, log)
		}
	}

	if len(logs) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return logs, nil
}
