// Package type23 implements SMBIOS Type 23 - System Reset
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type23

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for System Reset
const StructureType uint8 = 23

// SystemReset represents Type 23 - System Reset
type SystemReset struct {
	Header        gosmbios.Header
	Capabilities  Capabilities
	ResetCount    uint16
	ResetLimit    uint16
	TimerInterval uint16 // In minutes
	Timeout       uint16 // In minutes
}

// Capabilities represents system reset capabilities
type Capabilities uint8

// IsEnabled returns true if system reset is enabled
func (c Capabilities) IsEnabled() bool {
	return c&0x01 != 0
}

// BootOption returns the boot option on successful reset
func (c Capabilities) BootOption() BootOption {
	return BootOption((c >> 1) & 0x03)
}

// BootOptionOnLimit returns the boot option when reset limit is reached
func (c Capabilities) BootOptionOnLimit() BootOption {
	return BootOption((c >> 3) & 0x03)
}

// WatchdogTimerPresent returns true if watchdog timer is present
func (c Capabilities) WatchdogTimerPresent() bool {
	return c&0x20 != 0
}

func (c Capabilities) String() string {
	status := "Disabled"
	if c.IsEnabled() {
		status = "Enabled"
	}
	watchdog := "No"
	if c.WatchdogTimerPresent() {
		watchdog = "Yes"
	}
	return fmt.Sprintf("Status: %s, Boot Option: %s, Watchdog: %s",
		status, c.BootOption().String(), watchdog)
}

// BootOption identifies boot behavior
type BootOption uint8

const (
	BootOptionReserved      BootOption = 0x00
	BootOptionOS            BootOption = 0x01
	BootOptionSystemUtility BootOption = 0x02
	BootOptionDoNotReboot   BootOption = 0x03
)

func (b BootOption) String() string {
	switch b {
	case BootOptionReserved:
		return "Reserved"
	case BootOptionOS:
		return "Operating System"
	case BootOptionSystemUtility:
		return "System Utility"
	case BootOptionDoNotReboot:
		return "Do Not Reboot"
	default:
		return fmt.Sprintf("Unknown (0x%02X)", uint8(b))
	}
}

// Parse parses a System Reset structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*SystemReset, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 13 bytes
	if len(s.Data) < 13 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &SystemReset{
		Header:        s.Header,
		Capabilities:  Capabilities(s.GetByte(0x04)),
		ResetCount:    s.GetWord(0x05),
		ResetLimit:    s.GetWord(0x07),
		TimerInterval: s.GetWord(0x09),
		Timeout:       s.GetWord(0x0B),
	}

	return info, nil
}

// Get retrieves the System Reset from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*SystemReset, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// ResetCountString returns the reset count as a string
func (s *SystemReset) ResetCountString() string {
	if s.ResetCount == 0xFFFF {
		return "Unknown"
	}
	return fmt.Sprintf("%d", s.ResetCount)
}

// ResetLimitString returns the reset limit as a string
func (s *SystemReset) ResetLimitString() string {
	if s.ResetLimit == 0xFFFF {
		return "Unknown"
	}
	return fmt.Sprintf("%d", s.ResetLimit)
}

// TimerIntervalString returns the timer interval as a string
func (s *SystemReset) TimerIntervalString() string {
	if s.TimerInterval == 0xFFFF {
		return "Unknown"
	}
	return fmt.Sprintf("%d minutes", s.TimerInterval)
}

// TimeoutString returns the timeout as a string
func (s *SystemReset) TimeoutString() string {
	if s.Timeout == 0xFFFF {
		return "Unknown"
	}
	return fmt.Sprintf("%d minutes", s.Timeout)
}
