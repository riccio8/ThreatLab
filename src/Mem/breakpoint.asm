.section .text
.global isHardwareBreakpointSet
.isHardwareBreakpointSet:
    pushq %rbp
    movq %rsp, %rbp

    movq %dr7, %rax
