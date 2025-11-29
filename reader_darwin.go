//go:build darwin

package gosmbios

import (
	"bytes"
	"encoding/binary"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// readSMBIOS reads SMBIOS data on macOS systems
// macOS doesn't expose raw SMBIOS tables directly like Linux or Windows
// We synthesize SMBIOS-compatible structures from available system information
func readSMBIOS() (*SMBIOS, error) {
	var structures []Structure

	// Try IOPlatformExpertDevice first (works on all Macs including Apple Silicon)
	ioregStructures := readFromIOPlatformExpert()
	structures = append(structures, ioregStructures...)

	// Get additional info from system_profiler
	profilerStructures := readFromSystemProfiler()

	// Merge structures, avoiding duplicate types that are already present
	// (but allowing multiple structures of the same type from the same source)
	existingTypes := make(map[uint8]bool)
	for _, s := range structures {
		existingTypes[s.Header.Type] = true
	}

	// Track which types we've started adding from profiler
	addedFromProfiler := make(map[uint8]bool)
	for _, s := range profilerStructures {
		// Add if type doesn't exist in ioreg structures, OR if we've already started
		// adding this type from profiler (to allow multiple caches, memory devices, etc.)
		if !existingTypes[s.Header.Type] || addedFromProfiler[s.Header.Type] {
			structures = append(structures, s)
			addedFromProfiler[s.Header.Type] = true
		}
	}

	if len(structures) == 0 {
		return nil, ErrNotFound
	}

	// Add End-of-Table
	structures = append(structures, Structure{
		Header: Header{Type: 127, Length: 4, Handle: 0xFFFF},
		Data:   []byte{127, 4, 0xFF, 0xFF},
	})

	// Create a synthetic entry point since macOS doesn't expose the raw entry point
	entryPoint := EntryPoint{
		Type:         EntryPoint64Bit,
		MajorVersion: 3,
		MinorVersion: 0,
		Revision:     0,
	}

	return &SMBIOS{
		EntryPoint: entryPoint,
		Structures: structures,
	}, nil
}

// readFromIOPlatformExpert reads system info from IOPlatformExpertDevice
func readFromIOPlatformExpert() []Structure {
	var structures []Structure

	// Get platform info
	cmd := exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
	output, err := cmd.Output()
	if err != nil {
		return structures
	}

	props := parseIORegProperties(output)

	if len(props) > 0 {
		// Create Type 1 - System Information
		structures = append(structures, createSystemInfoFromIOProps(props))

		// Create Type 2 - Baseboard Information
		structures = append(structures, createBaseboardFromIOProps(props))

		// Create Type 3 - Chassis Information
		structures = append(structures, createChassisFromIOProps(props))
	}

	return structures
}

// readFromSystemProfiler reads additional info from system_profiler
func readFromSystemProfiler() []Structure {
	var structures []Structure

	// Get hardware info
	cmd := exec.Command("system_profiler", "SPHardwareDataType")
	output, err := cmd.Output()
	if err != nil {
		return structures
	}

	hwInfo := parseSystemProfilerOutput(string(output))

	// Create Type 0 - BIOS Information if we have data
	if len(hwInfo) > 0 {
		structures = append(structures, createBIOSInfoStructure(hwInfo))
	}

	// Create Type 4 - Processor Information
	procStructures := createProcessorStructures(hwInfo)
	structures = append(structures, procStructures...)

	// Create Type 7 - Cache Information
	cacheStructures := createCacheStructures()
	structures = append(structures, cacheStructures...)

	// Get memory info (works on Intel Macs, may be limited on Apple Silicon)
	cmd = exec.Command("system_profiler", "SPMemoryDataType")
	output, _ = cmd.Output()
	memStructures := parseSystemProfilerMemory(string(output))

	if len(memStructures) > 0 {
		// Add memory array first
		structures = append(structures, createMemoryArrayStructure(len(memStructures)))
		structures = append(structures, memStructures...)
	} else if len(hwInfo) > 0 {
		// Apple Silicon - create synthetic memory info from hardware info
		if memStr, ok := hwInfo["Memory"]; ok {
			structures = append(structures, createMemoryArrayStructure(1))
			structures = append(structures, createSyntheticMemoryDevice(memStr))
		}
	}

	return structures
}

// parseIORegProperties extracts key-value pairs from ioreg output
func parseIORegProperties(data []byte) map[string]string {
	props := make(map[string]string)

	patterns := map[string]*regexp.Regexp{
		"manufacturer":    regexp.MustCompile(`"manufacturer"\s*=\s*<?"?([^"<>]+)"?`),
		"model":           regexp.MustCompile(`"model"\s*=\s*<?"?([^"<>]+)"?`),
		"product-name":    regexp.MustCompile(`"product-name"\s*=\s*<?"?([^"<>]+)"?`),
		"serial-number":   regexp.MustCompile(`"IOPlatformSerialNumber"\s*=\s*"([^"]+)"`),
		"uuid":            regexp.MustCompile(`"IOPlatformUUID"\s*=\s*"([^"]+)"`),
		"board-id":        regexp.MustCompile(`"board-id"\s*=\s*<?"?([^"<>]+)"?`),
		"target-sub-type": regexp.MustCompile(`"target-sub-type"\s*=\s*"([^"]+)"`),
	}

	for key, pattern := range patterns {
		match := pattern.FindSubmatch(data)
		if len(match) > 1 {
			value := strings.TrimSpace(string(match[1]))
			// Clean up byte sequences like <"Mac..."
			value = strings.Trim(value, "<>\"")
			if value != "" && value != "0" {
				props[key] = value
			}
		}
	}

	return props
}

// parseSystemProfilerOutput parses system_profiler output into key-value pairs
func parseSystemProfilerOutput(output string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if idx := strings.Index(line, ":"); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			value := strings.TrimSpace(line[idx+1:])
			if value != "" {
				result[key] = value
			}
		}
	}

	return result
}

// parseSystemProfilerMemory parses memory information from system_profiler
func parseSystemProfilerMemory(output string) []Structure {
	var structures []Structure
	var handle uint16 = 0x1100

	lines := strings.Split(output, "\n")
	var currentSlot map[string]string
	inSlot := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Detect memory slot start
		if strings.HasPrefix(trimmed, "BANK") || strings.HasPrefix(trimmed, "DIMM") ||
			strings.Contains(trimmed, "Slot") {
			if currentSlot != nil && len(currentSlot) > 1 {
				structures = append(structures, createMemoryDeviceStructure(currentSlot, handle))
				handle++
			}
			currentSlot = make(map[string]string)
			currentSlot["slot"] = trimmed
			inSlot = true
			continue
		}

		if inSlot && strings.Contains(trimmed, ":") {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				if value != "" && value != "(Empty)" {
					currentSlot[key] = value
				}
			}
		}
	}

	// Don't forget the last slot
	if currentSlot != nil && len(currentSlot) > 1 {
		structures = append(structures, createMemoryDeviceStructure(currentSlot, handle))
	}

	return structures
}

// createSystemInfoFromIOProps creates Type 1 - System Information
func createSystemInfoFromIOProps(props map[string]string) Structure {
	var strTable []string
	var data bytes.Buffer

	// Write header
	data.WriteByte(1)                                        // Type
	data.WriteByte(27)                                       // Length (SMBIOS 2.4+)
	binary.Write(&data, binary.LittleEndian, uint16(0x0001)) // Handle

	// Manufacturer (offset 0x04)
	manufacturer := "Apple Inc."
	strTable = append(strTable, manufacturer)
	data.WriteByte(uint8(len(strTable)))

	// Product Name (offset 0x05)
	productName := ""
	if v, ok := props["model"]; ok {
		productName = v
	} else if v, ok := props["product-name"]; ok {
		productName = v
	}
	strTable = append(strTable, productName)
	data.WriteByte(uint8(len(strTable)))

	// Version (offset 0x06)
	version := "1.0"
	strTable = append(strTable, version)
	data.WriteByte(uint8(len(strTable)))

	// Serial Number (offset 0x07)
	serial := ""
	if v, ok := props["serial-number"]; ok {
		serial = v
	}
	strTable = append(strTable, serial)
	data.WriteByte(uint8(len(strTable)))

	// UUID (offset 0x08) - 16 bytes
	uuid := props["uuid"]
	uuidBytes := parseUUID(uuid)
	data.Write(uuidBytes)

	// Wake-up Type (offset 0x18)
	data.WriteByte(0x06) // Power Switch

	// SKU Number (offset 0x19)
	sku := ""
	if v, ok := props["board-id"]; ok {
		sku = v
	}
	strTable = append(strTable, sku)
	data.WriteByte(uint8(len(strTable)))

	// Family (offset 0x1A)
	family := "Mac"
	strTable = append(strTable, family)
	data.WriteByte(uint8(len(strTable)))

	return Structure{
		Header:  Header{Type: 1, Length: 27, Handle: 0x0001},
		Data:    data.Bytes(),
		Strings: strTable,
	}
}

// createBaseboardFromIOProps creates Type 2 - Baseboard Information
func createBaseboardFromIOProps(props map[string]string) Structure {
	var strTable []string
	var data bytes.Buffer

	// Write header
	data.WriteByte(2)                                        // Type
	data.WriteByte(15)                                       // Length
	binary.Write(&data, binary.LittleEndian, uint16(0x0002)) // Handle

	// Manufacturer (offset 0x04)
	strTable = append(strTable, "Apple Inc.")
	data.WriteByte(uint8(len(strTable)))

	// Product (offset 0x05)
	product := ""
	if v, ok := props["board-id"]; ok {
		product = v
	} else if v, ok := props["model"]; ok {
		product = v
	}
	strTable = append(strTable, product)
	data.WriteByte(uint8(len(strTable)))

	// Version (offset 0x06)
	strTable = append(strTable, "1.0")
	data.WriteByte(uint8(len(strTable)))

	// Serial Number (offset 0x07)
	serial := ""
	if v, ok := props["serial-number"]; ok {
		serial = v
	}
	strTable = append(strTable, serial)
	data.WriteByte(uint8(len(strTable)))

	// Asset Tag (offset 0x08)
	strTable = append(strTable, "")
	data.WriteByte(uint8(len(strTable)))

	// Feature Flags (offset 0x09)
	data.WriteByte(0x01) // Hosting board

	// Location in Chassis (offset 0x0A)
	strTable = append(strTable, "")
	data.WriteByte(uint8(len(strTable)))

	// Chassis Handle (offset 0x0B)
	binary.Write(&data, binary.LittleEndian, uint16(0x0003))

	// Board Type (offset 0x0D)
	data.WriteByte(0x0A) // Motherboard

	// Number of Contained Object Handles (offset 0x0E)
	data.WriteByte(0)

	return Structure{
		Header:  Header{Type: 2, Length: 15, Handle: 0x0002},
		Data:    data.Bytes(),
		Strings: strTable,
	}
}

// createChassisFromIOProps creates Type 3 - Chassis Information
func createChassisFromIOProps(props map[string]string) Structure {
	var strTable []string
	var data bytes.Buffer

	// Write header
	data.WriteByte(3)                                        // Type
	data.WriteByte(22)                                       // Length (SMBIOS 2.3+)
	binary.Write(&data, binary.LittleEndian, uint16(0x0003)) // Handle

	// Manufacturer (offset 0x04)
	strTable = append(strTable, "Apple Inc.")
	data.WriteByte(uint8(len(strTable)))

	// Chassis Type (offset 0x05) - determine from model
	chassisType := uint8(0x09) // Laptop default
	model := strings.ToLower(props["model"])
	if strings.Contains(model, "macpro") || strings.Contains(model, "mac pro") {
		chassisType = 0x07 // Tower
	} else if strings.Contains(model, "macmini") || strings.Contains(model, "mac mini") {
		chassisType = 0x23 // Mini PC
	} else if strings.Contains(model, "imac") {
		chassisType = 0x0D // All-in-One
	} else if strings.Contains(model, "macstudio") || strings.Contains(model, "mac studio") {
		chassisType = 0x23 // Mini PC
	} else if strings.Contains(model, "macbook") {
		chassisType = 0x0A // Notebook
	}
	data.WriteByte(chassisType)

	// Version (offset 0x06)
	strTable = append(strTable, "1.0")
	data.WriteByte(uint8(len(strTable)))

	// Serial Number (offset 0x07)
	serial := ""
	if v, ok := props["serial-number"]; ok {
		serial = v
	}
	strTable = append(strTable, serial)
	data.WriteByte(uint8(len(strTable)))

	// Asset Tag (offset 0x08)
	strTable = append(strTable, "")
	data.WriteByte(uint8(len(strTable)))

	// Boot-up State (offset 0x09)
	data.WriteByte(0x03) // Safe

	// Power Supply State (offset 0x0A)
	data.WriteByte(0x03) // Safe

	// Thermal State (offset 0x0B)
	data.WriteByte(0x03) // Safe

	// Security Status (offset 0x0C)
	data.WriteByte(0x02) // Unknown

	// OEM-defined (offset 0x0D) - 4 bytes
	binary.Write(&data, binary.LittleEndian, uint32(0))

	// Height (offset 0x11)
	data.WriteByte(0) // Unspecified

	// Number of Power Cords (offset 0x12)
	data.WriteByte(1)

	// Contained Element Count (offset 0x13)
	data.WriteByte(0)

	// Contained Element Record Length (offset 0x14)
	data.WriteByte(0)

	// SKU Number (offset 0x15)
	strTable = append(strTable, "")
	data.WriteByte(uint8(len(strTable)))

	return Structure{
		Header:  Header{Type: 3, Length: 22, Handle: 0x0003},
		Data:    data.Bytes(),
		Strings: strTable,
	}
}

// createBIOSInfoStructure creates a Type 0 structure from system_profiler data
func createBIOSInfoStructure(info map[string]string) Structure {
	var strTable []string
	var data bytes.Buffer

	// Write header
	data.WriteByte(0)                                        // Type 0 - BIOS Information
	data.WriteByte(26)                                       // Length (SMBIOS 3.1+)
	binary.Write(&data, binary.LittleEndian, uint16(0x0000)) // Handle

	// Vendor (offset 0x04)
	strTable = append(strTable, "Apple Inc.")
	data.WriteByte(uint8(len(strTable)))

	// BIOS Version (offset 0x05)
	biosVersion := ""
	if v, ok := info["System Firmware Version"]; ok {
		biosVersion = v
	} else if v, ok := info["Boot ROM Version"]; ok {
		biosVersion = v
	} else if v, ok := info["OS Loader Version"]; ok {
		biosVersion = v
	}
	strTable = append(strTable, biosVersion)
	data.WriteByte(uint8(len(strTable)))

	// BIOS Starting Address Segment (offset 0x06)
	binary.Write(&data, binary.LittleEndian, uint16(0xE000))

	// BIOS Release Date (offset 0x08)
	strTable = append(strTable, "")
	data.WriteByte(uint8(len(strTable)))

	// BIOS ROM Size (offset 0x09)
	data.WriteByte(0xFF) // Extended BIOS ROM Size used

	// BIOS Characteristics (offset 0x0A) - 8 bytes
	binary.Write(&data, binary.LittleEndian, uint64(0x000000007FFBDE90))

	// BIOS Characteristics Extension Bytes (offset 0x12) - 2 bytes
	data.WriteByte(0x03) // ACPI, USB Legacy
	data.WriteByte(0x0D) // UEFI, Boot spec, Target content distribution

	// System BIOS Major Release (offset 0x14)
	data.WriteByte(0xFF)

	// System BIOS Minor Release (offset 0x15)
	data.WriteByte(0xFF)

	// Embedded Controller Major Release (offset 0x16)
	data.WriteByte(0xFF)

	// Embedded Controller Minor Release (offset 0x17)
	data.WriteByte(0xFF)

	// Extended BIOS ROM Size (offset 0x18) - 2 bytes
	binary.Write(&data, binary.LittleEndian, uint16(0))

	return Structure{
		Header:  Header{Type: 0, Length: 26, Handle: 0x0000},
		Data:    data.Bytes(),
		Strings: strTable,
	}
}

// createProcessorStructures creates Type 4 - Processor Information structures
func createProcessorStructures(hwInfo map[string]string) []Structure {
	var structures []Structure

	// Get CPU info from sysctl
	cpuInfo := getSysctlCPUInfo()

	// Merge with hwInfo
	for k, v := range hwInfo {
		if _, exists := cpuInfo[k]; !exists {
			cpuInfo[k] = v
		}
	}

	if len(cpuInfo) == 0 {
		return structures
	}

	var strTable []string
	var data bytes.Buffer

	// Write header
	data.WriteByte(4)                                        // Type 4 - Processor
	data.WriteByte(48)                                       // Length (SMBIOS 3.0)
	binary.Write(&data, binary.LittleEndian, uint16(0x0004)) // Handle

	// Socket Designation (offset 0x04)
	socket := "U3E1"
	if _, ok := cpuInfo["Chip"]; ok {
		socket = "Apple Silicon"
	}
	strTable = append(strTable, socket)
	data.WriteByte(uint8(len(strTable)))

	// Processor Type (offset 0x05)
	data.WriteByte(0x03) // Central Processor

	// Processor Family (offset 0x06)
	family := uint8(0xFE) // Use Family2
	data.WriteByte(family)

	// Processor Manufacturer (offset 0x07)
	manufacturer := "Apple Inc."
	if v, ok := cpuInfo["Manufacturer"]; ok && v != "" {
		manufacturer = v
	}
	strTable = append(strTable, manufacturer)
	data.WriteByte(uint8(len(strTable)))

	// Processor ID (offset 0x08) - 8 bytes
	binary.Write(&data, binary.LittleEndian, uint64(0))

	// Processor Version (offset 0x10)
	version := ""
	if v, ok := cpuInfo["Chip"]; ok {
		version = v
	} else if v, ok := cpuInfo["Processor Name"]; ok {
		version = v
	} else if v, ok := cpuInfo["cpu_brand"]; ok {
		version = v
	}
	strTable = append(strTable, version)
	data.WriteByte(uint8(len(strTable)))

	// Voltage (offset 0x11)
	data.WriteByte(0x80 | 12) // 1.2V

	// External Clock (offset 0x12)
	binary.Write(&data, binary.LittleEndian, uint16(0)) // Unknown

	// Max Speed (offset 0x14)
	maxSpeed := uint16(0)
	if v, ok := cpuInfo["max_freq"]; ok {
		if freq, err := strconv.ParseUint(v, 10, 16); err == nil {
			maxSpeed = uint16(freq)
		}
	}
	binary.Write(&data, binary.LittleEndian, maxSpeed)

	// Current Speed (offset 0x16)
	currentSpeed := maxSpeed
	if v, ok := cpuInfo["current_freq"]; ok {
		if freq, err := strconv.ParseUint(v, 10, 16); err == nil {
			currentSpeed = uint16(freq)
		}
	}
	binary.Write(&data, binary.LittleEndian, currentSpeed)

	// Status (offset 0x18)
	data.WriteByte(0x41) // CPU Enabled, Populated

	// Processor Upgrade (offset 0x19)
	data.WriteByte(0x06) // None (soldered)

	// L1 Cache Handle (offset 0x1A)
	binary.Write(&data, binary.LittleEndian, uint16(0x0700))

	// L2 Cache Handle (offset 0x1C)
	binary.Write(&data, binary.LittleEndian, uint16(0x0701))

	// L3 Cache Handle (offset 0x1E)
	binary.Write(&data, binary.LittleEndian, uint16(0xFFFF)) // Not provided

	// Serial Number (offset 0x20)
	strTable = append(strTable, "")
	data.WriteByte(uint8(len(strTable)))

	// Asset Tag (offset 0x21)
	strTable = append(strTable, "")
	data.WriteByte(uint8(len(strTable)))

	// Part Number (offset 0x22)
	strTable = append(strTable, "")
	data.WriteByte(uint8(len(strTable)))

	// Core Count (offset 0x23)
	coreCount := uint8(0)
	if v, ok := cpuInfo["cores"]; ok {
		if c, err := strconv.ParseUint(v, 10, 8); err == nil {
			coreCount = uint8(c)
		}
	} else if v, ok := cpuInfo["Total Number of Cores"]; ok {
		// Parse "X (Y performance and Z efficiency)" format
		parts := strings.Split(v, " ")
		if len(parts) > 0 {
			if c, err := strconv.ParseUint(parts[0], 10, 8); err == nil {
				coreCount = uint8(c)
			}
		}
	}
	data.WriteByte(coreCount)

	// Core Enabled (offset 0x24)
	data.WriteByte(coreCount)

	// Thread Count (offset 0x25)
	threadCount := coreCount // Default to core count
	if v, ok := cpuInfo["threads"]; ok {
		if t, err := strconv.ParseUint(v, 10, 8); err == nil {
			threadCount = uint8(t)
		}
	}
	data.WriteByte(threadCount)

	// Processor Characteristics (offset 0x26)
	characteristics := uint16(0x04 | 0x08 | 0x80) // 64-bit, Multi-Core, Power/Performance Control
	binary.Write(&data, binary.LittleEndian, characteristics)

	// Processor Family 2 (offset 0x28)
	family2 := uint16(0x0102) // ARMv8
	if _, ok := cpuInfo["Chip"]; ok {
		family2 = 0x0102 // ARMv8 for Apple Silicon
	}
	binary.Write(&data, binary.LittleEndian, family2)

	// Core Count 2 (offset 0x2A)
	binary.Write(&data, binary.LittleEndian, uint16(coreCount))

	// Core Enabled 2 (offset 0x2C)
	binary.Write(&data, binary.LittleEndian, uint16(coreCount))

	// Thread Count 2 (offset 0x2E)
	binary.Write(&data, binary.LittleEndian, uint16(threadCount))

	structures = append(structures, Structure{
		Header:  Header{Type: 4, Length: 48, Handle: 0x0004},
		Data:    data.Bytes(),
		Strings: strTable,
	})

	return structures
}

// getSysctlCPUInfo gets CPU information from sysctl
func getSysctlCPUInfo() map[string]string {
	info := make(map[string]string)

	// Get various CPU properties
	sysctlKeys := map[string]string{
		"machdep.cpu.brand_string": "cpu_brand",
		"machdep.cpu.core_count":   "cores",
		"machdep.cpu.thread_count": "threads",
		"hw.cpufrequency_max":      "max_freq_hz",
		"hw.ncpu":                  "ncpu",
		"hw.physicalcpu":           "physical_cpu",
		"hw.logicalcpu":            "logical_cpu",
	}

	for key, name := range sysctlKeys {
		cmd := exec.Command("sysctl", "-n", key)
		output, err := cmd.Output()
		if err == nil {
			value := strings.TrimSpace(string(output))
			if value != "" && value != "0" {
				info[name] = value
			}
		}
	}

	// Convert max frequency from Hz to MHz
	if hz, ok := info["max_freq_hz"]; ok {
		if freq, err := strconv.ParseUint(hz, 10, 64); err == nil {
			info["max_freq"] = strconv.FormatUint(freq/1000000, 10)
			info["current_freq"] = info["max_freq"]
		}
	}

	// Use physical/logical CPU if core/thread count not available
	if _, ok := info["cores"]; !ok {
		if v, ok := info["physical_cpu"]; ok {
			info["cores"] = v
		}
	}
	if _, ok := info["threads"]; !ok {
		if v, ok := info["logical_cpu"]; ok {
			info["threads"] = v
		}
	}

	return info
}

// createCacheStructures creates Type 7 - Cache Information structures
func createCacheStructures() []Structure {
	var structures []Structure

	// Get cache info from sysctl
	cacheInfo := getSysctlCacheInfo()

	// Create L1 Data Cache
	if size, ok := cacheInfo["l1d"]; ok {
		structures = append(structures, createCacheStructure(0x0700, "L1 Data Cache", 1, size, 0x04)) // Data
	}

	// Create L1 Instruction Cache
	if size, ok := cacheInfo["l1i"]; ok {
		structures = append(structures, createCacheStructure(0x0701, "L1 Instruction Cache", 1, size, 0x03)) // Instruction
	}

	// Create L2 Cache
	if size, ok := cacheInfo["l2"]; ok {
		structures = append(structures, createCacheStructure(0x0702, "L2 Cache", 2, size, 0x05)) // Unified
	}

	// Create L3 Cache (if available)
	if size, ok := cacheInfo["l3"]; ok {
		structures = append(structures, createCacheStructure(0x0703, "L3 Cache", 3, size, 0x05)) // Unified
	}

	return structures
}

// getSysctlCacheInfo gets cache information from sysctl
func getSysctlCacheInfo() map[string]uint32 {
	info := make(map[string]uint32)

	cacheKeys := map[string]string{
		"hw.l1dcachesize": "l1d",
		"hw.l1icachesize": "l1i",
		"hw.l2cachesize":  "l2",
		"hw.l3cachesize":  "l3",
	}

	for key, name := range cacheKeys {
		cmd := exec.Command("sysctl", "-n", key)
		output, err := cmd.Output()
		if err == nil {
			value := strings.TrimSpace(string(output))
			if size, err := strconv.ParseUint(value, 10, 32); err == nil && size > 0 {
				// Convert to KB
				info[name] = uint32(size / 1024)
			}
		}
	}

	return info
}

// createCacheStructure creates a single Type 7 - Cache Information structure
func createCacheStructure(handle uint16, designation string, level int, sizeKB uint32, cacheType uint8) Structure {
	var strTable []string
	var data bytes.Buffer

	// Write header
	data.WriteByte(7)  // Type 7 - Cache Information
	data.WriteByte(27) // Length (SMBIOS 3.1)
	binary.Write(&data, binary.LittleEndian, handle)

	// Socket Designation (offset 0x04)
	strTable = append(strTable, designation)
	data.WriteByte(uint8(len(strTable)))

	// Cache Configuration (offset 0x05)
	// Bits 0-2: Level (0-based), Bit 3: Socketed, Bits 5-6: Location, Bit 7: Enabled, Bits 8-9: Mode
	config := uint16(level-1) | 0x0080 | 0x0100 // Level + Enabled + Write-back
	binary.Write(&data, binary.LittleEndian, config)

	// Maximum Cache Size (offset 0x07)
	maxSize := uint16(sizeKB)
	if sizeKB > 0x7FFF {
		maxSize = 0x8000 | uint16(sizeKB/64) // Use 64K granularity
	}
	binary.Write(&data, binary.LittleEndian, maxSize)

	// Installed Size (offset 0x09)
	binary.Write(&data, binary.LittleEndian, maxSize)

	// Supported SRAM Type (offset 0x0B)
	binary.Write(&data, binary.LittleEndian, uint16(0x0020)) // Synchronous

	// Current SRAM Type (offset 0x0D)
	binary.Write(&data, binary.LittleEndian, uint16(0x0020)) // Synchronous

	// Cache Speed (offset 0x0F)
	data.WriteByte(0) // Unknown

	// Error Correction Type (offset 0x10)
	data.WriteByte(0x05) // Single-bit ECC

	// System Cache Type (offset 0x11)
	data.WriteByte(cacheType)

	// Associativity (offset 0x12)
	assoc := uint8(0x06) // Fully associative (default)
	if level == 2 {
		assoc = 0x08 // 16-way
	} else if level == 3 {
		assoc = 0x09 // 12-way
	}
	data.WriteByte(assoc)

	// Maximum Cache Size 2 (offset 0x13) - SMBIOS 3.1
	binary.Write(&data, binary.LittleEndian, sizeKB)

	// Installed Cache Size 2 (offset 0x17) - SMBIOS 3.1
	binary.Write(&data, binary.LittleEndian, sizeKB)

	return Structure{
		Header:  Header{Type: 7, Length: 27, Handle: handle},
		Data:    data.Bytes(),
		Strings: strTable,
	}
}

// createMemoryArrayStructure creates Type 16 - Physical Memory Array
func createMemoryArrayStructure(numDevices int) Structure {
	var data bytes.Buffer

	// Write header
	data.WriteByte(16)                                       // Type
	data.WriteByte(23)                                       // Length (SMBIOS 2.7+)
	binary.Write(&data, binary.LittleEndian, uint16(0x1000)) // Handle

	// Location (offset 0x04)
	data.WriteByte(0x03) // System board

	// Use (offset 0x05)
	data.WriteByte(0x03) // System memory

	// Memory Error Correction (offset 0x06)
	data.WriteByte(0x03) // None

	// Maximum Capacity (offset 0x07) - 4 bytes, in KB
	binary.Write(&data, binary.LittleEndian, uint32(0x80000000)) // Use extended

	// Memory Error Information Handle (offset 0x0B)
	binary.Write(&data, binary.LittleEndian, uint16(0xFFFE))

	// Number of Memory Devices (offset 0x0D)
	binary.Write(&data, binary.LittleEndian, uint16(numDevices))

	// Extended Maximum Capacity (offset 0x0F) - 8 bytes
	binary.Write(&data, binary.LittleEndian, uint64(256*1024*1024*1024)) // 256GB

	return Structure{
		Header:  Header{Type: 16, Length: 23, Handle: 0x1000},
		Data:    data.Bytes(),
		Strings: nil,
	}
}

// createMemoryDeviceStructure creates a Type 17 structure from memory slot data
func createMemoryDeviceStructure(slot map[string]string, handle uint16) Structure {
	var strTable []string
	var data bytes.Buffer

	// Write header
	data.WriteByte(17) // Type 17 - Memory Device
	data.WriteByte(40) // Length (SMBIOS 2.8)
	binary.Write(&data, binary.LittleEndian, handle)

	// Physical Memory Array Handle (offset 0x04)
	binary.Write(&data, binary.LittleEndian, uint16(0x1000))

	// Memory Error Information Handle (offset 0x06)
	binary.Write(&data, binary.LittleEndian, uint16(0xFFFE))

	// Total Width (offset 0x08)
	binary.Write(&data, binary.LittleEndian, uint16(64))

	// Data Width (offset 0x0A)
	binary.Write(&data, binary.LittleEndian, uint16(64))

	// Size (offset 0x0C) - parse from slot data
	var sizeMB uint16 = 0
	if sizeStr, ok := slot["Size"]; ok {
		sizeMB = parseMemorySize(sizeStr)
	}
	binary.Write(&data, binary.LittleEndian, sizeMB)

	// Form Factor (offset 0x0E)
	data.WriteByte(0x09) // DIMM

	// Device Set (offset 0x0F)
	data.WriteByte(0)

	// Device Locator (offset 0x10)
	locator := slot["slot"]
	if locator == "" {
		locator = "DIMM"
	}
	strTable = append(strTable, locator)
	data.WriteByte(uint8(len(strTable)))

	// Bank Locator (offset 0x11)
	strTable = append(strTable, "")
	data.WriteByte(uint8(len(strTable)))

	// Memory Type (offset 0x12)
	memType := uint8(0x1A) // DDR4 default
	if typeStr, ok := slot["Type"]; ok {
		memType = parseMemoryType(typeStr)
	}
	data.WriteByte(memType)

	// Type Detail (offset 0x13)
	binary.Write(&data, binary.LittleEndian, uint16(0x0080)) // Synchronous

	// Speed (offset 0x15)
	var speed uint16 = 0
	if speedStr, ok := slot["Speed"]; ok {
		speed = parseMemorySpeed(speedStr)
	}
	binary.Write(&data, binary.LittleEndian, speed)

	// Manufacturer (offset 0x17)
	manufacturer := ""
	if v, ok := slot["Manufacturer"]; ok {
		manufacturer = v
	}
	strTable = append(strTable, manufacturer)
	data.WriteByte(uint8(len(strTable)))

	// Serial Number (offset 0x18)
	serial := ""
	if v, ok := slot["Serial Number"]; ok {
		serial = v
	}
	strTable = append(strTable, serial)
	data.WriteByte(uint8(len(strTable)))

	// Asset Tag (offset 0x19)
	strTable = append(strTable, "")
	data.WriteByte(uint8(len(strTable)))

	// Part Number (offset 0x1A)
	partNumber := ""
	if v, ok := slot["Part Number"]; ok {
		partNumber = v
	}
	strTable = append(strTable, partNumber)
	data.WriteByte(uint8(len(strTable)))

	// Attributes (offset 0x1B)
	data.WriteByte(0)

	// Extended Size (offset 0x1C)
	binary.Write(&data, binary.LittleEndian, uint32(0))

	// Configured Memory Speed (offset 0x20)
	binary.Write(&data, binary.LittleEndian, speed)

	// Minimum Voltage (offset 0x22)
	binary.Write(&data, binary.LittleEndian, uint16(1200))

	// Maximum Voltage (offset 0x24)
	binary.Write(&data, binary.LittleEndian, uint16(1200))

	// Configured Voltage (offset 0x26)
	binary.Write(&data, binary.LittleEndian, uint16(1200))

	return Structure{
		Header:  Header{Type: 17, Length: 40, Handle: handle},
		Data:    data.Bytes(),
		Strings: strTable,
	}
}

// createSyntheticMemoryDevice creates a memory device from Apple Silicon info
func createSyntheticMemoryDevice(memStr string) Structure {
	slot := make(map[string]string)
	slot["slot"] = "Unified Memory"
	slot["Size"] = memStr
	slot["Type"] = "LPDDR5" // Apple Silicon typically uses LPDDR
	return createMemoryDeviceStructure(slot, 0x1100)
}

// parseUUID parses a UUID string into 16 bytes
func parseUUID(uuid string) []byte {
	result := make([]byte, 16)
	uuid = strings.ReplaceAll(uuid, "-", "")

	if len(uuid) != 32 {
		return result
	}

	for i := 0; i < 16; i++ {
		b, _ := strconv.ParseUint(uuid[i*2:i*2+2], 16, 8)
		result[i] = byte(b)
	}

	return result
}

// parseMemorySize parses memory size string (e.g., "8 GB") to MB
func parseMemorySize(size string) uint16 {
	size = strings.ToUpper(size)
	size = strings.ReplaceAll(size, " ", "")

	var multiplier uint64 = 1
	if strings.HasSuffix(size, "TB") {
		multiplier = 1024 * 1024 // Result in MB
		size = strings.TrimSuffix(size, "TB")
	} else if strings.HasSuffix(size, "GB") {
		multiplier = 1024 // Result in MB
		size = strings.TrimSuffix(size, "GB")
	} else if strings.HasSuffix(size, "MB") {
		multiplier = 1
		size = strings.TrimSuffix(size, "MB")
	}

	val, err := strconv.ParseUint(size, 10, 64)
	if err != nil {
		return 0
	}

	result := val * multiplier
	if result > 0x7FFF {
		return 0x7FFF // Max value for 16-bit size field
	}
	return uint16(result)
}

// parseMemorySpeed parses memory speed string (e.g., "2400 MHz") to MT/s
func parseMemorySpeed(speed string) uint16 {
	speed = strings.ToUpper(speed)
	speed = strings.ReplaceAll(speed, " ", "")
	speed = strings.TrimSuffix(speed, "MHZ")
	speed = strings.TrimSuffix(speed, "MT/S")

	val, err := strconv.ParseUint(speed, 10, 16)
	if err != nil {
		return 0
	}
	return uint16(val)
}

// parseMemoryType parses memory type string to SMBIOS memory type code
func parseMemoryType(typeStr string) uint8 {
	typeStr = strings.ToUpper(typeStr)
	if strings.Contains(typeStr, "DDR5") {
		return 0x22
	} else if strings.Contains(typeStr, "LPDDR5") {
		return 0x23
	} else if strings.Contains(typeStr, "DDR4") {
		return 0x1A
	} else if strings.Contains(typeStr, "LPDDR4") {
		return 0x1E
	} else if strings.Contains(typeStr, "DDR3") {
		return 0x18
	} else if strings.Contains(typeStr, "LPDDR3") {
		return 0x1D
	}
	return 0x1A // DDR4 default
}
