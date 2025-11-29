package main

import (
	"fmt"
	"os"

	"github.com/earentir/gosmbios"
	"github.com/earentir/gosmbios/types/type4"
	"github.com/earentir/gosmbios/types/type7"
	"github.com/earentir/gosmbios/types/type17"
)

func main() {
	sm, err := gosmbios.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading SMBIOS: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("SMBIOS Version: %s\n", sm.EntryPoint.String())
	fmt.Printf("Total structures: %d\n\n", len(sm.Structures))

	// List all structure types present
	typeCounts := make(map[uint8]int)
	for _, s := range sm.Structures {
		typeCounts[s.Header.Type]++
	}

	fmt.Println("=== Structure Types Present ===")
	for t := uint8(0); t < 128; t++ {
		if count, ok := typeCounts[t]; ok {
			fmt.Printf("  Type %3d: %d structure(s)\n", t, count)
		}
	}
	for t := uint8(128); t > 0; t++ {
		if count, ok := typeCounts[t]; ok {
			fmt.Printf("  Type %3d: %d structure(s) (OEM)\n", t, count)
		}
		if t == 255 {
			break
		}
	}

	// Debug Type 4 (Processor)
	fmt.Println("\n=== Type 4 (Processor) Debug ===")
	type4Structs := sm.GetStructures(4)
	fmt.Printf("Found %d Type 4 structures\n", len(type4Structs))
	for i, s := range type4Structs {
		fmt.Printf("  [%d] Handle=0x%04X, Length=%d, DataLen=%d, Strings=%d\n",
			i, s.Header.Handle, s.Header.Length, len(s.Data), len(s.Strings))

		proc, err := type4.Parse(&s)
		if err != nil {
			fmt.Printf("      Parse error: %v\n", err)
		} else {
			fmt.Printf("      Parsed OK: %s by %s, %d cores, %d threads\n",
				proc.ProcessorVersion, proc.ProcessorManufacturer,
				proc.GetCoreCount(), proc.GetThreadCount())
		}
	}

	// Debug Type 7 (Cache)
	fmt.Println("\n=== Type 7 (Cache) Debug ===")
	type7Structs := sm.GetStructures(7)
	fmt.Printf("Found %d Type 7 structures\n", len(type7Structs))
	for i, s := range type7Structs {
		fmt.Printf("  [%d] Handle=0x%04X, Length=%d, DataLen=%d, Strings=%d\n",
			i, s.Header.Handle, s.Header.Length, len(s.Data), len(s.Strings))

		cache, err := type7.Parse(&s)
		if err != nil {
			fmt.Printf("      Parse error: %v\n", err)
		} else {
			fmt.Printf("      Parsed OK: %s, L%d, Size=%d KB\n",
				cache.SocketDesignation, cache.Level(), cache.InstalledSize)
		}
	}

	// Debug Type 17 (Memory Device)
	fmt.Println("\n=== Type 17 (Memory Device) Debug ===")
	type17Structs := sm.GetStructures(17)
	fmt.Printf("Found %d Type 17 structures\n", len(type17Structs))
	for i, s := range type17Structs {
		fmt.Printf("  [%d] Handle=0x%04X, Length=%d, DataLen=%d, Strings=%d\n",
			i, s.Header.Handle, s.Header.Length, len(s.Data), len(s.Strings))

		mem, err := type17.Parse(&s)
		if err != nil {
			fmt.Printf("      Parse error: %v\n", err)
		} else {
			fmt.Printf("      Parsed OK: %s, Size=%d MB, Type=%s\n",
				mem.DeviceLocator, mem.Size, mem.MemoryType.String())
		}
	}
}
