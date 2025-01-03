/*
 * Copyright 2023-2024 Riccardo Adami. All rights reserved.
 * License: https://github.com/riccio8/ThreatLab/blob/main/LICENSE
 */

 
package main

import (
	"debug/pe"
	"encoding/json"
	"fmt"
	"os"
)

func load(fileName string) (*pe.File, error) {
	file, err := pe.Open(fileName)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func lib(f *pe.File) ([]string, error) {
	libs, err := f.ImportedLibraries()
	if err != nil {
		return nil, err
	}
	return libs, nil
}

func sym(f *pe.File) ([]string, error) {
	symbols, err := f.ImportedSymbols()
	if err != nil {
		return nil, err
	}
	return symbols, nil
}

func sections(f *pe.File, name string) *pe.Section {
	return f.Section(name)
}

func info(f *pe.File) pe.FileHeader {
	return f.FileHeader
}

func optionalHeaders(f *pe.File) any {
	return f.OptionalHeader
}
func theFile(f *pe.File) any {
	return f.FileHeader
}

func coffNums(f *pe.File) any {
	return f.COFFSymbols
}

func machine(f *pe.File) any {
    return f.Machine 
}

func stringsTable(f *pe.File) any {
	return f.StringTable
}

func dates(f *pe.File) any {
	return f.TimeDateStamp
}

func dwarf(f *pe.File) any {
	res, err := f.DWARF()
	if err != nil {
		return err
	}
	return res
}

func pointerSymTables(f *pe.File) any {
	return f.PointerToSymbolTable
}

func Characteristics(f *pe.File) any {
    return f.Characteristics
}



func strings(f *pe.File) ([]string, error) {
	var strTable []string


	for _, section := range f.Sections {
		if section.Name == ".strtab" {  // tipic name
			data, err := section.Data()  
			if err != nil {
				return nil, err
			}
			start := 0
			for start < len(data) {
				end := start
				//(null terminator)
				for end < len(data) && data[end] != 0 {
					end++
				}
				if end > start {
					strTable = append(strTable, string(data[start:end]))
				}
				// next block
				start = end + 1
			}
			break
		}
	}

	if len(strTable) == 0 {
		return nil, fmt.Errorf("no string table found")
	}
	return strTable, nil
}



func prettyPrintJSON(v interface{}) {
	// Format the result as JSON with indentation for pretty printing
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("Error formatting output as JSON: %v\n", err)
		return
	}
	fmt.Println(string(data))
}


func help() {
	// ANSI color codes
	reset := "\033[0m"
	red := "\033[31m"
	green := "\033[32m"
	yellow := "\033[33m"
	blue := "\033[34m"
	cyan := "\033[36m"

	fmt.Printf("%sUsage:%s\n", yellow, reset)
	fmt.Printf("  %s<tool_name>%s %s<file>%s %s<command>%s [%s<sectionName>%s]\n", cyan, reset, green, reset, red, reset, blue, reset)

	fmt.Printf("\n%sCommands:%s\n", yellow, reset)
	fmt.Printf("  %slib%s               - List libraries required by the PE file.\n", green, reset)
	fmt.Printf("  %ssym%s               - Display symbols table.\n", green, reset)
	fmt.Printf("  %ssections%s          - Show details of a specific section (use with [sectionName]).\n", green, reset)
	fmt.Printf("  %sinfo%s              - Display basic information about the PE file.\n", green, reset)
	fmt.Printf("  %soptionalHeaders%s   - Show the optional headers of the PE file.\n", green, reset)
	fmt.Printf("  %sfileHeader%s        - Display the PE file header.\n", green, reset)
	fmt.Printf("  %scoffSymbols%s       - Display COFF symbols.\n", green, reset)
	fmt.Printf("  %smachine%s           - Show the target machine type.\n", green, reset)
	fmt.Printf("  %sstringTable%s       - Display the string table.\n", green, reset)
	fmt.Printf("  %stime%s              - Show the timestamp of the PE file.\n", green, reset)
	fmt.Printf("  %sdwarf%s             - Extract DWARF debug data (if available).\n", green, reset)
	fmt.Printf("  %spointerSymTables%s  - Show pointers to symbol tables.\n", green, reset)
	fmt.Printf("  %scharacteristics%s   - Display characteristics of the PE file.\n", green, reset)
	fmt.Printf("  %sstring%s            - Extract strings from the PE file.\n", green, reset)

	fmt.Printf("\n%sExample:%s\n", yellow, reset)
	fmt.Printf("  %s./peutils.exe%s %smyfile.exe%s %sinfo%s\n", cyan, reset, green, reset, red, reset)
	fmt.Printf("  %s./peutils.exe%s %smyfile.exe%s %ssections%s %s.text%s\n", cyan, reset, green, reset, red, reset, blue, reset)

	fmt.Printf("\n%sNote:%s For more details on specific fields, refer to the documentation or PE specification.\n", yellow, reset)
	fmt.Printf("  %shttps://learn.microsoft.com/en-us/windows/win32/debug/pe-format%s\n", cyan, reset)
}

func main() {
	if len(os.Args) < 3 {
		help()
		return
	}

	fileName := os.Args[1]
	command := os.Args[2]
	sectionName := ""
	if command == "sections" && len(os.Args) >= 4 {
		sectionName = os.Args[3]
	}

	peFile, err := load(fileName)
	if err != nil {
		fmt.Printf("Error loading file: %v\n", err)
		return
	}
	defer peFile.Close()

	var result interface{}
	switch command {
	case "lib":
		libs, err := lib(peFile)
		if err != nil {
			fmt.Printf("Error fetching libraries: %v\n", err)
			return
		}
		result = libs
	case "sym":
		symbols, err := sym(peFile)
		if err != nil {
			fmt.Printf("Error fetching symbols: %v\n", err)
			return
		}
		result = symbols
	case "sections":
		if sectionName == "" {
			fmt.Println("Please specify a section name for the 'sections' command.")
			return
		}
		section := sections(peFile, sectionName)
		if section == nil {
			fmt.Printf("Section %s not found.\n", sectionName)
			return
		}
		result = section
	case "info":
		header := info(peFile)
		result = header
	case "optionalHeaders":
		headers := optionalHeaders(peFile)
		result = headers
	case "fileHeader":
		result = theFile(peFile)
	case "coffSymbols":
		result = coffNums(peFile)
	case "machine":
		result = machine(peFile)
	case "stringTable":
		result = stringsTable(peFile)
	case "time":
		result = dates(peFile)
	case "dwarf":
		result = dwarf(peFile)
	case "pointerSymTables":
		result = pointerSymTables(peFile)
	case "characteristics":
		result = Characteristics(peFile)
	case "string":
		result, err = strings(peFile)
		if err!= nil {
            fmt.Printf("Error fetching strings: %v\n", err)
            return
        }
	default:
		fmt.Println("Unknown command. Valid commands are: lib, sym, sections, info, optionalHeaders, fileHeader, coffSymbols, machine, stringTable, timeDateStamp, dwarf, pointerSymTables, characteristics")
		help()
		return
	}

	// Pretty print the result in JSON format
	prettyPrintJSON(result)
	
	file, err := os.Create(fileName+".json")
	if err!= nil {
            panic(err)
   	 }
   	 defer file.Close()
    
	bs, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
	    panic(err)
	}
	file.Write(bs)
}
