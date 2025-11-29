// Package type25 implements SMBIOS Type 25 - System Power Controls
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type25

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for System Power Controls
const StructureType uint8 = 25

// SystemPowerControls represents Type 25 - System Power Controls
type SystemPowerControls struct {
	Header          gosmbios.Header
	NextScheduledPowerOnMonth   uint8
	NextScheduledPowerOnDay     uint8
	NextScheduledPowerOnHour    uint8
	NextScheduledPowerOnMinute  uint8
	NextScheduledPowerOnSecond  uint8
}

// Parse parses a System Power Controls structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*SystemPowerControls, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 9 bytes
	if len(s.Data) < 9 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &SystemPowerControls{
		Header:                     s.Header,
		NextScheduledPowerOnMonth:  s.GetByte(0x04),
		NextScheduledPowerOnDay:    s.GetByte(0x05),
		NextScheduledPowerOnHour:   s.GetByte(0x06),
		NextScheduledPowerOnMinute: s.GetByte(0x07),
		NextScheduledPowerOnSecond: s.GetByte(0x08),
	}

	return info, nil
}

// Get retrieves the System Power Controls from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*SystemPowerControls, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// NextPowerOnString returns the next scheduled power-on time as a string
func (s *SystemPowerControls) NextPowerOnString() string {
	// BCD format: 0xFF means unspecified
	month := bcdToInt(s.NextScheduledPowerOnMonth)
	day := bcdToInt(s.NextScheduledPowerOnDay)
	hour := bcdToInt(s.NextScheduledPowerOnHour)
	minute := bcdToInt(s.NextScheduledPowerOnMinute)
	second := bcdToInt(s.NextScheduledPowerOnSecond)

	if month == -1 && day == -1 && hour == -1 && minute == -1 && second == -1 {
		return "Not Scheduled"
	}

	result := ""
	if month != -1 && day != -1 {
		result = fmt.Sprintf("%02d/%02d", month, day)
	}
	if hour != -1 && minute != -1 {
		if result != "" {
			result += " "
		}
		if second != -1 {
			result += fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
		} else {
			result += fmt.Sprintf("%02d:%02d", hour, minute)
		}
	}

	if result == "" {
		return "Partially Specified"
	}
	return result
}

// bcdToInt converts BCD value to integer, returns -1 for 0xFF (unspecified)
func bcdToInt(bcd uint8) int {
	if bcd == 0xFF {
		return -1
	}
	return int((bcd>>4)*10 + (bcd & 0x0F))
}

// IsScheduled returns true if a power-on is scheduled
func (s *SystemPowerControls) IsScheduled() bool {
	return s.NextScheduledPowerOnMonth != 0xFF ||
		s.NextScheduledPowerOnDay != 0xFF ||
		s.NextScheduledPowerOnHour != 0xFF ||
		s.NextScheduledPowerOnMinute != 0xFF ||
		s.NextScheduledPowerOnSecond != 0xFF
}
