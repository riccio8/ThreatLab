BITS 64
default rel


section .data
    folderPath db "C:\\hidden_folder", 0    ; Path to the folder (Note: double backslashes in Windows paths)

section .text
    global _start
    extern CreateDirectoryW, SetFileAttributesW, ExitProcess

_start:
    ; Create Directory (CreateDirectoryW)
    lea rdi, [folderPath]    ; Address of folderPath (using rdi for 64-bit)
    xor rsi, rsi             ; Security attributes (NULL) in rsi (for 64-bit)
    call CreateDirectoryW    ; Call CreateDirectoryW function

    ; Check for success
    test rax, rax
    jz create_error          ; Jump to create_error if creation failed

    ; Set Hidden attribute (SetFileAttributesW)
    lea rdi, [folderPath]    ; Address of folderPath
    mov rdx, 0x2             ; FILE_ATTRIBUTE_HIDDEN (value for hidden attribute)
    call SetFileAttributesW  ; Call SetFileAttributesW function

    ; Check for success
    test rax, rax
    jz set_error             ; Jump to set_error if setting the attribute failed

    ; Exit program successfully (ExitProcess)
    xor rax, rax             ; Exit code 0
    call ExitProcess         ; Call ExitProcess to terminate the program

create_error:
    ; Handle error: directory creation failed
    mov rax, 1               ; Exit code 1
    call ExitProcess         ; Call ExitProcess to terminate with error code

set_error:
    ; Handle error: setting attribute failed
    mov rax, 1               ; Exit code 1
    call ExitProcess         ; Call ExitProcess to terminate with error code


; compilation: nasm -f win64 -o hide.obj hidden.asm
; linking: ld -o hide.exe hide.obj -lkernel32 -luser32 
