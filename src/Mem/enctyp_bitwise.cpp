#include <iostream>
#include <vector>
#include <cstdint>
#include <cstdlib>

std::vector<uint8_t> keyGen(){
    std::vector<uint8_t> key(16); 
    for (auto &k : key) {
        k = rand() % 256; 
    }
    return key;
}  

std::vector<uint8_t> encrypt(const std::vector<uint8_t> &data, const std::vector<int> &keys) {
    std::vector<uint8_t> encryptedData;
    for (size_t i = 0; i < data.size(); ++i) {
        uint8_t maskedByte = data[i];
        int key = keys[i % keys.size()];

        maskedByte = maskedByte ^ key;             // XOR with the key
        maskedByte = (maskedByte & ~key) | (key);  // AND+OR conditional
        maskedByte = maskedByte << (key % 8);      // Shift to left csugin the key
        encryptedData.push_back(maskedByte);
    }
    return encryptedData;
}

int main() {
    std::vector<uint8_t> not_data = { 0x45, 0x23, 0x56, 0x12 };
    std::vector<int> number = { 1, 2, 3, 4 };

    std::vector<uint8_t> encryptedData = encrypt(not_data, number);

    for (uint8_t byte : encryptedData) {
        std::cout << std::hex << static_cast<int>(byte) << " ";
    }
    std::cout << std::endl;

    return 0;
}
