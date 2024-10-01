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


const (
	DELETE          = "0x00010000L" 	//Required to delete the object.
	READ_CONTROL    = "0x00020000L" 	//Required to read information in the security descriptor for the object, not including the information in the SACL. To read or write the SACL, you must request the ACCESS_SYSTEM_SECURITY access right. For more information, see SACL Access Right.
	SYNCHRONIZE     = "0x00100000L" 	//The right to use the object for synchronization. This enables a thread to wait until the object is in the signaled state.
	WRITE_DAC       = "0x00040000L" 	//Required to modify the DACL in the security descriptor for the object.
	WRITE_OWNER     ="0x00080000L"      // Required to change the owner in the security descriptor for the object.

)

const (
	ACCESS = windows.PROCESS_SET_INFORMATION
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
    BELOW_NORMAL_PRIORITY_CLASS   = 0x00004000
    NORMAL_PRIORITY_CLASS         = 0x00000020
    ABOVE_NORMAL_PRIORITY_CLASS   = 0x00008000
    HIGH_PRIORITY_CLASS           = 0x00000080
    REALTIME_PRIORITY_CLASS       = 0x00000100
)

var ( 
	entry syscall.ProcessEntry32
)
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
			break // No more processes to enumerate
		}
	}

	entry = syscall.ProcessEntry32{}
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// Function to get detailed information about a specific process
func GetProcessInfo(name string) {
    // Trova i PIDs per il nome del processo fornito
    pids, err := FindPidByNamePowerShell(name)
    if err != nil {
        fmt.Println("\033[31mError finding process:\033[0m", err)
        return
    }

    if len(pids) > 0 {
        for _, pid := range pids {
            pidValue := uint32(pid) // Converti il PID in uint32
            fmt.Printf("\033[36mRetrieving information for PID: %d...\033[0m\n", pidValue)

            // Apri il processo con i diritti di accesso necessari
            hProcess, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ, false, pidValue)
            if err != nil {
                fmt.Println("\033[31mError opening process:\033[0m", err)
                continue // Salta al PID successivo se c'è un errore
            }
            defer windows.CloseHandle(hProcess) // Assicurati che l'handle venga chiuso quando fatto

            // Recupera il nome completo del processo
            var processName [windows.MAX_PATH]uint16
            processPathLength := uint32(len(processName))
            err = windows.QueryFullProcessImageName(hProcess, 0, &processName[0], &processPathLength)
            if err != nil {
                fmt.Println("\033[31mError retrieving process name:\033[0m", err)
                continue
            }

            // Converti il nome del processo da UTF16 a una stringa Go
            processNameStr := windows.UTF16ToString(processName[:])
            fmt.Printf("\033[32mPID: %d\tName: %s\033[0m\n", pidValue, processNameStr)

            // Ottieni informazioni sulla memoria per il processo
            var memInfo windows.MemoryBasicInformation
            addr := uintptr(0) // Inizia all'indirizzo 0

            for {
                // Query memory information for the specified process
                ret := windows.VirtualQueryEx(hProcess, addr, &memInfo, uintptr(unsafe.Sizeof(memInfo)))
                if ret != nil {
                    fmt.Println("\033[31mFinished querying memory regions.\033[0m")
                    break // Esci dal ciclo se la query fallisce
                }

                if memInfo.State == windows.MEM_COMMIT {
                    // Stampa informazioni sulla regione di memoria se è impegnata
                    fmt.Printf("\033[34mMemory Region: Base Address: %x, Region Size: %d bytes\033[0m\n", memInfo.BaseAddress, memInfo.RegionSize)
                }
                addr += memInfo.RegionSize // Passa alla regione successiva
            }

            // Ottieni informazioni sull'utilizzo della CPU
            var creationTime, exitTime, kernelTime, userTime windows.Filetime
            err = windows.GetProcessTimes(hProcess, &creationTime, &exitTime, &kernelTime, &userTime)
            if err != nil {
                fmt.Println("\033[31mError retrieving process times:\033[0m", err)
                continue
            }

            // Calcola il tempo totale della CPU utilizzato dal processo
            cpuTime := kernelTime.Nanoseconds() + userTime.Nanoseconds()
            fmt.Printf("\033[34mCPU Time: %d nanoseconds\033[0m\n", cpuTime)
        }
    } else {
        fmt.Println("\033[33mNo processes found with the given name.\033[0m")
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
func TerminateProcess(proc string) {
	pids, err := FindPidByNamePowerShell(proc)
	if err != nil {
		fmt.Println("\033[31mError finding process:\033[0m", err)
	}
	
	for _, pid := range pids {
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
}
// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// Function to set the priority of a process
func SetProcessPriority(proc string, priority uint32) error {
	fmt.Printf("\033[33mSetting priority for PROCESS: %s to %d\033[0m\n", proc, priority)
	
	pids, err := FindPidByNamePowerShell(proc)
    if err != nil{
        fmt.Println(err)
    }
    
    for _, hpid := range pids {
	// Open the process with required access
		pid := uint32(hpid)		
		handle, err := windows.OpenProcess(ACCESS, false, pid)
		if err != nil {
			return fmt.Errorf("error opening process: %v", err)
		}
		defer windows.CloseHandle(handle)
	
		// Set the priority class
		err = windows.SetPriorityClass(handle, priority)
		if err != nil {
			fmt.Printf("error setting priority class: %v", err)
			fmt.Println(windows.GetLastError())
		}
	
		return nil
	}
	return nil
	
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

//suspend te thread of the process
func SuspendProcess(proc string) {
	pids, err := FindPidByNamePowerShell(proc)
    if err != nil{
        fmt.Println(err)
    }
    
    for _, pid := range pids {
	
		fmt.Printf("\033[31mSuspending process with PID: %d...\033[0m\n", pid)
		
	
		// Open the process with required access
		hProcess, err := windows.OpenProcess(windows.PROCESS_SUSPEND_RESUME, false, uint32(pid))
		if err != nil {
			fmt.Println("\033[31mError opening process for suspending:\033[0m", err)
			return
		}
		defer windows.CloseHandle(hProcess)
	
		// Suspend the process
		// _, err = windows.SuspendThread(hProcess)
		// if err != nil {
		// 	fmt.Println("\033[31mError suspending process:\033[0m", err)
		// 	return
		// }
	
		fmt.Println("\033[32mProcess suspended successfully\033[0m")
	}
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// Close the handle
func closeHandle(handle syscall.Handle) {
	procCloseHandle.Call(uintptr(handle))
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// Function to resume a suspended process
func ResumeProcess(proc string) {
	pids, err := FindPidByNamePowerShell(proc)
    if err != nil{
        fmt.Println(err)
    }
    
    for _, pid := range pids {
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
	fmt.Println("\033[32m  info <proc_name>\033[0m             \033[37mRetrieve detailed information for a specific process by its PID.\033[0m")
	fmt.Println("\033[32m  kill <proc_name>\033[0m        \033[37mTerminate a process by its PID.\033[0m")
	fmt.Println("\033[32m  set-priority <process_name><priority>\033[0m \033[37mSet the priority for a process. Priority can be one of: low, normal, high, realtime.\033[0m")
	fmt.Println("\033[32m  suspend <proc_name>\033[0m          \033[37mSuspend a process by its PID.\033[0m")
	fmt.Println("\033[32m  resume <proc_name>\033[0m           \033[37mResume a suspended process by its PID.\033[0m")
	fmt.Println("\033[32m  read-memory <process_name><address> <size>\033[0m \033[37mRead memory at a specific address of a process.\033[0m")
	fmt.Println("\033[32m  write-memory <process_name><address> <data>\033[0m \033[37mWrite data to a specific memory address of a process.\033[0m")
}

func main() {
	if len(os.Args) < 2 {
		DisplayHelp()
		return
	}

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
			fmt.Println("\033[31mUsage: info <process_name>\033[0m")
			return
		}
		processName := os.Args[2]

		GetProcessInfo(processName)
		
	// -------------------------------------------------------------------------------------------------------------------------------------------------------------
		
	case "generic":
		generic()
		
	case "kill":
		if len(os.Args) < 3 {
			fmt.Println("\033[31mUsage: kill <process_name>\033[0m")
			return
		}
		processName := os.Args[2]

		TerminateProcess(processName)
		
	// -------------------------------------------------------------------------------------------------------------------------------------------------------------
		
case "set-priority":
    if len(os.Args) != 4 && len(os.Args) != 2{
        fmt.Println("\033[31mUsage: set-priority <process_name> <priority>\033[0m")
    }
    
    processName := os.Args[2]

    // If the user requests information about priority classes
    if processName == "info" {
        fmt.Println("Process Priority Classes and Corresponding Thread Priority Levels:")
        fmt.Println("----------------------------------------------------------------------")
        fmt.Println("IDLE_PRIORITY_CLASS          THREAD_PRIORITY_IDLE              1")
        fmt.Println("                             THREAD_PRIORITY_LOWEST            2")
        fmt.Println("                             THREAD_PRIORITY_BELOW_NORMAL      3")
        fmt.Println("                             THREAD_PRIORITY_NORMAL            4")
        fmt.Println("                             THREAD_PRIORITY_ABOVE_NORMAL      5")
        fmt.Println("                             THREAD_PRIORITY_HIGHEST           6")
        fmt.Println("                             THREAD_PRIORITY_TIME_CRITICAL     15")
        fmt.Println("BELOW_NORMAL_PRIORITY_CLASS  THREAD_PRIORITY_IDLE              1")
        fmt.Println("                             THREAD_PRIORITY_LOWEST            4")
        fmt.Println("                             THREAD_PRIORITY_BELOW_NORMAL      5")
        fmt.Println("                             THREAD_PRIORITY_NORMAL            6")
        fmt.Println("                             THREAD_PRIORITY_ABOVE_NORMAL      7")
        fmt.Println("                             THREAD_PRIORITY_HIGHEST           8")
        fmt.Println("                             THREAD_PRIORITY_TIME_CRITICAL     15")
        fmt.Println("NORMAL_PRIORITY_CLASS        THREAD_PRIORITY_IDLE              1")
        fmt.Println("                             THREAD_PRIORITY_LOWEST            6")
        fmt.Println("                             THREAD_PRIORITY_BELOW_NORMAL      7")
        fmt.Println("                             THREAD_PRIORITY_NORMAL            8")
        fmt.Println("                             THREAD_PRIORITY_ABOVE_NORMAL      9")
        fmt.Println("                             THREAD_PRIORITY_HIGHEST           10")
        fmt.Println("                             THREAD_PRIORITY_TIME_CRITICAL     15")
        fmt.Println("ABOVE_NORMAL_PRIORITY_CLASS  THREAD_PRIORITY_IDLE              1")
        fmt.Println("                             THREAD_PRIORITY_LOWEST            8")
        fmt.Println("                             THREAD_PRIORITY_BELOW_NORMAL      9")
        fmt.Println("                             THREAD_PRIORITY_NORMAL            10")
        fmt.Println("                             THREAD_PRIORITY_ABOVE_NORMAL      11")
        fmt.Println("                             THREAD_PRIORITY_HIGHEST           12")
        fmt.Println("                             THREAD_PRIORITY_TIME_CRITICAL     15")
        fmt.Println("HIGH_PRIORITY_CLASS          THREAD_PRIORITY_IDLE              1")
        fmt.Println("                             THREAD_PRIORITY_LOWEST            11")
        fmt.Println("                             THREAD_PRIORITY_BELOW_NORMAL      12")
        fmt.Println("                             THREAD_PRIORITY_NORMAL            13")
        fmt.Println("                             THREAD_PRIORITY_ABOVE_NORMAL      14")
        fmt.Println("                             THREAD_PRIORITY_HIGHEST           15")
        fmt.Println("                             THREAD_PRIORITY_TIME_CRITICAL     15")
        fmt.Println("REALTIME_PRIORITY_CLASS      THREAD_PRIORITY_IDLE              16")
        fmt.Println("                             THREAD_PRIORITY_LOWEST            22")
        fmt.Println("                             THREAD_PRIORITY_BELOW_NORMAL      23")
        fmt.Println("                             THREAD_PRIORITY_NORMAL            24")
        fmt.Println("                             THREAD_PRIORITY_ABOVE_NORMAL      25")
        fmt.Println("                             THREAD_PRIORITY_HIGHEST           26")
        fmt.Println("                             THREAD_PRIORITY_TIME_CRITICAL     31")
		fmt.Println("Priority Mappings:")
        fmt.Println("----------------------------------------------------------------------")
        fmt.Println("\"idle\":          windows.IDLE_PRIORITY_CLASS")
        fmt.Println("\"below_normal\":  windows.BELOW_NORMAL_PRIORITY_CLASS")
        fmt.Println("\"low\":           windows.BELOW_NORMAL_PRIORITY_CLASS")
        fmt.Println("\"normal\":        windows.NORMAL_PRIORITY_CLASS")
        fmt.Println("\"above_normal\":  windows.ABOVE_NORMAL_PRIORITY_CLASS")
        fmt.Println("\"high\":          windows.HIGH_PRIORITY_CLASS")
        fmt.Println("\"realtime\":      windows.REALTIME_PRIORITY_CLASS")
        return
    }

    // Get the priority argument from the command line
    priorityStr := os.Args[3]
    
    // Define a map to convert priority name to Windows constant
    priorityMap := map[string]uint32{
        "idle":          windows.IDLE_PRIORITY_CLASS,
        "below_normal":  windows.BELOW_NORMAL_PRIORITY_CLASS,
        "low":           windows.BELOW_NORMAL_PRIORITY_CLASS, // Added this line
        "normal":        windows.NORMAL_PRIORITY_CLASS,
        "above_normal":  windows.ABOVE_NORMAL_PRIORITY_CLASS,
        "high":          windows.HIGH_PRIORITY_CLASS,
        "realtime":      windows.REALTIME_PRIORITY_CLASS,
    }

    // Find the corresponding priority class
    priority, exists := priorityMap[priorityStr]
    if !exists {
        fmt.Println("\033[31mError: Invalid priority value\033[0m")
        return
    }

    // Set the process priority for the specified process using the uint32 priority value
    if err := SetProcessPriority(processName, priority); err != nil {
        fmt.Println("\033[31mError setting process priority:\033[0m", err)
    }

		
	// -----------------------------------------------------------------------------------------------------------------------------------------------------------------
	
	case "suspend":
		if len(os.Args) < 3 {
			fmt.Println("\033[31mUsage: suspend <process_name>\033[0m")
			return
		}
		processName := os.Args[2]

		SuspendProcess(processName)
		
	// -------------------------------------------------------------------------------------------------------------------------------------------------------------
		
	case "resume":
		if len(os.Args) < 3 {
			fmt.Println("\033[31mUsage: resume <process_name>\033[0m")
			return
		}
		processName := os.Args[2]

		ResumeProcess(processName)
		
	// -------------------------------------------------------------------------------------------------------------------------------------------------------------
		
	case "read-memory":
		if len(os.Args) < 5 {
			fmt.Println("\033[31mUsage: read-memory <process_name> <address> <size>\033[0m")
			return
		}
		processName := os.Args[2]
		address := os.Args[3]
		size, err := strconv.Atoi(os.Args[4])
		if err != nil {
			fmt.Println("\033[31mError: Invalid size value:\033[0m", os.Args[4])
			return
		}
		ReadMemory(processName, address, size)
		
	// -------------------------------------------------------------------------------------------------------------------------------------------------------------
		
	case "write-memory":
		if len(os.Args) < 5 {
			fmt.Println("\033[31mUsage: write-memory <process_name> <address> <data>\033[0m")
			return
		}
		processName := os.Args[2]
		address := os.Args[3]
		data := os.Args[4]
		WriteMemory(processName, address, data)
		
	// -------------------------------------------------------------------------------------------------------------------------------------------------------------
		
	default:
		fmt.Println("\033[31mError: Unknown command:\033[0m", command)
		DisplayHelp()
	}	
}

	// comming soon
// func GetThreadPriority(hThread syscall.Handle) {
// and
// netowrk visualizer for processes and anti debugger
