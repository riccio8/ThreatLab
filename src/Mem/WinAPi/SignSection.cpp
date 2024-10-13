// it will print the signature if found and the number of section, like .data, .text, .bss and so on i'll add another program that will print the export table and import table, the signature and the section

#include <windows.h>
#include <imagehlp.h>
#include <iostream>

int main() {
    char fileName[512];
    std::cout << "Paste the path to your PE file (.exe) using // or '' as separators: " << std::endl;
    std::cin >> fileName;
    
    if (stringLen(fileName) > 256) {
        std::cerr << "Buffer too big: " << fileName << std::endl;
    }

    HANDLE file = CreateFile(fileName, GENERIC_READ, FILE_SHARE_READ, NULL, OPEN_EXISTING, 0, NULL);
    if (file == INVALID_HANDLE_VALUE) {
        std::cerr << "Failed to open file. Error code: " << GetLastError() << std::endl;
        return 1;
    } else {
        std::cout << "File opened successfully!" << std::endl;
    }

    HANDLE mapping = CreateFileMapping(file, NULL, PAGE_READONLY, 0, 0, NULL);
    if (!mapping) {
        std::cerr << "Failed to create file mapping. Error code: " << GetLastError() << std::endl;
        CloseHandle(file);
        return 1;
    } else {
        std::cout << "File mapping created successfully!" << std::endl;
    }

    LPVOID baseAddress = MapViewOfFile(mapping, FILE_MAP_READ, 0, 0, 0);
    if (!baseAddress) {
        std::cerr << "Failed to map view of file. Error code: " << GetLastError() << std::endl;
        CloseHandle(mapping);
        CloseHandle(file);
        return 1;
    } else {
        std::cout << "File mapped successfully!" << std::endl;
    }

    PIMAGE_NT_HEADERS ntHeaders = ImageNtHeader(baseAddress);
    if (!ntHeaders) {
        std::cerr << "Failed to get NT headers. Error code: " << GetLastError() << std::endl;
        UnmapViewOfFile(baseAddress);
        CloseHandle(mapping);
        CloseHandle(file);
        return 1;
    } else {
        std::cout << "PE file signature: " << std::hex << ntHeaders->Signature << std::endl;
        std::cout << "Number of sections: " << ntHeaders->FileHeader.NumberOfSections << std::endl;
    }
    
    // Clean up all stuff that we opened before
    UnmapViewOfFile(baseAddress);
    CloseHandle(mapping);
    CloseHandle(file);
    return 0;
}


// compile using: g++ -o dbginfos .\ImageDebugInformation.cpp -limagehlp
