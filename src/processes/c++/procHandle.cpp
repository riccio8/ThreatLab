#include <windows.h>
#include <tlhelp32.h>
#include <iostream>
#include <psapi.h>
#include <cstdlib> 
#include <string>
#include <cstdio>
#include <memory>
#include <stdexcept>
#include <array>
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
    std::cout << "This is a tool for process and network analysis. Use the corresponding command as the first argument.\n"
              << "Note: Process names should be provided without '.exe' and without uppercase letters if applicable.\n\n"
              << "Usage: ProcHandle <command> [arguments]\n"
              << "Commands:\n"
              << "\nProcess Management:\n"
              << "  list                   - List all running processes on the system.\n"
              << "  info <proc_name>       - Retrieve path for a specific process by its name.\n"
              << "  suspend <proc_name>    - Suspend the process and its threads.\n"
              << "  resume <proc_name>     - Resume the process and its threads.\n"
              << "  kill <proc_name>       - Terminate the process.\n"
              << "  setpriority <proc_name> <priority>  - Set the priority of a process.\n"
              << "    Priority levels:\n"
              << "      idle, below_normal, normal, above_normal, high, realtime\n"
              << "      background_begin, background_end\n"
              << "\nNetwork Commands:\n"
              << "  net <command>      - Display network status and statistics.\n"
              << "    Available sub-commands for 'netstat':\n"
              << "      all                 - Show all TCP/UDP connections with numeric addresses.\n"
              << "      names               - Show all TCP/UDP connections with resolved names.\n"
              << "      programs            - Show connections along with the programs using them.\n"
              << "      processid           - Show connections with their respective process IDs.\n"
              << "      listening           - Show only listening connections.\n"
              << "      established         - Show only established connections.\n"
              << "      routing             - Show the routing table.\n"
              << "      stats               - Show network statistics.\n"
              << "\nARP Commands:\n"
              << "  arp <command>          - Display or modify the ARP (Address Resolution Protocol) table.\n"
              << "    Available sub-commands for 'arp':\n"
              << "      table               - Show the ARP table.\n"
              << "      cache               - Show statistics about the ARP cache.\n"
              << "      delete <ip_address> - Delete an ARP entry for the specified IP address.\n"
              << "\nExample usage:\n"
              << "  ProcHandle list\n"
              << "  ProcHandle netstat all\n"
              << "  ProcHandle arp delete 192.168.1.1\n";
}



class NetworkManager {
public:
    NetworkManager() {}
    ~NetworkManager() {}

    // Functions for executing netstat commands
    void ExecuteNetstatAllConnections();
    void ExecuteNetstatConnectionsAndNames();
    void ExecuteNetstatWithPrograms();
    void ExecuteNetstatWithProcessId();
    void ExecuteNetstatListeningConnections();
    void ExecuteNetstatEstablishedConnections();
    void ExecuteNetstatRoutingTable();
    void ExecuteNetstatNetworkStats();
    void NetProcess(const std::string& names);

    // Functions for executing arp commands
    void ExecuteArpTable();
    void ExecuteArpCacheStatistics();
    void ExecuteArpDeleteEntry(const std::string& ipAddress);
    

private:
    std::string exec(const std::string& cmd);
};

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
    
    std::vector<DWORD> FindPidByName(const std::string& processName);
    
    

private:
    void DisplayPriorityLevels();
};

// Executes a command and returns the output as a string
std::string NetworkManager::exec(const std::string& cmd) {
    std::shared_ptr<FILE> pipe(popen(cmd.c_str(), "r"), [](FILE* f) { if (f) fclose(f); });
    if (!pipe) return "";

    char buffer[128];
    std::string result = "";
    while (fgets(buffer, sizeof(buffer), pipe.get()) != nullptr) {
        result += buffer;
    }
    return result;
}



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
    if (pids.empty()) {
        std::cerr << "No processes found with name: " << name << std::endl;
        return;
    }
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


// Execute "netstat -an": Shows TCP/UDP connections with numeric addresses
void NetworkManager::ExecuteNetstatAllConnections() {
    std::string cmd = "powershell.exe netstat -an";
    std::string output = exec(cmd);
    std::cout << "TCP/UDP connections with numeric addresses:\n" << output << std::endl;
}

// Execute "netstat -a": Shows TCP/UDP connections with resolved names
void NetworkManager::ExecuteNetstatConnectionsAndNames() {
    std::string cmd = "powershell.exe netstat -a";
    std::string output = exec(cmd);
    std::cout << "TCP/UDP connections with resolved names:\n" << output << std::endl;
}

// Execute "netstat -b": Shows connections with the programs using them
void NetworkManager::ExecuteNetstatWithPrograms() {
    std::string cmd = "powershell.exe netstat -b";
    std::string output = exec(cmd);
    std::cout << "Connections with the programs using them:\n" << output << std::endl;
}

// Execute "netstat -o": Shows connections with the process ID
void NetworkManager::ExecuteNetstatWithProcessId() {
    std::string cmd = "powershell.exe netstat -o";
    std::string output = exec(cmd);
    std::cout << "Connections with the process ID:\n" << output << std::endl;
}

// Execute "netstat -an | findstr LISTENING": Shows only listening connections
void NetworkManager::ExecuteNetstatListeningConnections() {
    std::string cmd = "powershell.exe netstat -an | findstr LISTENING";
    std::string output = exec(cmd);
    std::cout << "Listening connections:\n" << output << std::endl;
}

// Execute "netstat -an | findstr ESTABLISHED": Shows only established connections
void NetworkManager::ExecuteNetstatEstablishedConnections() {
    std::string cmd = "powershell.exe netstat -an | findstr ESTABLISHED";
    std::string output = exec(cmd);
    std::cout << "Established connections:\n" << output << std::endl;
}

// Execute "netstat -r": Shows the routing table
void NetworkManager::ExecuteNetstatRoutingTable() {
    std::string cmd = "powershell.exe netstat -r";
    std::string output = exec(cmd);
    std::cout << "Routing table:\n" << output << std::endl;
}

// Execute "netstat -e": Shows network statistics
void NetworkManager::ExecuteNetstatNetworkStats() {
    std::string cmd = "powershell.exe netstat -e";
    std::string output = exec(cmd);
    std::cout << "Network statistics:\n" << output << std::endl;
}

// Execute "arp -a": Shows the ARP table
void NetworkManager::ExecuteArpTable() {
    std::string cmd = "powershell.exe arp -a";
    std::string output = exec(cmd);
    std::cout << "ARP table:\n" << output << std::endl;
}

// Execute "arp -s <ip> <mac>": Adds a static ARP entry
void NetworkManager::ExecuteArpCacheStatistics() {
    std::string cmd = "powershell.exe arp -s";
    std::string output = exec(cmd);
    std::cout << "ARP cache statistics:\n" << output << std::endl;
}

// Execute "arp -d <ip>": Deletes an ARP entry by IP address
void NetworkManager::ExecuteArpDeleteEntry(const std::string& ipAddress) {
    std::string cmd = "powershell.exe arp -d " + ipAddress;
    std::string output = exec(cmd);
    std::cout << "Deleted ARP entry for IP address " << ipAddress << ":\n" << output << std::endl;
}

void NetworkManager::NetProcess(const std::string& names) {
    ProcessManager pm;
    std::vector<DWORD> pids = pm.FindPidByName(names);
        if (pids.empty()) {
        std::cout << "No processes found with the given name." << std::endl;
        return;
    }
    
    std::vector<std::string> options = {"-t", "-u", "-l", "-n", "-p"};
    
    for (auto id : pids) {
        for (const auto& option : options) {
            std::string cmd = "netstat " + option + " | grep " + std::to_string(id);
            std::string output = exec(cmd); 
            std::cout << "Connection for process\t" << names 
                      << " with specific id\t" << id 
                      << " option\t" << option 
                      << "\noutput:\n" << output << std::endl; 
        }
    }
}


int main(int argc, char* argv[]) {
    if (argc < 2) {
        DisplayHelp();
        return 1;
    }

    ProcessManager pm;
    NetworkManager nm;  
    std::string command = argv[1];

    if (command == "list") {
        pm.ListProcesses();
    } else if (command == "help") {
        DisplayHelp();
    } else if (command == "info") {
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
    } else if (command == "net") { 
        if (argc < 3) {
            std::cout << "Usage: net <command>" << std::endl;
            return 1;
        }

        std::string netstatCommand = argv[2];
        if (netstatCommand == "all") {
            nm.ExecuteNetstatAllConnections();
        } else if (netstatCommand == "names") {
            nm.ExecuteNetstatConnectionsAndNames();
        } else if (netstatCommand == "programs") {
            nm.ExecuteNetstatWithPrograms();
        } else if (netstatCommand == "processid") {
            nm.ExecuteNetstatWithProcessId();
        } else if (netstatCommand == "listening") {
            nm.ExecuteNetstatListeningConnections();
        } else if (netstatCommand == "established") {
            nm.ExecuteNetstatEstablishedConnections();
        } else if (netstatCommand == "routing") {
            nm.ExecuteNetstatRoutingTable();
        } else if (netstatCommand == "stats") {
            nm.ExecuteNetstatNetworkStats();
            
        } else if (netstatCommand == "process"){
            if (argc < 4) {
                std::cout << "Usage: net process <process_name>" << std::endl;
                return 1;
            }
            nm.NetProcess(argv[3]);
        }
        
        else {
            std::cout << "Unknown netstat command: " << netstatCommand << std::endl;
            DisplayHelp();
        }
    } else if (command == "arp") {  
        if (argc < 3) {
            std::cout << "Usage: arp <command>" << std::endl;
            return 1;
        }

        std::string arpCommand = argv[2];
        if (arpCommand == "table") {
            nm.ExecuteArpTable();
        } else if (arpCommand == "cache") {
            nm.ExecuteArpCacheStatistics();
        } else if (arpCommand == "delete") {
            if (argc < 4) {
                std::cout << "Usage: arp delete <ip_address>" << std::endl;
                return 1;
            }
            nm.ExecuteArpDeleteEntry(argv[3]);
        } else {
            std::cout << "Unknown arp command: " << arpCommand << std::endl;
            DisplayHelp();
        }
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
*/ 


//     const std::string name = names+".exe";
