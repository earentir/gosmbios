# gosmbios

Pure Go implementation for reading SMBIOS/DMI data across Windows, Linux, and macOS on both AMD64 and ARM64 architectures.

Implements **DSP0134 SMBIOS Reference Specification Version 3.9.0**.

## Features

- **Pure Go** - No external dependencies or tools required
- **Cross-platform** - Supports Windows, Linux, and macOS
- **Multi-architecture** - Works on AMD64 and ARM64
- **Comprehensive** - Implements major SMBIOS structure types
- **Modular** - Each SMBIOS type is in its own package
- **Type-safe** - Full Go type definitions with constants and methods

## Supported SMBIOS Types

| Type | Package | Description |
|------|---------|-------------|
| 0 | `types/type0` | BIOS Information |
| 1 | `types/type1` | System Information |
| 2 | `types/type2` | Baseboard (Motherboard) Information |
| 3 | `types/type3` | System Enclosure/Chassis |
| 4 | `types/type4` | Processor Information |
| 7 | `types/type7` | Cache Information |
| 8 | `types/type8` | Port Connector Information |
| 9 | `types/type9` | System Slots |
| 11 | `types/type11` | OEM Strings |
| 16 | `types/type16` | Physical Memory Array |
| 17 | `types/type17` | Memory Device |
| 32 | `types/type32` | System Boot Information |

## Installation

```bash
go get github.com/earentir/gosmbios
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

### Checking System Type

```go
import "github.com/earentir/gosmbios/types/type3"

chassis, _ := type3.Get(sm)
if chassis.Type.IsPortable() {
    fmt.Println("This is a portable device (laptop/tablet)")
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

## Building

```bash
# Build for current platform
go build ./...

# Cross-compile examples
GOOS=linux GOARCH=amd64 go build -o example-linux-amd64 ./examples/
GOOS=linux GOARCH=arm64 go build -o example-linux-arm64 ./examples/
GOOS=windows GOARCH=amd64 go build -o example-windows.exe ./examples/
GOOS=darwin GOARCH=amd64 go build -o example-darwin-amd64 ./examples/
GOOS=darwin GOARCH=arm64 go build -o example-darwin-arm64 ./examples/
```

## Testing

```bash
go test ./...
```

## Requirements

- Go 1.21 or later
- Linux: Read access to `/sys/firmware/dmi/tables/` (needs root)
- Windows: No special requirements
- macOS: Command-line tools available (`ioreg`, `system_profiler`)

## License
GNU General Public License v2.0

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## References

- [DMTF DSP0134 - SMBIOS Reference Specification](https://www.dmtf.org/standards/smbios)
- [SMBIOS on Wikipedia](https://en.wikipedia.org/wiki/System_Management_BIOS)
