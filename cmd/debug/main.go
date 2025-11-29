// Debug tool for SMBIOS data - displays raw structure information for all types
package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/earentir/gosmbios"
	"github.com/earentir/gosmbios/types"
	"github.com/earentir/gosmbios/types/type0"
	"github.com/earentir/gosmbios/types/type1"
	"github.com/earentir/gosmbios/types/type11"
	"github.com/earentir/gosmbios/types/type16"
	"github.com/earentir/gosmbios/types/type17"
	"github.com/earentir/gosmbios/types/type2"
	"github.com/earentir/gosmbios/types/type3"
	"github.com/earentir/gosmbios/types/type32"
	"github.com/earentir/gosmbios/types/type4"
	"github.com/earentir/gosmbios/types/type7"
	"github.com/earentir/gosmbios/types/type8"
	"github.com/earentir/gosmbios/types/type9"
)

func main() {
	sm, err := gosmbios.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading SMBIOS: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("================================================================================")
	fmt.Println("                           SMBIOS DEBUG INFORMATION")
	fmt.Println("================================================================================")
	fmt.Printf("\nSMBIOS Version: %s\n", sm.EntryPoint.String())
	fmt.Printf("Entry Point Type: %s\n", entryPointTypeString(sm.EntryPoint.Type))
	fmt.Printf("Table Address: 0x%016X\n", sm.EntryPoint.TableAddress)
	fmt.Printf("Table Length: %d bytes\n", sm.EntryPoint.TableLength)
	fmt.Printf("Total Structures: %d\n\n", len(sm.Structures))

	// List all structure types present
	typeCounts := make(map[uint8]int)
	for _, s := range sm.Structures {
		typeCounts[s.Header.Type]++
	}

	fmt.Println("================================================================================")
	fmt.Println("                            STRUCTURE TYPE SUMMARY")
	fmt.Println("================================================================================")
	for t := uint8(0); t <= 127; t++ {
		if count, ok := typeCounts[t]; ok {
			fmt.Printf("  Type %3d: %2d structure(s) - %s\n", t, count, types.TypeName(t))
		}
	}
	for t := uint8(128); t > 0; t++ {
		if count, ok := typeCounts[t]; ok {
			fmt.Printf("  Type %3d: %2d structure(s) - %s (OEM)\n", t, count, types.TypeName(t))
		}
		if t == 255 {
			break
		}
	}

	// Debug each structure type
	fmt.Println("\n================================================================================")
	fmt.Println("                          DETAILED STRUCTURE DEBUG")
	fmt.Println("================================================================================")

	debugType0(sm)
	debugType1(sm)
	debugType2(sm)
	debugType3(sm)
	debugType4(sm)
	debugType7(sm)
	debugType8(sm)
	debugType9(sm)
	debugType11(sm)
	debugType16(sm)
	debugType17(sm)
	debugType32(sm)

	// Debug any remaining/unknown types with raw data
	debugRemainingTypes(sm, typeCounts)
}

func entryPointTypeString(t gosmbios.EntryPointType) string {
	if t == gosmbios.EntryPoint64Bit {
		return "64-bit (SMBIOS 3.x)"
	}
	return "32-bit (SMBIOS 2.x)"
}

func printStructureHeader(s *gosmbios.Structure) {
	fmt.Printf("  Handle: 0x%04X | Type: %d | Length: %d | Data Size: %d | Strings: %d\n",
		s.Header.Handle, s.Header.Type, s.Header.Length, len(s.Data), len(s.Strings))
}

func printHexDump(data []byte, indent string) {
	if len(data) == 0 {
		return
	}
	fmt.Printf("%sRaw Data:\n", indent)
	for i := 0; i < len(data); i += 16 {
		end := i + 16
		if end > len(data) {
			end = len(data)
		}
		fmt.Printf("%s  %04X: %s\n", indent, i, hex.EncodeToString(data[i:end]))
	}
}

func printStrings(strings []string, indent string) {
	if len(strings) == 0 {
		return
	}
	fmt.Printf("%sStrings:\n", indent)
	for i, s := range strings {
		fmt.Printf("%s  [%d]: %q\n", indent, i+1, s)
	}
}

func debugType0(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(0)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 0: BIOS Information ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		bios, err := type0.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Vendor:          %q\n", bios.Vendor)
			fmt.Printf("  Version:         %q\n", bios.Version)
			fmt.Printf("  Release Date:    %q\n", bios.ReleaseDate)
			fmt.Printf("  ROM Size:        %s (%d bytes)\n", bios.ROMSizeString(), bios.ROMSizeBytes)
			fmt.Printf("  BIOS Release:    %d.%d\n", bios.SystemBIOSMajorRelease, bios.SystemBIOSMinorRelease)
			fmt.Printf("  EC Release:      %s\n", bios.ECVersionString())
			fmt.Printf("  Characteristics: 0x%016X\n", uint64(bios.Characteristics))
			fmt.Printf("  CharExt1:        0x%02X\n", uint8(bios.CharacteristicsExt1))
			fmt.Printf("  CharExt2:        0x%02X\n", uint8(bios.CharacteristicsExt2))
			fmt.Printf("  UEFI:            %v\n", bios.IsUEFI())
			fmt.Printf("  Virtual Machine: %v\n", bios.IsVirtualMachine())
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType1(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(1)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 1: System Information ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		sys, err := type1.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Manufacturer:    %q\n", sys.Manufacturer)
			fmt.Printf("  Product Name:    %q\n", sys.ProductName)
			fmt.Printf("  Version:         %q\n", sys.Version)
			fmt.Printf("  Serial Number:   %q\n", sys.SerialNumber)
			fmt.Printf("  UUID:            %s\n", sys.UUID.String())
			fmt.Printf("  UUID (hex):      %s\n", sys.UUID.Hex())
			fmt.Printf("  UUID Zero:       %v\n", sys.UUID.IsZero())
			fmt.Printf("  UUID Invalid:    %v\n", sys.UUID.IsInvalid())
			fmt.Printf("  Wake-up Type:    %s (0x%02X)\n", sys.WakeUpType.String(), uint8(sys.WakeUpType))
			fmt.Printf("  SKU Number:      %q\n", sys.SKUNumber)
			fmt.Printf("  Family:          %q\n", sys.Family)
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType2(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(2)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 2: Baseboard Information ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		board, err := type2.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Manufacturer:    %q\n", board.Manufacturer)
			fmt.Printf("  Product:         %q\n", board.Product)
			fmt.Printf("  Version:         %q\n", board.Version)
			fmt.Printf("  Serial Number:   %q\n", board.SerialNumber)
			fmt.Printf("  Asset Tag:       %q\n", board.AssetTag)
			fmt.Printf("  Location:        %q\n", board.LocationInChassis)
			fmt.Printf("  Board Type:      %s (0x%02X)\n", board.BoardType.String(), uint8(board.BoardType))
			fmt.Printf("  Features:        0x%02X\n", uint8(board.FeatureFlags))
			fmt.Printf("  Chassis Handle:  0x%04X\n", board.ChassisHandle)
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType3(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(3)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 3: Chassis Information ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		chassis, err := type3.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Manufacturer:    %q\n", chassis.Manufacturer)
			fmt.Printf("  Type:            %s (0x%02X)\n", chassis.Type.String(), uint8(chassis.Type))
			fmt.Printf("  Version:         %q\n", chassis.Version)
			fmt.Printf("  Serial Number:   %q\n", chassis.SerialNumber)
			fmt.Printf("  Asset Tag:       %q\n", chassis.AssetTag)
			fmt.Printf("  Boot-up State:   %s (0x%02X)\n", chassis.BootUpState.String(), uint8(chassis.BootUpState))
			fmt.Printf("  Power State:     %s (0x%02X)\n", chassis.PowerSupplyState.String(), uint8(chassis.PowerSupplyState))
			fmt.Printf("  Thermal State:   %s (0x%02X)\n", chassis.ThermalState.String(), uint8(chassis.ThermalState))
			fmt.Printf("  Security:        %s (0x%02X)\n", chassis.SecurityStatus.String(), uint8(chassis.SecurityStatus))
			fmt.Printf("  Height:          %s\n", chassis.HeightString())
			fmt.Printf("  Power Cords:     %d\n", chassis.NumberOfPowerCords)
			fmt.Printf("  Portable:        %v\n", chassis.Type.IsPortable())
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType4(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(4)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 4: Processor Information ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		proc, err := type4.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Socket:          %q\n", proc.SocketDesignation)
			fmt.Printf("  Type:            %s (0x%02X)\n", proc.ProcessorType.String(), uint8(proc.ProcessorType))
			fmt.Printf("  Manufacturer:    %q\n", proc.ProcessorManufacturer)
			fmt.Printf("  Version:         %q\n", proc.ProcessorVersion)
			fmt.Printf("  Family:          %s (0x%04X)\n", proc.ProcessorFamily.String(), uint16(proc.ProcessorFamily))
			fmt.Printf("  ID:              0x%016X\n", proc.ProcessorID)
			fmt.Printf("  Voltage:         %s\n", proc.Voltage.String())
			fmt.Printf("  External Clock:  %d MHz\n", proc.ExternalClock)
			fmt.Printf("  Max Speed:       %d MHz\n", proc.MaxSpeed)
			fmt.Printf("  Current Speed:   %d MHz\n", proc.CurrentSpeed)
			fmt.Printf("  Status:          %s (0x%02X)\n", proc.Status.String(), uint8(proc.Status))
			fmt.Printf("  Populated:       %v\n", proc.Status.IsPopulated())
			fmt.Printf("  Upgrade:         %s (0x%02X)\n", proc.ProcessorUpgrade.String(), uint8(proc.ProcessorUpgrade))
			fmt.Printf("  Core Count:      %d\n", proc.GetCoreCount())
			fmt.Printf("  Core Enabled:    %d\n", proc.GetCoreEnabled())
			fmt.Printf("  Thread Count:    %d\n", proc.GetThreadCount())
			fmt.Printf("  Characteristics: 0x%04X\n", uint16(proc.ProcessorCharacteristics))
			fmt.Printf("  64-bit:          %v\n", proc.ProcessorCharacteristics.Is64Bit())
			fmt.Printf("  Multi-Core:      %v\n", proc.ProcessorCharacteristics.IsMultiCore())
			fmt.Printf("  L1 Cache:        0x%04X\n", proc.L1CacheHandle)
			fmt.Printf("  L2 Cache:        0x%04X\n", proc.L2CacheHandle)
			fmt.Printf("  L3 Cache:        0x%04X\n", proc.L3CacheHandle)
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType7(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(7)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 7: Cache Information ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		cache, err := type7.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Socket:          %q\n", cache.SocketDesignation)
			fmt.Printf("  Configuration:   0x%04X\n", uint16(cache.Configuration))
			fmt.Printf("  Level:           L%d\n", cache.Level())
			fmt.Printf("  Enabled:         %v\n", cache.Configuration.IsEnabled())
			fmt.Printf("  Socketed:        %v\n", cache.Configuration.IsSocketed())
			fmt.Printf("  Location:        %s\n", cache.Configuration.Location().String())
			fmt.Printf("  Mode:            %s\n", cache.Configuration.OperationalMode().String())
			fmt.Printf("  Max Size:        %s (%d KB)\n", cache.MaximumSizeString(), cache.MaximumCacheSize)
			fmt.Printf("  Installed Size:  %s (%d KB)\n", cache.InstalledSizeString(), cache.InstalledSize)
			fmt.Printf("  Supported SRAM:  0x%04X\n", uint16(cache.SupportedSRAMType))
			fmt.Printf("  Current SRAM:    0x%04X\n", uint16(cache.CurrentSRAMType))
			fmt.Printf("  Speed:           %d ns\n", cache.CacheSpeed)
			fmt.Printf("  ECC:             %s (0x%02X)\n", cache.ErrorCorrectionType.String(), uint8(cache.ErrorCorrectionType))
			fmt.Printf("  Type:            %s (0x%02X)\n", cache.SystemCacheType.String(), uint8(cache.SystemCacheType))
			fmt.Printf("  Associativity:   %s (0x%02X)\n", cache.Associativity.String(), uint8(cache.Associativity))
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType8(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(8)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 8: Port Connector Information ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		port, err := type8.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Internal Ref:    %q\n", port.InternalReferenceDesignator)
			fmt.Printf("  Internal Conn:   %s (0x%02X)\n", port.InternalConnectorType.String(), uint8(port.InternalConnectorType))
			fmt.Printf("  External Ref:    %q\n", port.ExternalReferenceDesignator)
			fmt.Printf("  External Conn:   %s (0x%02X)\n", port.ExternalConnectorType.String(), uint8(port.ExternalConnectorType))
			fmt.Printf("  Port Type:       %s (0x%02X)\n", port.PortType.String(), uint8(port.PortType))
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType9(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(9)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 9: System Slots ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		slot, err := type9.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Designation:     %q\n", slot.Designation)
			fmt.Printf("  Slot Type:       %s (0x%02X)\n", slot.SlotType.String(), uint8(slot.SlotType))
			fmt.Printf("  Data Bus Width:  %s (0x%02X)\n", slot.SlotDataBusWidth.String(), uint8(slot.SlotDataBusWidth))
			fmt.Printf("  Current Usage:   %s (0x%02X)\n", slot.CurrentUsage.String(), uint8(slot.CurrentUsage))
			fmt.Printf("  Slot Length:     %s (0x%02X)\n", slot.SlotLength.String(), uint8(slot.SlotLength))
			fmt.Printf("  Slot ID:         0x%04X\n", slot.SlotID)
			fmt.Printf("  Char 1:          0x%02X\n", uint8(slot.Characteristics1))
			fmt.Printf("  Char 2:          0x%02X\n", uint8(slot.Characteristics2))
			fmt.Printf("  PCI Address:     %s\n", slot.PCIAddress())
			fmt.Printf("  In Use:          %v\n", slot.IsInUse())
			fmt.Printf("  Hot-Plug:        %v\n", slot.SupportsHotPlug())
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType11(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(11)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 11: OEM Strings ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		oem, err := type11.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Count:           %d\n", oem.Count)
			fmt.Printf("  Strings:\n")
			for j, str := range oem.Strings {
				fmt.Printf("    [%d]: %q\n", j+1, str)
			}
		}
	}
}

func debugType16(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(16)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 16: Physical Memory Array ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		arr, err := type16.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Location:        %s (0x%02X)\n", arr.Location.String(), uint8(arr.Location))
			fmt.Printf("  Use:             %s (0x%02X)\n", arr.Use.String(), uint8(arr.Use))
			fmt.Printf("  ECC:             %s (0x%02X)\n", arr.ErrorCorrection.String(), uint8(arr.ErrorCorrection))
			fmt.Printf("  Max Capacity:    %s\n", arr.MaximumCapacityString())
			fmt.Printf("  Error Handle:    0x%04X\n", arr.MemoryErrorInfoHandle)
			fmt.Printf("  Num Devices:     %d\n", arr.NumberOfMemoryDevices)
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType17(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(17)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 17: Memory Device ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		mem, err := type17.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Array Handle:    0x%04X\n", mem.PhysicalMemoryArrayHandle)
			fmt.Printf("  Error Handle:    0x%04X\n", mem.MemoryErrorInfoHandle)
			fmt.Printf("  Total Width:     %d bits\n", mem.TotalWidth)
			fmt.Printf("  Data Width:      %d bits\n", mem.DataWidth)
			fmt.Printf("  Size:            %s (%d MB)\n", mem.SizeString(), mem.Size)
			fmt.Printf("  Form Factor:     %s (0x%02X)\n", mem.FormFactor.String(), uint8(mem.FormFactor))
			fmt.Printf("  Device Set:      0x%02X\n", mem.DeviceSet)
			fmt.Printf("  Device Locator:  %q\n", mem.DeviceLocator)
			fmt.Printf("  Bank Locator:    %q\n", mem.BankLocator)
			fmt.Printf("  Memory Type:     %s (0x%02X)\n", mem.MemoryType.String(), uint8(mem.MemoryType))
			fmt.Printf("  Type Detail:     0x%04X\n", uint16(mem.TypeDetail))
			fmt.Printf("  Speed:           %s (%d MT/s)\n", mem.SpeedString(), mem.Speed)
			fmt.Printf("  Configured:      %d MT/s\n", mem.GetConfiguredSpeed())
			fmt.Printf("  Manufacturer:    %q\n", mem.Manufacturer)
			fmt.Printf("  Serial Number:   %q\n", mem.SerialNumber)
			fmt.Printf("  Asset Tag:       %q\n", mem.AssetTag)
			fmt.Printf("  Part Number:     %q\n", mem.PartNumber)
			fmt.Printf("  Ranks:           %d\n", mem.Ranks())
			fmt.Printf("  Voltage:         %s\n", mem.VoltageString())
			fmt.Printf("  Is Populated:    %v\n", mem.IsPopulated())
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType32(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(32)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 32: System Boot Information ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		boot, err := type32.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Boot Status:     %s (0x%02X)\n", boot.BootStatus.String(), uint8(boot.BootStatus))
			fmt.Printf("  Success:         %v\n", boot.BootStatus.IsSuccess())
			fmt.Printf("  Failure:         %v\n", boot.BootStatus.IsFailure())
			fmt.Printf("  Reserved:        %s\n", hex.EncodeToString(boot.Reserved[:]))
		}
		printStrings(s.Strings, "  ")
	}
}

func debugRemainingTypes(sm *gosmbios.SMBIOS, typeCounts map[uint8]int) {
	// Types we've already handled
	handled := map[uint8]bool{
		0: true, 1: true, 2: true, 3: true, 4: true,
		7: true, 8: true, 9: true, 11: true,
		16: true, 17: true, 32: true, 127: true,
	}

	for structType := range typeCounts {
		if handled[structType] {
			continue
		}

		structs := sm.GetStructures(structType)
		if len(structs) == 0 {
			continue
		}

		fmt.Printf("\n--- Type %d: %s (Raw) ---\n", structType, types.TypeName(structType))
		for i, s := range structs {
			fmt.Printf("[%d]\n", i)
			printStructureHeader(&s)
			printHexDump(s.Data, "  ")
			printStrings(s.Strings, "  ")
		}
	}
}
