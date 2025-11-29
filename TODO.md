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

## Other Future Enhancements

- [ ] Add more SMBIOS types (Type 10, 13, 14, 15, 18-31, 33-46)
- [ ] Add unit tests
- [ ] Add benchmarks
- [ ] Consider adding a CLI tool for dumping SMBIOS data
