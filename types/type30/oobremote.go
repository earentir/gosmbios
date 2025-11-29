// Package type30 implements SMBIOS Type 30 - Out-of-Band Remote Access
// Per DSP0134 SMBIOS Reference Specification 3.9.0
package type30

import (
	"fmt"

	"github.com/earentir/gosmbios"
)

// StructureType is the SMBIOS structure type for Out-of-Band Remote Access
const StructureType uint8 = 30

// OutOfBandRemoteAccess represents Type 30 - Out-of-Band Remote Access
type OutOfBandRemoteAccess struct {
	Header             gosmbios.Header
	ManufacturerName   string
	Connections        Connections
}

// Connections represents remote access connection information
type Connections uint8

// InboundEnabled returns true if inbound connection is enabled
func (c Connections) InboundEnabled() bool {
	return c&0x01 != 0
}

// OutboundEnabled returns true if outbound connection is enabled
func (c Connections) OutboundEnabled() bool {
	return c&0x02 != 0
}

func (c Connections) String() string {
	inbound := "Disabled"
	if c.InboundEnabled() {
		inbound = "Enabled"
	}
	outbound := "Disabled"
	if c.OutboundEnabled() {
		outbound = "Enabled"
	}
	return fmt.Sprintf("Inbound: %s, Outbound: %s", inbound, outbound)
}

// Parse parses an Out-of-Band Remote Access structure from raw SMBIOS data
func Parse(s *gosmbios.Structure) (*OutOfBandRemoteAccess, error) {
	if s == nil || s.Header.Type != StructureType {
		return nil, gosmbios.ErrInvalidStructure
	}

	// Minimum length is 6 bytes
	if len(s.Data) < 6 {
		return nil, gosmbios.ErrInvalidStructure
	}

	info := &OutOfBandRemoteAccess{
		Header:           s.Header,
		ManufacturerName: s.GetString(s.GetByte(0x04)),
		Connections:      Connections(s.GetByte(0x05)),
	}

	return info, nil
}

// Get retrieves the Out-of-Band Remote Access from SMBIOS data
func Get(sm *gosmbios.SMBIOS) (*OutOfBandRemoteAccess, error) {
	s := sm.GetStructure(StructureType)
	if s == nil {
		return nil, gosmbios.ErrNotFound
	}
	return Parse(s)
}
