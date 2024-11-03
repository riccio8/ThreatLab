---

# README: x64 Assembly Programming on Windows

## Overview

This guide covers:
- Argument passing in x64 Assembly using NASM and FASM.
- Calling conventions: stdcall and fastcall.
- Creating a DLL that checks for a debugger.
- Using the DLL in a C++ program.

## Calling Conventions in x64

### Microsoft x64 Calling Convention
- **First Four Arguments:** Passed in the registers:
  - `RCX`: First argument
  - `RDX`: Second argument
  - `R8`: Third argument
  - `R9`: Fourth argument
- **Additional Arguments:** Pushed onto the stack.
- **Return Value:** The result goes into `RAX`.

### Registers Summary Table

| Register | Purpose                                     |
|----------|---------------------------------------------|
| `RCX`    | 1st integer/pointer argument                |
| `RDX`    | 2nd integer/pointer argument                |
| `R8`     | 3rd integer/pointer argument                |
| `R9`     | 4th integer/pointer argument                |
| `RSP`    | Stack Pointer; points to the top of the stack |
| `RAX`    | Return value of a function                  |
| `RBP`    | Base Pointer; used for stack frame access   |
| `R10-R15`| General-purpose registers                    |

## Explanation of `stdcall` and `fastcall`

- **Stdcall:** 
  - A calling convention where arguments are passed on the stack, right to left. The caller cleans up the stack after the function call. Commonly used in WinAPI.

- **Fastcall:** 
  - A calling convention that passes the first two arguments in `ECX` and `EDX` registers and the rest on the stack. The callee cleans up the stack. It’s used for performance.

## Example: Calling a Function Using NASM

Here’s how to call `IsDebuggerPresent` and save its return value in a variable using NASM:

### NASM Example

```asm
section .data
    is_debugger db 0   ; Variable to store the debugger state

section .text
    extern IsDebuggerPresent
    global main

main:
    ; Call IsDebuggerPresent and store the return value
    call IsDebuggerPresent
    test rax, rax      ; Check if the result is non-zero
    jz .not_debugger    ; Jump if zero (not present)

    mov byte [is_debugger], 1  ; Set to 1 if debugger is present
    jmp .done

.not_debugger:
    mov byte [is_debugger], 0  ; Set to 0 if not present

.done:
    ret
```

## Example: Calling a Function Using FASM

Here’s the same example using FASM syntax, utilizing `invoke` for cleaner calls:

### FASM Example

```asm
section .data
    is_debugger db 0   ; Variable to store the debugger state

section .text
    extern IsDebuggerPresent
    global main

main:
    ; Call IsDebuggerPresent and store the return value
    invoke IsDebuggerPresent
    test rax, rax      ; Check if the result is non-zero
    jz .not_debugger    ; Jump if zero (not present)

    mov byte [is_debugger], 1  ; Set to 1 if debugger is present
    jmp .done

.not_debugger:
    mov byte [is_debugger], 0  ; Set to 0 if not present

.done:
    ret
```

## Creating a DLL that Checks for a Debugger

### NASM DLL Example

Here’s how to create a DLL that checks if a debugger is present using NASM:

```asm
section .data
    is_debugger db 0   ; Variable to store the debugger state

section .text
    global DllMain
    global CheckDebugger

; DLL Entry Point
DllMain:
    ; We don’t need to do anything here for this example
    mov eax, 1         ; Return TRUE
    ret

; Function to check if a debugger is present
CheckDebugger:
    call IsDebuggerPresent
    test rax, rax      ; Check if the result is non-zero
    jz .not_debugger

    mov byte [is_debugger], 1  ; Set to 1 if debugger is present
    mov rax, 1                  ; Return 1 (true)
    ret

.not_debugger:
    mov byte [is_debugger], 0  ; Set to 0 if not present
    xor rax, rax                ; Return 0 (false)
    ret
```

### FASM DLL Example

And here’s how to do the same with FASM:

```asm
section .data
    is_debugger db 0   ; Variable to store the debugger state

section .text
    global DllMain
    global CheckDebugger

; DLL Entry Point
DllMain:
    mov eax, 1         ; Return TRUE
    ret

; Function to check if a debugger is present
CheckDebugger:
    invoke IsDebuggerPresent
    test rax, rax      ; Check if the result is non-zero
    jz .not_debugger

    mov byte [is_debugger], 1  ; Set to 1 if debugger is present
    mov rax, 1                  ; Return 1 (true)
    ret

.not_debugger:
    mov byte [is_debugger], 0  ; Set to 0 if not present
    xor rax, rax                ; Return 0 (false)
    ret
```

### Compiling Your DLL

- For **NASM**:
  ```bash
  nasm -f win64 mydll.asm -o mydll.obj
  link /DLL /OUT:mydll.dll mydll.obj
  ```

- For **FASM**:
  ```bash
  fasm mydll.asm mydll.dll
  ```

## Using the DLL in C++

Here’s how to use the DLL in a C++ application:

### C++ Example

```cpp
#include <windows.h>
#include <iostream>

typedef int(*CheckDebuggerType)();

int main() {
    HMODULE hModule = LoadLibraryA("mydll.dll");
    if (hModule) {
        CheckDebuggerType CheckDebugger = (CheckDebuggerType)GetProcAddress(hModule, "CheckDebugger");
        if (CheckDebugger) {
            int isDebug = CheckDebugger(); // Call the function from the DLL
            if (isDebug) {
                std::cout << "Debugger is present." << std::endl;
            } else {
                std::cout << "No debugger detected." << std::endl;
            }
        }
        FreeLibrary(hModule);
    } else {
        std::cout << "Failed to load DLL." << std::endl;
    }
    return 0;
}
```
