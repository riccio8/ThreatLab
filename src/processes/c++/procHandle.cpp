#include <windows.h>
#include <tlhelp32.h>
#include <iostream>
#include <psapi.h>
#include <cstdlib> 
#include <string>
#include <vector>


#ifndef ABOVE_NORMAL_PRIORITY_CLASS
#define ABOVE_NORMAL_PRIORITY_CLASS 0x00008000  // Priority above NORMAL, below HIGH
#endif

#ifndef BELOW_NORMAL_PRIORITY_CLASS
#define BELOW_NORMAL_PRIORITY_CLASS 0x00004000  // Priority above IDLE, below NORMAL
#endif

#ifndef HIGH_PRIORITY_CLASS
#define HIGH_PRIORITY_CLASS 0x00000080  // For time-critical tasks; preempts NORMAL and IDLE
#endif

#ifndef IDLE_PRIORITY_CLASS
#define IDLE_PRIORITY_CLASS 0x00000040  // Runs only when the system is idle
#endif

#ifndef NORMAL_PRIORITY_CLASS
#define NORMAL_PRIORITY_CLASS 0x00000020  // No special scheduling needs
#endif

#ifndef PROCESS_MODE_BACKGROUND_BEGIN
#define PROCESS_MODE_BACKGROUND_BEGIN 0x00100000  // Start background processing mode
#endif

#ifndef PROCESS_MODE_BACKGROUND_END
#define PROCESS_MODE_BACKGROUND_END 0x00200000  // End background processing mode
#endif

#ifndef REALTIME_PRIORITY_CLASS
#define REALTIME_PRIORITY_CLASS 0x00000100  // Highest possible priority, preempts all others
#endif


void DisplayHelp() {
    std::cout << "This is a tool for process analysis. Use 'list' as the first argument.\n"
              << "Note: Process names should be provided without '.exe' and without uppercase letters if applicable.\n\n"
              << "Usage: ProcHandle <command> [arguments]\n"
              << "Commands:\n"
              << "  list                   - List all running processes on the system.\n"
              << "  info <proc_name>       - Retrieve path for a specific process by its name.\n"
              << "  suspend <proc_name>    - Suspend the process and its threads.\n"
              << "  resume <proc_name>     - Resume the process and its threads.\n"
              << "  kill <proc_name>       - Terminate the process.\n"
              << "  setpriority <proc_name> <priority>  - Set the priority of a process.\n"
              << "    Priority levels:\n"
              << "      idle, below_normal, normal, above_normal, high, realtime\n"
              << "      background_begin, background_end\n";
}



class ProcessManager {
public:
    ProcessManager(){
    }
    ~ProcessManager() {
    }

    void ListProcesses();
    void GetProcessInfo(const std::string& name);
    void SuspendProcess(const std::string& name);
    void ResumeProcess(const std::string& name);
    void Kill(const std::string& name);
    void SetProcessPriority(const std::string& name, const std::string& priority);

private:
    std::vector<DWORD> FindPidByName(const std::string& processName);
    void DisplayPriorityLevels();
};


void ProcessManager::DisplayPriorityLevels() {
    std::cout << "Available Priority Levels:\n"
              << "-----------------------------------------\n"
              << "| Priority Level     | Description      |\n"
              << "-----------------------------------------\n"
              << "| idle               | Runs only when the system is idle                 |\n"
              << "| below_normal       | Priority above IDLE, below NORMAL                 |\n"
              << "| normal             | No special scheduling needs                       |\n"
              << "| above_normal       | Priority above NORMAL, below HIGH                 |\n"
              << "| high               | For time-critical tasks; preempts NORMAL and IDLE |\n"
              << "| realtime           | Highest possible priority, preempts all others    |\n"
              << "| background_begin   | Start background processing mode                  |\n"
              << "| background_end     | End background processing mode                    |\n"
              << "-----------------------------------------\n";
}


// This function use find the pid (process id) by it's name creating a tool that makes a snapshot of the processes
std::vector<DWORD> ProcessManager::FindPidByName(const std::string& processName) {
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

void ProcessManager::ListProcesses() {
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

void ProcessManager::GetProcessInfo(const std::string& names) {
    std::vector<DWORD> pids = FindPidByName(names);
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

void ProcessManager::SuspendProcess(const std::string& name) {
     std::vector<DWORD> pids = FindPidByName(name);

    if (pids.empty()) {
        std::cerr << "No processes found with name: " << name << std::endl;
        return;
    }

    HANDLE hThreadSnapshot = CreateToolhelp32Snapshot(TH32CS_SNAPTHREAD, 0);
    if (hThreadSnapshot == INVALID_HANDLE_VALUE) {
        std::cerr << "Error: Unable to create thread snapshot." << std::endl;
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
                            std::cerr << "Error suspending thread with ID: " 
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


void ProcessManager::ResumeProcess(const std::string& name) {
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
                        DWORD suspendResult = ResumeThread(hThread);
                        if (suspendResult == (DWORD)-1) {
                            std::cout << "Error resuming thread with ID: " 
                                      << threadEntry.th32ThreadID 
                                      << " in process ID: " << pid << std::endl;
                        } else {
                            std::cout << "Successfully resumed thread with ID: " 
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

void ProcessManager::Kill(const std::string& name) {
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

void ProcessManager::SetProcessPriority(const std::string& name, const std::string& priority) {
    std::vector<DWORD> pids = FindPidByName(name);
    if (pids.empty()) {
        std::cerr << "No processes found with name: " << name << std::endl;
        return;
    }
    
    DWORD priorityClass;
    if (priority == "idle") priorityClass = IDLE_PRIORITY_CLASS;
    else if (priority == "below_normal") priorityClass = BELOW_NORMAL_PRIORITY_CLASS;
    else if (priority == "normal") priorityClass = NORMAL_PRIORITY_CLASS;
    else if (priority == "above_normal") priorityClass = ABOVE_NORMAL_PRIORITY_CLASS;
    else if (priority == "high") priorityClass = HIGH_PRIORITY_CLASS;
    else if (priority == "realtime") priorityClass = REALTIME_PRIORITY_CLASS;
    else if (priority == "background_begin") priorityClass = PROCESS_MODE_BACKGROUND_BEGIN;
    else if (priority == "background_end") priorityClass = PROCESS_MODE_BACKGROUND_END;
    else {
        std::cerr << "Invalid priority level: " << priority << std::endl;
        DisplayPriorityLevels();
        return;
    }

    std::cout << "Setting process priority to " << priority << std::endl;
    
    for (DWORD pid : pids) {
        HANDLE hProcess = OpenProcess(PROCESS_ALL_ACCESS, FALSE, pid);
        if (hProcess == NULL) {
            std::cerr << "Error opening process: " << GetLastError() << std::endl;
            continue;
        }
        
        if (!SetPriorityClass(hProcess, priorityClass)) {
            std::cerr << "Error setting priority: " << GetLastError() << std::endl;
        }
        
        CloseHandle(hProcess);
    }
}


int main(int argc, char* argv[]) {
    if (argc < 2) {
        DisplayHelp();
        return 1;
    }

    ProcessManager pm;
    std::string command = argv[1];

    if (command == "list") {
        pm.ListProcesses();
    } else if (command == "help") {
        DisplayHelp();
    }else if (command == "info") {
        if (argc < 3) {
            std::cout << "Usage: info <process_name>" << std::endl;
            return 1;
        }
        pm.GetProcessInfo(argv[2]);
    } else if (command == "kill") {
        if (argc < 3) {
            std::cout << "Usage: kill <process_name>" << std::endl;
            return 1;
        }
        pm.Kill(argv[2]);
    } else if (command == "suspend") {
        if (argc < 3) {
            std::cout << "Usage: suspend <process name>" << std::endl;
            return 1;
        }
        pm.SuspendProcess(argv[2]);
    } else if (command == "resume") {
        if (argc < 3) {
            std::cout << "Usage: resume <process name>" << std::endl;
            return 1;
        }
        pm.ResumeProcess(argv[2]);
    } else if (command == "setpriority") {
        if (argc < 4) {
            std::cout << "Usage: setpriority <process_name> <priority>" << std::endl;
            return 1;
        }
        pm.SetProcessPriority(argv[2], argv[3]); 
        
    }else {
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


//     const std::string name = names+".exe";
