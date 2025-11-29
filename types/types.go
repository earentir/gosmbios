// Package types provides a convenient import for all SMBIOS type packages
package types

// SMBIOS Structure Type Constants
const (
	BIOSInformation          uint8 = 0   // Type 0
	SystemInformation        uint8 = 1   // Type 1
	BaseboardInformation     uint8 = 2   // Type 2
	SystemEnclosure          uint8 = 3   // Type 3
	ProcessorInformation     uint8 = 4   // Type 4
	MemoryController         uint8 = 5   // Type 5 (obsolete)
	MemoryModule             uint8 = 6   // Type 6 (obsolete)
	CacheInformation         uint8 = 7   // Type 7
	PortConnector            uint8 = 8   // Type 8
	SystemSlots              uint8 = 9   // Type 9
	OnBoardDevices           uint8 = 10  // Type 10 (obsolete, use Type 41)
	OEMStrings               uint8 = 11  // Type 11
	SystemConfig             uint8 = 12  // Type 12
	BIOSLanguage             uint8 = 13  // Type 13
	GroupAssociations        uint8 = 14  // Type 14
	SystemEventLog           uint8 = 15  // Type 15
	PhysicalMemoryArray      uint8 = 16  // Type 16
	MemoryDevice             uint8 = 17  // Type 17
	MemoryError32Bit         uint8 = 18  // Type 18
	MemoryArrayMappedAddr    uint8 = 19  // Type 19
	MemoryDeviceMappedAddr   uint8 = 20  // Type 20
	BuiltInPointingDevice    uint8 = 21  // Type 21
	PortableBattery          uint8 = 22  // Type 22
	SystemReset              uint8 = 23  // Type 23
	HardwareSecurity         uint8 = 24  // Type 24
	SystemPowerControls      uint8 = 25  // Type 25
	VoltageProbe             uint8 = 26  // Type 26
	CoolingDevice            uint8 = 27  // Type 27
	TemperatureProbe         uint8 = 28  // Type 28
	ElectricalCurrentProbe   uint8 = 29  // Type 29
	OutOfBandRemoteAccess    uint8 = 30  // Type 30
	BootIntegrityServices    uint8 = 31  // Type 31
	SystemBoot               uint8 = 32  // Type 32
	MemoryError64Bit         uint8 = 33  // Type 33
	ManagementDevice         uint8 = 34  // Type 34
	ManagementDeviceComp     uint8 = 35  // Type 35
	ManagementDeviceThresh   uint8 = 36  // Type 36
	MemoryChannel            uint8 = 37  // Type 37
	IPMIDevice               uint8 = 38  // Type 38
	SystemPowerSupply        uint8 = 39  // Type 39
	Additional               uint8 = 40  // Type 40
	OnboardDevicesExtended   uint8 = 41  // Type 41
	ManagementControllerHost uint8 = 42  // Type 42
	TPMDevice                uint8 = 43  // Type 43
	ProcessorAdditional      uint8 = 44  // Type 44
	FirmwareInventory        uint8 = 45  // Type 45
	StringProperty           uint8 = 46  // Type 46
	Inactive                 uint8 = 126 // Type 126
	EndOfTable               uint8 = 127 // Type 127
)

// TypeName returns the human-readable name for an SMBIOS structure type
func TypeName(structType uint8) string {
	names := map[uint8]string{
		0:   "BIOS Information",
		1:   "System Information",
		2:   "Baseboard Information",
		3:   "System Enclosure",
		4:   "Processor Information",
		5:   "Memory Controller Information",
		6:   "Memory Module Information",
		7:   "Cache Information",
		8:   "Port Connector Information",
		9:   "System Slots",
		10:  "On Board Devices Information",
		11:  "OEM Strings",
		12:  "System Configuration Options",
		13:  "BIOS Language Information",
		14:  "Group Associations",
		15:  "System Event Log",
		16:  "Physical Memory Array",
		17:  "Memory Device",
		18:  "32-bit Memory Error Information",
		19:  "Memory Array Mapped Address",
		20:  "Memory Device Mapped Address",
		21:  "Built-in Pointing Device",
		22:  "Portable Battery",
		23:  "System Reset",
		24:  "Hardware Security",
		25:  "System Power Controls",
		26:  "Voltage Probe",
		27:  "Cooling Device",
		28:  "Temperature Probe",
		29:  "Electrical Current Probe",
		30:  "Out-of-Band Remote Access",
		31:  "Boot Integrity Services Entry Point",
		32:  "System Boot Information",
		33:  "64-bit Memory Error Information",
		34:  "Management Device",
		35:  "Management Device Component",
		36:  "Management Device Threshold Data",
		37:  "Memory Channel",
		38:  "IPMI Device Information",
		39:  "System Power Supply",
		40:  "Additional Information",
		41:  "Onboard Devices Extended Information",
		42:  "Management Controller Host Interface",
		43:  "TPM Device",
		44:  "Processor Additional Information",
		45:  "Firmware Inventory Information",
		46:  "String Property",
		126: "Inactive",
		127: "End-of-Table",
	}

	if name, ok := names[structType]; ok {
		return name
	}
	if structType >= 128 {
		return "OEM-specific"
	}
	return "Unknown"
}
