BITS 64                 ; Define the target architecture as 64-bit
default rel            ; Set default addressing mode to RIP-relative

section .data
    file db "hide.txt", 0       ; String "hide.txt" with null terminator (used as file name)
    hTemplateFile dq 0          ; Initialize hTemplateFile (NULL)
    message db "error", 0       ; Error message
    msg_len equ $ - message     ; Calculate the length of the error message

section .text
    global _start              
    extern SetFileAttributesW, ExitProcess, CreateFileA, WriteFile  

_start:
    ; Call CreateFileA to create/open the file
    lea rdi, [rel file]                 ; Load the address of the file name (RIP-relative addressing)
    mov rsi, 0x40000000                 ; dwDesiredAccess: GENERIC_WRITE (0x40000000)
    mov rdx, 0x00000001                 ; dwShareMode: FILE_SHARE_READ (0x00000001)
    xor rcx, rcx                        ; lpSecurityAttributes: NULL (0x0)
    mov r8, 0x00000004                  ; dwCreationDisposition: OPEN_ALWAYS (0x00000004)
    mov r9, 0x80                        ; dwFlagsAndAttributes: FILE_ATTRIBUTE_NORMAL (0x80)
    mov qword [rsp+8], 0                ; hTemplateFile: NULL (0x0)
    call CreateFileA            

    ; Check if the file creation/opening succeeded
    test rax, rax              
    jz set_error                

    ; Set file attribute to hidden using SetFileAttributesW (Unicode string required)
    lea rdi, [rel file]         ; Load the address of the file name (RIP-relative addressing)
    mov rdx, 0x2                ; dwFileAttributes: FILE_ATTRIBUTE_HIDDEN (0x2)
    call SetFileAttributesW     

    ; Check if setting the attribute succeeded
    test rax, rax               ; Test if rax is zero (failure condition)
    jz set_error                ; Jump to set_error if SetFileAttributesW failed

    ; Exit the program successfully (exit code 0)
    xor rax, rax                
    call ExitProcess            ; Call ExitProcess to exit the program

set_error:

    mov rdi, 0xFFFFFFF5         ; File handle for stdout (negative value: console output)
    lea rsi, [rel message]      
    mov rdx, msg_len            
    call WriteFile              ; Call WriteFile to print the error message

    ; Exit with error code 1
    mov rax, 1                  ; Set the exit code to 1
    call ExitProcess            ; Call ExitProcess to exit with error code 1


; nasm -f win64 -o hide_file.obj .\hide_file.asm
; ld -o hide_file.exe hide_file.obj -lkernel32
; hide_file.exe 
