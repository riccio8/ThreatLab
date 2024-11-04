#include <iostream>
#include <fstream>
#include <iomanip>
#include <algorithm> 

std::string trim_quotes(const std::string& path) {
    std::string trimmed = path;
    if (trimmed.front() == '"') {
        trimmed.erase(trimmed.begin());
    }
    if (trimmed.back() == '"') {
        trimmed.pop_back();
    }
    return trimmed;
}

int main() {
    std::string filepath; 
    std::cout << "Enter your file path: " << std::endl;
    std::getline(std::cin, filepath);
    
    filepath = trim_quotes(filepath); 
    
    std::ifstream file(filepath, std::ios::binary);
    
    if (!file) {
        std::cerr << "Error opening file: " << filepath << std::endl; 
    } else {
        std::cout << "File opened successfully!" << std::endl;

    }
    char byte;
    while (file.get(byte)) {
        std::cout << std::hex << std::setw(2) << std::setfill('0') << (static_cast<unsigned char>(byte) & 0xFF) << " ";
    }

    file.close();
    return 0;
}
