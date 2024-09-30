package main

import (
	"fmt"
	"os"
	"os/exec"
	"bytes"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)


var (
	kernel32              = syscall.NewLazyDLL("kernel32.dll")
	procOpenProcess       = kernel32.NewProc("OpenProcess")
	procEnumProcesses     = kernel32.NewProc("EnumProcesses")
	procEnumProcessThreads = kernel32.NewProc("EnumProcessThreads")
	procSuspendThread      = kernel32.NewProc("SuspendThread")
	procCloseHandle        = kernel32.NewProc("CloseHandle")
)


const (
	PROCESS_QUERY_INFORMATION = 0x0400
	PROCESS_SUSPEND_RESUME   = 0x0800
)

const (
	INVALID_HANDLE_VALUE = ^uintptr(0)
	PROCESS_ALL_ACCESS        = 0x1F0FFF
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

//************************************************************************************************************************************************************************************************


func FindPidByNamePowerShell(processName string) ([]int, error) {
    // Build the PowerShell command
    // `Get-Process -Name %s | Select-Object -ExpandProperty Id` finds processes with the given name and gets their IDs.
    powershellCommand := fmt.Sprintf(`Get-Process -Name %s | Where-Object { $_.Parent -eq $null } | Select-Object -ExpandProperty Id`, processName)

    // Execute the command using `powershell -Command <command>`
    // `exec.Command` creates a command to be executed. Here, it runs PowerShell with the specified command.
    cmd := exec.Command("powershell", "-Command", powershellCommand)

    // Capture the output using a buffer
    // `bytes.Buffer` is used to capture the output of the command.
    var out bytes.Buffer
    cmd.Stdout = &out // Set `cmd.Stdout` to point to `out` so output is stored in this buffer.

    // Run the command and check for errors
    // `cmd.Run()` executes the command. If an error occurs, it returns an error.
    err := cmd.Run()
    if err != nil {
        return nil, fmt.Errorf("error executing PowerShell command: %v", err) // Format error if command fails.
    }

    // Parse the output to get the PIDs
    // `out.String()` converts the captured output buffer to a string. `strings.TrimSpace()` removes leading/trailing whitespace.
    output := strings.TrimSpace(out.String())
    if output == "" {
        return nil, fmt.Errorf("no process found with the name: %s", processName) // Return error if no process is found.
    }

    // Split the output to get individual PIDs
    // `strings.Split()` splits the output string by newline to get individual PIDs.
    lines := strings.Split(output, "\n")
    var pids []int

    // Loop through each line, convert to int, and collect valid PIDs
    // `strings.TrimSpace()` removes any extra whitespace, and `strconv.Atoi()` converts strings to integers.
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if pid, err := strconv.Atoi(line); err == nil { // Convert to integer if valid.
            pids = append(pids, pid) // Add the PID to the list of PIDs.
        }
    }

    // Check if we found any PIDs
    if len(pids) == 0 {
        return nil, fmt.Errorf("no process found with the name: %s", processName) // Return error if no PIDs were found.
    }

    return pids, nil // Return the list of PIDs.
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

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
			// No more processes to enumerate
		}
	}

	entry = syscall.ProcessEntry32{}
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// Function to get detailed information about a specific process
func GetProcessInfo(name string) {
    for i := 0; i < 1; i++ { 
        pids, err := FindPidByNamePowerShell(name)
        if err != nil {
            fmt.Println("\033[31mError finding process:\033[0m", err)
            return
        }

        if len(pids) > 0 {
            for _, pid := range pids {
                pidValue := uint32(pid) // Convert the PID to uint32
                fmt.Printf("\033[36mRetrieving information for PID: %d...\033[0m\n", pidValue)

                hProcess, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ, false, pidValue)
                if err != nil {
                    fmt.Println("\033[31mError opening process:\033[0m", err)
                    continue // Skip to the next PID if there's an error
                }
                defer windows.CloseHandle(hProcess)

                var processName [windows.MAX_PATH]uint16
                processPathLength := uint32(len(processName))
                err = windows.QueryFullProcessImageName(hProcess, 0, &processName[0], &processPathLength)
                if err != nil {
                    fmt.Println("\033[31mError retrieving process name:\033[0m", err)
                    continue
                }

                processNameStr := windows.UTF16ToString(processName[:])
                fmt.Printf("\033[32mPID: %d\tName: %s\033[0m\n", pidValue, processNameStr)
            }
        } else {
            fmt.Println("\033[33mNo processes found with the given name.\033[0m")
        }
    }
}

	
	
	// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	
func generic() error {
	cmd := exec.Command("powershell", "-Command", "Get-Process | Where-Object { $_.Parent -eq $null }")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	fmt.Println(string(output))
	return nil
}


// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// Function to terminate a process by its PID
func TerminateProcess(pid int) {
	fmt.Printf("\033[31mTerminating process with PID: %d...\033[0m\n", pid)

	// Execute the taskkill command to terminate the process
	cmd := exec.Command("taskkill", "/PID", fmt.Sprint(pid), "/F")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("\033[31mError terminating process:\033[0m %s\n", err)
		return
	}

	fmt.Printf("\033[32mProcess terminated successfully:\033[0m %s\n", string(output))
}
// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// Function to set the priority of a process
func SetProcessPriority(proc string, priority uint32) error {
	fmt.Printf("\033[33mSetting priority for PROCESS: %s to %d\033[0m\n", proc, priority)
	
	pid, err := FindPidByNamePowerShell(proc)
    if err != nil{
        fmt.Println(err)
    }
	

	// Open the process with required access
	handle, err := windows.OpenProcess(windows.PROCESS_SET_INFORMATION|windows.PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if err != nil {
		return fmt.Errorf("error opening process: %v", err)
	}
	defer windows.CloseHandle(handle)

	// Set the priority class
	if err := windows.SetPriorityClass(handle, priority); err != nil {
		return fmt.Errorf("error setting priority class: %v", err)
	}

	return nil
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

//suspend te thread of the process
func SuspendProcess(proc string) {
	pid, err := FindPidByNamePowerShell(proc)
    if err != nil{
        fmt.Println(err)
    }
	
	fmt.Printf("\033[31mSuspending process with PID: %d...\033[0m\n", pid)
	

	// Open the process with required access
	hProcess, err := windows.OpenProcess(windows.PROCESS_SUSPEND_RESUME, false, uint32(pid))
	if err != nil {
		fmt.Println("\033[31mError opening process for suspending:\033[0m", err)
		return
	}
	defer windows.CloseHandle(hProcess)

	// Suspend the process
	_, err = windows.SuspendThread(hProcess)
	if err != nil {
		fmt.Println("\033[31mError suspending process:\033[0m", err)
		return
	}

	fmt.Println("\033[32mProcess suspended successfully\033[0m")
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// Close the handle
func closeHandle(handle syscall.Handle) {
	procCloseHandle.Call(uintptr(handle))
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// Function to resume a suspended process
func ResumeProcess(proc string) {
	pid, err := FindPidByNamePowerShell(proc)
    if err != nil{
        fmt.Println(err)
    }
	fmt.Printf("\033[36mResuming process with PID: %d\033[0m\n", pid)

	// Open the file descriptor
	hProcess, err := windows.OpenProcess(windows.PROCESS_SUSPEND_RESUME, false, uint32(pid))
	if err != nil {
		fmt.Println("\033[31mError opening process for resuming:\033[0m", err)
		return
	}
	defer windows.CloseHandle(hProcess)

	// Resume the process
	_, err = windows.ResumeThread(hProcess)
	if err != nil {
		fmt.Println("\033[31mError resuming process:\033[0m", err)
		return
	}

	fmt.Println("\033[32mProcess resumed successfully\033[0m")

}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// Function to read memory from a specific process
func ReadMemory(proc string, address string, size int) {
	pid, err := FindPidByNamePowerShell(proc)
    if err != nil{
        fmt.Println(err)
    }
	fmt.Printf("\033[37mReading memory at address: %s for PID: %d\033[0m\n", address, pid)
	// Logic to read memory
}


// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// Function to write data to a specific memory address of a process
func WriteMemory(proc string, address string, data string) {
	pid, err := FindPidByNamePowerShell(proc)
    if err != nil{
        fmt.Println(err)
    }
	fmt.Printf("\033[32mWriting data to memory address: %s for PID: %d\033[0m\n", address, pid)
	// Logic to write memory
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func DisplayHelp() {
	fmt.Println("\033[36mThis is a tool for process analysis, is suggested to use the 'generic' args as first one... \033[0m")
	fmt.Println("Use help for display this massage...")
	fmt.Println("\033[36mUsage: ProcHandle <command> [arguments]\033[0m")
	fmt.Println("\033[33mCommands:\033[0m")
	fmt.Println("\033[32m  list\033[0m                   \033[37mList all running processes on the system.\033[0m")
	fmt.Println("\033[32m  info <pid>\033[0m             \033[37mRetrieve detailed information for a specific process by its PID.\033[0m")
	fmt.Println("\033[32m  kill <pid>\033[0m        \033[37mTerminate a process by its PID.\033[0m")
	fmt.Println("\033[32m  set-priority <pid> <priority>\033[0m \033[37mSet the priority for a process. Priority can be one of: low, normal, high, realtime.\033[0m")
	fmt.Println("\033[32m  suspend <pid>\033[0m          \033[37mSuspend a process by its PID.\033[0m")
	fmt.Println("\033[32m  resume <pid>\033[0m           \033[37mResume a suspended process by its PID.\033[0m")
	fmt.Println("\033[32m  read-memory <pid> <address> <size>\033[0m \033[37mRead memory at a specific address of a process.\033[0m")
	fmt.Println("\033[32m  write-memory <pid> <address> <data>\033[0m \033[37mWrite data to a specific memory address of a process.\033[0m")
}

func main() {
	for{
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
			
		//comming soon	
		// case "thread": 
		// 	GetThreadPriority()
			
		case "info":
			if len(os.Args) < 3 {
				fmt.Println("\033[31mUsage: info <pid>\033[0m")
				return
			}
			pid := os.Args[2]
	
			GetProcessInfo(pid)
			
			// -------------------------------------------------------------------------------------------------------------------------------------------------------------
			
		case "kill":
			if len(os.Args) < 3 {
				fmt.Println("\033[31mUsage: kill <pid>\033[0m")
				return
			}
			pid, err := strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Println("\033[31mError: Invalid PID value:\033[0m", os.Args[2])
				return
			}
			
			hpid, err := os.FindProcess(pid) 
	        if err != nil {
	            fmt.Println("\033[31mError: \t \033[0m", os.Args[2])
	            return
	        }
	
			TerminateProcess(hpid)
			
			// -------------------------------------------------------------------------------------------------------------------------------------------------------------
			
		case "set-priority":
			if len(os.Args) != 4 {
				fmt.Println("\033[31mUsage: set-priority <pid> <priority>\033[0m")
				return
			}
			
			pid, err := strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Println("\033[31mError: Invalid PID value \033[0m")
			}
			
			hpid, err := os.FindProcess(pid) 
	        if err != nil {
	            fmt.Println("\033[31mError: \t \033[0m", os.Args[2])
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
		
			priorities := map[string]uint32{
				"low": BELOW_NORMAL_PRIORITY_CLASS,
				"normal": NORMAL_PRIORITY_CLASS,
				"high": HIGH_PRIORITY_CLASS,
				"realtime": REALTIME_PRIORITY_CLASS,
			}
			priority, exists := priorities[priorityStr]
			if !exists {
				fmt.Println("\033[31mError: Invalid priority value\033[0m")
				return
			}
		
			if err := SetProcessPriority(hpid, priority); err != nil {
				fmt.Println("\033[31mError setting process priority:\033[0m", err)
			}
		
		
		// -----------------------------------------------------------------------------------------------------------------------------------------------------------------
		
		case "suspend":
			if len(os.Args) < 3 {
				fmt.Println("\033[31mUsage: suspend <pid>\033[0m")
				return
			}
			pid, err := strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Println("\033[31mError: Invalid PID value:\033[0m", os.Args[2])
				return
			}
			
			hpid, err := os.FindProcess(pid) 
	        if err != nil {
	            fmt.Println("\033[31mError: \t \033[0m", os.Args[2])
	            return
	        }
	
	
			SuspendProcess(hpid)
			
			// -------------------------------------------------------------------------------------------------------------------------------------------------------------
			
		case "resume":
			if len(os.Args) < 3 {
				fmt.Println("\033[31mUsage: resume <pid>\033[0m")
				return
			}
			pid, err := strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Println("\033[31mError: Invalid PID value:\033[0m", os.Args[2])
				return
			}
			
			hpid, err := os.FindProcess(pid) 
	        if err != nil {
	            fmt.Println("\033[31mError: \t \033[0m", os.Args[2])
	            return
	        }
	
	
			ResumeProcess(hpid)
			
			// -------------------------------------------------------------------------------------------------------------------------------------------------------------
			
		case "read-memory":
			if len(os.Args) < 5 {
				fmt.Println("\033[31mUsage: read-memory <pid> <address> <size>\033[0m")
				return
			}
			pid, err := strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Println("\033[31mError: Invalid PID value:\033[0m", os.Args[2])
				return
			}
			hpid, err := os.FindProcess(pid) 
	        if err != nil {
	            fmt.Println("\033[31mError: \t \033[0m", os.Args[2])
	            return
	        }
	
	
			address := os.Args[3]
			size, err := strconv.Atoi(os.Args[4])
			if err != nil {
				fmt.Println("\033[31mError: Invalid size value:\033[0m", os.Args[4])
				return
			}
			ReadMemory(hpid, address, size)
			
			// -------------------------------------------------------------------------------------------------------------------------------------------------------------
			
		case "write-memory":
			if len(os.Args) < 5 {
				fmt.Println("\033[31mUsage: write-memory <pid> <address> <data>\033[0m")
				return
			}
			pid, err := strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Println("\033[31mError: Invalid PID value:\033[0m", os.Args[2])
				return
			}
			
			hpid, err := os.FindProcess(pid) 
	        if err != nil {
	            fmt.Println("\033[31mError: \t \033[0m", os.Args[2])
	            return
	        }
	
			address := os.Args[3]
			data := os.Args[4]
			WriteMemory(hpid, address, data)
			
			// -------------------------------------------------------------------------------------------------------------------------------------------------------------
			
		default:
			fmt.Println("\033[31mError: Unknown command:\033[0m", command)
			DisplayHelp()
		}
	}
	
	// comming soon
	// func GetThreadPriority(hThread syscall.Handle) {
}
	
