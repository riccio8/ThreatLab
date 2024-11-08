#include <windows.h>
#include <tlhelp32.h>
#include <iostream>
#include <psapi.h>
#include <string>
#include <vector>

void DisplayHelp() {
    std::cout << "This is a tool for process analysis, it is suggested to use the 'generic' args as first one..." << std::endl;
    std::cout << "Usage: ProcHandle <command> [arguments]" << std::endl;
    std::cout << "Commands:" << std::endl;
    std::cout << "  list                    List all running processes on the system." << std::endl;
    std::cout << "  info <proc_name>       Retrieve detailed information for a specific process by its PID." << std::endl;
}

std::vector<DWORD> FindPidByName(const std::string& processName) {
    std::vector<DWORD> pids;
    PROCESSENTRY32 pe32;
    pe32.dwSize = sizeof(PROCESSENTRY32);

    HANDLE snapshot = CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0);
    if (snapshot == INVALID_HANDLE_VALUE) return pids;

    if (Process32First(snapshot, &pe32)) {
        do {
            if (processName == pe32.szExeFile) {
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

void GetProcessInfo(const std::string& name) {
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
    } else {
        std::cout << "Error: Unknown command: " << command << std::endl;
        DisplayHelp();
    }

    return 0;
}
