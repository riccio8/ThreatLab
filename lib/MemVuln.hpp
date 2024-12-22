/*
 * Copyright 2023-2024 Riccardo Adami. All rights reserved.
 * License: https://github.com/riccio8/ThreatLab/blob/main/LICENSE
 */


#ifndef MemVuln
#define MemVuln

//#pragma once

#include <iostream>
#include <cstring>

class MemType {
public:
    long number;
    char* buffer;

    MemType(long n, const char* b);
    ~MemType();

    void MemLeak();
    void heapLeak();
};

class VulnerableVoid {
public:
    void stack_overflow(const char* input);
    void heap_overflow(const char* big);
};

#endif 
