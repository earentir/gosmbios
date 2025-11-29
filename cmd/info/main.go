// smbiosinfo - Simple tool to dump all SMBIOS information in a human-readable format
package main

import (
	"fmt"
	"os"

	"github.com/earentir/gosmbios"
	"github.com/earentir/gosmbios/types"
	"github.com/earentir/gosmbios/types/type0"
	"github.com/earentir/gosmbios/types/type1"
	"github.com/earentir/gosmbios/types/type11"
	"github.com/earentir/gosmbios/types/type127"
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

	// Print header
	fmt.Println("================================================================================")
	fmt.Println("                              SMBIOS INFORMATION")
	fmt.Println("================================================================================")
	fmt.Printf("\nSMBIOS Version: %s\n", sm.EntryPoint.String())
	fmt.Printf("Total Structures: %d\n", len(sm.Structures))

	// Print structure summary
	typeCounts := make(map[uint8]int)
	for _, s := range sm.Structures {
		typeCounts[s.Header.Type]++
	}
	fmt.Printf("Unique Types: %d\n\n", len(typeCounts))

	// Print all information
	printBIOS(sm)
	printSystem(sm)
	printBaseboard(sm)
	printChassis(sm)
	printProcessors(sm)
	printCaches(sm)
	printPorts(sm)
	printSlots(sm)
	printOEMStrings(sm)
	printMemoryArrays(sm)
	printMemoryDevices(sm)
	printBootInfo(sm)
	printEndOfTable(sm)
	printUnknownTypes(sm, typeCounts)

	fmt.Println("\n================================================================================")
	fmt.Println("                                   END OF DUMP")
	fmt.Println("================================================================================")
}

func printBIOS(sm *gosmbios.SMBIOS) {
	bios, err := type0.Get(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 0: BIOS Information")
	fmt.Println("================================================================================")
	fmt.Printf("  Vendor:                 %s\n", bios.Vendor)
	fmt.Printf("  Version:                %s\n", bios.Version)
	fmt.Printf("  Release Date:           %s\n", bios.ReleaseDate)
	fmt.Printf("  BIOS Revision:          %s\n", bios.BIOSVersionString())
	fmt.Printf("  EC Revision:            %s\n", bios.ECVersionString())
	fmt.Printf("  ROM Size:               %s\n", bios.ROMSizeString())
	fmt.Printf("  Address:                0x%04X0\n", bios.StartingAddressSegment)
	fmt.Printf("  Runtime Size:           %d kB\n", (0x10000-int(bios.StartingAddressSegment))*16/1024)
	fmt.Println("  Characteristics:")
	if bios.Characteristics.Has(type0.CharBIOSUpgradeable) {
		fmt.Println("    - BIOS is upgradeable")
	}
	if bios.Characteristics.Has(type0.CharBIOSShadowingAllowed) {
		fmt.Println("    - BIOS shadowing is allowed")
	}
	if bios.Characteristics.Has(type0.CharBootFromCDSupported) {
		fmt.Println("    - Boot from CD is supported")
	}
	if bios.Characteristics.Has(type0.CharSelectableBootSupported) {
		fmt.Println("    - Selectable boot is supported")
	}
	if bios.Characteristics.Has(type0.CharEDDSupported) {
		fmt.Println("    - EDD is supported")
	}
	if bios.CharacteristicsExt1.Has(type0.CharExt1ACPISupported) {
		fmt.Println("    - ACPI is supported")
	}
	if bios.CharacteristicsExt1.Has(type0.CharExt1USBLegacySupported) {
		fmt.Println("    - USB legacy is supported")
	}
	if bios.CharacteristicsExt2.Has(type0.CharExt2UEFISupported) {
		fmt.Println("    - UEFI is supported")
	}
	if bios.CharacteristicsExt2.Has(type0.CharExt2VirtualMachine) {
		fmt.Println("    - System is a virtual machine")
	}
	fmt.Println()
}

func printSystem(sm *gosmbios.SMBIOS) {
	sys, err := type1.Get(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 1: System Information")
	fmt.Println("================================================================================")
	fmt.Printf("  Manufacturer:           %s\n", sys.Manufacturer)
	fmt.Printf("  Product Name:           %s\n", sys.ProductName)
	fmt.Printf("  Version:                %s\n", sys.Version)
	fmt.Printf("  Serial Number:          %s\n", sys.SerialNumber)
	fmt.Printf("  UUID:                   %s\n", sys.UUID.String())
	fmt.Printf("  Wake-up Type:           %s\n", sys.WakeUpType.String())
	fmt.Printf("  SKU Number:             %s\n", sys.SKUNumber)
	fmt.Printf("  Family:                 %s\n", sys.Family)
	fmt.Println()
}

func printBaseboard(sm *gosmbios.SMBIOS) {
	boards := sm.GetStructures(2)
	if len(boards) == 0 {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 2: Baseboard Information")
	fmt.Println("================================================================================")
	for i, s := range boards {
		board, err := type2.Parse(&s)
		if err != nil {
			continue
		}
		if len(boards) > 1 {
			fmt.Printf("  Board %d:\n", i+1)
		}
		fmt.Printf("  Manufacturer:           %s\n", board.Manufacturer)
		fmt.Printf("  Product Name:           %s\n", board.Product)
		fmt.Printf("  Version:                %s\n", board.Version)
		fmt.Printf("  Serial Number:          %s\n", board.SerialNumber)
		fmt.Printf("  Asset Tag:              %s\n", board.AssetTag)
		fmt.Printf("  Location in Chassis:    %s\n", board.LocationInChassis)
		fmt.Printf("  Chassis Handle:         0x%04X\n", board.ChassisHandle)
		fmt.Printf("  Type:                   %s\n", board.BoardType.String())
		fmt.Printf("  Features:\n")
		if board.FeatureFlags.IsHostingBoard() {
			fmt.Println("    - Hosting board")
		}
		if board.FeatureFlags.RequiresDaughterBoard() {
			fmt.Println("    - Requires daughter board")
		}
		if board.FeatureFlags.IsRemovable() {
			fmt.Println("    - Removable")
		}
		if board.FeatureFlags.IsReplaceable() {
			fmt.Println("    - Replaceable")
		}
		if board.FeatureFlags.IsHotSwappable() {
			fmt.Println("    - Hot swappable")
		}
		fmt.Println()
	}
}

func printChassis(sm *gosmbios.SMBIOS) {
	chassis, err := type3.Get(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 3: Chassis Information")
	fmt.Println("================================================================================")
	fmt.Printf("  Manufacturer:           %s\n", chassis.Manufacturer)
	fmt.Printf("  Type:                   %s\n", chassis.Type.String())
	fmt.Printf("  Lock:                   %s\n", chassis.LockString())
	fmt.Printf("  Version:                %s\n", chassis.Version)
	fmt.Printf("  Serial Number:          %s\n", chassis.SerialNumber)
	fmt.Printf("  Asset Tag:              %s\n", chassis.AssetTag)
	fmt.Printf("  Boot-up State:          %s\n", chassis.BootUpState.String())
	fmt.Printf("  Power Supply State:     %s\n", chassis.PowerSupplyState.String())
	fmt.Printf("  Thermal State:          %s\n", chassis.ThermalState.String())
	fmt.Printf("  Security Status:        %s\n", chassis.SecurityStatus.String())
	fmt.Printf("  OEM Information:        0x%08X\n", chassis.OEMDefined)
	fmt.Printf("  Height:                 %s\n", chassis.HeightString())
	fmt.Printf("  Number of Power Cords:  %d\n", chassis.NumberOfPowerCords)
	fmt.Printf("  SKU Number:             %s\n", chassis.SKUNumber)
	fmt.Println()
}

func printProcessors(sm *gosmbios.SMBIOS) {
	processors, err := type4.GetAll(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 4: Processor Information")
	fmt.Println("================================================================================")
	for i, proc := range processors {
		fmt.Printf("  Processor %d:\n", i+1)
		fmt.Printf("    Socket Designation:   %s\n", proc.SocketDesignation)
		fmt.Printf("    Type:                 %s\n", proc.ProcessorType.String())
		fmt.Printf("    Family:               %s\n", proc.ProcessorFamily.String())
		fmt.Printf("    Manufacturer:         %s\n", proc.ProcessorManufacturer)
		fmt.Printf("    ID:                   %016X\n", proc.ProcessorID)
		fmt.Printf("    Version:              %s\n", proc.ProcessorVersion)
		fmt.Printf("    Voltage:              %s\n", proc.Voltage.String())
		fmt.Printf("    External Clock:       %d MHz\n", proc.ExternalClock)
		fmt.Printf("    Max Speed:            %d MHz\n", proc.MaxSpeed)
		fmt.Printf("    Current Speed:        %d MHz\n", proc.CurrentSpeed)
		fmt.Printf("    Status:               %s\n", proc.Status.String())
		fmt.Printf("    Upgrade:              %s\n", proc.ProcessorUpgrade.String())
		fmt.Printf("    L1 Cache Handle:      0x%04X\n", proc.L1CacheHandle)
		fmt.Printf("    L2 Cache Handle:      0x%04X\n", proc.L2CacheHandle)
		fmt.Printf("    L3 Cache Handle:      0x%04X\n", proc.L3CacheHandle)
		fmt.Printf("    Serial Number:        %s\n", proc.SerialNumber)
		fmt.Printf("    Asset Tag:            %s\n", proc.AssetTag)
		fmt.Printf("    Part Number:          %s\n", proc.PartNumber)
		fmt.Printf("    Core Count:           %d\n", proc.GetCoreCount())
		fmt.Printf("    Core Enabled:         %d\n", proc.GetCoreEnabled())
		fmt.Printf("    Thread Count:         %d\n", proc.GetThreadCount())
		fmt.Printf("    Characteristics:\n")
		if proc.ProcessorCharacteristics.Is64Bit() {
			fmt.Println("      - 64-bit capable")
		}
		if proc.ProcessorCharacteristics.IsMultiCore() {
			fmt.Println("      - Multi-Core")
		}
		if proc.ProcessorCharacteristics.IsHWThread() {
			fmt.Println("      - Hardware Thread")
		}
		if proc.ProcessorCharacteristics.IsVirtualizationCapable() {
			fmt.Println("      - Execute Protection")
		}
		if proc.ProcessorCharacteristics.IsPowerPerformanceControl() {
			fmt.Println("      - Power/Performance Control")
		}
		fmt.Println()
	}
}

func printCaches(sm *gosmbios.SMBIOS) {
	caches, err := type7.GetAll(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 7: Cache Information")
	fmt.Println("================================================================================")
	for _, cache := range caches {
		fmt.Printf("  %s:\n", cache.SocketDesignation)
		fmt.Printf("    Configuration:        %s, %s, Level %d\n",
			cache.Configuration.OperationalMode().String(),
			cache.Configuration.Location().String(),
			cache.Level())
		fmt.Printf("    Operational Mode:     %s\n", cache.Configuration.OperationalMode().String())
		fmt.Printf("    Location:             %s\n", cache.Configuration.Location().String())
		fmt.Printf("    Installed Size:       %s\n", cache.InstalledSizeString())
		fmt.Printf("    Maximum Size:         %s\n", cache.MaximumSizeString())
		fmt.Printf("    Supported SRAM Types: %s\n", cache.SupportedSRAMType.String())
		fmt.Printf("    Installed SRAM Type:  %s\n", cache.CurrentSRAMType.String())
		fmt.Printf("    Speed:                %d ns\n", cache.CacheSpeed)
		fmt.Printf("    Error Correction:     %s\n", cache.ErrorCorrectionType.String())
		fmt.Printf("    System Type:          %s\n", cache.SystemCacheType.String())
		fmt.Printf("    Associativity:        %s\n", cache.Associativity.String())
		fmt.Println()
	}
}

func printPorts(sm *gosmbios.SMBIOS) {
	ports, err := type8.GetAll(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 8: Port Connector Information")
	fmt.Println("================================================================================")
	for _, port := range ports {
		name := port.DisplayName()
		fmt.Printf("  %s:\n", name)
		if port.InternalReferenceDesignator != "" {
			fmt.Printf("    Internal Designator:  %s\n", port.InternalReferenceDesignator)
			fmt.Printf("    Internal Connector:   %s\n", port.InternalConnectorType.String())
		}
		if port.ExternalReferenceDesignator != "" {
			fmt.Printf("    External Designator:  %s\n", port.ExternalReferenceDesignator)
			fmt.Printf("    External Connector:   %s\n", port.ExternalConnectorType.String())
		}
		fmt.Printf("    Port Type:            %s\n", port.PortType.String())
		fmt.Println()
	}
}

func printSlots(sm *gosmbios.SMBIOS) {
	slots, err := type9.GetAll(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 9: System Slots")
	fmt.Println("================================================================================")
	for _, slot := range slots {
		fmt.Printf("  %s:\n", slot.Designation)
		fmt.Printf("    Type:                 %s\n", slot.SlotType.String())
		fmt.Printf("    Current Usage:        %s\n", slot.CurrentUsage.String())
		fmt.Printf("    Length:               %s\n", slot.SlotLength.String())
		fmt.Printf("    ID:                   %d\n", slot.SlotID)
		fmt.Printf("    Data Bus Width:       %s\n", slot.SlotDataBusWidth.String())
		fmt.Printf("    Bus Address:          %s\n", slot.PCIAddress())
		fmt.Printf("    Characteristics:\n")
		if slot.Characteristics1.Has(type9.SlotChar1Provides5V) {
			fmt.Println("      - 5.0 V is provided")
		}
		if slot.Characteristics1.Has(type9.SlotChar1Provides3_3V) {
			fmt.Println("      - 3.3 V is provided")
		}
		if slot.Characteristics1.Has(type9.SlotChar1Shared) {
			fmt.Println("      - Opening is shared")
		}
		if slot.Characteristics2.Has(type9.SlotChar2PMESignal) {
			fmt.Println("      - PME signal is supported")
		}
		if slot.Characteristics2.Has(type9.SlotChar2HotPlugDevices) {
			fmt.Println("      - Hot-plug devices are supported")
		}
		if slot.Characteristics2.Has(type9.SlotChar2SMBusSignal) {
			fmt.Println("      - SMBus signal is supported")
		}
		if slot.Characteristics2.Has(type9.SlotChar2Bifurcation) {
			fmt.Println("      - PCIe slot bifurcation is supported")
		}
		fmt.Println()
	}
}

func printOEMStrings(sm *gosmbios.SMBIOS) {
	oems, err := type11.GetAll(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 11: OEM Strings")
	fmt.Println("================================================================================")
	for _, oem := range oems {
		for i, str := range oem.Strings {
			fmt.Printf("  String %d: %s\n", i+1, str)
		}
	}
	fmt.Println()
}

func printMemoryArrays(sm *gosmbios.SMBIOS) {
	arrays, err := type16.GetAll(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 16: Physical Memory Array")
	fmt.Println("================================================================================")
	for i, arr := range arrays {
		if len(arrays) > 1 {
			fmt.Printf("  Array %d:\n", i+1)
		}
		fmt.Printf("  Location:               %s\n", arr.Location.String())
		fmt.Printf("  Use:                    %s\n", arr.Use.String())
		fmt.Printf("  Error Correction:       %s\n", arr.ErrorCorrection.String())
		fmt.Printf("  Maximum Capacity:       %s\n", arr.MaximumCapacityString())
		fmt.Printf("  Error Info Handle:      0x%04X\n", arr.MemoryErrorInfoHandle)
		fmt.Printf("  Number of Devices:      %d\n", arr.NumberOfMemoryDevices)
		fmt.Println()
	}
}

func printMemoryDevices(sm *gosmbios.SMBIOS) {
	devices, err := type17.GetAll(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 17: Memory Device")
	fmt.Println("================================================================================")
	var totalSize uint64
	for _, dev := range devices {
		fmt.Printf("  %s:\n", dev.DeviceLocator)
		fmt.Printf("    Array Handle:         0x%04X\n", dev.PhysicalMemoryArrayHandle)
		fmt.Printf("    Error Info Handle:    0x%04X\n", dev.MemoryErrorInfoHandle)
		fmt.Printf("    Total Width:          %d bits\n", dev.TotalWidth)
		fmt.Printf("    Data Width:           %d bits\n", dev.DataWidth)
		fmt.Printf("    Size:                 %s\n", dev.SizeString())
		fmt.Printf("    Form Factor:          %s\n", dev.FormFactor.String())
		fmt.Printf("    Set:                  0x%02X\n", dev.DeviceSet)
		fmt.Printf("    Locator:              %s\n", dev.DeviceLocator)
		fmt.Printf("    Bank Locator:         %s\n", dev.BankLocator)
		fmt.Printf("    Type:                 %s\n", dev.MemoryType.String())
		fmt.Printf("    Type Detail:          %s\n", dev.TypeDetail.String())
		fmt.Printf("    Speed:                %s\n", dev.SpeedString())
		fmt.Printf("    Manufacturer:         %s\n", dev.Manufacturer)
		fmt.Printf("    Serial Number:        %s\n", dev.SerialNumber)
		fmt.Printf("    Asset Tag:            %s\n", dev.AssetTag)
		fmt.Printf("    Part Number:          %s\n", dev.PartNumber)
		fmt.Printf("    Rank:                 %d\n", dev.Ranks())
		fmt.Printf("    Configured Speed:     %d MT/s\n", dev.GetConfiguredSpeed())
		fmt.Printf("    Voltage:              %s\n", dev.VoltageString())
		fmt.Println()
		totalSize += dev.Size
	}
	fmt.Printf("  Total Installed Memory: %d MB (%d GB)\n", totalSize, totalSize/1024)
	fmt.Println()
}

func printBootInfo(sm *gosmbios.SMBIOS) {
	boot, err := type32.Get(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 32: System Boot Information")
	fmt.Println("================================================================================")
	fmt.Printf("  Status:                 %s\n", boot.BootStatus.String())
	fmt.Println()
}

func printEndOfTable(sm *gosmbios.SMBIOS) {
	eot, err := type127.Get(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 127: End Of Table")
	fmt.Println("================================================================================")
	fmt.Printf("  Handle:                 0x%04X\n", eot.Header.Handle)
	fmt.Println()
}

func printUnknownTypes(sm *gosmbios.SMBIOS, typeCounts map[uint8]int) {
	// Types we handle
	handled := map[uint8]bool{
		0: true, 1: true, 2: true, 3: true, 4: true,
		7: true, 8: true, 9: true, 11: true,
		16: true, 17: true, 32: true, 127: true,
	}

	hasUnknown := false
	for t := range typeCounts {
		if !handled[t] {
			hasUnknown = true
			break
		}
	}

	if !hasUnknown {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Other Structures (not parsed)")
	fmt.Println("================================================================================")
	for t := uint8(0); t <= 255; t++ {
		if count, ok := typeCounts[t]; ok && !handled[t] {
			fmt.Printf("  Type %3d: %d structure(s) - %s\n", t, count, types.TypeName(t))
		}
		if t == 255 {
			break
		}
	}
	fmt.Println()
}
