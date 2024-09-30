package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	PROCESS_ALL_ACCESS = 0x1F0FFF
)

const (
	IDLE_PRIORITY_CLASS           = 0x00000040
	BELOW_NORMAL_PRIORITY_CLASS   = 0x00040000
	NORMAL_PRIORITY_CLASS         = 0x00000020
	ABOVE_NORMAL_PRIORITY_CLASS   = 0x00080000
	HIGH_PRIORITY_CLASS           = 0x00000080
	REALTIME_PRIORITY_CLASS       = 0x00000100
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

	entry = syscall.ProcessEntry32{}
}

// Function to get detailed information about a specific process
func GetProcessInfo(pid int) {
	fmt.Printf("\033[36mRetrieving information for PID: %d...\033[0m\n", pid)

	// Open the process with the necessary permissions
	hProcess, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ, false, uint32(pid))
	if err != nil {
		fmt.Println("\033[31mError opening process:\033[0m", err)
		return
	}
	defer windows.CloseHandle(hProcess)

	// Get the full name of the process
	var processName [windows.MAX_PATH]uint16
	processPathLength := uint32(len(processName))
	err = windows.QueryFullProcessImageName(hProcess, 0, &processName[0], &processPathLength)
	if err != nil {
		fmt.Println("\033[31mError retrieving process name:\033[0m", err)
		return
	}

	// Convert the process name to a string
	name := windows.UTF16ToString(processName[:])

	// Print the process information
	fmt.Printf("\033[32mPID: %d\tName: %s\033[0m\n", pid, name)

	entry = syscall.ProcessEntry32{}
}

func generic() error {
	cmd := exec.Command("powershell", "-Command", "Get-Process") // working in powershell, not cmd
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	entry = syscall.ProcessEntry32{}
	fmt.Println(string(output))
	return nil
}

// Function to terminate a process by its PID
func TerminateProcess(pid int) {
	fmt.Printf("\033[31mTerminating process with PID: %d...\033[0m\n", pid)

	// Esegui il comando taskkill per terminare il processo
	cmd := exec.Command("taskkill", "/PID", fmt.Sprint(pid), "/F")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("\033[31mError terminating process:\033[0m %s\n", err)
		return
	}

	fmt.Printf("\033[32mProcess terminated successfully:\033[0m %s\n", string(output))
}

// Function to set the priority of a process
func SetProcessPriority(pid int, priority uint32) error {
	fmt.Printf("\033[33mSetting priority for PID: %d to %d\033[0m\n", pid, priority)

	// Open process with required access
	handle, err := windows.OpenProcess(windows.PROCESS_SET_INFORMATION, false, uint32(pid))
	if err != nil {
		return fmt.Errorf("failed to open process: %v", err)
	}
	defer windows.CloseHandle(handle)

	// Set the priority class
	if err := windows.SetPriorityClass(handle, priority); err != nil {
		return fmt.Errorf("failed to set priority class: %v", err)
	}

	return nil
}

func GetThreadPriority(hThread syscall.Handle) {
	fmt.Println("\033[1;34mThread possibilities:\033[0m")
	fmt.Println("\033[1;32mReturn code/value \t Description\033[0m")
	fmt.Println("--------------------------------------------------")
	fmt.Println("\033[1;33mTHREAD_PRIORITY_ABOVE_NORMAL \t 1 \t\t Priority 1 point above the priority class.\033[0m")
	fmt.Println("\033[1;33mTHREAD_PRIORITY_BELOW_NORMAL \t -1 \t\t Priority 1 point below the priority class.\033[0m")
	fmt.Println("\033[1;33mTHREAD_PRIORITY_HIGHEST \t 2 \t\t Priority 2 points above the priority class.\033[0m")
	fmt.Println("\033[1;33mTHREAD_PRIORITY_IDLE \t\t -15 \t\t Base priority of 1 for IDLE_PRIORITY_CLASS, BELOW_NORMAL_PRIORITY_CLASS, NORMAL_PRIORITY_CLASS, ABOVE_NORMAL_PRIORITY_CLASS, or HIGH_PRIORITY_CLASS processes, and a base priority of 16 for REALTIME_PRIORITY_CLASS processes.\033[0m")
	fmt.Println("\033[1;33mTHREAD_PRIORITY_LOWEST \t\t -2 \t\t Priority 2 points below the priority class.\033[0m")
	fmt.Println("\033[1;33mTHREAD_PRIORITY_NORMAL \t\t 0 \t\t Normal priority for the priority class.\033[0m")
	fmt.Println("\033[1;33mTHREAD_PRIORITY_TIME_CRITICAL \t 15 \t\t Base-priority level of 15 for IDLE_PRIORITY_CLASS, BELOW_NORMAL_PRIORITY_CLASS, NORMAL_PRIORITY_CLASS, ABOVE_NORMAL_PRIORITY_CLASS, or HIGH_PRIORITY_CLASS processes, and a base-priority level of 31 for REALTIME_PRIORITY_CLASS processes.\033[0m")
	fmt.Println("\033[1;34mThread priority: \033[0m", hThread)
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

func DisplayHelp() {
	fmt.Println("\033[36mThis is a tool for process analysis, is suggested to use the 'generic' args as first one... \033[0m")
	fmt.Println("Use help for display this massage...")
	fmt.Println("\033[36mUsage: ProcHandle <command> [arguments]\033[0m")
	fmt.Println("\033[33mCommands:\033[0m")
	fmt.Println("\033[32m  list\033[0m                   \033[37mList all running processes on the system.\033[0m")
	fmt.Println("\033[32m  info <pid>\033[0m             \033[37mRetrieve detailed information for a specific process by its PID.\033[0m")
	fmt.Println("\033[32m  terminate <pid>\033[0m        \033[37mTerminate a process by its PID.\033[0m")
	fmt.Println("\033[32m  set-priority <pid> <priority>\033[0m \033[37mSet the priority for a process. Priority can be one of: low, normal, high, realtime.\033[0m")
	fmt.Println("\033[32m  suspend <pid>\033[0m          \033[37mSuspend a process by its PID.\033[0m")
	fmt.Println("\033[32m  resume <pid>\033[0m           \033[37mResume a suspended process by its PID.\033[0m")
	fmt.Println("\033[32m  read-memory <pid> <address> <size>\033[0m \033[37mRead memory at a specific address of a process.\033[0m")
	fmt.Println("\033[32m  write-memory <pid> <address> <data>\033[0m \033[37mWrite data to a specific memory address of a process.\033[0m")
}

func main() {
	if len(os.Args) < 2 {
		DisplayHelp()
		return
	}
// -------------------------------------------------------------------------------------------------------------------------------------------------------------
	command := os.Args[1]

	switch command {
	case "list":
		ListInfoProcesses()
		
		// -------------------------------------------------------------------------------------------------------------------------------------------------------------
		
	case "info":
		if len(os.Args) < 3 {
			fmt.Println("\033[31mError: Missing PID argument for 'info' command.\033[0m")
			return
		}
		pid, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("\033[31mError: Invalid PID value:\033[0m", os.Args[2])
			return
		}
		GetProcessInfo(pid)
		
		// -------------------------------------------------------------------------------------------------------------------------------------------------------------
		
	case "terminate":
		if len(os.Args) < 3 {
			fmt.Println("\033[31mError: Missing PID argument for 'terminate' command.\033[0m")
			return
		}
		pid, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("\033[31mError: Invalid PID value:\033[0m", os.Args[2])
			return
		}
		TerminateProcess(pid)
		
		// -------------------------------------------------------------------------------------------------------------------------------------------------------------
		
	case "set-priority":
		if len(os.Args) < 3 {
			fmt.Println("\033[31mError: Missing PID or priority argument for 'set-priority' command.\033[0m")
			return
		}
	
		if os.Args[2] == "info" {
			fmt.Println("Process priority class\tThread priority level\tBase priority")
			fmt.Println("IDLE_PRIORITY_CLASS\t\tTHREAD_PRIORITY_IDLE\t1")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_LOWEST\t2")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_BELOW_NORMAL\t3")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_NORMAL\t4")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_ABOVE_NORMAL\t5")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_HIGHEST\t6")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_TIME_CRITICAL\t15")
			fmt.Println("BELOW_NORMAL_PRIORITY_CLASS\tTHREAD_PRIORITY_IDLE\t1")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_LOWEST\t4")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_BELOW_NORMAL\t5")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_NORMAL\t6")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_ABOVE_NORMAL\t7")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_HIGHEST\t8")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_TIME_CRITICAL\t15")
			fmt.Println("NORMAL_PRIORITY_CLASS\t\tTHREAD_PRIORITY_IDLE\t1")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_LOWEST\t6")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_BELOW_NORMAL\t7")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_NORMAL\t8")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_ABOVE_NORMAL\t9")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_HIGHEST\t10")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_TIME_CRITICAL\t15")
			fmt.Println("ABOVE_NORMAL_PRIORITY_CLASS\tTHREAD_PRIORITY_IDLE\t1")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_LOWEST\t8")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_BELOW_NORMAL\t9")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_NORMAL\t10")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_ABOVE_NORMAL\t11")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_HIGHEST\t12")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_TIME_CRITICAL\t15")
			fmt.Println("HIGH_PRIORITY_CLASS\t\tTHREAD_PRIORITY_IDLE\t1")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_LOWEST\t11")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_BELOW_NORMAL\t12")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_NORMAL\t13")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_ABOVE_NORMAL\t14")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_HIGHEST\t15")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_TIME_CRITICAL\t15")
			fmt.Println("REALTIME_PRIORITY_CLASS\t\tTHREAD_PRIORITY_IDLE\t16")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_LOWEST\t22")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_BELOW_NORMAL\t23")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_NORMAL\t24")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_ABOVE_NORMAL\t25")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_HIGHEST\t26")
			fmt.Println("\t\t\t\tTHREAD_PRIORITY_TIME_CRITICAL\t31")
			return // return and exit
		}
	
		if len(os.Args) < 4 {
			fmt.Println("\033[31mError: Missing priority argument for 'set-priority' command.\033[0m")
			return
		}
	
		priorityStr := os.Args[3]
		var priority uint32
	
		switch priorityStr {
		case "low":
			priority = BELOW_NORMAL_PRIORITY_CLASS
		case "normal":
			priority = NORMAL_PRIORITY_CLASS
		case "high":
			priority = HIGH_PRIORITY_CLASS
		case "realtime":
			priority = REALTIME_PRIORITY_CLASS
		default:
			fmt.Println("\033[31mError: Invalid priority value. Choose from: low, normal, high, realtime.\033[0m")
			return
		}
	
		pid, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("\033[31mError: Invalid PID value:\033[0m", os.Args[2])
			return
		}
	
		if err := SetProcessPriority(pid, priority); err != nil {
			fmt.Println("\033[31mError setting process priority:\033[0m", err)
		}
	
	
	// -----------------------------------------------------------------------------------------------------------------------------------------------------------------
	
	case "suspend":
		if len(os.Args) < 3 {
			fmt.Println("\033[31mError: Missing PID argument for 'suspend' command.\033[0m")
			return
		}
		pid, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("\033[31mError: Invalid PID value:\033[0m", os.Args[2])
			return
		}
		SuspendProcess(pid)
		
		// -------------------------------------------------------------------------------------------------------------------------------------------------------------
		
	case "resume":
		if len(os.Args) < 3 {
			fmt.Println("\033[31mError: Missing PID argument for 'resume' command.\033[0m")
			return
		}
		pid, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("\033[31mError: Invalid PID value:\033[0m", os.Args[2])
			return
		}
		ResumeProcess(pid)
		
		// -------------------------------------------------------------------------------------------------------------------------------------------------------------
		
	case "read-memory":
		if len(os.Args) < 5 {
			fmt.Println("\033[31mError: Missing PID, address or size argument for 'read-memory' command.\033[0m")
			return
		}
		pid, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("\033[31mError: Invalid PID value:\033[0m", os.Args[2])
			return
		}
		address := os.Args[3]
		size, err := strconv.Atoi(os.Args[4])
		if err != nil {
			fmt.Println("\033[31mError: Invalid size value:\033[0m", os.Args[4])
			return
		}
		ReadMemory(pid, address, size)
		
		// -------------------------------------------------------------------------------------------------------------------------------------------------------------
		
	case "write-memory":
		if len(os.Args) < 5 {
			fmt.Println("\033[31mError: Missing PID, address or data argument for 'write-memory' command.\033[0m")
			return
		}
		pid, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("\033[31mError: Invalid PID value:\033[0m", os.Args[2])
			return
		}
		address := os.Args[3]
		data := os.Args[4]
		WriteMemory(pid, address, data)
		
		// -------------------------------------------------------------------------------------------------------------------------------------------------------------
		
	default:
		fmt.Println("\033[31mError: Unknown command:\033[0m", command)
		DisplayHelp()
	}
}
