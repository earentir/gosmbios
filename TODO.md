# gosmbios TODO

## macOS Implementation - External Tools Issue

The current macOS implementation uses external command-line tools (`ioreg`, `system_profiler`, `sysctl`) via `exec.Command`, which violates the "pure Go, no external tools" requirement.

### Current Status

| Platform | Method | Pure Go? |
|----------|--------|----------|
| Linux | `/sys/firmware/dmi/tables/` file I/O | ✅ Yes |
| Windows | `GetSystemFirmwareTable` syscall | ✅ Yes |
| macOS | `ioreg`, `system_profiler`, `sysctl` commands | ❌ No |

### Options

#### Option 1: Accept the Limitation
Keep the current implementation where macOS uses command-line tools as a workaround.

**Pros:**
- Full functionality on macOS
- Works on both Intel and Apple Silicon Macs
- Already implemented and tested

**Cons:**
- Violates "no external tools" requirement
- Depends on external binaries being available
- Slower due to process spawning

---

#### Option 2: Remove macOS Support
Only support Linux and Windows with pure Go. macOS returns `ErrUnsupportedOS`.

**Pros:**
- Strict adherence to "pure Go, no external tools"
- Simpler codebase
- No external dependencies

**Cons:**
- No macOS support at all
- Limits library usefulness for macOS developers

---

#### Option 3: Pure Go IOKit Implementation
Implement direct IOKit framework access via raw Mach syscalls without CGO.

**Pros:**
- True pure Go implementation
- No external tool dependencies
- Potentially faster than exec.Command

**Cons:**
- Extremely complex to implement
- Requires reimplementing parts of IOKit in pure Go
- Mach syscalls are poorly documented
- May break with macOS updates
- Significant development effort (weeks/months)
- Apple Silicon vs Intel differences to handle

**Technical Requirements:**
- Mach port manipulation via syscalls
- IOKit service matching and property access
- Parsing IOKit registry data structures
- Handle both Intel SMBIOS and Apple Silicon equivalents

---

## Decision

- [ ] Option 1: Accept limitation (external tools on macOS)
- [ ] Option 2: Remove macOS support
- [ ] Option 3: Pure Go IOKit (major undertaking)

## Completed Tasks

- [x] **All 47 SMBIOS types implemented** (Types 0-46, 126-127)
  - Core types: 0, 1, 2, 3, 4, 7, 8, 9
  - Memory types: 5, 6, 16, 17, 18, 19, 20, 33, 37
  - Device types: 10, 21, 22, 41
  - Management types: 34, 35, 36, 38, 42
  - Power/thermal types: 23, 24, 25, 26, 27, 28, 29, 39
  - Boot/firmware types: 31, 32, 43, 44, 45
  - Other types: 11, 12, 13, 14, 15, 30, 40, 46, 127

- [x] **CLI tools implemented**
  - `cmd/info` - Human-readable SMBIOS information display
  - `cmd/debug` - Raw debug output with hex dumps
  - `cmd/dump` - Export to text, JSON, or raw hex formats
  - `cmd/examples` - Basic usage examples

## Future Enhancements

- [ ] Add unit tests for all type parsers
- [ ] Add benchmarks
- [ ] Add JSON output to info tool
- [ ] Support for reading from file (for offline analysis)
- [ ] Support for writing SMBIOS data (for testing/simulation)
