// smbiosinfo - Simple tool to dump all SMBIOS information in a human-readable format
package main

import (
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
	"github.com/earentir/gosmbios/types/type127"
	"github.com/earentir/gosmbios/types/type16"
	"github.com/earentir/gosmbios/types/type17"
	"github.com/earentir/gosmbios/types/type2"
	"github.com/earentir/gosmbios/types/type21"
	"github.com/earentir/gosmbios/types/type22"
	"github.com/earentir/gosmbios/types/type26"
	"github.com/earentir/gosmbios/types/type27"
	"github.com/earentir/gosmbios/types/type28"
	"github.com/earentir/gosmbios/types/type3"
	"github.com/earentir/gosmbios/types/type32"
	"github.com/earentir/gosmbios/types/type38"
	"github.com/earentir/gosmbios/types/type39"
	"github.com/earentir/gosmbios/types/type4"
	"github.com/earentir/gosmbios/types/type41"
	"github.com/earentir/gosmbios/types/type43"
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
	printOnboardDevices(sm)
	printOnboardDevicesExtended(sm)
	printOEMStrings(sm)
	printSystemConfig(sm)
	printBIOSLanguage(sm)
	printMemoryArrays(sm)
	printMemoryDevices(sm)
	printPointingDevices(sm)
	printBatteries(sm)
	printVoltageProbes(sm)
	printCoolingDevices(sm)
	printTemperatureProbes(sm)
	printBootInfo(sm)
	printIPMI(sm)
	printPowerSupplies(sm)
	printTPM(sm)
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
	fmt.Println("  Characteristics:")
	if bios.Characteristics.Has(type0.CharBIOSUpgradeable) {
		fmt.Println("    - BIOS is upgradeable")
	}
	if bios.Characteristics.Has(type0.CharBootFromCDSupported) {
		fmt.Println("    - Boot from CD is supported")
	}
	if bios.Characteristics.Has(type0.CharSelectableBootSupported) {
		fmt.Println("    - Selectable boot is supported")
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
		fmt.Printf("  Type:                   %s\n", board.BoardType.String())
		fmt.Printf("  Features:\n")
		if board.FeatureFlags.IsHostingBoard() {
			fmt.Println("    - Hosting board")
		}
		if board.FeatureFlags.Has(type2.FeatureRequiresDaughter) {
			fmt.Println("    - Requires daughter board")
		}
		if board.FeatureFlags.Has(type2.FeatureRemovable) {
			fmt.Println("    - Removable")
		}
		if board.FeatureFlags.Has(type2.FeatureReplaceable) {
			fmt.Println("    - Replaceable")
		}
		if board.FeatureFlags.Has(type2.FeatureHotSwappable) {
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
	lockStatus := "Not Present"
	if chassis.TypeLocked {
		lockStatus = "Present"
	}
	fmt.Printf("  Lock:                   %s\n", lockStatus)
	fmt.Printf("  Version:                %s\n", chassis.Version)
	fmt.Printf("  Serial Number:          %s\n", chassis.SerialNumber)
	fmt.Printf("  Asset Tag:              %s\n", chassis.AssetTag)
	fmt.Printf("  Boot-up State:          %s\n", chassis.BootUpState.String())
	fmt.Printf("  Power Supply State:     %s\n", chassis.PowerSupplyState.String())
	fmt.Printf("  Thermal State:          %s\n", chassis.ThermalState.String())
	fmt.Printf("  Security Status:        %s\n", chassis.SecurityStatus.String())
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
		fmt.Printf("    Version:              %s\n", proc.ProcessorVersion)
		fmt.Printf("    Voltage:              %s\n", proc.Voltage.String())
		fmt.Printf("    External Clock:       %d MHz\n", proc.ExternalClock)
		fmt.Printf("    Max Speed:            %d MHz\n", proc.MaxSpeed)
		fmt.Printf("    Current Speed:        %d MHz\n", proc.CurrentSpeed)
		fmt.Printf("    Status:               %s\n", proc.Status.String())
		fmt.Printf("    Upgrade:              %s\n", proc.ProcessorUpgrade.String())
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
		fmt.Printf("    Level:                L%d\n", cache.Level())
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
			fmt.Printf("    Internal Connector:   %s\n", port.InternalConnectorType.String())
		}
		if port.ExternalReferenceDesignator != "" {
			fmt.Printf("    External Connector:   %s\n", port.ExternalConnectorType.String())
		}
		fmt.Printf("    Port Type:            %s\n", port.PortType.String())
	}
	fmt.Println()
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
		fmt.Printf("    Data Bus Width:       %s\n", slot.SlotDataBusWidth.String())
		fmt.Printf("    Bus Address:          %s\n", slot.PCIAddress())
	}
	fmt.Println()
}

func printOnboardDevices(sm *gosmbios.SMBIOS) {
	devices, err := type10.GetAllDevices(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 10: On Board Devices (Obsolete)")
	fmt.Println("================================================================================")
	for _, dev := range devices {
		status := "Disabled"
		if dev.Enabled {
			status = "Enabled"
		}
		fmt.Printf("  %s: %s (%s)\n", dev.Description, dev.DeviceType.String(), status)
	}
	fmt.Println()
}

func printOnboardDevicesExtended(sm *gosmbios.SMBIOS) {
	devices, err := type41.GetAll(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 41: Onboard Devices Extended Information")
	fmt.Println("================================================================================")
	for _, dev := range devices {
		fmt.Printf("  %s:\n", dev.ReferenceDesignation)
		fmt.Printf("    Type:                 %s\n", dev.TypeString())
		fmt.Printf("    Status:               %s\n", dev.StatusString())
		fmt.Printf("    Instance:             %d\n", dev.DeviceTypeInstance)
		fmt.Printf("    Bus Address:          %s\n", dev.PCIAddress())
	}
	fmt.Println()
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

func printSystemConfig(sm *gosmbios.SMBIOS) {
	configs, err := type12.GetAll(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 12: System Configuration Options")
	fmt.Println("================================================================================")
	for _, cfg := range configs {
		for i, opt := range cfg.Options {
			fmt.Printf("  Option %d: %s\n", i+1, opt)
		}
	}
	fmt.Println()
}

func printBIOSLanguage(sm *gosmbios.SMBIOS) {
	lang, err := type13.Get(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 13: BIOS Language Information")
	fmt.Println("================================================================================")
	fmt.Printf("  Language Format:        %s\n", lang.Flags.String())
	fmt.Printf("  Current Language:       %s\n", lang.CurrentLanguage)
	fmt.Printf("  Installable Languages:  %d\n", lang.InstallableLanguages)
	for i, l := range lang.Languages {
		fmt.Printf("    [%d]: %s\n", i+1, l)
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
		fmt.Printf("    Size:                 %s\n", dev.SizeString())
		fmt.Printf("    Form Factor:          %s\n", dev.FormFactor.String())
		fmt.Printf("    Locator:              %s\n", dev.DeviceLocator)
		fmt.Printf("    Bank Locator:         %s\n", dev.BankLocator)
		fmt.Printf("    Type:                 %s\n", dev.MemoryType.String())
		fmt.Printf("    Type Detail:          %s\n", dev.TypeDetail.String())
		fmt.Printf("    Speed:                %s\n", dev.SpeedString())
		fmt.Printf("    Manufacturer:         %s\n", dev.Manufacturer)
		fmt.Printf("    Serial Number:        %s\n", dev.SerialNumber)
		fmt.Printf("    Part Number:          %s\n", dev.PartNumber)
		fmt.Printf("    Configured Speed:     %d MT/s\n", dev.GetConfiguredSpeed())
		fmt.Printf("    Voltage:              %s\n", dev.VoltageString())
		fmt.Println()
		totalSize += dev.Size
	}
	fmt.Printf("  Total Installed Memory: %d MB (%d GB)\n", totalSize, totalSize/1024)
	fmt.Println()
}

func printPointingDevices(sm *gosmbios.SMBIOS) {
	devices, err := type21.GetAll(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 21: Built-in Pointing Device")
	fmt.Println("================================================================================")
	for _, dev := range devices {
		fmt.Printf("  Type:                   %s\n", dev.DeviceType.String())
		fmt.Printf("  Interface:              %s\n", dev.Interface.String())
		fmt.Printf("  Buttons:                %d\n", dev.NumberOfButtons)
	}
	fmt.Println()
}

func printBatteries(sm *gosmbios.SMBIOS) {
	batteries, err := type22.GetAll(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 22: Portable Battery")
	fmt.Println("================================================================================")
	for _, bat := range batteries {
		fmt.Printf("  %s:\n", bat.DeviceName)
		fmt.Printf("    Location:             %s\n", bat.Location)
		fmt.Printf("    Manufacturer:         %s\n", bat.Manufacturer)
		fmt.Printf("    Chemistry:            %s\n", bat.DeviceChemistry.String())
		fmt.Printf("    Design Capacity:      %s\n", bat.DesignCapacityString())
		fmt.Printf("    Design Voltage:       %s\n", bat.DesignVoltageString())
		fmt.Printf("    Serial Number:        %s\n", bat.SerialNumber)
		fmt.Printf("    Manufacture Date:     %s\n", bat.ManufactureDate)
		fmt.Println()
	}
}

func printVoltageProbes(sm *gosmbios.SMBIOS) {
	probes, err := type26.GetAll(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 26: Voltage Probe")
	fmt.Println("================================================================================")
	for _, probe := range probes {
		fmt.Printf("  %s:\n", probe.Description)
		fmt.Printf("    Location:             %s\n", probe.LocationAndStatus.Location().String())
		fmt.Printf("    Status:               %s\n", probe.LocationAndStatus.Status().String())
		fmt.Printf("    Minimum Value:        %s\n", probe.MinimumValueString())
		fmt.Printf("    Maximum Value:        %s\n", probe.MaximumValueString())
		fmt.Printf("    Nominal Value:        %s\n", probe.NominalValueString())
	}
	fmt.Println()
}

func printCoolingDevices(sm *gosmbios.SMBIOS) {
	devices, err := type27.GetAll(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 27: Cooling Device")
	fmt.Println("================================================================================")
	for _, dev := range devices {
		name := dev.Description
		if name == "" {
			name = dev.DeviceTypeAndStatus.DeviceType().String()
		}
		fmt.Printf("  %s:\n", name)
		fmt.Printf("    Type:                 %s\n", dev.DeviceTypeAndStatus.DeviceType().String())
		fmt.Printf("    Status:               %s\n", dev.DeviceTypeAndStatus.Status().String())
		fmt.Printf("    Nominal Speed:        %s\n", dev.NominalSpeedString())
		fmt.Printf("    Cooling Unit Group:   %s\n", dev.CoolingUnitGroupString())
	}
	fmt.Println()
}

func printTemperatureProbes(sm *gosmbios.SMBIOS) {
	probes, err := type28.GetAll(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 28: Temperature Probe")
	fmt.Println("================================================================================")
	for _, probe := range probes {
		fmt.Printf("  %s:\n", probe.Description)
		fmt.Printf("    Location:             %s\n", probe.LocationAndStatus.Location().String())
		fmt.Printf("    Status:               %s\n", probe.LocationAndStatus.Status().String())
		fmt.Printf("    Minimum Value:        %s\n", probe.MinimumValueString())
		fmt.Printf("    Maximum Value:        %s\n", probe.MaximumValueString())
		fmt.Printf("    Nominal Value:        %s\n", probe.NominalValueString())
	}
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

func printIPMI(sm *gosmbios.SMBIOS) {
	ipmi, err := type38.Get(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 38: IPMI Device Information")
	fmt.Println("================================================================================")
	fmt.Printf("  Interface Type:         %s\n", ipmi.InterfaceType.String())
	fmt.Printf("  Specification Version:  %s\n", ipmi.SpecificationRevisionString())
	fmt.Printf("  I2C Slave Address:      %s\n", ipmi.I2CAddressString())
	fmt.Printf("  Base Address:           %s\n", ipmi.BaseAddressString())
	fmt.Printf("  Interrupt:              %s\n", ipmi.InterruptNumberString())
	fmt.Println()
}

func printPowerSupplies(sm *gosmbios.SMBIOS) {
	supplies, err := type39.GetAll(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 39: System Power Supply")
	fmt.Println("================================================================================")
	for _, psu := range supplies {
		name := psu.DeviceName
		if name == "" {
			name = psu.Location
		}
		fmt.Printf("  %s:\n", name)
		fmt.Printf("    Location:             %s\n", psu.Location)
		fmt.Printf("    Manufacturer:         %s\n", psu.Manufacturer)
		fmt.Printf("    Serial Number:        %s\n", psu.SerialNumber)
		fmt.Printf("    Model:                %s\n", psu.ModelPartNumber)
		fmt.Printf("    Max Power:            %s\n", psu.MaxPowerCapacityString())
		fmt.Printf("    Status:               %s\n", psu.Characteristics.Status().String())
		fmt.Printf("    Type:                 %s\n", psu.Characteristics.Type().String())
		fmt.Printf("    Hot Replaceable:      %v\n", psu.Characteristics.IsHotReplaceable())
	}
	fmt.Println()
}

func printTPM(sm *gosmbios.SMBIOS) {
	tpm, err := type43.Get(sm)
	if err != nil {
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("Type 43: TPM Device")
	fmt.Println("================================================================================")
	fmt.Printf("  Vendor ID:              %s\n", tpm.VendorIDString())
	fmt.Printf("  Specification Version:  %s\n", tpm.SpecVersionString())
	fmt.Printf("  Firmware Version:       %s\n", tpm.FirmwareVersionString())
	fmt.Printf("  Description:            %s\n", tpm.Description)
	fmt.Printf("  Family:                 %s\n", tpm.Family())
	fmt.Printf("  Supported:              %v\n", tpm.IsSupported())
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
		7: true, 8: true, 9: true, 10: true, 11: true,
		12: true, 13: true, 16: true, 17: true, 21: true,
		22: true, 26: true, 27: true, 28: true, 32: true,
		38: true, 39: true, 41: true, 43: true, 127: true,
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
	fmt.Println("Other Structures (not displayed in detail)")
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
