#include <windows.h>
#include <iostream>

int main() {
    char fileName[512];
    
    std::cout << "Paste the path to your PE file (.exe) using // or '' as separators: " << std::endl;
    std::cin.getline(fileName, sizeof(fileName));

    if (strlen(fileName) >= 512) { 
        std::cerr << "Buffer too big: " << fileName << std::endl;
        return 1;
    } else {
        std::cout << "File retrieving: " << fileName << std::endl;
    }

    DWORD lengthNeeded = 0;

    BOOL success = GetFileSecurityA(
        fileName,
        DACL_SECURITY_INFORMATION | OWNER_SECURITY_INFORMATION,
        NULL,
        0,
        &lengthNeeded
    );

    std::cout << "Buffer size needed: " << lengthNeeded << std::endl;

    if (!success && GetLastError() == ERROR_INSUFFICIENT_BUFFER) {

        PSECURITY_DESCRIPTOR pSecDesc = (PSECURITY_DESCRIPTOR)new char[lengthNeeded];

        if (pSecDesc == nullptr) {
            std::cerr << "Memory allocation failed." << std::endl;
            return 1;
        }

        success = GetFileSecurityA(
            fileName,
            DACL_SECURITY_INFORMATION | OWNER_SECURITY_INFORMATION,
            pSecDesc,
            lengthNeeded,
            &lengthNeeded
        );

        if (success) {
            std::cout << "Security descriptor retrieved successfully." << std::endl;
            std::cout << "Security Descriptor Address: " << pSecDesc << std::endl;
        } else {
            std::cout << "Error retrieving security descriptor: " << GetLastError() << std::endl;
        }

        delete[] pSecDesc;  
    } else {
        std::cout << "Error calculating buffer size: " << GetLastError() << std::endl;
    }
    
    return 0;
}

// Tells where the file descriptor is allocated in the memory
// compile using: g++ -o file_desc FileDescriptor.cpp
