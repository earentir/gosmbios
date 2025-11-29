// smbiosdump - Tool to dump all SMBIOS data to a file
package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/earentir/gosmbios"
	"github.com/earentir/gosmbios/types"
)

// OutputFormat represents the output format type
type OutputFormat string

const (
	FormatText OutputFormat = "text"
	FormatJSON OutputFormat = "json"
	FormatRaw  OutputFormat = "raw"
)

// SMBIOSDump represents the complete SMBIOS dump for JSON export
type SMBIOSDump struct {
	Timestamp   string           `json:"timestamp"`
	Version     string           `json:"version"`
	EntryPoint  EntryPointInfo   `json:"entry_point"`
	Structures  []StructureDump  `json:"structures"`
	Summary     SummaryInfo      `json:"summary"`
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
	format := flag.String("f", "text", "Output format: text, json, raw")
	showHelp := flag.Bool("h", false, "Show help")
	flag.Parse()

	if *showHelp {
		printUsage()
		os.Exit(0)
	}

	// Read SMBIOS data
	sm, err := gosmbios.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading SMBIOS: %v\n", err)
		os.Exit(1)
	}

	// Determine output writer
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
	fmt.Println("  -o <file>   Output file path (default: stdout)")
	fmt.Println("  -f <format> Output format: text, json, raw (default: text)")
	fmt.Println("  -h          Show this help message")
	fmt.Println()
	fmt.Println("Formats:")
	fmt.Println("  text        Human-readable text format")
	fmt.Println("  json        JSON format with all structure data")
	fmt.Println("  raw         Raw hexadecimal dump of all structures")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  smbiosdump -o smbios.txt")
	fmt.Println("  smbiosdump -o smbios.json -f json")
	fmt.Println("  smbiosdump -f raw > smbios.hex")
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

	// All structures
	fmt.Fprintln(w, "================================================================================")
	fmt.Fprintln(w, "                            STRUCTURE DETAILS")
	fmt.Fprintln(w, "================================================================================")

	for i, s := range sm.Structures {
		fmt.Fprintf(w, "\n--- Structure %d: Type %d (%s) ---\n", i, s.Header.Type, types.TypeName(s.Header.Type))
		fmt.Fprintf(w, "Handle: 0x%04X\n", s.Header.Handle)
		fmt.Fprintf(w, "Length: %d bytes\n", s.Header.Length)
		fmt.Fprintf(w, "Data Length: %d bytes\n", len(s.Data))

		// Hex dump of data
		fmt.Fprintln(w, "Data (hex):")
		for j := 0; j < len(s.Data); j += 16 {
			end := j + 16
			if end > len(s.Data) {
				end = len(s.Data)
			}
			// Print offset
			fmt.Fprintf(w, "  %04X: ", j)
			// Print hex bytes
			hexStr := hex.EncodeToString(s.Data[j:end])
			for k := 0; k < len(hexStr); k += 2 {
				fmt.Fprintf(w, "%s ", hexStr[k:k+2])
			}
			// Pad if needed
			for k := end - j; k < 16; k++ {
				fmt.Fprint(w, "   ")
			}
			// Print ASCII
			fmt.Fprint(w, " |")
			for k := j; k < end; k++ {
				if s.Data[k] >= 32 && s.Data[k] < 127 {
					fmt.Fprintf(w, "%c", s.Data[k])
				} else {
					fmt.Fprint(w, ".")
				}
			}
			fmt.Fprintln(w, "|")
		}

		// Strings
		if len(s.Strings) > 0 {
			fmt.Fprintln(w, "Strings:")
			for j, str := range s.Strings {
				fmt.Fprintf(w, "  [%d]: %q\n", j+1, str)
			}
		}
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "================================================================================")
	fmt.Fprintln(w, "                               END OF DUMP")
	fmt.Fprintln(w, "================================================================================")

	return nil
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
