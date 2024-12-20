package main

import (
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

func sections(f *elf.File) any{
	return f.Sections
}

func symbolTable(f *elf.File) (any, error) {
	symbols, err := f.Symbols()
	if err != nil {
		return nil, err
	}
	return symbols, nil
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

func relocs(f *elf.File) *elf.Section{
	return f.Section(".strtab")
}

func sectionsInfo(f *elf.File) []string {
	var sectionsInfo []string
	for _, section := range f.Sections {
		sectionsInfo = append(sectionsInfo, fmt.Sprintf("Section name: %s, Type: %s", section.Name, section.Type))
	}
	return sectionsInfo
}

func ImportedSymbols(f *elf.File) any{
	return f.ImportedSymbols
}

func symbols(f *elf.File) any{
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

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: elfutils.exe <file> <command> [sectionName]")
		fmt.Println("Commands: sections, sym, class, symbols, headers, machine, entryPoint, importSym, fileHeader, stringTable, dynamicSymbols, relocs, sectionsInfo")
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
		result = ImportedSymbols(elfFile)
	case "stringTable":
		result = stringTable(elfFile)
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
		fmt.Println("Unknown command. Valid commands are: sections, sym, class, machine, entryPoint, fileHeader, stringTable, dynamicSymbols, relocs, sectionsInfo")
		return
	}

	// Pretty print the result in JSON format
	prettyPrintJSON(result)
}
