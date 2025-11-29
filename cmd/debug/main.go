// Debug tool for SMBIOS data - displays raw structure information for all types
package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/earentir/gosmbios"
	"github.com/earentir/gosmbios/types"
	"github.com/earentir/gosmbios/types/type0"
	"github.com/earentir/gosmbios/types/type1"
	"github.com/earentir/gosmbios/types/type10"
	"github.com/earentir/gosmbios/types/type11"
	"github.com/earentir/gosmbios/types/type12"
	"github.com/earentir/gosmbios/types/type13"
	"github.com/earentir/gosmbios/types/type14"
	"github.com/earentir/gosmbios/types/type15"
	"github.com/earentir/gosmbios/types/type16"
	"github.com/earentir/gosmbios/types/type17"
	"github.com/earentir/gosmbios/types/type18"
	"github.com/earentir/gosmbios/types/type19"
	"github.com/earentir/gosmbios/types/type2"
	"github.com/earentir/gosmbios/types/type20"
	"github.com/earentir/gosmbios/types/type21"
	"github.com/earentir/gosmbios/types/type22"
	"github.com/earentir/gosmbios/types/type23"
	"github.com/earentir/gosmbios/types/type24"
	"github.com/earentir/gosmbios/types/type25"
	"github.com/earentir/gosmbios/types/type26"
	"github.com/earentir/gosmbios/types/type27"
	"github.com/earentir/gosmbios/types/type28"
	"github.com/earentir/gosmbios/types/type29"
	"github.com/earentir/gosmbios/types/type3"
	"github.com/earentir/gosmbios/types/type30"
	"github.com/earentir/gosmbios/types/type31"
	"github.com/earentir/gosmbios/types/type32"
	"github.com/earentir/gosmbios/types/type33"
	"github.com/earentir/gosmbios/types/type34"
	"github.com/earentir/gosmbios/types/type35"
	"github.com/earentir/gosmbios/types/type36"
	"github.com/earentir/gosmbios/types/type37"
	"github.com/earentir/gosmbios/types/type38"
	"github.com/earentir/gosmbios/types/type39"
	"github.com/earentir/gosmbios/types/type4"
	"github.com/earentir/gosmbios/types/type40"
	"github.com/earentir/gosmbios/types/type41"
	"github.com/earentir/gosmbios/types/type42"
	"github.com/earentir/gosmbios/types/type43"
	"github.com/earentir/gosmbios/types/type44"
	"github.com/earentir/gosmbios/types/type45"
	"github.com/earentir/gosmbios/types/type46"
	"github.com/earentir/gosmbios/types/type5"
	"github.com/earentir/gosmbios/types/type6"
	"github.com/earentir/gosmbios/types/type7"
	"github.com/earentir/gosmbios/types/type8"
	"github.com/earentir/gosmbios/types/type9"
)

func main() {
	inputFile := flag.String("i", "", "Input file (gosmbios dump format)")
	showHelp := flag.Bool("h", false, "Show help")
	flag.Parse()

	if *showHelp {
		fmt.Println("smbiosdebug - Debug SMBIOS data")
		fmt.Println()
		fmt.Println("Usage: smbiosdebug [options]")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -i <file>   Read from gosmbios dump file instead of system")
		fmt.Println("  -h          Show this help message")
		os.Exit(0)
	}

	var sm *gosmbios.SMBIOS
	var err error

	if *inputFile != "" {
		sm, err = gosmbios.ReadFromFile(*inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading dump file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("(Reading from dump file: %s)\n\n", *inputFile)
	} else {
		sm, err = gosmbios.Read()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading SMBIOS: %v\n", err)
			os.Exit(1)
		}
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
	debugType5(sm)
	debugType6(sm)
	debugType7(sm)
	debugType8(sm)
	debugType9(sm)
	debugType10(sm)
	debugType11(sm)
	debugType12(sm)
	debugType13(sm)
	debugType14(sm)
	debugType15(sm)
	debugType16(sm)
	debugType17(sm)
	debugType18(sm)
	debugType19(sm)
	debugType20(sm)
	debugType21(sm)
	debugType22(sm)
	debugType23(sm)
	debugType24(sm)
	debugType25(sm)
	debugType26(sm)
	debugType27(sm)
	debugType28(sm)
	debugType29(sm)
	debugType30(sm)
	debugType31(sm)
	debugType32(sm)
	debugType33(sm)
	debugType34(sm)
	debugType35(sm)
	debugType36(sm)
	debugType37(sm)
	debugType38(sm)
	debugType39(sm)
	debugType40(sm)
	debugType41(sm)
	debugType42(sm)
	debugType43(sm)
	debugType44(sm)
	debugType45(sm)
	debugType46(sm)

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

func debugType5(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(5)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 5: Memory Controller Information (Obsolete) ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		mc, err := type5.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Error Detecting: %s (0x%02X)\n", mc.ErrorDetectingMethod.String(), uint8(mc.ErrorDetectingMethod))
			fmt.Printf("  Error Correcting:%s (0x%02X)\n", mc.ErrorCorrectingCapability.String(), uint8(mc.ErrorCorrectingCapability))
			fmt.Printf("  Supported Interleave: %s (0x%02X)\n", mc.SupportedInterleave.String(), uint8(mc.SupportedInterleave))
			fmt.Printf("  Current Interleave:   %s (0x%02X)\n", mc.CurrentInterleave.String(), uint8(mc.CurrentInterleave))
			fmt.Printf("  Max Module Size: %d MB\n", mc.MaxModuleSizeMB())
			fmt.Printf("  Supported Speeds:%s\n", mc.SupportedSpeeds.String())
			fmt.Printf("  Voltage:         %s\n", mc.MemoryModuleVoltage.String())
			fmt.Printf("  Num Slots:       %d\n", mc.NumberOfAssociatedMemorySlots)
			for j, h := range mc.MemoryModuleConfigHandles {
				fmt.Printf("    Slot %d Handle: 0x%04X\n", j, h)
			}
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType6(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(6)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 6: Memory Module Information (Obsolete) ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		mm, err := type6.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Socket:          %q\n", mm.SocketDesignation)
			fmt.Printf("  Bank Connection: %s\n", mm.BankConnectionString())
			fmt.Printf("  Current Speed:   %d ns\n", mm.CurrentSpeed)
			fmt.Printf("  Memory Type:     %s\n", mm.CurrentMemoryType.String())
			fmt.Printf("  Installed Size:  %s\n", mm.InstalledSize.String())
			fmt.Printf("  Enabled Size:    %s\n", mm.EnabledSize.String())
			fmt.Printf("  Error Status:    %s\n", mm.ErrorStatus.String())
			fmt.Printf("  Is Installed:    %v\n", mm.IsInstalled())
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
			fmt.Printf("  Max Size:        %s (%d KB)\n", cache.MaximumSizeString(), cache.MaximumSize)
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

func debugType10(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(10)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 10: On Board Devices (Obsolete) ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		obd, err := type10.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Device Count:    %d\n", len(obd.Devices))
			for j, dev := range obd.Devices {
				fmt.Printf("  Device %d:\n", j)
				fmt.Printf("    Type:          %s (0x%02X)\n", dev.DeviceType.String(), uint8(dev.DeviceType))
				fmt.Printf("    Enabled:       %v\n", dev.Enabled)
				fmt.Printf("    Description:   %q\n", dev.Description)
			}
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

func debugType12(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(12)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 12: System Configuration Options ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		cfg, err := type12.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Count:           %d\n", cfg.Count)
			fmt.Printf("  Options:\n")
			for j, opt := range cfg.Options {
				fmt.Printf("    [%d]: %q\n", j+1, opt)
			}
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType13(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(13)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 13: BIOS Language Information ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		lang, err := type13.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Installable:     %d\n", lang.InstallableLanguages)
			fmt.Printf("  Flags:           %s (0x%02X)\n", lang.Flags.String(), uint8(lang.Flags))
			fmt.Printf("  Current:         %q\n", lang.CurrentLanguage)
			fmt.Printf("  Languages:\n")
			for j, l := range lang.Languages {
				fmt.Printf("    [%d]: %q\n", j+1, l)
			}
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType14(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(14)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 14: Group Associations ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		grp, err := type14.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Group Name:      %q\n", grp.GroupName)
			fmt.Printf("  Items:           %d\n", len(grp.Items))
			for j, item := range grp.Items {
				fmt.Printf("    [%d]: Type %d, Handle 0x%04X\n", j, item.ItemType, item.ItemHandle)
			}
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType15(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(15)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 15: System Event Log ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		log, err := type15.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Log Area Length: %d bytes\n", log.LogAreaLength)
			fmt.Printf("  Log Header Start:0x%04X\n", log.LogHeaderStartOffset)
			fmt.Printf("  Log Data Start:  0x%04X\n", log.LogDataStartOffset)
			fmt.Printf("  Access Method:   %s (0x%02X)\n", log.AccessMethod.String(), uint8(log.AccessMethod))
			fmt.Printf("  Log Status:      0x%02X\n", uint8(log.LogStatus))
			fmt.Printf("  Log Full:        %v\n", log.LogStatus.IsFull())
			fmt.Printf("  Log Valid:       %v\n", log.LogStatus.IsValid())
			fmt.Printf("  Change Token:    0x%08X\n", log.LogChangeToken)
			fmt.Printf("  Log Data Format: 0x%02X\n", log.LogHeaderFormat)
			fmt.Printf("  Type Descriptors:%d\n", log.NumberOfSupportedLogTypes)
		}
		printStrings(s.Strings, "  ")
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
			fmt.Printf("  Error Handle:    0x%04X\n", arr.ErrorInformationHandle)
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
			fmt.Printf("  Error Handle:    0x%04X\n", mem.MemoryErrorInformationHandle)
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

func debugType18(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(18)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 18: 32-Bit Memory Error Information ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		err32, err := type18.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Error Type:      %s (0x%02X)\n", err32.ErrorType.String(), uint8(err32.ErrorType))
			fmt.Printf("  Error Granular:  %s (0x%02X)\n", err32.ErrorGranularity.String(), uint8(err32.ErrorGranularity))
			fmt.Printf("  Error Operation: %s (0x%02X)\n", err32.ErrorOperation.String(), uint8(err32.ErrorOperation))
			fmt.Printf("  Vendor Syndrome: 0x%08X\n", err32.VendorSyndrome)
			fmt.Printf("  Memory Address:  0x%08X (Unknown: %v)\n", err32.MemoryArrayErrorAddress, err32.IsAddressUnknown())
			fmt.Printf("  Device Address:  0x%08X (Unknown: %v)\n", err32.DeviceErrorAddress, err32.IsDeviceAddressUnknown())
			fmt.Printf("  Error Resolution:0x%08X\n", err32.ErrorResolution)
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType19(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(19)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 19: Memory Array Mapped Address ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		map19, err := type19.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Starting Address:0x%X (%d KB)\n", map19.GetStartingAddressBytes(), map19.StartingAddress)
			fmt.Printf("  Ending Address:  0x%X (%d KB)\n", map19.GetEndingAddressBytes(), map19.EndingAddress)
			fmt.Printf("  Array Handle:    0x%04X\n", map19.MemoryArrayHandle)
			fmt.Printf("  Partition Width: %d\n", map19.PartitionWidth)
			fmt.Printf("  Size:            %s\n", map19.GetSizeString())
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType20(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(20)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 20: Memory Device Mapped Address ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		map20, err := type20.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Starting Address:0x%X (%d KB)\n", map20.GetStartingAddressBytes(), map20.StartingAddress)
			fmt.Printf("  Ending Address:  0x%X (%d KB)\n", map20.GetEndingAddressBytes(), map20.EndingAddress)
			fmt.Printf("  Device Handle:   0x%04X\n", map20.MemoryDeviceHandle)
			fmt.Printf("  Array Map Handle:0x%04X\n", map20.MemoryArrayMappedAddressHandle)
			fmt.Printf("  Partition Row:   %s\n", map20.PartitionRowPositionString())
			fmt.Printf("  Interleave Pos:  %s\n", map20.InterleavePositionString())
			fmt.Printf("  Interleave Depth:%d\n", map20.InterleavedDataDepth)
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType21(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(21)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 21: Built-in Pointing Device ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		dev, err := type21.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Type:            %s (0x%02X)\n", dev.DeviceType.String(), uint8(dev.DeviceType))
			fmt.Printf("  Interface:       %s (0x%02X)\n", dev.Interface.String(), uint8(dev.Interface))
			fmt.Printf("  Buttons:         %d\n", dev.NumberOfButtons)
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType22(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(22)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 22: Portable Battery ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		bat, err := type22.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Location:        %q\n", bat.Location)
			fmt.Printf("  Manufacturer:    %q\n", bat.Manufacturer)
			fmt.Printf("  Manufacture Date:%q\n", bat.ManufactureDate)
			fmt.Printf("  Serial Number:   %q\n", bat.SerialNumber)
			fmt.Printf("  Device Name:     %q\n", bat.DeviceName)
			fmt.Printf("  Chemistry:       %s (0x%02X)\n", bat.DeviceChemistry.String(), uint8(bat.DeviceChemistry))
			fmt.Printf("  Design Capacity: %s\n", bat.DesignCapacityString())
			fmt.Printf("  Design Voltage:  %s\n", bat.DesignVoltageString())
			fmt.Printf("  Max Error:       %s\n", bat.MaximumErrorString())
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType23(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(23)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 23: System Reset ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		rst, err := type23.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Capabilities:    0x%02X\n", uint8(rst.Capabilities))
			fmt.Printf("  Enabled:         %v\n", rst.Capabilities.IsEnabled())
			fmt.Printf("  Boot Option:     %s\n", rst.Capabilities.BootOption().String())
			fmt.Printf("  Boot Option On Limit: %s\n", rst.Capabilities.BootOptionOnLimit().String())
			fmt.Printf("  Watchdog Timer:  %v\n", rst.Capabilities.WatchdogTimerPresent())
			fmt.Printf("  Reset Count:     %s\n", rst.ResetCountString())
			fmt.Printf("  Reset Limit:     %s\n", rst.ResetLimitString())
			fmt.Printf("  Timer Interval:  %s\n", rst.TimerIntervalString())
			fmt.Printf("  Timeout:         %s\n", rst.TimeoutString())
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType24(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(24)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 24: Hardware Security ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		sec, err := type24.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Settings:        0x%02X\n", uint8(sec.HardwareSettings))
			fmt.Printf("  Power-on Passwd: %s\n", sec.HardwareSettings.PowerOnPasswordStatus().String())
			fmt.Printf("  Keyboard Passwd: %s\n", sec.HardwareSettings.KeyboardPasswordStatus().String())
			fmt.Printf("  Admin Passwd:    %s\n", sec.HardwareSettings.AdministratorPasswordStatus().String())
			fmt.Printf("  Front Panel:     %s\n", sec.HardwareSettings.FrontPanelResetStatus().String())
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType25(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(25)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 25: System Power Controls ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		pwr, err := type25.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Next Power On:   %s\n", pwr.NextPowerOnString())
			fmt.Printf("  Is Scheduled:    %v\n", pwr.IsScheduled())
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType26(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(26)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 26: Voltage Probe ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		vp, err := type26.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Description:     %q\n", vp.Description)
			fmt.Printf("  Location:        %s\n", vp.LocationAndStatus.Location().String())
			fmt.Printf("  Status:          %s\n", vp.LocationAndStatus.Status().String())
			fmt.Printf("  Maximum Value:   %s\n", vp.MaximumValueString())
			fmt.Printf("  Minimum Value:   %s\n", vp.MinimumValueString())
			fmt.Printf("  Resolution:      %s\n", vp.ResolutionString())
			fmt.Printf("  Tolerance:       %d mV\n", vp.Tolerance)
			fmt.Printf("  Accuracy:        %s\n", vp.AccuracyString())
			fmt.Printf("  Nominal Value:   %s\n", vp.NominalValueString())
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType27(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(27)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 27: Cooling Device ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		cd, err := type27.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Temp Probe Handle: 0x%04X\n", cd.TemperatureProbeHandle)
			fmt.Printf("  Device Type:     %s\n", cd.DeviceTypeAndStatus.DeviceType().String())
			fmt.Printf("  Status:          %s\n", cd.DeviceTypeAndStatus.Status().String())
			fmt.Printf("  Cooling Unit:    %s\n", cd.CoolingUnitGroupString())
			fmt.Printf("  OEM-Defined:     0x%08X\n", cd.OEMDefined)
			fmt.Printf("  Nominal Speed:   %s\n", cd.NominalSpeedString())
			fmt.Printf("  Description:     %q\n", cd.Description)
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType28(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(28)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 28: Temperature Probe ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		tp, err := type28.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Description:     %q\n", tp.Description)
			fmt.Printf("  Location:        %s\n", tp.LocationAndStatus.Location().String())
			fmt.Printf("  Status:          %s\n", tp.LocationAndStatus.Status().String())
			fmt.Printf("  Maximum Value:   %s\n", tp.MaximumValueString())
			fmt.Printf("  Minimum Value:   %s\n", tp.MinimumValueString())
			fmt.Printf("  Resolution:      %s\n", tp.ResolutionString())
			fmt.Printf("  Tolerance:       %s\n", tp.ToleranceString())
			fmt.Printf("  Accuracy:        %s\n", tp.AccuracyString())
			fmt.Printf("  Nominal Value:   %s\n", tp.NominalValueString())
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType29(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(29)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 29: Electrical Current Probe ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		cp, err := type29.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Description:     %q\n", cp.Description)
			fmt.Printf("  Location:        %s\n", cp.LocationAndStatus.Location().String())
			fmt.Printf("  Status:          %s\n", cp.LocationAndStatus.Status().String())
			fmt.Printf("  Maximum Value:   %s\n", cp.MaximumValueString())
			fmt.Printf("  Minimum Value:   %s\n", cp.MinimumValueString())
			fmt.Printf("  Resolution:      %s\n", cp.ResolutionString())
			fmt.Printf("  Tolerance:       %d mA\n", cp.Tolerance)
			fmt.Printf("  Accuracy:        %s\n", cp.AccuracyString())
			fmt.Printf("  Nominal Value:   %s\n", cp.NominalValueString())
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType30(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(30)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 30: Out-of-Band Remote Access ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		oob, err := type30.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Manufacturer:    %q\n", oob.ManufacturerName)
			fmt.Printf("  Connections:     0x%02X\n", uint8(oob.Connections))
			fmt.Printf("  Inbound Enabled: %v\n", oob.Connections.InboundEnabled())
			fmt.Printf("  Outbound Enabled:%v\n", oob.Connections.OutboundEnabled())
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType31(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(31)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 31: Boot Integrity Services Entry Point ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		bis, err := type31.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Checksum:        0x%02X\n", bis.Checksum)
			fmt.Printf("  Reserved1:       0x%02X\n", bis.Reserved1)
			fmt.Printf("  Reserved2:       0x%04X\n", bis.Reserved2)
			fmt.Printf("  BIS Entry Point: 0x%08X\n", bis.BISEntryPoint)
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

func debugType33(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(33)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 33: 64-Bit Memory Error Information ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		err64, err := type33.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Error Type:      %s (0x%02X)\n", err64.ErrorType.String(), uint8(err64.ErrorType))
			fmt.Printf("  Error Granular:  %s (0x%02X)\n", err64.ErrorGranularity.String(), uint8(err64.ErrorGranularity))
			fmt.Printf("  Error Operation: %s (0x%02X)\n", err64.ErrorOperation.String(), uint8(err64.ErrorOperation))
			fmt.Printf("  Vendor Syndrome: 0x%08X\n", err64.VendorSyndrome)
			fmt.Printf("  Memory Address:  0x%016X (Unknown: %v)\n", err64.MemoryArrayErrorAddress, err64.IsAddressUnknown())
			fmt.Printf("  Device Address:  0x%016X (Unknown: %v)\n", err64.DeviceErrorAddress, err64.IsDeviceAddressUnknown())
			fmt.Printf("  Error Resolution:0x%08X\n", err64.ErrorResolution)
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType34(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(34)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 34: Management Device ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		md, err := type34.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Description:     %q\n", md.Description)
			fmt.Printf("  Device Type:     %s (0x%02X)\n", md.DeviceType.String(), uint8(md.DeviceType))
			fmt.Printf("  Address:         0x%08X\n", md.Address)
			fmt.Printf("  Address Type:    %s (0x%02X)\n", md.AddressType.String(), uint8(md.AddressType))
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType35(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(35)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 35: Management Device Component ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		mdc, err := type35.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Description:     %q\n", mdc.Description)
			fmt.Printf("  Mgmt Device:     0x%04X\n", mdc.ManagementDeviceHandle)
			fmt.Printf("  Component:       0x%04X\n", mdc.ComponentHandle)
			fmt.Printf("  Threshold:       0x%04X\n", mdc.ThresholdHandle)
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType36(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(36)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 36: Management Device Threshold Data ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		mdt, err := type36.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Lower Non-Critical:  %s\n", mdt.LowerNonCriticalString())
			fmt.Printf("  Upper Non-Critical:  %s\n", mdt.UpperNonCriticalString())
			fmt.Printf("  Lower Critical:      %s\n", mdt.LowerCriticalString())
			fmt.Printf("  Upper Critical:      %s\n", mdt.UpperCriticalString())
			fmt.Printf("  Lower Non-Recoverable:%s\n", mdt.LowerNonRecoverableString())
			fmt.Printf("  Upper Non-Recoverable:%s\n", mdt.UpperNonRecoverableString())
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType37(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(37)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 37: Memory Channel ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		mc, err := type37.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Channel Type:    %s (0x%02X)\n", mc.ChannelType.String(), uint8(mc.ChannelType))
			fmt.Printf("  Max Channel Load:%d\n", mc.MaximumChannelLoad)
			fmt.Printf("  Memory Device Count: %d\n", mc.MemoryDeviceCount)
			for j, dev := range mc.MemoryDevices {
				fmt.Printf("    Device %d: Handle 0x%04X, Load %d\n", j, dev.MemoryDeviceHandle, dev.MemoryDeviceLoad)
			}
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType38(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(38)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 38: IPMI Device Information ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		ipmi, err := type38.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Interface Type:  %s (0x%02X)\n", ipmi.InterfaceType.String(), uint8(ipmi.InterfaceType))
			fmt.Printf("  IPMI Spec Rev:   %s\n", ipmi.SpecificationRevisionString())
			fmt.Printf("  I2C Address:     %s\n", ipmi.I2CAddressString())
			fmt.Printf("  NV Storage:      0x%02X\n", ipmi.NVStorageDeviceAddress)
			fmt.Printf("  Base Address:    %s\n", ipmi.BaseAddressString())
			fmt.Printf("  Base Addr Mod:   0x%02X\n", ipmi.BaseAddressModifier)
			fmt.Printf("  Interrupt:       %s\n", ipmi.InterruptNumberString())
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType39(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(39)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 39: System Power Supply ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		psu, err := type39.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Power Unit Group:%d\n", psu.PowerUnitGroup)
			fmt.Printf("  Location:        %q\n", psu.Location)
			fmt.Printf("  Device Name:     %q\n", psu.DeviceName)
			fmt.Printf("  Manufacturer:    %q\n", psu.Manufacturer)
			fmt.Printf("  Serial Number:   %q\n", psu.SerialNumber)
			fmt.Printf("  Asset Tag:       %q\n", psu.AssetTagNumber)
			fmt.Printf("  Model Part:      %q\n", psu.ModelPartNumber)
			fmt.Printf("  Revision Level:  %q\n", psu.RevisionLevel)
			fmt.Printf("  Max Power:       %s\n", psu.MaxPowerCapacityString())
			fmt.Printf("  Characteristics: 0x%04X\n", uint16(psu.Characteristics))
			fmt.Printf("  Status:          %s\n", psu.Characteristics.Status().String())
			fmt.Printf("  Type:            %s\n", psu.Characteristics.Type().String())
			fmt.Printf("  Hot Replaceable: %v\n", psu.Characteristics.IsHotReplaceable())
			fmt.Printf("  Present:         %v\n", psu.Characteristics.IsPresent())
			fmt.Printf("  Voltage Handle:  0x%04X\n", psu.InputVoltageProbeHandle)
			fmt.Printf("  Cooling Handle:  0x%04X\n", psu.CoolingDeviceHandle)
			fmt.Printf("  Current Handle:  0x%04X\n", psu.InputCurrentProbeHandle)
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType40(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(40)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 40: Additional Information ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		ai, err := type40.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Entries:         %d\n", ai.NumberOfEntries)
			for j, entry := range ai.Entries {
				fmt.Printf("  Entry %d:\n", j)
				fmt.Printf("    Entry Length:  %d\n", entry.EntryLength)
				fmt.Printf("    Ref Handle:    0x%04X\n", entry.ReferencedHandle)
				fmt.Printf("    Ref Offset:    0x%02X\n", entry.ReferencedOffset)
				fmt.Printf("    String:        %q\n", entry.String)
			}
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType41(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(41)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 41: Onboard Devices Extended Information ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		dev, err := type41.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Reference:       %q\n", dev.ReferenceDesignation)
			fmt.Printf("  Device Type:     %s (0x%02X)\n", dev.TypeString(), dev.DeviceType)
			fmt.Printf("  Status:          %s\n", dev.StatusString())
			fmt.Printf("  Instance:        %d\n", dev.DeviceTypeInstance)
			fmt.Printf("  Segment Group:   %d\n", dev.SegmentGroupNumber)
			fmt.Printf("  Bus Number:      %d\n", dev.BusNumber)
			fmt.Printf("  Device/Function: %d/%d\n", (dev.DeviceFunctionNumber>>3)&0x1F, dev.DeviceFunctionNumber&0x07)
			fmt.Printf("  PCI Address:     %s\n", dev.PCIAddress())
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType42(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(42)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 42: Management Controller Host Interface ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		mchi, err := type42.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Interface Type:  %s (0x%02X)\n", mchi.InterfaceType.String(), uint8(mchi.InterfaceType))
			fmt.Printf("  IF Data Length:  %d bytes\n", len(mchi.InterfaceTypeSpecificData))
			if len(mchi.InterfaceTypeSpecificData) > 0 {
				fmt.Printf("  IF Data:         %s\n", hex.EncodeToString(mchi.InterfaceTypeSpecificData))
			}
			fmt.Printf("  Protocol Records:%d\n", len(mchi.ProtocolRecords))
			for j, pr := range mchi.ProtocolRecords {
				fmt.Printf("    Protocol %d: %s (0x%02X)\n", j, pr.ProtocolType.String(), uint8(pr.ProtocolType))
				if len(pr.ProtocolTypeSpecific) > 0 {
					fmt.Printf("      Data: %s\n", hex.EncodeToString(pr.ProtocolTypeSpecific))
				}
			}
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType43(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(43)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 43: TPM Device ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		tpm, err := type43.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Vendor ID:       %s\n", tpm.VendorIDString())
			fmt.Printf("  Spec Version:    %s\n", tpm.SpecVersionString())
			fmt.Printf("  Firmware Version:%s\n", tpm.FirmwareVersionString())
			fmt.Printf("  Description:     %q\n", tpm.Description)
			fmt.Printf("  Characteristics: 0x%016X\n", tpm.Characteristics)
			fmt.Printf("  OEM-Defined:     0x%08X\n", tpm.OEMDefined)
			fmt.Printf("  Family:          %s\n", tpm.Family())
			fmt.Printf("  Supported:       %v\n", tpm.IsSupported())
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType44(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(44)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 44: Processor Additional Information ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		pai, err := type44.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Referenced Handle: 0x%04X\n", pai.ReferencedHandle)
			fmt.Printf("  Block Length:    %d\n", pai.ProcessorSpecificBlock.Length)
			fmt.Printf("  Processor Type:  %s (0x%02X)\n", pai.ProcessorSpecificBlock.ProcessorType.String(), uint8(pai.ProcessorSpecificBlock.ProcessorType))
			if len(pai.ProcessorSpecificBlock.Data) > 0 {
				fmt.Printf("  Block Data:      %s\n", hex.EncodeToString(pai.ProcessorSpecificBlock.Data))
			}
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType45(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(45)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 45: Firmware Inventory Information ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		fw, err := type45.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Component Name:  %q\n", fw.FirmwareComponentName)
			fmt.Printf("  Version:         %q\n", fw.FirmwareVersion)
			fmt.Printf("  Version Format:  %s (0x%02X)\n", fw.VersionFormat.String(), uint8(fw.VersionFormat))
			fmt.Printf("  ID:              %q\n", fw.FirmwareID)
			fmt.Printf("  ID Format:       %s (0x%02X)\n", fw.FirmwareIDFormat.String(), uint8(fw.FirmwareIDFormat))
			fmt.Printf("  Release Date:    %q\n", fw.ReleaseDate)
			fmt.Printf("  Manufacturer:    %q\n", fw.Manufacturer)
			fmt.Printf("  Lowest Version:  %q\n", fw.LowestSupportedVersion)
			fmt.Printf("  Image Size:      %s\n", fw.ImageSizeString())
			fmt.Printf("  Characteristics: 0x%04X\n", uint16(fw.Characteristics))
			fmt.Printf("  State:           %s (0x%02X)\n", fw.State.String(), uint8(fw.State))
			fmt.Printf("  Assoc Components:%d\n", fw.AssociatedComponentCount)
		}
		printStrings(s.Strings, "  ")
	}
}

func debugType46(sm *gosmbios.SMBIOS) {
	structs := sm.GetStructures(46)
	if len(structs) == 0 {
		return
	}

	fmt.Println("\n--- Type 46: String Property ---")
	for i, s := range structs {
		fmt.Printf("[%d]\n", i)
		printStructureHeader(&s)

		sp, err := type46.Parse(&s)
		if err != nil {
			fmt.Printf("  Parse Error: %v\n", err)
			printHexDump(s.Data, "  ")
		} else {
			fmt.Printf("  Property ID:     %s (0x%04X)\n", sp.StringPropertyID.String(), uint16(sp.StringPropertyID))
			fmt.Printf("  Property Value:  %q\n", sp.StringPropertyValue)
			fmt.Printf("  Parent Handle:   0x%04X\n", sp.ParentHandle)
		}
		printStrings(s.Strings, "  ")
	}
}

func debugRemainingTypes(sm *gosmbios.SMBIOS, typeCounts map[uint8]int) {
	// Types we've already handled
	handled := map[uint8]bool{
		0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true,
		7: true, 8: true, 9: true, 10: true, 11: true,
		12: true, 13: true, 14: true, 15: true, 16: true,
		17: true, 18: true, 19: true, 20: true, 21: true,
		22: true, 23: true, 24: true, 25: true, 26: true,
		27: true, 28: true, 29: true, 30: true, 31: true, 32: true,
		33: true, 34: true, 35: true, 36: true, 37: true,
		38: true, 39: true, 40: true, 41: true, 42: true, 43: true,
		44: true, 45: true, 46: true, 127: true,
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
