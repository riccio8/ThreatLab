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
		fmt.Println("Usage: go run main.go <file> <command> [sectionName]")
		fmt.Println("Commands: lib, sym, sections, info, optionalHeaders")
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
	default:
		fmt.Println("Unknown command. Valid commands are: lib, sym, sections, info, optionalHeaders")
		return
	}

	// Pretty print the result in JSON format
	prettyPrintJSON(result)
}
