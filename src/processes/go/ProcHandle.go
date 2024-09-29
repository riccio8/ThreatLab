package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var entry syscall.ProcessEntry32

// Function to list all running processes
func ListInfoProcesses() {
	fmt.Println("\033[36mListing all processes...\033[0m")
	snapshot, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPALL, 0)
	if err != nil {
		fmt.Println("\033[31mError creating process snapshot:\033[0m", err)
		return
	}
	defer syscall.CloseHandle(snapshot)

	entry.Size = uint32(unsafe.Sizeof(entry))

	err = syscall.Process32First(snapshot, &entry)
	if err != nil {
		fmt.Println("\033[31mError retrieving first process:\033[0m", err)
		return
	}

	fmt.Println("\033[32mProcesses:\033[0m")
	for {
		processName := syscall.UTF16ToString(entry.ExeFile[:]) // convert the exefile name string from UTF-16 to UTF-8

		fmt.Printf("\033[32mPid: %d\tFile Name: %s\tThread: %d\tHeap Allocation: %d\tProcess Flags: %d\033[0m\n",
			entry.ProcessID, processName, entry.Threads, entry.DefaultHeapID, entry.Flags)

		err = syscall.Process32Next(snapshot, &entry)
		if err != nil {
			fmt.Println("\033[33mNo more processes...\033[0m")
			break // No more processes to enumerate
		}
	}
}

// Function to get detailed information about a specific process

func GetProcessInfo(pid int) {
	fmt.Printf("\033[36mRetrieving information for PID: %d...\033[0m\n", pid) // Cyan for info retrieval

	// Open the process with the necessary permissions
	hProcess, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ, false, uint32(pid))
	if err != nil {
		fmt.Println("\033[31mError opening process:\033[0m", err) // Red for error messages
		return
	}
	defer windows.CloseHandle(hProcess)

	// Get the full name of the process
	var processName [windows.MAX_PATH]uint16
	processPathLength := uint32(len(processName))
	err = windows.QueryFullProcessImageName(hProcess, 0, &processName[0], &processPathLength)
	if err != nil {
		fmt.Println("\033[31mError retrieving process name:\033[0m", err) // Red for error messages
		return
	}

	// Convert the process name to a string
	name := windows.UTF16ToString(processName[:])

	// Print the process information
	fmt.Printf("\033[32mPID: %d\tName: %s\033[0m\n", pid, name) // Green for process information
}

func generic() error {
	proc := exec.Command("ps")
	// Utilizziamo Output() per ottenere l'output del comando
	info, err := proc.Output()

	if err != nil {
		fmt.Println("Error while getting generic infos: ", err)
		return err
	}

	// Stampiamo l'output del comando
	fmt.Println(string(info))
	return nil
}

// Function to terminate a process by its PID
func TerminateProcess(pid int) {
	fmt.Printf("\033[31mTerminating process with PID: %d\033[0m\n", pid)
	// Logic to terminate process
}

// Function to set the priority of a process
func SetProcessPriority(pid int, priority string) {
	fmt.Printf("\033[33mSetting priority for PID: %d to %s\033[0m\n", pid, priority)
	// Logic to set process priority
}

// Function to suspend a process by its PID
func SuspendProcess(pid int) {
	fmt.Printf("\033[35mSuspending process with PID: %d\033[0m\n", pid)
	// Logic to suspend process
}

// Function to resume a suspended process
func ResumeProcess(pid int) {
	fmt.Printf("\033[36mResuming process with PID: %d\033[0m\n", pid)
	// Logic to resume process
}

// Function to read memory from a specific process
func ReadMemory(pid int, address string, size int) {
	fmt.Printf("\033[37mReading memory at address: %s for PID: %d\033[0m\n", address, pid)
	// Logic to read memory
}

// Function to write data to a specific memory address of a process
func WriteMemory(pid int, address string, data string) {
	fmt.Printf("\033[32mWriting data to memory address: %s for PID: %d\033[0m\n", address, pid)
	// Logic to write memory
}

func main() {
	fmt.Println("\033[36mThis is a tool for process analysis, is suggested to use the 'general' args as first one... \033[0m")
	
	time.Sleep(1 * time.Second) 
	
	if len(os.Args) < 2 {
		fmt.Println("\033[31mError: No command provided. Please use 'list', 'info <pid>', etc.\033[0m")
		return
	}

	command := os.Args[1]

	switch command {
	case "list":
		ListInfoProcesses()
	case "info":
		if len(os.Args) < 3 {
			fmt.Println("\033[31mError: PID required for info command.\033[0m")
			return
		}
		// Convert PID from string to int and call GetProcessInfo
	case "terminate":
		if len(os.Args) < 3 {
			fmt.Println("\033[31mError: PID required for terminate command.\033[0m")
			return
		}
		// Convert PID from string to int and call TerminateProcess
	case "set-priority":
		if len(os.Args) < 4 {
			fmt.Println("\033[31mError: PID and priority required for set-priority command.\033[0m")
			return
		}
		// Convert PID from string to int and call SetProcessPriority
	case "suspend":
		if len(os.Args) < 3 {
			fmt.Println("\033[31mError: PID required for suspend command.\033[0m")
			return
		}
		// Convert PID from string to int and call SuspendProcess
	case "resume":
		if len(os.Args) < 3 {
			fmt.Println("\033[31mError: PID required for resume command.\033[0m")
			return
		}
		// Convert PID from string to int and call ResumeProcess
	case "read-memory":
		if len(os.Args) < 5 {
			fmt.Println("\033[31mError: PID, address, and size required for read-memory command.\033[0m")
			return
		}
		// Convert PID from string to int, address and size, and call ReadMemory
	case "write-memory":
		if len(os.Args) < 5 {
			fmt.Println("\033[31mError: PID, address, and data required for write-memory command.\033[0m")
			return
		}

	case "generic":
		generic()
		// Convert PID from string to int and call WriteMemory
	default:
		fmt.Println("\033[31mError: Unknown command. Please use 'list', 'info', etc.\033[0m")
	}
}
