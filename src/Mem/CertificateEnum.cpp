#include <windows.h>
#include <imagehlp.h>
#include <iostream>


int main(){
    char fileName[512];
    
    std::cout << "Paste the path to your PE file (.exe) using // or '' as separators: " << std::endl;
    std::cin.getline(fileName, sizeof(fileName));

    if (strlen(fileName) >= 512) { 
        std::cerr << "Buffer too big: " << fileName << std::endl;
        return 1;
    }
    
    
    HANDLE file = CreateFile(fileName, GENERIC_READ, FILE_SHARE_READ | FILE_SHARE_WRITE, NULL, OPEN_EXISTING, 0, NULL);
    if (file == INVALID_HANDLE_VALUE) {
        std::cerr << "Failed to open file. Error code: " << GetLastError() << std::endl;
        return 1;
    } else {
        std::cout << "File opened successfully!" << std::endl;
    }
    
    DWORD certificates[99];
    DWORD count = sizeof(certificates) / sizeof(certificates[0]);
    DWORD type = CERT_SECTION_TYPE_ANY;
    DWORD flags = 0;
    
    bool enumerated = ImageEnumerateCertificates(file, type, certificates, &count, flags);
    if (enumerated == false || certificates == NULL) {
    std::cerr << "Failed to enumerate certificates. Error: " << GetLastError() << std::endl;
        return 1;
    }else{
        std::cout << "Enumerated succesfully certificates\n";
    }
    if (count <= 0){
        std::cerr << "No certificates were found. Error: " << GetLastError() << std::endl;
    }
    for (int i = 0; i < *certificates; i++) {
        std::cout << certificates[i] << "\t at index \t" << i << "\n"<<  std::endl;
    }
    std::cout << "No certificates were found. Error: " << GetLastError() << std::endl;
    CloseHandle(file);
    return 0;
}
