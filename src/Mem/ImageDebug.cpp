#include <windows.h>
#include <imagehlp.h>
#include <iostream>
#include <dbghelp.h>

int main() {
    char fileName[512];
    char debugFilePath[MAX_PATH];
    std::cout << "Paste the path to your PE file (.exe) using // or '' as separators: " << std::endl;
    std::cin >> fileName;

    if (strlen(fileName) > 512) {
        std::cerr << "Buffer too big: " << fileName << std::endl;
        return 1;
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

    const char* symbolPath = "C:\\Symbols"; 

    HANDLE debugHandle = FindDebugInfoFile(fileName, symbolPath, debugFilePath);
    if (debugHandle == NULL) {
        std::cerr << "Failed to find debug info file. Error code: " << GetLastError() << std::endl;
        if (GetLastError() == 0){
            std::cerr << "Failed to open debug info file. Ensure that it exists, or try adjusting the symblo path (default: C:\\Symbols)" << std::endl;
        }
    } else {
        std::cout << "Debug info file found at: " << debugFilePath << std::endl;
        CloseHandle(debugHandle); // Close the handle after use
    }

    // Clean up
    UnmapViewOfFile(baseAddress);
    CloseHandle(mapping);
    CloseHandle(file);

    return 0;
}


// compile it using:  g++ -o dbginfos .\ImageDebugInformation.cpp -limagehlp -ldbghelp
