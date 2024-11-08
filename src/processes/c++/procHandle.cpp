#include <windows.h>
#include <tlhelp32.h>
#include <iostream>
#include <psapi.h>
#include <cstdlib> 
#include <string>
#include <vector>

void DisplayHelp() {
    std::cout << "This is a tool for process analysis, it is suggested to use the 'generic' args as first one..." << std::endl;
    std::cout << "Usage: ProcHandle <command> [arguments]" << std::endl;
    std::cout << "Commands:" << std::endl;
    std::cout << "  list                    List all running processes on the system." << std::endl;
    std::cout << "  info <proc_name>       Retrieve detailed information for a specific process by its PID." << std::endl;

}

// This function use find the pid (process id) by it's name creating a tool that makes a snapshot of the processes
std::vector<DWORD> FindPidByName(const std::string& processName) {
    std::vector<DWORD> pids;
    PROCESSENTRY32 pe32;
    pe32.dwSize = sizeof(PROCESSENTRY32);

    HANDLE snapshot = CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0);
    if (snapshot == INVALID_HANDLE_VALUE) return pids;


    if (Process32First(snapshot, &pe32)) {
        do {
            std::string exeFile = pe32.szExeFile;

            if (exeFile.compare(0, processName.length(), processName) == 0 &&
                exeFile.length() == processName.length() + 4 &&
                exeFile.compare(processName.length(), 4, ".exe") == 0) {
                pids.push_back(pe32.th32ProcessID);
            }
        } while (Process32Next(snapshot, &pe32));
    }


    CloseHandle(snapshot);
    return pids;
}

void ListInfoProcesses() {
    std::cout << "Listing all processes..." << std::endl;
    PROCESSENTRY32 pe32;
    pe32.dwSize = sizeof(PROCESSENTRY32);

    HANDLE snapshot = CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0);
    if (snapshot == INVALID_HANDLE_VALUE) {
        std::cerr << "Error creating process snapshot:" << std::endl;
        return;
    }

    if (Process32First(snapshot, &pe32)) {
        do {
            std::cout << "Pid: " << pe32.th32ProcessID
                      << "\tFile Name: " << pe32.szExeFile
                      << "\tThread: " << pe32.cntThreads
                      << "\tProcess Flags: " << pe32.dwFlags
                      << "\tUsage Count: " << pe32.cntUsage 
                      << std::endl;
        } while (Process32Next(snapshot, &pe32));
    } else {
        std::cerr << "Error retrieving processes:" << std::endl;
    }

    CloseHandle(snapshot);
}

void GetProcessInfo(std::string names) {
    const std::string name = names+".exe";
    std::vector<DWORD> pids = FindPidByName(name);
    if (pids.empty()) {
        std::cout << "No processes found with the given name." << std::endl;
        return;
    }

    for (const auto& pid : pids) {
        std::cout << "Retrieving information for PID: " << pid << "..." << std::endl;
        HANDLE hProcess = OpenProcess(PROCESS_QUERY_INFORMATION | PROCESS_VM_READ, FALSE, pid);
        if (!hProcess) {
            std::cerr << "Error opening process:" << GetLastError() << std::endl;
            continue;
        }
        

        char processName[MAX_PATH];
        if (GetModuleFileNameExA(hProcess, NULL, processName, sizeof(processName) / sizeof(char))) {
            std::cout << "PID: " << pid << "\tName: " << processName << std::endl;
        } else {
            std::cerr << "Error retrieving process name:" << GetLastError() << std::endl;
        }

        CloseHandle(hProcess);
    }
}

void suspend(const std::string& name)
{

    std::vector<DWORD> pids = FindPidByName(name);

    if (pids.empty()) {
        std::cout << "No processes found with name: " << name << std::endl;
        return;
    }

    HANDLE hThreadSnapshot = CreateToolhelp32Snapshot(TH32CS_SNAPTHREAD, 0);
    if (hThreadSnapshot == INVALID_HANDLE_VALUE) {
        std::cout << "Error: Unable to create thread snapshot." << std::endl;
        return;
    }

    THREADENTRY32 threadEntry;
    threadEntry.dwSize = sizeof(THREADENTRY32);

    if (Thread32First(hThreadSnapshot, &threadEntry)) {
        do {
            for (DWORD pid : pids) {
                if (threadEntry.th32OwnerProcessID == pid) {
                    HANDLE hThread = OpenThread(THREAD_SUSPEND_RESUME, FALSE, threadEntry.th32ThreadID);
                    if (hThread != NULL) {
                        DWORD suspendResult = SuspendThread(hThread);
                        if (suspendResult == (DWORD)-1) {
                            std::cout << "Error suspending thread with ID: " 
                                      << threadEntry.th32ThreadID 
                                      << " in process ID: " << pid << std::endl;
                        } else {
                            std::cout << "Successfully suspended thread with ID: " 
                                      << threadEntry.th32ThreadID 
                                      << " in process ID: " << pid << std::endl;
                        }
                        CloseHandle(hThread);
                    } else {
                        std::cout << "Failed to open thread with ID: " 
                                  << threadEntry.th32ThreadID << std::endl;
                    }
                }
            }
        } while (Thread32Next(hThreadSnapshot, &threadEntry));
    }

    CloseHandle(hThreadSnapshot);
}

void TerminateProcessByName(const std::string& name) {
    std::vector<DWORD> pids = FindPidByName(name);
    for (const auto& pid : pids) {
        HANDLE hProcess = OpenProcess(PROCESS_TERMINATE, FALSE, pid);
        if (hProcess) {
            TerminateProcess(hProcess, 0);
            std::cout << "Process " << pid << " terminated successfully." << std::endl;
            CloseHandle(hProcess);
        } else {
            std::cerr << "Error terminating process " << pid << ":" << GetLastError() << std::endl;
        }
    }
}

int main(int argc, char* argv[]) {
    if (argc < 2) {
        DisplayHelp();
        return 1;
    }

    std::string command = argv[1];

    if (command == "list") {
        ListInfoProcesses();
    } else if (command == "info") {
        if (argc < 3) {
            std::cout << "Usage: info <process_name>" << std::endl;
            return 1;
        }
        GetProcessInfo(argv[2]);
    } else if (command == "kill") {
        if (argc < 3) {
            std::cout << "Usage: kill <process_name>" << std::endl;
            return 1;
        }
        TerminateProcessByName(argv[2]);
    } else if (command == "suspend") {
        if (argc < 3) {
            std::cout << "Usage: suspend <process name>" << std::endl;
            return 1;
        }
        suspend(argv[2]);
    } else {
        std::cout << "Error: Unknown command: " << command << std::endl;
        DisplayHelp();
    }

    return 0;
}

/*
Resources:

- https://learn.microsoft.com/en-us/windows/win32/api/tlhelp32/ns-tlhelp32-processentry32
- https://learn.microsoft.com/en-us/windows/win32/api/tlhelp32/nf-tlhelp32-createtoolhelp32snapshot
- https://learn.microsoft.com/en-us/windows/win32/api/tlhelp32/nf-tlhelp32-process32first
- https://learn.microsoft.com/en-us/windows/win32/api/handleapi/nf-handleapi-closehandle
- https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-openprocess
- https://learn.microsoft.com/en-us/windows/win32/api/errhandlingapi/nf-errhandlingapi-getlasterror
- https://learn.microsoft.com/en-us/windows/win32/api/psapi/nf-psapi-getmodulefilenameexa
- https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-terminateprocess
- 
- 
- 
- 
- 
- 
*/ 
