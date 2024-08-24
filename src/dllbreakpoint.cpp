#include <iostream>

using namespace std;

extern "C" bool __attribute__((visibility("default"))) isHardwareBreakpointSet(){
        unsigned long dr0, dr1, dr2, dr3, dr7;
        asm volatile (
            "mov %%dr0, %0\n\t"
            "mov %%dr1, %1\n\t"
            "mov %%dr2, %2\n\t"
            "mov %%dr3, %3\n\t"
            "mov %%dr7, %4\n\t"
            : "=r" (dr0), "=r" (dr1), "=r" (dr2), "=r" (dr3), "=r" (dr7)
        );

        
        if (dr7 & 0xFF) {
            std::cout << "Hardware breakpoint detected!" << std::endl;
            return true;
        } else {
            std::cout << "No hardware breakpoint detected." << std::endl;
            return false;
        }
    }


/*
for a dll file: 

[1] g++ -c (scriptname).cpp -fPIC

[2] gcc -shared -fPIC -o brekdll.dll (scriptname).o
*/