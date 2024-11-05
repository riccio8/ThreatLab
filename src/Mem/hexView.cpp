#include <iostream>
#include <fstream>
#include <iomanip>
#include <vector>
#include <string>

std::string trim_quotes(const std::string& path) {
    std::string trimmed = path;
    if (!trimmed.empty() && trimmed.front() == '"') {
        trimmed.erase(trimmed.begin());
    }
    if (!trimmed.empty() && trimmed.back() == '"') {
        trimmed.pop_back();
    }
    return trimmed;
}


std::string get_replacement(const std::string& hex_value) {
    if (hex_value == "00") {
        return "NULL";   
    } else if (hex_value == "01") {
        return "START"; 
    } else if (hex_value == "02") {
        return "OPEN"; 
    } else if (hex_value == "03") {
        return "CLOSE"; 
    } else if (hex_value == "04") {
        return "READ"; 
    } else if (hex_value == "05") {
        return "WRITE"; 
    } else if (hex_value == "90") {
        return "NOP";   
    } else if (hex_value == "ff") {
        return "INVALID"; 
    } else if (hex_value == "C3") {
        return "RET"; 
    } else if (hex_value == "E8") {
        return "CALL"; 
    } else if (hex_value == "B8") {
        return "MOV"; 
    } else if (hex_value == "6A") {
        return "PUSH"; 
    } else if (hex_value == "68") {
        return "PUSH_IMM"; 
    } else if (hex_value == "C7") {
        return "MOV_MEM"; 
    }
    return hex_value; 
}



int main() {
    std::vector<unsigned char> bytes;  
    std::vector<int> formatted_bytes;  

    std::string filepath; 
    std::cout << "Enter your file path: " << std::endl;
    std::getline(std::cin, filepath);
    
    filepath = trim_quotes(filepath); 
    
    std::ifstream file(filepath, std::ios::binary);
    
    if (!file) {
        std::cerr << "Error opening file: " << filepath << std::endl; 
        return 1;
    } else {
        std::cout << "File opened successfully!" << std::endl;
    }
    
    unsigned char byte;
    while (file.read(reinterpret_cast<char*>(&byte), sizeof(byte))) {
        bytes.push_back(byte);               
        formatted_bytes.push_back((int)byte); 
    }
    
        
    for (const auto b : formatted_bytes) {
        std::stringstream ss;
        ss << std::hex << std::setw(2) << std::setfill('0') << b;
        std::string hex_value = ss.str();

        std::string output_value = get_replacement(hex_value);
        std::cout << "Hex value: " << hex_value << " -> " << output_value << std::endl;
    }

  

    return 0;
}
