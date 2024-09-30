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
	handle, err := windows.OpenProcess(PROCESS_ALL_ACCESS, false, uint32(pid))
	if err != nil {
		return err
	}
	defer windows.CloseHandle(handle)

	err = windows.SetPriorityClass(handle, priority)
	if err != nil {
		return err
	}

	return nil
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
	fmt.Println("\033[32m  read-memory <pid> <address> <size>\033[0m \033[37mRead memory at a specific address of a given process. Address should be in hexadecimal format.\033[0m")
	fmt.Println("\033[32m  write-memory <pid> <address> <data>\033[0m \033[37mWrite data to a specific memory address of a process. Address should be in hexadecimal format.\033[0m")
	fmt.Println("\033[32m  generic\033[0m                \033[37mExecute a generic command to retrieve process information.\033[0m")
	fmt.Println("\033[35mNote:\033[0m")
	fmt.Println("\033[37m- PID is the Process ID you want to perform the operation on.\033[0m")
	fmt.Println("\033[37m- Addresses should be provided in hexadecimal format (e.g., 0x7ffde000).\033[0m")
	fmt.Println("\033[37m- The priority can be specified as low, normal, high, or realtime.\033[0m")
	fmt.Println("\033[37m- Ensure you have sufficient permissions to execute commands on the processes you target.\033[0m")
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("\033[36mThis is a tool for process analysis, is suggested to use the 'generic' args as first one... \033[0m")
		fmt.Println("\033[31mError: No command provided. Please use 'list', 'info <pid>', etc.\033[0m")
		return
	}

	command := os.Args[1]

	switch command {
	case "list":
		fmt.Println("\033[36mThis is a tool for process analysis, is suggested to use the 'generic' args as first one... \033[0m")
		ListInfoProcesses()
	case "info":
		if len(os.Args) < 3 {
			fmt.Println("\033[36mThis is a tool for process analysis, is suggested to use the 'generic' args as first one... \033[0m")
			fmt.Println("\033[31mError: PID required for info command.\033[0m")
			return
		}
		// Convert PID from string to int and call GetProcessInfo
	case "kill":
		if len(os.Args) < 3 {
			fmt.Println("\033[36mThis is a tool for process analysis, is suggested to use the 'generic' args as first one... \033[0m")
			fmt.Println("\033[31mError: PID required for terminate command.\033[0m")
			return
		}

		pid, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("\033[31mError: Invalid PID value. It must be an integer.\033[0m")
			return
		}
		TerminateProcess(pid)

		// Convert PID from string to int and call TerminateProcess
	case "set-priority":
		if len(os.Args) < 4 {
			fmt.Println("\033[36mThis is a tool for process analysis, it is suggested to use the 'generic' args as the first one...\033[0m")
			fmt.Println("\033[31mError: PID and priority required for set-priority command.\033[0m")
			return
		}
	
		fmt.Println("\033[33mThose are the following priority classes:\033[0m")
		fmt.Println("\033[32mIDLE_PRIORITY_CLASS(1) \n" +
			"BELOW_NORMAL_PRIORITY_CLASS(2)\n" +
			"NORMAL_PRIORITY_CLASS(3)\n" +
			"ABOVE_NORMAL_PRIORITY_CLASS(4)\n" +
			"HIGH_PRIORITY_CLASS(5)\n" +
			"REALTIME_PRIORITY_CLASS(6)\n" +
			"INFOS(0)\n\033[0m")
	
		pid, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("\033[31mError: Invalid PID. Please enter a numeric value.\033[0m")
			return
		}
	
		priorityArg := os.Args[3]
		var priority uint32
	
		switch priorityArg {
		case "1":
			priority = IDLE_PRIORITY_CLASS
		case "2":
			priority = BELOW_NORMAL_PRIORITY_CLASS
		case "3":
			priority = NORMAL_PRIORITY_CLASS
		case "4":
			priority = ABOVE_NORMAL_PRIORITY_CLASS
		case "5":
			priority = HIGH_PRIORITY_CLASS
		case "6":
			priority = REALTIME_PRIORITY_CLASS
		case "0":
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

		default:
			fmt.Println("\033[31mError: Invalid priority. Please select a value between 1 and 6.\033[0m")
			return
		}
	
		err = SetProcessPriority(pid, priority)
		if err != nil {
			fmt.Println("\033[31mError setting priority:", err, "\033[0m")
		} else {
			fmt.Println("\033[32mPriority set successfully.\033[0m")
		}
	
	
						
		// Convert PID from string to int and call SetProcessPriority
	case "suspend":
		if len(os.Args) < 3 {
			fmt.Println("\033[36mThis is a tool for process analysis, is suggested to use the 'generic' args as first one... \033[0m")
			fmt.Println("\033[31mError: PID required for suspend command.\033[0m")
			return
		}
		// Convert PID from string to int and call SuspendProcess
	case "resume":
		if len(os.Args) < 3 {
			fmt.Println("\033[36mThis is a tool for process analysis, is suggested to use the 'generic' args as first one... \033[0m")
			fmt.Println("\033[31mError: PID required for resume command.\033[0m")
			return
		}
		// Convert PID from string to int and call ResumeProcess
	case "read-memory":
		if len(os.Args) < 5 {
			fmt.Println("\033[36mThis is a tool for process analysis, is suggested to use the 'generic' args as first one... \033[0m")
			fmt.Println("\033[31mError: PID, address, and size required for read-memory command.\033[0m")
			return
		}
		// Convert PID from string to int, address and size, and call ReadMemory
	case "write-memory":
		if len(os.Args) < 5 {
			fmt.Println("\033[36mThis is a tool for process analysis, is suggested to use the 'generic' args as first one... \033[0m")
			fmt.Println("\033[31mError: PID, address, and data required for write-memory command.\033[0m")
			return
		}
		// Convert PID from string to int, address and data, and call WriteMemory
	case "generic":
		err := generic()
		if err != nil {
			fmt.Println("\033[31mError executing generic command:\033[0m", err)
		}
	default:
		DisplayHelp()
	}
}
