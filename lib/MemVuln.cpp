#include "MemVuln.hpp"

using namespace std;

MemType::MemType(long n, const char* b) {
    this->number = n;
    this->buffer = new char[strlen(b) + 1];
    strcpy(this->buffer, b);
}


MemType::~MemType() {
    delete[] buffer;
}


void MemType::MemLeak() {
    char buffer[20];
    strcpy(buffer, "TooLargeForBuffer");
    cout << "Buffer contains: " << buffer << endl;
}


void MemType::heapLeak() {
    long long* num1 = new long long(1000000000000LL); 
    long long num2 = (*num1) * 10000;
    cout << "Result of heap operation: " << num2 << endl;
    delete num1; 
}

void VulnerableVoid::stack_overflow(const char* input) {
    char buffer[8];  
    cout << "Simulating a stack overflow: " << endl;
    strcpy(buffer, input);  
    cout << "Data in the buffer: " << buffer << endl;
}

void VulnerableVoid::heap_overflow(const char* big) {
    cout << "Simulating heap overflow..." << endl;
    char* buffer1 = new char[4];
    char* buffer2 = new char[64];

    strcpy(buffer1, big); 
    cout << "Buffer1: " << buffer1 << endl;
    cout << "Buffer2: " << buffer2 << endl;

}
