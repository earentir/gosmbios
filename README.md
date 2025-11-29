# gosmbios

Pure Go implementation for reading SMBIOS/DMI data across Windows, Linux, and macOS on both AMD64 and ARM64 architectures.

Implements **DSP0134 SMBIOS Reference Specification Version 3.9.0** with **100% type coverage**.

## Features

- **Pure Go** - No external dependencies or tools required
- **Cross-platform** - Supports Windows, Linux, and macOS
- **Multi-architecture** - Works on AMD64 and ARM64
- **Complete Coverage** - Implements all 47 SMBIOS structure types (including obsolete)
- **Modular** - Each SMBIOS type is in its own package
- **Type-safe** - Full Go type definitions with constants and methods

## Supported SMBIOS Types

### Core System Information
| Type | Package | Description |
|------|---------|-------------|
| 0 | `types/type0` | BIOS Information |
| 1 | `types/type1` | System Information |
| 2 | `types/type2` | Baseboard (Motherboard) Information |
| 3 | `types/type3` | System Enclosure/Chassis |

### Processor & Cache
| Type | Package | Description |
|------|---------|-------------|
| 4 | `types/type4` | Processor Information |
| 7 | `types/type7` | Cache Information |
| 44 | `types/type44` | Processor Additional Information |

### Memory
| Type | Package | Description |
|------|---------|-------------|
| 5 | `types/type5` | Memory Controller Information *(obsolete)* |
| 6 | `types/type6` | Memory Module Information *(obsolete)* |
| 16 | `types/type16` | Physical Memory Array |
| 17 | `types/type17` | Memory Device |
| 18 | `types/type18` | 32-Bit Memory Error Information |
| 19 | `types/type19` | Memory Array Mapped Address |
| 20 | `types/type20` | Memory Device Mapped Address |
| 33 | `types/type33` | 64-Bit Memory Error Information |
| 37 | `types/type37` | Memory Channel |

### Ports & Slots
| Type | Package | Description |
|------|---------|-------------|
| 8 | `types/type8` | Port Connector Information |
| 9 | `types/type9` | System Slots |

### On-Board Devices
| Type | Package | Description |
|------|---------|-------------|
| 10 | `types/type10` | On Board Devices Information *(obsolete)* |
| 41 | `types/type41` | Onboard Devices Extended Information |

### System Configuration
| Type | Package | Description |
|------|---------|-------------|
| 11 | `types/type11` | OEM Strings |
| 12 | `types/type12` | System Configuration Options |
| 13 | `types/type13` | BIOS Language Information |
| 14 | `types/type14` | Group Associations |
| 15 | `types/type15` | System Event Log |

### Input Devices
| Type | Package | Description |
|------|---------|-------------|
| 21 | `types/type21` | Built-in Pointing Device |

### Power & Battery
| Type | Package | Description |
|------|---------|-------------|
| 22 | `types/type22` | Portable Battery |
| 23 | `types/type23` | System Reset |
| 25 | `types/type25` | System Power Controls |
| 39 | `types/type39` | System Power Supply |

### Security
| Type | Package | Description |
|------|---------|-------------|
| 24 | `types/type24` | Hardware Security |
| 43 | `types/type43` | TPM Device |

### Probes & Thermal
| Type | Package | Description |
|------|---------|-------------|
| 26 | `types/type26` | Voltage Probe |
| 27 | `types/type27` | Cooling Device |
| 28 | `types/type28` | Temperature Probe |
| 29 | `types/type29` | Electrical Current Probe |

### Management
| Type | Package | Description |
|------|---------|-------------|
| 30 | `types/type30` | Out-of-Band Remote Access |
| 34 | `types/type34` | Management Device |
| 35 | `types/type35` | Management Device Component |
| 36 | `types/type36` | Management Device Threshold Data |
| 38 | `types/type38` | IPMI Device Information |
| 42 | `types/type42` | Management Controller Host Interface |

### Boot & Firmware
| Type | Package | Description |
|------|---------|-------------|
| 31 | `types/type31` | Boot Integrity Services Entry Point |
| 32 | `types/type32` | System Boot Information |
| 45 | `types/type45` | Firmware Inventory Information |

### Miscellaneous
| Type | Package | Description |
|------|---------|-------------|
| 40 | `types/type40` | Additional Information |
| 46 | `types/type46` | String Property |
| 127 | `types/type127` | End-of-Table |

## Installation

```bash
go get github.com/earentir/gosmbios
```

## Command Line Tools

The package includes several command-line tools in the `cmd/` directory:

### smbiosinfo (`cmd/info`)
Displays comprehensive SMBIOS information in human-readable format.

```bash
go run ./cmd/info
```

### smbiosdebug (`cmd/debug`)
Debug tool showing raw structure information with hex dumps.

```bash
go run ./cmd/debug
```

### smbiosdump (`cmd/dump`)
Dumps SMBIOS data to file in various formats (text, JSON, raw hex).

```bash
# Text format (default)
go run ./cmd/dump -o smbios.txt

# JSON format
go run ./cmd/dump -o smbios.json -f json

# Raw hex format
go run ./cmd/dump -f raw > smbios.hex
```

### examples (`cmd/examples`)
Basic example demonstrating library usage.

```bash
go run ./cmd/examples
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/earentir/gosmbios"
    "github.com/earentir/gosmbios/types/type0"
    "github.com/earentir/gosmbios/types/type1"
    "github.com/earentir/gosmbios/types/type17"
)

func main() {
    // Read SMBIOS data
    sm, err := gosmbios.Read()
    if err != nil {
        panic(err)
    }

    fmt.Printf("SMBIOS Version: %s\n", sm.EntryPoint.String())

    // Get BIOS Information
    bios, _ := type0.Get(sm)
    fmt.Printf("BIOS: %s %s\n", bios.Vendor, bios.Version)

    // Get System Information
    sys, _ := type1.Get(sm)
    fmt.Printf("System: %s %s\n", sys.Manufacturer, sys.ProductName)
    fmt.Printf("UUID: %s\n", sys.UUID.String())

    // Get Memory Information
    memory, _ := type17.GetPopulated(sm)
    for _, m := range memory {
        fmt.Printf("Memory: %s %s %s\n",
            m.DeviceLocator,
            m.SizeString(),
            m.MemoryType.String())
    }
}
```

## Platform-Specific Behavior

### Linux
Reads from `/sys/firmware/dmi/tables/smbios_entry_point` and `/sys/firmware/dmi/tables/DMI`. May require root privileges on some systems.

### Windows
Uses `GetSystemFirmwareTable` API to read SMBIOS data. Works without administrative privileges.

### macOS
Uses `ioreg` and `system_profiler` to gather system information. Some data may be synthesized as macOS doesn't expose raw SMBIOS tables directly.

## Advanced Usage

### Working with Raw Structures

```go
sm, _ := gosmbios.Read()

// Get all structures of a specific type
for _, s := range sm.GetStructures(17) { // Type 17 = Memory Device
    fmt.Printf("Handle: 0x%04X\n", s.Header.Handle)
    fmt.Printf("Type: %d\n", s.Header.Type)
    fmt.Printf("Length: %d\n", s.Header.Length)

    // Access raw data
    fmt.Printf("Raw byte at 0x04: 0x%02X\n", s.GetByte(0x04))
    fmt.Printf("Raw word at 0x08: 0x%04X\n", s.GetWord(0x08))

    // Access strings
    for i, str := range s.Strings {
        fmt.Printf("String %d: %s\n", i+1, str)
    }
}
```

### Getting All Processors

```go
import "github.com/earentir/gosmbios/types/type4"

processors, _ := type4.GetAll(sm)
for _, proc := range processors {
    if proc.Status.IsPopulated() {
        fmt.Printf("%s: %d cores, %d threads @ %d MHz\n",
            proc.DisplayName(),
            proc.GetCoreCount(),
            proc.GetThreadCount(),
            proc.CurrentSpeed)
    }
}
```

### Checking TPM Status

```go
import "github.com/earentir/gosmbios/types/type43"

tpm, err := type43.Get(sm)
if err == nil && tpm.IsSupported() {
    fmt.Printf("TPM: %s (Spec %s)\n", tpm.Family(), tpm.SpecVersionString())
    fmt.Printf("TPM Vendor: %s\n", tpm.VendorIDString())
}
```

### Getting Battery Information (Laptops)

```go
import "github.com/earentir/gosmbios/types/type22"

batteries, _ := type22.GetAll(sm)
for _, bat := range batteries {
    fmt.Printf("Battery: %s\n", bat.DeviceName)
    fmt.Printf("  Chemistry: %s\n", bat.DeviceChemistry.String())
    fmt.Printf("  Capacity: %s\n", bat.DesignCapacityString())
    fmt.Printf("  Voltage: %s\n", bat.DesignVoltageString())
}
```

### Reading Temperature Probes

```go
import "github.com/earentir/gosmbios/types/type28"

probes, _ := type28.GetAll(sm)
for _, probe := range probes {
    fmt.Printf("Probe: %s\n", probe.Description)
    fmt.Printf("  Location: %s\n", probe.LocationAndStatus.Location().String())
    fmt.Printf("  Status: %s\n", probe.LocationAndStatus.Status().String())
    fmt.Printf("  Range: %s to %s\n", probe.MinimumValueString(), probe.MaximumValueString())
}
```

### Checking if Virtual Machine

```go
import "github.com/earentir/gosmbios/types/type0"

bios, _ := type0.Get(sm)
if bios.IsVirtualMachine() {
    fmt.Println("Running in a virtual machine")
}
```

### Getting Total Memory

```go
import "github.com/earentir/gosmbios/types/type17"

devices, _ := type17.GetPopulated(sm)
var totalMB uint64
for _, dev := range devices {
    totalMB += dev.Size
}
fmt.Printf("Total Memory: %d GB\n", totalMB/1024)
```

## API Reference

### Main Package

| Function | Description |
|----------|-------------|
| `Read()` | Reads and parses SMBIOS data from the system |
| `GetStructure(type)` | Returns first structure of given type |
| `GetStructures(type)` | Returns all structures of given type |

### Structure Methods

| Method | Description |
|--------|-------------|
| `GetString(index)` | Get string by 1-based index |
| `GetByte(offset)` | Get byte at offset |
| `GetWord(offset)` | Get 16-bit value at offset |
| `GetDWord(offset)` | Get 32-bit value at offset |
| `GetQWord(offset)` | Get 64-bit value at offset |

### Type Constants

The `types` package provides constants for all SMBIOS structure types:

```go
import "github.com/earentir/gosmbios/types"

types.BIOSInformation          // 0
types.SystemInformation        // 1
types.BaseboardInformation     // 2
types.ProcessorInformation     // 4
types.MemoryDevice             // 17
types.TPMDevice                // 43
// ... and all others
```

## Building

```bash
# Build all packages and tools
go build ./...

# Build specific tool
go build -o smbiosinfo ./cmd/info
go build -o smbiosdebug ./cmd/debug
go build -o smbiosdump ./cmd/dump

# Cross-compile
GOOS=linux GOARCH=amd64 go build -o smbiosinfo-linux-amd64 ./cmd/info
GOOS=linux GOARCH=arm64 go build -o smbiosinfo-linux-arm64 ./cmd/info
GOOS=windows GOARCH=amd64 go build -o smbiosinfo.exe ./cmd/info
GOOS=darwin GOARCH=amd64 go build -o smbiosinfo-darwin-amd64 ./cmd/info
GOOS=darwin GOARCH=arm64 go build -o smbiosinfo-darwin-arm64 ./cmd/info
```

## Testing

```bash
go test ./...
```

## Requirements

- Go 1.21 or later
- Linux: Read access to `/sys/firmware/dmi/tables/` (may need root)
- Windows: No special requirements
- macOS: Command-line tools available (`ioreg`, `system_profiler`)

## License
GNU General Public License v2.0

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## References

- [DMTF DSP0134 - SMBIOS Reference Specification](https://www.dmtf.org/standards/smbios)
- [SMBIOS on Wikipedia](https://en.wikipedia.org/wiki/System_Management_BIOS)
