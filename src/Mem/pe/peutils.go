/*
 * Copyright 2023-2024 Riccardo Adami. All rights reserved.
 * License: https://github.com/riccio8/ThreatLab/blob/main/LICENSE
 */


package main

import (
	"debug/pe"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"runtime"

	"gopkg.in/yaml.v3"
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

func stringss(f *pe.File) ([]string, error) {
	var strTable []string

	for _, section := range f.Sections {
		if section.Name == ".strtab" { // tipic name
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

func saveResult(fileName string, result interface{}, format string) (int, error) {
	so := runtime.GOOS

	switch format {
	case "json":
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return 0, err
		}
		if so == "linux" {
			file, err := os.Create("/var/log/" + fileName + ".json")
			if err != nil {
				panic(err)
			}
			defer file.Close()

			fmt.Println("Logged successfully at /var/log/", fileName+".json")
			return file.Write(data)

		} else if so == "windows" {

			file, err := os.Create(fileName + ".json")
			if err != nil {
				panic(err)
			}
			defer file.Close()

			fmt.Println("Logged successfully at current directory: ", fileName+".json")
			return file.Write(data)
		} else {
			err = fmt.Errorf("Unsupported operating system for logging.")
			return 0, err
		}

	case "xml":
		data, err := xml.MarshalIndent(result, "", "  ")
		if err != nil {
			return 0, err
		}
		if so == "linux" {
			file, err := os.Create("/var/log/" + fileName + ".xml")
			if err != nil {
				panic(err)
			}
			defer file.Close()

			fmt.Println("Logged successfully at /var/log/", fileName+".xml")
			return file.Write(data)

		} else if so == "windows" {

			file, err := os.Create(fileName + ".xml")
			if err != nil {
				panic(err)
			}
			defer file.Close()

			fmt.Println("Logged successfully at current directory: ", fileName+".xml")
			return file.Write(data)
		} else {
			err = fmt.Errorf("Unsupported operating system for logging.")
			return 0, err
		}

	case "yaml":
		data, err := yaml.Marshal(result)
		if err != nil {
			return 0, err
		}
		if so == "linux" {
			file, err := os.Create("/var/log/" + fileName + ".yaml")
			if err != nil {
				panic(err)
			}
			defer file.Close()

			fmt.Println("Logged successfully at /var/log/", fileName+".yaml")
			return file.Write(data)

		} else if so == "windows" {

			file, err := os.Create(fileName + ".yaml")
			if err != nil {
				panic(err)
			}
			defer file.Close()

			fmt.Println("Logged successfully at current directory: ", fileName+".yaml")
			return file.Write(data)
		} else {
			err = fmt.Errorf("Unsupported operating system for logging.")
			return 0, err
		}

	default:
		return 0, fmt.Errorf("unsupported format: %s", format)
	}
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
	fmt.Printf("  %speutils(.exe)%s %s<file>%s %s<command>%s [%s<sectionName>%s]\n", cyan, reset, green, reset, red, reset, blue, reset)

	fmt.Printf("\n%sCommands:%s\n", yellow, reset)
	fmt.Printf("  %s-l%s             - List all imported libraries.\n", green, reset)
	fmt.Printf("  %s-sy%s            - List all imported symbols.\n", green, reset)
	fmt.Printf("  %s-s%s             - Display a specific section by name (use with [sectionName]).\n", green, reset)
	fmt.Printf("  %s-i%s             - Display the PE file header.\n", green, reset)
	fmt.Printf("  %s-oh%s            - Display the optional headers.\n", green, reset)
	fmt.Printf("  %s-fh%s            - Display the file header (basic details about the PE file).\n", green, reset)
	fmt.Printf("  %s-cs%s            - Display COFF symbols.\n", green, reset)
	fmt.Printf("  %s-ma%s            - Display the machine architecture type.\n", green, reset)
	fmt.Printf("  %s-sT%s            - Display the string table.\n", green, reset)
	fmt.Printf("  %s-tm%s            - Display the timestamp of the PE file.\n", green, reset)
	fmt.Printf("  %s-dw%s            - Extract and display DWARF debug data (if present).\n", green, reset)
	fmt.Printf("  %s-ps%s            - Display the pointer to the symbol table.\n", green, reset)
	fmt.Printf("  %s-ch%s            - Display the characteristics of the PE file.\n", green, reset)
	fmt.Printf("  %s-st%s            - Extract and display strings from the `.strtab` section.\n", green, reset)

	fmt.Printf("\n%sExamples:%s\n", yellow, reset)
	fmt.Printf("  %speutils.exe%s %smyfile.exe%s %s-l%s\n", cyan, reset, green, reset, red, reset)
	fmt.Printf("  %speutils.exe%s %smyfile.exe%s %s-s%s %s.text%s\n", cyan, reset, green, reset, red, reset, blue, reset)
	fmt.Printf("  %speutils.exe%s %smyfile.exe%s %s-sy%s\n", cyan, reset, green, reset, red, reset)

	fmt.Printf("\n%sLogging:%s\n", yellow, reset)
	fmt.Printf("  Use the %s--log%s flag to save the output to a file (supported formats: json, xml, yaml).\n", blue, reset)
	fmt.Printf("  Example: %speutils.exe%s %smyfile.exe%s %s-l%s %s--log json%s\n", cyan, reset, green, reset, red, reset, blue, reset)

	fmt.Printf("\n%sNote:%s The program will overwrite any existing log file in the same directory.\n", yellow, reset)
	fmt.Printf("%sFor detailed documentation on PE file analysis, refer to:%s\n", yellow, reset)
	fmt.Printf("  %shttps://pkg.go.dev/debug/pe%s\n", cyan, reset)
}


func main() {
	if len(os.Args) < 3 {
		help()
		return
	}

	logEnabled := false
	logFormat := ""

	fileName := os.Args[1]
	command := os.Args[2]
	sectionName := ""
	if command == "-s" {
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
	case "-l":
		libs, err := lib(peFile)
		if err != nil {
			fmt.Printf("Error fetching libraries: %v\n", err)
		}
		result = libs
	case "-sy":
		symbols, err := sym(peFile)
		if err != nil {
			fmt.Printf("Error fetching symbols: %v\n", err)
		}
		result = symbols
	case "-s":
		if sectionName == "" {
			fmt.Println("Please specify a section name for the 'sections' command.")
			return
		}
		section := sections(peFile, sectionName)
		if section == nil {
			fmt.Printf("Section %s not found.\n", sectionName)
		}
		result = section
	case "-i":
		header := info(peFile)
		result = header
	case "-oh":
		headers := optionalHeaders(peFile)
		result = headers
	case "-fh":
		result = theFile(peFile)
	case "-cs":
		result = coffNums(peFile)
	case "-ma":
		result = machine(peFile)
	case "-sT":
		result = stringsTable(peFile)
	case "-tm":
		result = dates(peFile)
	case "-dw":
		result = dwarf(peFile)
	case "-ps":
		result = pointerSymTables(peFile)
	case "-ch":
		result = Characteristics(peFile)
	case "-st":
		result, err = stringss(peFile)
		if err != nil {
			fmt.Printf("Error fetching strings: %v\n", err)
		}
	default:
		fmt.Println("Unknown command. Valid commands are: lib, sym, sections, info, optionalHeaders, fileHeader, coffSymbols, machine, stringTable, timeDateStamp, dwarf, pointerSymTables, characteristics")
		help()
	}

	for i := 3; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "--log" {
			logEnabled = true
			if i+1 < len(os.Args) {
				logFormat = os.Args[i+1]
				i++ // next arg
			} else {
				fmt.Println("Error: File type, json, xml or yaml")
				panic("Error: File type, json, xml or yaml needed")
			}
		}
	}

	if logEnabled {
		if logFormat == "" {
			panic("Error:  File type, json, xml or yaml")
		}
		_, err := saveResult(fileName, result, logFormat)
		if err != nil {
			fmt.Printf("Error while saving log: %v\n", err)
			panic(err)
		}
	}

	// Pretty print the result in JSON format
	prettyPrintJSON(result)
}
