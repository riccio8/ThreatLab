#include "Offensive-defensive-tolls-/lib/MemVuln.hpp"

using namespace std;

class MemType {
public:
    long number;
    char* buffer;

    MemType(long n, const char* b) {
        this->number = n;
        this->buffer = new char[strlen(b) + 1];
        strcpy(this->buffer, b);
    }

    ~MemType() {
        delete[] buffer;
    }

    void MemLeak() {
        char buffer[2];
        strcpy(buffer, "TooLargeForBuffer");
        cout << "Buffer contains: " << sizeof(buffer) << endl;
    }

    void heapLeak() {
        long* num1 = new long(1000000000000);
        long num2 = (*num1) * 10000;
        cout << "Result of heap operation: " << num2 << endl;
    }
};

class VulnerableVoid {
public:
    void stack_overflow(const char* input) {
        char buffer[8];  
        cout << "Simulating a stack overflow: " << endl;
        strcpy(buffer, input);  
        cout << "Data in the buffer: " << buffer << endl;
    }

    void heap_overflow(const char* big) {
        cout << "Simulating heap overflow..." << endl;
        char* buffer1 = new char[4];
        char* buffer2 = new char[64];

        strcpy(buffer1, big); 
        cout << "Buffer1: " << buffer1 << endl;
        cout << "Buffer2: " << buffer2 << endl;
        
    }
};

int main() {
    VulnerableVoid vul;
    MemType mem(1000, "Hello");

    vul.stack_overflow("ThisIsAVeryLongStringThatWillOverflowTheBufferInStackOverflowMethod");
    vul.heap_overflow("ThisIsALongStringThatWillCauseHeapOverflow");

    mem.MemLeak();
    mem.heapLeak();

    return 0;
}
