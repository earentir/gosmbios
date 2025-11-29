// Package type36 implements SMBIOS Type 36 - Management Device Threshold Data
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type36

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Management Device Threshold Data
const StructureType uint8 = 36

// ManagementDeviceThreshold represents Type 36 - Management Device Threshold Data
type ManagementDeviceThreshold struct {
	Header                   gosmbios.Header
	LowerThresholdNonCritical uint16
	UpperThresholdNonCritical uint16
	LowerThresholdCritical    uint16
	UpperThresholdCritical    uint16
	LowerThresholdNonRecoverable uint16
	UpperThresholdNonRecoverable uint16
}

// Parse parses a Management Device Threshold Data structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*ManagementDeviceThreshold, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 16 bytes
	if len(s.Data) < 16 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &ManagementDeviceThreshold{
		Header:                       s.Header,
		LowerThresholdNonCritical:    s.GetWord(0x04),
		UpperThresholdNonCritical:    s.GetWord(0x06),
		LowerThresholdCritical:       s.GetWord(0x08),
		UpperThresholdCritical:       s.GetWord(0x0A),
		LowerThresholdNonRecoverable: s.GetWord(0x0C),
		UpperThresholdNonRecoverable: s.GetWord(0x0E),
	}

	return info, nil
}

// Get retrieves the first Management Device Threshold from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*ManagementDeviceThreshold, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}

// GetAll retrieves all Management Device Threshold structures from SMBIOS data
func GetAll(sm *gosmbios.SMBIOS) ([]*ManagementDeviceThreshold, error) {
	structures := sm.GetStructures(StructureType)
	if len(structures) == 0 {
		return nil, gosmbios.ErrNotFound
	}

	var thresholds []*ManagementDeviceThreshold
	for i := range structures {
		thresh, err := Parse(&structures[i])
		if err == nil {
			thresholds = append(thresholds, thresh)
		}
	}

	if len(thresholds) == 0 {
		return nil, gosmbios.ErrNotFound
	}
	return thresholds, nil
}

// thresholdString returns a threshold value as a string
func thresholdString(value uint16) string {
	if value == 0x8000 {
		return "Not Supported"
	}
	return fmt.Sprintf("%d", value)
}

// LowerNonCriticalString returns the lower non-critical threshold as a string
func (m *ManagementDeviceThreshold) LowerNonCriticalString() string {
	return thresholdString(m.LowerThresholdNonCritical)
}

// UpperNonCriticalString returns the upper non-critical threshold as a string
func (m *ManagementDeviceThreshold) UpperNonCriticalString() string {
	return thresholdString(m.UpperThresholdNonCritical)
}

// LowerCriticalString returns the lower critical threshold as a string
func (m *ManagementDeviceThreshold) LowerCriticalString() string {
	return thresholdString(m.LowerThresholdCritical)
}

// UpperCriticalString returns the upper critical threshold as a string
func (m *ManagementDeviceThreshold) UpperCriticalString() string {
	return thresholdString(m.UpperThresholdCritical)
}

// LowerNonRecoverableString returns the lower non-recoverable threshold as a string
func (m *ManagementDeviceThreshold) LowerNonRecoverableString() string {
	return thresholdString(m.LowerThresholdNonRecoverable)
}

// UpperNonRecoverableString returns the upper non-recoverable threshold as a string
func (m *ManagementDeviceThreshold) UpperNonRecoverableString() string {
	return thresholdString(m.UpperThresholdNonRecoverable)
}
