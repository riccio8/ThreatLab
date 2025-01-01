/*
 * Copyright 2023-2024 Riccardo Adami. All rights reserved.
 * License: https://github.com/riccio8/ThreatLab/blob/main/LICENSE
 */



package main

import (
	"debug/dwarf"
	"debug/elf"
	"encoding/json"
	"fmt"
	"os"
)

func load(fileName string) (*elf.File, error) {
	file, err := elf.Open(fileName)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func DWARF(f *elf.File) (*dwarf.Data, error) {
	data, err := f.DWARF() 
	if err != nil {
		return nil, err
	}
	return data, nil
}

func sections(f *elf.File) any {
	return f.Sections
}

func symbolTable(f *elf.File) (any, error) {
	symbols, err := f.Symbols()
	if err != nil {
		return nil, err
	}
	return symbols, nil
}

func lib(f *elf.File) (any, error) {
	return f.ImportedLibraries()
}

func class(f *elf.File) string {
	return f.Class.String()
}

func machine(f *elf.File) string {
	return f.Machine.String()
}

func entryPoint(f *elf.File) uint64 {
	return f.Entry
}

func sectionByName(f *elf.File, name string) any {
	for _, section := range f.Sections {
		if section.Name == name {
			return &section
		}
	}
	return nil
}

func headers(f *elf.File) any {
	fmt.Println("\nProgram Headers:")
	for _, progHeader := range f.Progs {
		return fmt.Sprintf("Type: %v, Offset: 0x%x, Virtual Address: 0x%x\n", progHeader.Type, progHeader.Off, progHeader.Vaddr, progHeader.ProgHeader, progHeader.Flags, progHeader.Memsz, progHeader.Filesz)
	}
	return f.Progs
}

func fileHeader(f *elf.File) any {
	return f.FileHeader
}

func stringTable(f *elf.File) *elf.Section {
	return f.Section(".strtab")
}

func dynamicSymbols(f *elf.File) ([]elf.Symbol, error) {
	symbols, err := f.DynamicSymbols()
	if err != nil {

		return nil, err
	}
	return symbols, nil
}

func relocs(f *elf.File) *elf.Section {
	return f.Section(".strtab")
}

func sectionsInfo(f *elf.File) []string {
	var sectionsInfo []string
	for _, section := range f.Sections {
		sectionsInfo = append(sectionsInfo, fmt.Sprintf("Section name: %s, Type: %s", section.Name, section.Type))
	}
	return sectionsInfo
}

func ImportedSymbols(f *elf.File) (any, error) {
	sym, err := f.ImportedSymbols()
	if err!= nil {
        return nil, err
    }
    return sym, nil
}

func symbols(f *elf.File) any {
	sym, err := f.Symbols()
	if err != nil {
		fmt.Println("Error while getting data", err)
		return err
	}
	if len(sym) > 0 {
		fmt.Println("\nSymbols:")
		for _, symbol := range sym {
			return fmt.Sprintf("Symbol Name: %s, Value: 0x%x\n", symbol.Name, symbol.Value)
		}
	} else {
		return fmt.Sprintf("No symbols found")
	}
	return nil
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
	fmt.Printf("  %ssections%s       - List all sections or a specific section by name (use with [sectionName]).\n", green, reset)
	fmt.Printf("  %ssym%s            - List the symbol table.\n", green, reset)
	fmt.Printf("  %sclass%s          - Show the ELF class (e.g., 32-bit, 64-bit).\n", green, reset)
	fmt.Printf("  %ssymbols%s        - Show all symbols in the ELF file.\n", green, reset)
	fmt.Printf("  %sdwarf%s          - Extract DWARF debug data.\n", green, reset)
	fmt.Printf("  %smachine%s        - Display machine architecture details.\n", green, reset)
	fmt.Printf("  %sentryPoint%s     - Show the program's entry point address.\n", green, reset)
	fmt.Printf("  %sfileHeader%s     - Display the ELF file header.\n", green, reset)
	fmt.Printf("  %sheaders%s        - Show all headers in the ELF file.\n", green, reset)
	fmt.Printf("  %simportSym%s      - List imported symbols.\n", green, reset)
	fmt.Printf("  %sstringTable%s    - Show the string table.\n", green, reset)
	fmt.Printf("  %slib%s            - List dynamic libraries required by the ELF file.\n", green, reset)
	fmt.Printf("  %sdynamicSymbols%s - Show dynamic symbols in the ELF file.\n", green, reset)
	fmt.Printf("  %srelocs%s         - List relocation entries.\n", green, reset)
	fmt.Printf("  %ssectionsInfo%s   - Display detailed section information.\n", green, reset)

	fmt.Printf("\n%sExample:%s\n", yellow, reset)
	fmt.Printf("  %s./elfutils.exe%s %smyfile.elf%s %ssections%s\n", cyan, reset, green, reset, red, reset)
	fmt.Printf("  %s./elfutils.exe%s %smyfile.elf%s %ssections%s %s.text%s\n", cyan, reset, green, reset, red, reset, blue, reset)

	fmt.Printf("\n%sNote:%s For more details on specific fields, refer to the documentation:\n", yellow, reset)
	fmt.Printf("  %shttps://pkg.go.dev/debug/elf%s\n", cyan, reset)
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

	elfFile, err := load(fileName)
	if err != nil {
		fmt.Printf("Error loading file: %v\n", err)
		return
	}
	defer elfFile.Close()

	var result interface{}
	switch command {
	case "sections":
		if sectionName == "" {
			sections := sections(elfFile)
			result = sections
		} else {
			section := sectionByName(elfFile, sectionName)
			if section == nil {
				fmt.Printf("Section %s not found.\n", sectionName)
				return
			}
			result = section
		}
	case "sym":
		symbols, err := symbolTable(elfFile)
		if err != nil {
			fmt.Printf("Error fetching symbols: %v\n", err)
			return
		}
		result = symbols
	case "class":
		result = class(elfFile)
	case "symbols":
		result = symbols(elfFile)
	case "dwarf":
		result, err = DWARF(elfFile)
		if err!= nil {
            fmt.Printf("Error fetching DWARF data: %v\n", err)
            return
        }
	case "machine":
		result = machine(elfFile)
		fmt.Sprintln("for further infos look at: https://pkg.go.dev/debug/elf#Machine")
	case "entryPoint":
		result = entryPoint(elfFile)
	case "fileHeader":
		result = fileHeader(elfFile)
	case "headers":
		result = headers(elfFile)
	case "importSym":
		result, err = ImportedSymbols(elfFile)
		if err!= nil {
            fmt.Printf("Error fetching imported symbols: %v\n", err)
            return
        }
	case "stringTable":
		result = stringTable(elfFile)
	case "lib":
		result, err = lib(elfFile)
		if err != nil {
			fmt.Printf("Error fetching libraries: %v\n", err)
			return
		}
	case "dynamicSymbols":
		symbols, err := dynamicSymbols(elfFile)
		if err != nil {
			fmt.Printf("Error fetching dynamic symbols: %v\n", err)
			return
		}
		result = symbols
	case "relocs":
		result = relocs(elfFile)
	case "sectionsInfo":
		result = sectionsInfo(elfFile)
	default:
		fmt.Println("Unknown command. Valid commands are: dwarf, lib, sections, sym, class, machine, entryPoint, fileHeader, stringTable, dynamicSymbols, relocs, sectionsInfo")
		help()
		return
	}

	// Pretty print the result in JSON format
	prettyPrintJSON(result)
}
