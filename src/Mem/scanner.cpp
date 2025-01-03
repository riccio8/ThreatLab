/*
 * Copyright 2023-2024 Riccardo Adami. All rights reserved.
 * License: https://github.com/riccio8/ThreatLab/blob/main/LICENSE
 */



/*
This file is for check
*/

#include <iostream>

const int dangerous_syscalls[] = {

    // Syscalls that handle memory

    9,   // sys_mmap: Maps files or devices into memory

    10,  // sys_mprotect: Changes protection of memory regions

    12,  // sys_brk: Changes the end of the heap (dynamic memory)



    // Syscalls that execute or control processes

    11,  // sys_execve: Executes a program

    39,  // sys_setuid: Sets the user ID (can change process privileges)

    40,  // sys_setgid: Sets the group ID (can change process privileges)

    57,  // sys_fork: Creates a new process

    56,  // sys_vfork: Creates a new process (similar to fork but less costly)

    120, // sys_clone: Creates a new process (similar to fork with more control)

    31,  // sys_wait4: Waits for a child process to terminate



    // Syscalls that handle network resources

    102, // sys_socketcall: Manages socket calls



    // Syscalls that handle inter-process communication

    186, // sys_rt_sigprocmask: Changes signal mask (potentially dangerous)

    192, // sys_kexec_load: Loads a new kernel (potentially dangerous if misused)



    // Syscalls that modify files or disks

    43,  // sys_truncate: Changes the size of a file

    45,  // sys_ftruncate: Changes the size of a file via a file descriptor

    52,  // sys_fsync: Synchronizes a file descriptor's data with disk



    // Syscalls that handle shared memory

    25,  // sys_shmat: Attaches a shared memory segment

    26,  // sys_shmdt: Detaches a shared memory segment

    27,  // sys_shmctl: Controls shared memory (create, modify, remove)

};



bool isDangerousSyscall(int* syscall_number) {

    for (int i = 0; i < sizeof(dangerous_syscalls) / sizeof(dangerous_syscalls[0]); ++i) {

        if (dangerous_syscalls[i] == *syscall_number) {

            return true;
        }
    }
    return false;
}

bool check() {

    int *r0, *r1, *r2, *r3, *r4, *r5 = new int(64);

    int *rmain = new int;

    asm volatile(

        "mov %%eax, %0" : "=g"(rmain) 

    );


    if (isDangerousSyscall(rmain)) {

        asm volatile (

            "mov %%ebx, %0\n\t"

            "mov %%ecx, %1\n\t"

            "mov %%edx, %2\n\t"

            "mov %%esi, %3\n\t"

            "mov %%edi, %4\n\t"

            "mov %%ebp, %5\n\t"

            : "=g" (r0), "=g" (r1), "=g" (r2), "=g" (r3), "=g" (r4), "=g" (r5) 

        );

        std::cout << "----------------------------------------------------------------" 

                  << "\n" << "Those are the arguments of the syscall number " << rmain << ": \n" 

                  << r0 << "\n" << r1 << "\n" << r2 << "\n" << r3 << "\n" << r4 << "\n" << r5 

                  << "----------------------------------------------------------------" << std::endl;
        return true;
    }
    return false;
}


int main() {
    check();
    return 0;
}


/* 
for compile that, need to set to 32 bit version, like gcc or g++ -m32 -0 ...
*/
