/*
 * Copyright 2023-2024 Riccardo Adami. All rights reserved.
 * License: https://github.com/riccio8/ThreatLab/blob/main/LICENSE
 */



#include "Offensive-defensive-tolls-/lib/MemVuln.hpp"


using namespace std;

int main() {

    VulnerableVoid vul;
    MemType mem(1000, "Hello");


    vul.stack_overflow("ThisIsAVeryLongStringThatWillOverflowTheBufferInStackOverflowMethod");
    vul.heap_overflow("ThisIsALongStringThatWillCauseHeapOverflow");

    mem.MemLeak();
    mem.heapLeak();

    return 0;
}
