.section .text
.global isHardwareBreakpointSet
.isHardwareBreakpointSet:
    pushq %rbp
    movq %rsp, %rbp

    movq %dr7, %rax
    testq $0xFF, %rax
    jz no_breakpoint

    movl $1, %eax
    jmp end

no_breakpoint:
    xorl %eax, %eax

end:
    popq %rbp
    ret
