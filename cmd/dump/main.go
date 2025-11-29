// smbiosdump - Tool to dump all SMBIOS data to a file
package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

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

// getSystemIdentifier returns a unique identifier for the system (UUID or serial)
// suitable for use in a filename
func getSystemIdentifier(sm *gosmbios.SMBIOS) string {
	sys, err := type1.Get(sm)
	if err != nil {
		return ""
	}

	// Try UUID first (most unique)
	if !sys.UUID.IsZero() && !sys.UUID.IsInvalid() {
		// Return UUID without dashes for cleaner filenames
		return strings.ReplaceAll(sys.UUID.String(), "-", "")
	}

	// Fall back to serial number
	if sys.SerialNumber != "" && sys.SerialNumber != "To Be Filled By O.E.M." {
		return sanitizeFilename(sys.SerialNumber)
	}

	return ""
}

// sanitizeFilename removes/replaces characters that are invalid in filenames
func sanitizeFilename(s string) string {
	// Remove or replace characters that are problematic in filenames
	re := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)
	s = re.ReplaceAllString(s, "_")
	// Trim spaces and dots from ends
	s = strings.Trim(s, " .")
	// Limit length
	if len(s) > 64 {
		s = s[:64]
	}
	return s
}

// OutputFormat represents the output format type
type OutputFormat string

const (
	FormatText OutputFormat = "text"
	FormatJSON OutputFormat = "json"
	FormatRaw  OutputFormat = "raw"
	FormatBin  OutputFormat = "bin"
)

// SMBIOSDump represents the complete SMBIOS dump for JSON export
type SMBIOSDump struct {
	Timestamp  string          `json:"timestamp"`
	Version    string          `json:"version"`
	EntryPoint EntryPointInfo  `json:"entry_point"`
	Structures []StructureDump `json:"structures"`
	Summary    SummaryInfo     `json:"summary"`
}

// EntryPointInfo represents entry point information
type EntryPointInfo struct {
	Type         string `json:"type"`
	MajorVersion uint8  `json:"major_version"`
	MinorVersion uint8  `json:"minor_version"`
	Revision     uint8  `json:"revision,omitempty"`
	TableAddress string `json:"table_address"`
	TableLength  uint32 `json:"table_length"`
}

// StructureDump represents a single SMBIOS structure
type StructureDump struct {
	Type     uint8    `json:"type"`
	TypeName string   `json:"type_name"`
	Handle   string   `json:"handle"`
	Length   uint8    `json:"length"`
	Data     string   `json:"data"`
	Strings  []string `json:"strings,omitempty"`
}

// SummaryInfo contains summary statistics
type SummaryInfo struct {
	TotalStructures int            `json:"total_structures"`
	TypeCounts      map[string]int `json:"type_counts"`
}

func main() {
	// Command line flags
	outputFile := flag.String("o", "", "Output file path (default: stdout)")
	inputFile := flag.String("i", "", "Input file (gosmbios dump format) - read from dump instead of system")
	format := flag.String("f", "text", "Output format: text, json, raw, bin")
	showHelp := flag.Bool("h", false, "Show help")
	flag.Parse()

	if *showHelp {
		printUsage()
		os.Exit(0)
	}

	// Read SMBIOS data
	var sm *gosmbios.SMBIOS
	var err error

	if *inputFile != "" {
		sm, err = gosmbios.ReadFromFile(*inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading dump file: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "(Reading from dump file: %s)\n", *inputFile)
	} else {
		sm, err = gosmbios.Read()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading SMBIOS: %v\n", err)
			os.Exit(1)
		}
	}

	// Handle binary format specially (uses WriteToFile directly)
	if OutputFormat(strings.ToLower(*format)) == FormatBin {
		binFile := *outputFile

		// Auto-generate filename from system UUID if not specified
		if binFile == "" {
			identifier := getSystemIdentifier(sm)
			if identifier == "" {
				// Fall back to timestamp if no unique identifier
				identifier = time.Now().Format("20060102-150405")
			}
			binFile = identifier + ".smbios"
		}

		// Ensure .smbios extension
		if !strings.HasSuffix(strings.ToLower(binFile), ".smbios") {
			binFile = binFile + ".smbios"
		}

		// Create directory if needed
		dir := filepath.Dir(binFile)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
				os.Exit(1)
			}
		}

		if err := sm.WriteToFile(binFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing binary dump: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "Raw SMBIOS dump written to: %s\n", binFile)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "This file contains the raw SMBIOS table data and can be read back\n")
		fmt.Fprintf(os.Stderr, "by gosmbios tools as if reading from the system. All structure types\n")
		fmt.Fprintf(os.Stderr, "(including unknown/future types) are preserved.\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Read with: smbiosdump -i %s\n", binFile)
		fmt.Fprintf(os.Stderr, "      or:  smbiosinfo -i %s\n", binFile)
		fmt.Fprintf(os.Stderr, "      or:  smbiosdebug -i %s\n", binFile)
		return
	}

	// Determine output writer for text-based formats
	var output *os.File
	if *outputFile != "" {
		// Create directory if needed
		dir := filepath.Dir(*outputFile)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
				os.Exit(1)
			}
		}

		output, err = os.Create(*outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer output.Close()
	} else {
		output = os.Stdout
	}

	// Generate output based on format
	switch OutputFormat(strings.ToLower(*format)) {
	case FormatJSON:
		err = dumpJSON(sm, output)
	case FormatRaw:
		err = dumpRaw(sm, output)
	default:
		err = dumpText(sm, output)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}

	if *outputFile != "" {
		fmt.Fprintf(os.Stderr, "SMBIOS data written to: %s\n", *outputFile)
	}
}

func printUsage() {
	fmt.Println("smbiosdump - Dump SMBIOS data to file")
	fmt.Println()
	fmt.Println("Usage: smbiosdump [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -o <file>   Output file path (default: auto-named for bin, stdout for others)")
	fmt.Println("  -i <file>   Input file (.smbios dump) - read from dump instead of system")
	fmt.Println("  -f <format> Output format: text, json, raw, bin (default: text)")
	fmt.Println("  -h          Show this help message")
	fmt.Println()
	fmt.Println("Formats:")
	fmt.Println("  text        Human-readable text format with all structure details")
	fmt.Println("  json        JSON format with all structure data")
	fmt.Println("  raw         Raw hexadecimal dump of all structures (text-based)")
	fmt.Println("  bin         Raw binary dump - stores SMBIOS table exactly as in memory")
	fmt.Println("              Auto-names file as <UUID>.smbios if -o not specified")
	fmt.Println("              Preserves ALL data including unknown/future types")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  smbiosdump -f bin                        # Auto-named: <UUID>.smbios")
	fmt.Println("  smbiosdump -o mypc.smbios -f bin         # Custom name with .smbios extension")
	fmt.Println("  smbiosdump -o smbios.txt                 # Dump as text")
	fmt.Println("  smbiosdump -o smbios.json -f json        # Dump as JSON")
	fmt.Println("  smbiosdump -i 4C4C4544.smbios            # Read from dump file")
	fmt.Println("  smbiosdump -i dump.smbios -f json        # Convert dump to JSON")
	fmt.Println()
	fmt.Println("The binary format (-f bin) is recommended for archiving SMBIOS data.")
	fmt.Println("Files are named using the system's UUID for easy identification.")
	fmt.Println("The .smbios extension is automatically added if not present.")
}

func dumpText(sm *gosmbios.SMBIOS, w *os.File) error {
	// Header
	fmt.Fprintln(w, "================================================================================")
	fmt.Fprintln(w, "                              SMBIOS DATA DUMP")
	fmt.Fprintln(w, "================================================================================")
	fmt.Fprintf(w, "Timestamp:     %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(w, "SMBIOS Version: %s\n", sm.EntryPoint.String())
	fmt.Fprintf(w, "Entry Point:   %s\n", entryPointTypeString(sm.EntryPoint.Type))
	fmt.Fprintf(w, "Table Address: 0x%016X\n", sm.EntryPoint.TableAddress)
	fmt.Fprintf(w, "Table Length:  %d bytes\n", sm.EntryPoint.TableLength)
	fmt.Fprintf(w, "Structures:    %d\n\n", len(sm.Structures))

	// Structure summary
	typeCounts := make(map[uint8]int)
	for _, s := range sm.Structures {
		typeCounts[s.Header.Type]++
	}

	fmt.Fprintln(w, "================================================================================")
	fmt.Fprintln(w, "                            STRUCTURE SUMMARY")
	fmt.Fprintln(w, "================================================================================")
	for t := uint8(0); t <= 255; t++ {
		if count, ok := typeCounts[t]; ok {
			fmt.Fprintf(w, "Type %3d: %2d structure(s) - %s\n", t, count, types.TypeName(t))
		}
		if t == 255 {
			break
		}
	}
	fmt.Fprintln(w)

	// Detailed structure information
	fmt.Fprintln(w, "================================================================================")
	fmt.Fprintln(w, "                            STRUCTURE DETAILS")
	fmt.Fprintln(w, "================================================================================")

	// Print all structure types
	printType0Text(sm, w)
	printType1Text(sm, w)
	printType2Text(sm, w)
	printType3Text(sm, w)
	printType4Text(sm, w)
	printType5Text(sm, w)
	printType6Text(sm, w)
	printType7Text(sm, w)
	printType8Text(sm, w)
	printType9Text(sm, w)
	printType10Text(sm, w)
	printType11Text(sm, w)
	printType12Text(sm, w)
	printType13Text(sm, w)
	printType14Text(sm, w)
	printType15Text(sm, w)
	printType16Text(sm, w)
	printType17Text(sm, w)
	printType18Text(sm, w)
	printType19Text(sm, w)
	printType20Text(sm, w)
	printType21Text(sm, w)
	printType22Text(sm, w)
	printType23Text(sm, w)
	printType24Text(sm, w)
	printType25Text(sm, w)
	printType26Text(sm, w)
	printType27Text(sm, w)
	printType28Text(sm, w)
	printType29Text(sm, w)
	printType30Text(sm, w)
	printType31Text(sm, w)
	printType32Text(sm, w)
	printType33Text(sm, w)
	printType34Text(sm, w)
	printType35Text(sm, w)
	printType36Text(sm, w)
	printType37Text(sm, w)
	printType38Text(sm, w)
	printType39Text(sm, w)
	printType40Text(sm, w)
	printType41Text(sm, w)
	printType42Text(sm, w)
	printType43Text(sm, w)
	printType44Text(sm, w)
	printType45Text(sm, w)
	printType46Text(sm, w)

	// Print unknown types
	printUnknownTypesText(sm, w, typeCounts)

	fmt.Fprintln(w)
	fmt.Fprintln(w, "================================================================================")
	fmt.Fprintln(w, "                               END OF DUMP")
	fmt.Fprintln(w, "================================================================================")

	return nil
}

// Type-specific text dump functions
func printType0Text(sm *gosmbios.SMBIOS, w *os.File) {
	bios, err := type0.Get(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 0: BIOS Information ---")
	fmt.Fprintf(w, "Vendor:           %s\n", bios.Vendor)
	fmt.Fprintf(w, "Version:          %s\n", bios.Version)
	fmt.Fprintf(w, "Release Date:     %s\n", bios.ReleaseDate)
	fmt.Fprintf(w, "ROM Size:         %s\n", bios.ROMSizeString())
	fmt.Fprintf(w, "BIOS Revision:    %s\n", bios.BIOSVersionString())
	fmt.Fprintf(w, "EC Revision:      %s\n", bios.ECVersionString())
	fmt.Fprintf(w, "UEFI Capable:     %v\n", bios.IsUEFI())
	fmt.Fprintf(w, "Virtual Machine:  %v\n", bios.IsVirtualMachine())
}

func printType1Text(sm *gosmbios.SMBIOS, w *os.File) {
	sys, err := type1.Get(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 1: System Information ---")
	fmt.Fprintf(w, "Manufacturer:     %s\n", sys.Manufacturer)
	fmt.Fprintf(w, "Product Name:     %s\n", sys.ProductName)
	fmt.Fprintf(w, "Version:          %s\n", sys.Version)
	fmt.Fprintf(w, "Serial Number:    %s\n", sys.SerialNumber)
	fmt.Fprintf(w, "UUID:             %s\n", sys.UUID.String())
	fmt.Fprintf(w, "Wake-up Type:     %s\n", sys.WakeUpType.String())
	fmt.Fprintf(w, "SKU Number:       %s\n", sys.SKUNumber)
	fmt.Fprintf(w, "Family:           %s\n", sys.Family)
}

func printType2Text(sm *gosmbios.SMBIOS, w *os.File) {
	boards, err := type2.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 2: Baseboard Information ---")
	for i, board := range boards {
		if len(boards) > 1 {
			fmt.Fprintf(w, "Board %d:\n", i+1)
		}
		fmt.Fprintf(w, "Manufacturer:     %s\n", board.Manufacturer)
		fmt.Fprintf(w, "Product:          %s\n", board.Product)
		fmt.Fprintf(w, "Version:          %s\n", board.Version)
		fmt.Fprintf(w, "Serial Number:    %s\n", board.SerialNumber)
		fmt.Fprintf(w, "Asset Tag:        %s\n", board.AssetTag)
		fmt.Fprintf(w, "Type:             %s\n", board.BoardType.String())
	}
}

func printType3Text(sm *gosmbios.SMBIOS, w *os.File) {
	chassis, err := type3.Get(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 3: Chassis Information ---")
	fmt.Fprintf(w, "Manufacturer:     %s\n", chassis.Manufacturer)
	fmt.Fprintf(w, "Type:             %s\n", chassis.Type.String())
	fmt.Fprintf(w, "Version:          %s\n", chassis.Version)
	fmt.Fprintf(w, "Serial Number:    %s\n", chassis.SerialNumber)
	fmt.Fprintf(w, "Asset Tag:        %s\n", chassis.AssetTag)
	fmt.Fprintf(w, "Height:           %s\n", chassis.HeightString())
	fmt.Fprintf(w, "Power Cords:      %d\n", chassis.NumberOfPowerCords)
}

func printType4Text(sm *gosmbios.SMBIOS, w *os.File) {
	procs, err := type4.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 4: Processor Information ---")
	for i, proc := range procs {
		fmt.Fprintf(w, "Processor %d:\n", i+1)
		fmt.Fprintf(w, "  Socket:         %s\n", proc.SocketDesignation)
		fmt.Fprintf(w, "  Type:           %s\n", proc.ProcessorType.String())
		fmt.Fprintf(w, "  Family:         %s\n", proc.ProcessorFamily.String())
		fmt.Fprintf(w, "  Manufacturer:   %s\n", proc.ProcessorManufacturer)
		fmt.Fprintf(w, "  Version:        %s\n", proc.ProcessorVersion)
		fmt.Fprintf(w, "  Max Speed:      %d MHz\n", proc.MaxSpeed)
		fmt.Fprintf(w, "  Current Speed:  %d MHz\n", proc.CurrentSpeed)
		fmt.Fprintf(w, "  Core Count:     %d\n", proc.GetCoreCount())
		fmt.Fprintf(w, "  Thread Count:   %d\n", proc.GetThreadCount())
	}
}

func printType5Text(sm *gosmbios.SMBIOS, w *os.File) {
	controllers, err := type5.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 5: Memory Controller Information (Obsolete) ---")
	for i, mc := range controllers {
		if len(controllers) > 1 {
			fmt.Fprintf(w, "Controller %d:\n", i+1)
		}
		fmt.Fprintf(w, "Error Detecting:  %s\n", mc.ErrorDetectingMethod.String())
		fmt.Fprintf(w, "Interleave:       %s\n", mc.CurrentInterleave.String())
		fmt.Fprintf(w, "Max Module Size:  %d MB\n", mc.MaxModuleSizeMB())
	}
}

func printType6Text(sm *gosmbios.SMBIOS, w *os.File) {
	modules, err := type6.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 6: Memory Module Information (Obsolete) ---")
	for _, mm := range modules {
		fmt.Fprintf(w, "%s:\n", mm.SocketDesignation)
		fmt.Fprintf(w, "  Installed Size: %s\n", mm.InstalledSize.String())
		fmt.Fprintf(w, "  Memory Type:    %s\n", mm.CurrentMemoryType.String())
	}
}

func printType7Text(sm *gosmbios.SMBIOS, w *os.File) {
	caches, err := type7.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 7: Cache Information ---")
	for _, cache := range caches {
		fmt.Fprintf(w, "%s (L%d):\n", cache.SocketDesignation, cache.Level())
		fmt.Fprintf(w, "  Max Size:       %s\n", cache.MaximumSizeString())
		fmt.Fprintf(w, "  Installed Size: %s\n", cache.InstalledSizeString())
		fmt.Fprintf(w, "  Type:           %s\n", cache.SystemCacheType.String())
	}
}

func printType8Text(sm *gosmbios.SMBIOS, w *os.File) {
	ports, err := type8.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 8: Port Connector Information ---")
	for _, port := range ports {
		fmt.Fprintf(w, "%s: %s\n", port.DisplayName(), port.PortType.String())
	}
}

func printType9Text(sm *gosmbios.SMBIOS, w *os.File) {
	slots, err := type9.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 9: System Slots ---")
	for _, slot := range slots {
		fmt.Fprintf(w, "%s:\n", slot.Designation)
		fmt.Fprintf(w, "  Type:           %s\n", slot.SlotType.String())
		fmt.Fprintf(w, "  Usage:          %s\n", slot.CurrentUsage.String())
		fmt.Fprintf(w, "  Address:        %s\n", slot.PCIAddress())
	}
}

func printType10Text(sm *gosmbios.SMBIOS, w *os.File) {
	devices, err := type10.GetAllDevices(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 10: On Board Devices (Obsolete) ---")
	for _, dev := range devices {
		status := "Disabled"
		if dev.Enabled {
			status = "Enabled"
		}
		fmt.Fprintf(w, "%s: %s (%s)\n", dev.Description, dev.DeviceType.String(), status)
	}
}

func printType11Text(sm *gosmbios.SMBIOS, w *os.File) {
	oems, err := type11.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 11: OEM Strings ---")
	for _, oem := range oems {
		for i, str := range oem.Strings {
			fmt.Fprintf(w, "[%d]: %s\n", i+1, str)
		}
	}
}

func printType12Text(sm *gosmbios.SMBIOS, w *os.File) {
	configs, err := type12.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 12: System Configuration Options ---")
	for _, cfg := range configs {
		for i, opt := range cfg.Options {
			fmt.Fprintf(w, "[%d]: %s\n", i+1, opt)
		}
	}
}

func printType13Text(sm *gosmbios.SMBIOS, w *os.File) {
	lang, err := type13.Get(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 13: BIOS Language Information ---")
	fmt.Fprintf(w, "Current Language: %s\n", lang.CurrentLanguage)
	fmt.Fprintf(w, "Installable:      %d\n", lang.InstallableLanguages)
}

func printType14Text(sm *gosmbios.SMBIOS, w *os.File) {
	groups, err := type14.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 14: Group Associations ---")
	for _, grp := range groups {
		fmt.Fprintf(w, "%s: %d items\n", grp.GroupName, len(grp.Items))
	}
}

func printType15Text(sm *gosmbios.SMBIOS, w *os.File) {
	log, err := type15.Get(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 15: System Event Log ---")
	fmt.Fprintf(w, "Log Area Length:  %d bytes\n", log.LogAreaLength)
	fmt.Fprintf(w, "Access Method:    %s\n", log.AccessMethod.String())
	fmt.Fprintf(w, "Log Full:         %v\n", log.LogStatus.IsFull())
}

func printType16Text(sm *gosmbios.SMBIOS, w *os.File) {
	arrays, err := type16.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 16: Physical Memory Array ---")
	for i, arr := range arrays {
		if len(arrays) > 1 {
			fmt.Fprintf(w, "Array %d:\n", i+1)
		}
		fmt.Fprintf(w, "Location:         %s\n", arr.Location.String())
		fmt.Fprintf(w, "Use:              %s\n", arr.Use.String())
		fmt.Fprintf(w, "Error Correction: %s\n", arr.ErrorCorrection.String())
		fmt.Fprintf(w, "Max Capacity:     %s\n", arr.MaximumCapacityString())
		fmt.Fprintf(w, "Num Devices:      %d\n", arr.NumberOfMemoryDevices)
	}
}

func printType17Text(sm *gosmbios.SMBIOS, w *os.File) {
	devices, err := type17.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 17: Memory Device ---")
	for _, dev := range devices {
		fmt.Fprintf(w, "%s:\n", dev.DeviceLocator)
		fmt.Fprintf(w, "  Size:           %s\n", dev.SizeString())
		fmt.Fprintf(w, "  Form Factor:    %s\n", dev.FormFactor.String())
		fmt.Fprintf(w, "  Type:           %s\n", dev.MemoryType.String())
		fmt.Fprintf(w, "  Speed:          %s\n", dev.SpeedString())
		fmt.Fprintf(w, "  Manufacturer:   %s\n", dev.Manufacturer)
		fmt.Fprintf(w, "  Part Number:    %s\n", dev.PartNumber)
	}
}

func printType18Text(sm *gosmbios.SMBIOS, w *os.File) {
	errors, err := type18.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 18: 32-Bit Memory Error Information ---")
	for i, me := range errors {
		fmt.Fprintf(w, "Error %d: %s\n", i+1, me.ErrorType.String())
	}
}

func printType19Text(sm *gosmbios.SMBIOS, w *os.File) {
	maps, err := type19.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 19: Memory Array Mapped Address ---")
	for _, m := range maps {
		fmt.Fprintf(w, "Array 0x%04X: %s\n", m.MemoryArrayHandle, m.GetSizeString())
	}
}

func printType20Text(sm *gosmbios.SMBIOS, w *os.File) {
	maps, err := type20.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 20: Memory Device Mapped Address ---")
	for _, m := range maps {
		fmt.Fprintf(w, "Device 0x%04X: 0x%X - 0x%X\n", m.MemoryDeviceHandle, m.GetStartingAddressBytes(), m.GetEndingAddressBytes())
	}
}

func printType21Text(sm *gosmbios.SMBIOS, w *os.File) {
	devices, err := type21.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 21: Built-in Pointing Device ---")
	for _, dev := range devices {
		fmt.Fprintf(w, "Type: %s, Interface: %s, Buttons: %d\n", dev.DeviceType.String(), dev.Interface.String(), dev.NumberOfButtons)
	}
}

func printType22Text(sm *gosmbios.SMBIOS, w *os.File) {
	batteries, err := type22.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 22: Portable Battery ---")
	for _, bat := range batteries {
		fmt.Fprintf(w, "%s:\n", bat.DeviceName)
		fmt.Fprintf(w, "  Location:       %s\n", bat.Location)
		fmt.Fprintf(w, "  Chemistry:      %s\n", bat.DeviceChemistry.String())
		fmt.Fprintf(w, "  Capacity:       %s\n", bat.DesignCapacityString())
	}
}

func printType23Text(sm *gosmbios.SMBIOS, w *os.File) {
	rst, err := type23.Get(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 23: System Reset ---")
	fmt.Fprintf(w, "Enabled: %v, Watchdog: %v\n", rst.Capabilities.IsEnabled(), rst.Capabilities.WatchdogTimerPresent())
}

func printType24Text(sm *gosmbios.SMBIOS, w *os.File) {
	sec, err := type24.Get(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 24: Hardware Security ---")
	fmt.Fprintf(w, "Power-on Password:  %s\n", sec.HardwareSettings.PowerOnPasswordStatus().String())
	fmt.Fprintf(w, "Admin Password:     %s\n", sec.HardwareSettings.AdministratorPasswordStatus().String())
}

func printType25Text(sm *gosmbios.SMBIOS, w *os.File) {
	pwr, err := type25.Get(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 25: System Power Controls ---")
	fmt.Fprintf(w, "Next Power On: %s\n", pwr.NextPowerOnString())
}

func printType26Text(sm *gosmbios.SMBIOS, w *os.File) {
	probes, err := type26.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 26: Voltage Probe ---")
	for _, probe := range probes {
		fmt.Fprintf(w, "%s: %s, Status: %s\n", probe.Description, probe.LocationAndStatus.Location().String(), probe.LocationAndStatus.Status().String())
	}
}

func printType27Text(sm *gosmbios.SMBIOS, w *os.File) {
	devices, err := type27.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 27: Cooling Device ---")
	for _, dev := range devices {
		fmt.Fprintf(w, "%s: %s, Speed: %s\n", dev.Description, dev.DeviceTypeAndStatus.DeviceType().String(), dev.NominalSpeedString())
	}
}

func printType28Text(sm *gosmbios.SMBIOS, w *os.File) {
	probes, err := type28.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 28: Temperature Probe ---")
	for _, probe := range probes {
		fmt.Fprintf(w, "%s: %s, Status: %s\n", probe.Description, probe.LocationAndStatus.Location().String(), probe.LocationAndStatus.Status().String())
	}
}

func printType29Text(sm *gosmbios.SMBIOS, w *os.File) {
	probes, err := type29.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 29: Electrical Current Probe ---")
	for _, probe := range probes {
		fmt.Fprintf(w, "%s: %s, Status: %s\n", probe.Description, probe.LocationAndStatus.Location().String(), probe.LocationAndStatus.Status().String())
	}
}

func printType30Text(sm *gosmbios.SMBIOS, w *os.File) {
	oob, err := type30.Get(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 30: Out-of-Band Remote Access ---")
	fmt.Fprintf(w, "Manufacturer: %s\n", oob.ManufacturerName)
	fmt.Fprintf(w, "Inbound: %v, Outbound: %v\n", oob.Connections.InboundEnabled(), oob.Connections.OutboundEnabled())
}

func printType31Text(sm *gosmbios.SMBIOS, w *os.File) {
	bis, err := type31.Get(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 31: Boot Integrity Services Entry Point ---")
	fmt.Fprintf(w, "Entry Point: 0x%08X\n", bis.BISEntryPoint)
}

func printType32Text(sm *gosmbios.SMBIOS, w *os.File) {
	boot, err := type32.Get(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 32: System Boot Information ---")
	fmt.Fprintf(w, "Status: %s\n", boot.BootStatus.String())
}

func printType33Text(sm *gosmbios.SMBIOS, w *os.File) {
	errors, err := type33.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 33: 64-Bit Memory Error Information ---")
	for i, me := range errors {
		fmt.Fprintf(w, "Error %d: %s\n", i+1, me.ErrorType.String())
	}
}

func printType34Text(sm *gosmbios.SMBIOS, w *os.File) {
	devices, err := type34.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 34: Management Device ---")
	for _, dev := range devices {
		fmt.Fprintf(w, "%s: %s, Address: 0x%08X\n", dev.Description, dev.DeviceType.String(), dev.Address)
	}
}

func printType35Text(sm *gosmbios.SMBIOS, w *os.File) {
	components, err := type35.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 35: Management Device Component ---")
	for _, comp := range components {
		fmt.Fprintf(w, "%s: Device 0x%04X, Component 0x%04X\n", comp.Description, comp.ManagementDeviceHandle, comp.ComponentHandle)
	}
}

func printType36Text(sm *gosmbios.SMBIOS, w *os.File) {
	thresholds, err := type36.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 36: Management Device Threshold Data ---")
	for i := range thresholds {
		fmt.Fprintf(w, "Threshold %d present\n", i+1)
	}
}

func printType37Text(sm *gosmbios.SMBIOS, w *os.File) {
	channels, err := type37.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 37: Memory Channel ---")
	for i, ch := range channels {
		fmt.Fprintf(w, "Channel %d: %s, Devices: %d\n", i+1, ch.ChannelType.String(), ch.MemoryDeviceCount)
	}
}

func printType38Text(sm *gosmbios.SMBIOS, w *os.File) {
	ipmi, err := type38.Get(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 38: IPMI Device Information ---")
	fmt.Fprintf(w, "Interface Type: %s\n", ipmi.InterfaceType.String())
	fmt.Fprintf(w, "Spec Revision:  %s\n", ipmi.SpecificationRevisionString())
	fmt.Fprintf(w, "Base Address:   %s\n", ipmi.BaseAddressString())
}

func printType39Text(sm *gosmbios.SMBIOS, w *os.File) {
	supplies, err := type39.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 39: System Power Supply ---")
	for _, psu := range supplies {
		fmt.Fprintf(w, "%s:\n", psu.DeviceName)
		fmt.Fprintf(w, "  Location:       %s\n", psu.Location)
		fmt.Fprintf(w, "  Manufacturer:   %s\n", psu.Manufacturer)
		fmt.Fprintf(w, "  Max Power:      %s\n", psu.MaxPowerCapacityString())
	}
}

func printType40Text(sm *gosmbios.SMBIOS, w *os.File) {
	info, err := type40.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 40: Additional Information ---")
	for i, ai := range info {
		fmt.Fprintf(w, "Info %d: %d entries\n", i+1, ai.NumberOfEntries)
	}
}

func printType41Text(sm *gosmbios.SMBIOS, w *os.File) {
	devices, err := type41.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 41: Onboard Devices Extended Information ---")
	for _, dev := range devices {
		fmt.Fprintf(w, "%s: %s, Status: %s, Address: %s\n", dev.ReferenceDesignation, dev.TypeString(), dev.StatusString(), dev.PCIAddress())
	}
}

func printType42Text(sm *gosmbios.SMBIOS, w *os.File) {
	mchis, err := type42.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 42: Management Controller Host Interface ---")
	for i, mchi := range mchis {
		fmt.Fprintf(w, "Interface %d: %s, Protocols: %d\n", i+1, mchi.InterfaceType.String(), len(mchi.ProtocolRecords))
	}
}

func printType43Text(sm *gosmbios.SMBIOS, w *os.File) {
	tpm, err := type43.Get(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 43: TPM Device ---")
	fmt.Fprintf(w, "Vendor ID:     %s\n", tpm.VendorIDString())
	fmt.Fprintf(w, "Spec Version:  %s\n", tpm.SpecVersionString())
	fmt.Fprintf(w, "Firmware:      %s\n", tpm.FirmwareVersionString())
	fmt.Fprintf(w, "Family:        %s\n", tpm.Family())
}

func printType44Text(sm *gosmbios.SMBIOS, w *os.File) {
	infos, err := type44.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 44: Processor Additional Information ---")
	for i, info := range infos {
		fmt.Fprintf(w, "Info %d: Handle 0x%04X, Type: %s\n", i+1, info.ReferencedHandle, info.ProcessorSpecificBlock.ProcessorType.String())
	}
}

func printType45Text(sm *gosmbios.SMBIOS, w *os.File) {
	firmwares, err := type45.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 45: Firmware Inventory Information ---")
	for _, fw := range firmwares {
		fmt.Fprintf(w, "%s:\n", fw.FirmwareComponentName)
		fmt.Fprintf(w, "  Version:        %s\n", fw.FirmwareVersion)
		fmt.Fprintf(w, "  Manufacturer:   %s\n", fw.Manufacturer)
		fmt.Fprintf(w, "  State:          %s\n", fw.State.String())
	}
}

func printType46Text(sm *gosmbios.SMBIOS, w *os.File) {
	props, err := type46.GetAll(sm)
	if err != nil {
		return
	}
	fmt.Fprintln(w, "\n--- Type 46: String Property ---")
	for _, prop := range props {
		fmt.Fprintf(w, "%s: %s (Parent: 0x%04X)\n", prop.StringPropertyID.String(), prop.StringPropertyValue, prop.ParentHandle)
	}
}

func printUnknownTypesText(sm *gosmbios.SMBIOS, w *os.File, typeCounts map[uint8]int) {
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

	fmt.Fprintln(w, "\n--- Unknown/OEM Types ---")
	for t := uint8(0); t <= 255; t++ {
		if count, ok := typeCounts[t]; ok && !handled[t] {
			fmt.Fprintf(w, "Type %3d: %d structure(s) - %s\n", t, count, types.TypeName(t))
			// Print raw data for unknown types
			structs := sm.GetStructures(t)
			for i, s := range structs {
				fmt.Fprintf(w, "  [%d] Handle: 0x%04X, Length: %d\n", i, s.Header.Handle, s.Header.Length)
				fmt.Fprintf(w, "      Data: %s\n", hex.EncodeToString(s.Data))
				if len(s.Strings) > 0 {
					for j, str := range s.Strings {
						fmt.Fprintf(w, "      String[%d]: %q\n", j+1, str)
					}
				}
			}
		}
		if t == 255 {
			break
		}
	}
}

func dumpJSON(sm *gosmbios.SMBIOS, w *os.File) error {
	// Build type counts
	typeCounts := make(map[string]int)
	for _, s := range sm.Structures {
		key := fmt.Sprintf("type_%d", s.Header.Type)
		typeCounts[key]++
	}

	// Build structure list
	structures := make([]StructureDump, 0, len(sm.Structures))
	for _, s := range sm.Structures {
		dump := StructureDump{
			Type:     s.Header.Type,
			TypeName: types.TypeName(s.Header.Type),
			Handle:   fmt.Sprintf("0x%04X", s.Header.Handle),
			Length:   s.Header.Length,
			Data:     hex.EncodeToString(s.Data),
			Strings:  s.Strings,
		}
		structures = append(structures, dump)
	}

	// Build complete dump
	epType := "32-bit"
	if sm.EntryPoint.Type == gosmbios.EntryPoint64Bit {
		epType = "64-bit"
	}

	dump := SMBIOSDump{
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   sm.EntryPoint.String(),
		EntryPoint: EntryPointInfo{
			Type:         epType,
			MajorVersion: sm.EntryPoint.MajorVersion,
			MinorVersion: sm.EntryPoint.MinorVersion,
			Revision:     sm.EntryPoint.Revision,
			TableAddress: fmt.Sprintf("0x%016X", sm.EntryPoint.TableAddress),
			TableLength:  sm.EntryPoint.TableLength,
		},
		Structures: structures,
		Summary: SummaryInfo{
			TotalStructures: len(sm.Structures),
			TypeCounts:      typeCounts,
		},
	}

	// Marshal to JSON with indentation
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(dump)
}

func dumpRaw(sm *gosmbios.SMBIOS, w *os.File) error {
	// Write header comment
	fmt.Fprintf(w, "# SMBIOS Raw Dump\n")
	fmt.Fprintf(w, "# Timestamp: %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(w, "# Version: %s\n", sm.EntryPoint.String())
	fmt.Fprintf(w, "# Structures: %d\n\n", len(sm.Structures))

	// Write each structure
	for i, s := range sm.Structures {
		fmt.Fprintf(w, "# Structure %d: Type %d (%s), Handle 0x%04X, Length %d\n",
			i, s.Header.Type, types.TypeName(s.Header.Type), s.Header.Handle, s.Header.Length)

		// Write hex data
		for j := 0; j < len(s.Data); j += 16 {
			end := j + 16
			if end > len(s.Data) {
				end = len(s.Data)
			}
			fmt.Fprintln(w, hex.EncodeToString(s.Data[j:end]))
		}

		// Write strings as comments
		if len(s.Strings) > 0 {
			fmt.Fprintf(w, "# Strings: %d\n", len(s.Strings))
			for j, str := range s.Strings {
				fmt.Fprintf(w, "# [%d]: %s\n", j+1, str)
			}
		}
		fmt.Fprintln(w)
	}

	return nil
}

func entryPointTypeString(t gosmbios.EntryPointType) string {
	if t == gosmbios.EntryPoint64Bit {
		return "64-bit (SMBIOS 3.x)"
	}
	return "32-bit (SMBIOS 2.x)"
}
