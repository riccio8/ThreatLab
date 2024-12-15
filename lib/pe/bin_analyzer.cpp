#include <windows.h>
#include <fstream>
#include <iostream>
#include <stdexcept>
#include <cstring>

#include "bin_analyzer.hpp"

// DOS Header structure, the first part of the PE file
struct DOSHeader {
    WORD e_magic;   // Magic number (0x5A4D - "MZ")
    WORD e_cblp;    // Bytes on last page of file
    // Additional fields are present, but these are the essential ones for the initial parsing
};

// PE Header structure, following the DOS Header
struct PEHeader {
    DWORD signature; // PE Signature: 0x00004550 ("PE\0\0")
    IMAGE_FILE_HEADER fileHeader;
    IMAGE_OPTIONAL_HEADER optionalHeader;
};

// Section Header structure, describing a section in the PE file
struct SectionHeader {
    char name[8];      // Section name (null-terminated string)
    DWORD virtualSize; // Virtual size of the section
    DWORD virtualAddress; // RVA (Relative Virtual Address)
    DWORD sizeOfRawData;  // Size of section data in file
    DWORD pointerToRawData; // Pointer to the raw data in file
    // Other fields are optional, but these are the key ones
};

// Class for analyzing the binary (PE) file
class BinaryAnalyzer {
public:
    explicit BinaryAnalyzer(const std::string& filePath) : filePath(filePath) {}

    // Function to parse the headers and retrieve section information
    void parseHeaders() {
        std::ifstream file(filePath, std::ios::binary);
        if (!file.is_open()) {
            throw std::runtime_error("Failed to open file. \n File open in another application ?");
        }

        // Read DOS Header
        DOSHeader dosHeader;
        file.read(reinterpret_cast<char*>(&dosHeader), sizeof(dosHeader));

        if (dosHeader.e_magic != 0x5A4D) { // Check if it's a PE file by verifying the magic number
            throw std::runtime_error("Invalid DOS Header.");
        }

        // Move to the PE Header (pointed by e_lfanew in the DOS header)
        file.seekg(dosHeader.e_lfanew, std::ios::beg); // e_lfanew points to where the PE Header starts
        PEHeader peHeader;
        file.read(reinterpret_cast<char*>(&peHeader), sizeof(peHeader)); // Searching on google i saw the cast is required

        if (peHeader.signature != 0x00004550) { // Check for the PE signature
            throw std::runtime_error("Invalid PE Signature.");
        }

        // Now, parse the sections
        parseSections(file, peHeader.fileHeader.NumberOfSections);
    }

    // Function to get a section by its name
    Section getSection(const std::string& name) {
        for (const auto& section : sections) {
            if (section.name == name) {
                return section;
            }
        }
        throw std::runtime_error("Section not found.");
    }

    // Function to save the sections to an output file
    void save(const std::string& outputPath) {
        std::ofstream outputFile(outputPath, std::ios::binary);
        if (!outputFile.is_open()) {
            throw std::runtime_error("Failed to open output file.");
        }

        // Write each section's data to the output file
        for (const auto& section : sections) {
            outputFile.write(reinterpret_cast<const char*>(section.data.data()), section.data.size());
        }
    }

private:
    std::string filePath;
    std::vector<Section> sections;

    // Function to parse section headers and load section data
    void parseSections(std::ifstream& file, WORD numSections) {
        // Each section is described by a Section Header (40 bytes)
        for (WORD i = 0; i < numSections; ++i) {
            SectionHeader sectionHeader;
            file.read(reinterpret_cast<char*>(&sectionHeader), sizeof(sectionHeader));

            // Create a new Section object and store its name
            Section section;
            section.name = std::string(sectionHeader.name, 8); // Convert the section name to a string
            section.data.resize(sectionHeader.sizeOfRawData); // Allocate space for section data

            // Move the file pointer to where the raw section data is stored
            file.seekg(sectionHeader.pointerToRawData, std::ios::beg);
            file.read(reinterpret_cast<char*>(section.data.data()), sectionHeader.sizeOfRawData); // Read the section data, again, the cast is required

            // Add the section to the list of sections
            sections.push_back(section);
        }
    }
};
