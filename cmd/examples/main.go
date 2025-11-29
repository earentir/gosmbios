// Example usage of the gosmbios package
package main

import (
	"fmt"
	"os"

	"github.com/earentir/gosmbios"
	"github.com/earentir/gosmbios/types/type0"
	"github.com/earentir/gosmbios/types/type1"
	"github.com/earentir/gosmbios/types/type16"
	"github.com/earentir/gosmbios/types/type17"
	"github.com/earentir/gosmbios/types/type2"
	"github.com/earentir/gosmbios/types/type22"
	"github.com/earentir/gosmbios/types/type3"
	"github.com/earentir/gosmbios/types/type4"
	"github.com/earentir/gosmbios/types/type43"
	"github.com/earentir/gosmbios/types/type7"
)

func main() {
	// Read SMBIOS data from the system
	sm, err := gosmbios.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading SMBIOS: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("SMBIOS Version: %s\n", sm.EntryPoint.String())
	fmt.Printf("Number of Structures: %d\n\n", len(sm.Structures))

	// BIOS Information (Type 0)
	printBIOSInfo(sm)

	// System Information (Type 1)
	printSystemInfo(sm)

	// Baseboard Information (Type 2)
	printBaseboardInfo(sm)

	// Chassis Information (Type 3)
	printChassisInfo(sm)

	// Processor Information (Type 4)
	printProcessorInfo(sm)

	// Cache Information (Type 7)
	printCacheInfo(sm)

	// Memory Information (Type 16/17)
	printMemoryInfo(sm)

	// TPM Information (Type 43)
	printTPMInfo(sm)

	// Battery Information (Type 22)
	printBatteryInfo(sm)
}

func printBIOSInfo(sm *gosmbios.SMBIOS) {
	bios, err := type0.Get(sm)
	if err != nil {
		fmt.Printf("BIOS Information: Not available\n\n")
		return
	}

	fmt.Println("=== BIOS Information ===")
	fmt.Printf("  Vendor:       %s\n", bios.Vendor)
	fmt.Printf("  Version:      %s\n", bios.Version)
	fmt.Printf("  Release Date: %s\n", bios.ReleaseDate)
	fmt.Printf("  ROM Size:     %s\n", bios.ROMSizeString())
	fmt.Printf("  UEFI:         %v\n", bios.IsUEFI())
	fmt.Printf("  Virtual:      %v\n", bios.IsVirtualMachine())
	fmt.Println()
}

func printSystemInfo(sm *gosmbios.SMBIOS) {
	sys, err := type1.Get(sm)
	if err != nil {
		fmt.Printf("System Information: Not available\n\n")
		return
	}

	fmt.Println("=== System Information ===")
	fmt.Printf("  Manufacturer: %s\n", sys.Manufacturer)
	fmt.Printf("  Product:      %s\n", sys.ProductName)
	fmt.Printf("  Version:      %s\n", sys.Version)
	fmt.Printf("  Serial:       %s\n", sys.SerialNumber)
	fmt.Printf("  UUID:         %s\n", sys.UUID.String())
	fmt.Printf("  SKU:          %s\n", sys.SKUNumber)
	fmt.Printf("  Family:       %s\n", sys.Family)
	fmt.Printf("  Wake-up Type: %s\n", sys.WakeUpType.String())
	fmt.Println()
}

func printBaseboardInfo(sm *gosmbios.SMBIOS) {
	board, err := type2.Get(sm)
	if err != nil {
		fmt.Printf("Baseboard Information: Not available\n\n")
		return
	}

	fmt.Println("=== Baseboard Information ===")
	fmt.Printf("  Manufacturer: %s\n", board.Manufacturer)
	fmt.Printf("  Product:      %s\n", board.Product)
	fmt.Printf("  Version:      %s\n", board.Version)
	fmt.Printf("  Serial:       %s\n", board.SerialNumber)
	fmt.Printf("  Asset Tag:    %s\n", board.AssetTag)
	fmt.Printf("  Board Type:   %s\n", board.BoardType.String())
	fmt.Println()
}

func printChassisInfo(sm *gosmbios.SMBIOS) {
	chassis, err := type3.Get(sm)
	if err != nil {
		fmt.Printf("Chassis Information: Not available\n\n")
		return
	}

	fmt.Println("=== Chassis Information ===")
	fmt.Printf("  Manufacturer:   %s\n", chassis.Manufacturer)
	fmt.Printf("  Type:           %s\n", chassis.Type.String())
	fmt.Printf("  Version:        %s\n", chassis.Version)
	fmt.Printf("  Serial:         %s\n", chassis.SerialNumber)
	fmt.Printf("  Asset Tag:      %s\n", chassis.AssetTag)
	fmt.Printf("  Boot-up State:  %s\n", chassis.BootUpState.String())
	fmt.Printf("  Power State:    %s\n", chassis.PowerSupplyState.String())
	fmt.Printf("  Thermal State:  %s\n", chassis.ThermalState.String())
	fmt.Printf("  Security:       %s\n", chassis.SecurityStatus.String())
	fmt.Printf("  Height:         %s\n", chassis.HeightString())
	fmt.Printf("  Portable:       %v\n", chassis.Type.IsPortable())
	fmt.Println()
}

func printProcessorInfo(sm *gosmbios.SMBIOS) {
	processors, err := type4.GetAll(sm)
	if err != nil {
		fmt.Printf("Processor Information: Not available\n\n")
		return
	}

	fmt.Println("=== Processor Information ===")
	for i, proc := range processors {
		if !proc.Status.IsPopulated() {
			continue
		}
		fmt.Printf("  Processor %d:\n", i+1)
		fmt.Printf("    Socket:       %s\n", proc.SocketDesignation)
		fmt.Printf("    Name:         %s\n", proc.DisplayName())
		fmt.Printf("    Manufacturer: %s\n", proc.ProcessorManufacturer)
		fmt.Printf("    Family:       %s\n", proc.ProcessorFamily.String())
		fmt.Printf("    Cores:        %d\n", proc.GetCoreCount())
		fmt.Printf("    Enabled:      %d\n", proc.GetCoreEnabled())
		fmt.Printf("    Threads:      %d\n", proc.GetThreadCount())
		fmt.Printf("    Max Speed:    %d MHz\n", proc.MaxSpeed)
		fmt.Printf("    Current:      %d MHz\n", proc.CurrentSpeed)
		fmt.Printf("    Status:       %s\n", proc.Status.String())
		fmt.Printf("    64-bit:       %v\n", proc.ProcessorCharacteristics.Is64Bit())
		fmt.Printf("    Upgrade:      %s\n", proc.ProcessorUpgrade.String())
		fmt.Printf("    Voltage:      %s\n", proc.Voltage.String())
		fmt.Println()
	}
}

func printCacheInfo(sm *gosmbios.SMBIOS) {
	caches, err := type7.GetAll(sm)
	if err != nil {
		fmt.Printf("Cache Information: Not available\n\n")
		return
	}

	fmt.Println("=== Cache Information ===")
	for _, cache := range caches {
		if !cache.Configuration.IsEnabled() {
			continue
		}
		fmt.Printf("  %s (L%d):\n", cache.SocketDesignation, cache.Level())
		fmt.Printf("    Max Size:     %s\n", cache.MaximumSizeString())
		fmt.Printf("    Installed:    %s\n", cache.InstalledSizeString())
		fmt.Printf("    Type:         %s\n", cache.SystemCacheType.String())
		fmt.Printf("    Associativity: %s\n", cache.Associativity.String())
		fmt.Printf("    Location:     %s\n", cache.Configuration.Location().String())
		fmt.Printf("    Mode:         %s\n", cache.Configuration.OperationalMode().String())
		fmt.Printf("    ECC:          %s\n", cache.ErrorCorrectionType.String())
	}
	fmt.Println()
}

func printMemoryInfo(sm *gosmbios.SMBIOS) {
	// Physical Memory Array
	arrays, err := type16.GetAll(sm)
	if err == nil {
		fmt.Println("=== Physical Memory Array ===")
		for i, arr := range arrays {
			fmt.Printf("  Array %d:\n", i+1)
			fmt.Printf("    Location:     %s\n", arr.Location.String())
			fmt.Printf("    Use:          %s\n", arr.Use.String())
			fmt.Printf("    Max Capacity: %s\n", arr.MaximumCapacityString())
			fmt.Printf("    ECC:          %s\n", arr.ErrorCorrection.String())
			fmt.Printf("    Devices:      %d\n", arr.NumberOfMemoryDevices)
		}
		fmt.Println()
	}

	// Memory Devices
	devices, err := type17.GetPopulated(sm)
	if err != nil {
		fmt.Printf("Memory Devices: Not available\n\n")
		return
	}

	fmt.Println("=== Memory Devices ===")
	var totalMemory uint64
	for _, dev := range devices {
		fmt.Printf("  %s:\n", dev.DeviceLocator)
		fmt.Printf("    Size:         %s\n", dev.SizeString())
		fmt.Printf("    Type:         %s\n", dev.MemoryType.String())
		fmt.Printf("    Form Factor:  %s\n", dev.FormFactor.String())
		fmt.Printf("    Speed:        %s\n", dev.SpeedString())
		fmt.Printf("    Configured:   %d MT/s\n", dev.GetConfiguredSpeed())
		fmt.Printf("    Manufacturer: %s\n", dev.Manufacturer)
		fmt.Printf("    Part Number:  %s\n", dev.PartNumber)
		fmt.Printf("    Serial:       %s\n", dev.SerialNumber)
		fmt.Printf("    Voltage:      %s\n", dev.VoltageString())
		fmt.Printf("    Data Width:   %d bits\n", dev.DataWidth)
		fmt.Printf("    Total Width:  %d bits\n", dev.TotalWidth)
		if dev.Ranks() > 0 {
			fmt.Printf("    Ranks:        %d\n", dev.Ranks())
		}
		totalMemory += dev.Size
		fmt.Println()
	}

	fmt.Printf("  Total Memory: %d GB\n", totalMemory/1024)
	fmt.Println()
}

func printTPMInfo(sm *gosmbios.SMBIOS) {
	tpm, err := type43.Get(sm)
	if err != nil {
		fmt.Printf("TPM Information: Not available\n\n")
		return
	}

	fmt.Println("=== TPM Device ===")
	fmt.Printf("  Vendor ID:      %s\n", tpm.VendorIDString())
	fmt.Printf("  Spec Version:   %s\n", tpm.SpecVersionString())
	fmt.Printf("  Firmware:       %s\n", tpm.FirmwareVersionString())
	fmt.Printf("  Family:         %s\n", tpm.Family())
	fmt.Printf("  Description:    %s\n", tpm.Description)
	fmt.Printf("  Supported:      %v\n", tpm.IsSupported())
	fmt.Println()
}

func printBatteryInfo(sm *gosmbios.SMBIOS) {
	batteries, err := type22.GetAll(sm)
	if err != nil {
		// Battery info is only available on laptops
		return
	}

	fmt.Println("=== Portable Battery ===")
	for i, bat := range batteries {
		fmt.Printf("  Battery %d:\n", i+1)
		fmt.Printf("    Name:         %s\n", bat.DeviceName)
		fmt.Printf("    Location:     %s\n", bat.Location)
		fmt.Printf("    Manufacturer: %s\n", bat.Manufacturer)
		fmt.Printf("    Chemistry:    %s\n", bat.DeviceChemistry.String())
		fmt.Printf("    Capacity:     %s\n", bat.DesignCapacityString())
		fmt.Printf("    Voltage:      %s\n", bat.DesignVoltageString())
		fmt.Printf("    Serial:       %s\n", bat.SerialNumber)
		if bat.ManufactureDate != "" {
			fmt.Printf("    Mfg Date:     %s\n", bat.ManufactureDate)
		}
	}
	fmt.Println()
}
