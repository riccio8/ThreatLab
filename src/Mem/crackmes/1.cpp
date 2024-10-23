#include <iostream>
#include <chrono>
#include <windows.h>
#include <thread>

using namespace std;

bool IsDebugge(){
    return IsDebuggerPresent();
}

void decrypt(int* buffer, int size, int key) {
    for (int i = 0; i < size; i++) {
        buffer[i] ^= key;  // Decrypt with XOR
    }
}

void delay() {
    std::this_thread::sleep_for(std::chrono::seconds(3));  // 3-second delay
}

int* ydobio() {

    if (IsDebugge()) {
        std::cout << "Debugger detected!" << std::endl;
        exit(1);  // Terminate if debugging
    }
    delay();
    int* buffer = new int[7];
    buffer[0] = 0x22 ^ 0x41; // Encrypted 'b'
    buffer[1] = 0x35 ^ 0x41; // Encrypted 'v'
    buffer[2] = 0x2d ^ 0x41; // Encrypted 'l'
    buffer[3] = 0x2e ^ 0x41; // Encrypted 'n'
    buffer[4] = 0x2d ^ 0x41; // Encrypted 'l'
    buffer[5] = 0x2e ^ 0x41; // Encrypted 'n'
    buffer[6] = 0x60 ^ 0x41; // Encrypted '!'
    
    decrypt(buffer, 7, 0x41); // XOR decryption
    return buffer;
}

int main(){
    std::cout << "Find the hidden world" << std::endl;
    int* result = ydobio();
    return 0;
}
