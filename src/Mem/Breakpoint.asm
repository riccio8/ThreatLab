;-----------------------------------------------------
; Copyright 2023-2024 Riccardo Adami. All rights reserved.
; License: https://github.com/riccio8/ThreatLab/blob/main/LICENSE
;-----------------------------------------------------


section .text
global get_debug_registers

get_debug_registers:
    ; Input: 
    ; %0 -> address to store DR0
    ; %1 -> address to store DR1
    ; %2 -> address to store DR2
    ; %3 -> address to store DR3
    ; %4 -> address to store DR7
    mov      %dr0, rax        ; move DR0 into RAX
    mov      rax, [rdi]       ; store RAX at the first address (DR0)
    
    mov      %dr1, rax        ; move DR1 into RAX
    mov      rax, [rsi]       ; store RAX at the second address (DR1)
    
    mov      %dr2, rax        ; move DR2 into RAX
    mov      rax, [rdx]       ; store RAX at the third address (DR2)
    
    mov      %dr3, rax        ; move DR3 into RAX
    mov      rax, [rcx]       ; store RAX at the fourth address (DR3)
    
    mov      %dr7, rax        ; move DR7 into RAX
    mov      rax, [r8]        ; store RAX at the fifth address (DR7)

    ret                       ; return the value of the check 
