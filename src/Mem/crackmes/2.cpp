/*
 * Copyright 2023-2024 Riccardo Adami. All rights reserved.
 * License: https://github.com/riccio8/ThreatLab/blob/main/LICENSE
 */


#include <iostream>
#include <chrono>
#include <windows.h>
#include <thread>

using namespace std;

bool IsDebugge() {
    return IsDebuggerPresent();
}

void decrypt(int* buffer, int size, int key) {
    for (int i = 0; i < size; i++) {
        buffer[i] ^= key;
    }
}

void delay() {
    std::this_thread::sleep_for(std::chrono::seconds(3));
}

int* ydobio() {
    static int buffer[7]; 
    if (IsDebugge()) {
        std::cout << "Debugger detected!" << std::endl;
        exit(1);
    }
    delay();

    buffer[0] = 0x22 ^ 0x41;
    buffer[1] = 0x35 ^ 0x41;
    buffer[2] = 0x2d ^ 0x41;
    buffer[3] = 0x2e ^ 0x41;
    buffer[4] = 0x2d ^ 0x41;
    buffer[5] = 0x2e ^ 0x41;
    buffer[6] = 0x60 ^ 0x41;

    decrypt(buffer, 7, 0x41);
    return buffer;
}

int main() {
    std::cout << "Find the hidden world" << std::endl;
    int* result = ydobio();
    std::cout << "maybe something cool at: " << std::oct << reinterpret_cast<uintptr_t>(result) << std::endl;
    delay();
    return 0;
}
