#include <iostream>
#include <fstream>
#include <iomanip>
#include <vector>
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
    std::vector<unsigned char> bytes;  // Declare a variable to store each byte read from the file.
    
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
    
 unsigned char byte;
    while (file.read(reinterpret_cast<char*>(&byte), sizeof(byte))) {
        bytes.push_back(byte);
    }

    
    std::cout << "Bytes read from file:" << std::endl;
    for (unsigned char b : bytes) {
        std::cout << std::hex << static_cast<int>(b) << " ";  
    }

    return 0;
}
