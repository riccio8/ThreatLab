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


// ALL THOSE FUNCTIONS WILL BE APLIED ON THE PARENT PROCESS AND ON THE CHILD TOO



const (
	DELETE          = "0x00010000L" 	//Required to delete the object.
	READ_CONTROL    = "0x00020000L" 	//Required to read information in the security descriptor for the object, not including the information in the SACL. To read or write the SACL, you must request the ACCESS_SYSTEM_SECURITY access right. For more information, see SACL Access Right.
	SYNCHRONIZE     = "0x00100000L" 	//The right to use the object for synchronization. This enables a thread to wait until the object is in the signaled state.
	WRITE_DAC       = "0x00040000L" 	//Required to modify the DACL in the security descriptor for the object.
	WRITE_OWNER     ="0x00080000L"      // Required to change the owner in the security descriptor for the object.

)


// Access rights for thread.
const(
	THREAD_DIRECT_IMPERSONATION      = 0x0200
	THREAD_GET_CONTEXT               = 0x0008
	THREAD_IMPERSONATE               = 0x0100
	THREAD_QUERY_INFORMATION         = 0x0040
	THREAD_QUERY_LIMITED_INFORMATION = 0x0800
	THREAD_SET_CONTEXT               = 0x0010
	THREAD_SET_INFORMATION           = 0x0020
	THREAD_SET_LIMITED_INFORMATION   = 0x0400
	THREAD_SET_THREAD_TOKEN          = 0x0080
	THREAD_SUSPEND_RESUME            = 0x0002
	THREAD_TERMINATE                 = 0x0001
)
                


// Access rights for process.
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

var (
	kernel32               = syscall.NewLazyDLL("kernel32.dll")
	procEnumProcessThreads = kernel32.NewProc("EnumProcessThreads")
	procSuspendThread      = kernel32.NewProc("SuspendThread")
	procResumeThread       = kernel32.NewProc("ResumeThread")
	procCloseHandle        = kernel32.NewProc("CloseHandle")
	procVirtualProtectEx   = kernel32.NewProc("VirtualProtectEx")
)

var (
    iphlpapi               = syscall.NewLazyDLL("Iphlpapi.dll")
    procGetExtendedTcpTable = iphlpapi.NewProc("GetExtendedTcpTable")
)

const INVALID_HANDLE_VALUE = ^uintptr(0)

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
func ListInfoProcesses() ([]string, error) {
	var output []string
	output = append(output, "isting all processes...")

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

	output = append(output, "Processes:")
	for {
		processName := syscall.UTF16ToString(entry.ExeFile[:])
		output = append(output, fmt.Sprintf("Pid: %d\tFile Name: %s\tThread: %d\tHeap Allocation: %d\tProcess Flags: %d",
			entry.ProcessID, processName, entry.Threads, entry.DefaultHeapID, entry.Flags))

		err = syscall.Process32Next(snapshot, &entry)
		if err != nil {
			output = append(output, "No more processes...")
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
			output += fmt.Sprintf("Retrieving information for PID: %d...\n", pidValue)

			hProcess, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ, false, pidValue)
			if err != nil {
				output += fmt.Sprintf("Error opening process: %v\n", err)
				continue
			}
			defer windows.CloseHandle(hProcess)

			var processName [windows.MAX_PATH]uint16
			processPathLength := uint32(len(processName))
			err = windows.QueryFullProcessImageName(hProcess, 0, &processName[0], &processPathLength)
			if err != nil {
				output += fmt.Sprintf("Error retrieving process name: %v\n", err)
				continue
			}

			processNameStr := windows.UTF16ToString(processName[:])
			output += fmt.Sprintf("PID: %d\tName: %s\n", pidValue, processNameStr)

			var memInfo windows.MemoryBasicInformation
			addr := uintptr(0)

			for {
				ret := windows.VirtualQueryEx(hProcess, addr, &memInfo, uintptr(unsafe.Sizeof(memInfo)))
				if ret != nil { 
					output += "Finished querying memory regions.\n"
					break
				}

				if memInfo.State == windows.MEM_COMMIT {
					output += fmt.Sprintf("Memory Region: Base Address: %x, Region Size: %d bytes\n", memInfo.BaseAddress, memInfo.RegionSize)
				}
				addr += memInfo.RegionSize
			}

			var creationTime, exitTime, kernelTime, userTime windows.Filetime
			err = windows.GetProcessTimes(hProcess, &creationTime, &exitTime, &kernelTime, &userTime)
			if err != nil {
				output += fmt.Sprintf("Error retrieving process times: %v\n", err)
				continue
			}

			cpuTime := kernelTime.Nanoseconds() + userTime.Nanoseconds()
			output += fmt.Sprintf("CPU Time: %d nanoseconds\n", cpuTime)
		}
	} else {
		output += "No processes found with the given name.\n"
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

    // Find process IDs by name
    pids, err := FindPidByNamePowerShell(name)
    if err != nil {
        return "", fmt.Errorf("error retrieving information for PID: %w", err)
    }

    // Iterate through each PID to attempt memory protection changes
    for _, hpid := range pids {
        pid := uint32(hpid)
        handle, err := windows.OpenProcess(ACCESS, false, pid)
        if err != nil {
            output += fmt.Sprintf("Error opening process (PID: %d): %v\n", pid, err) // Include PID in the output
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
            output += fmt.Sprintf("Error while calling VirtualProtectEx on PID %d: %v\n", pid, err)
        } else {
            output += fmt.Sprintf("Successfully changed protection for PID %d\n", pid) // Success message
        }
    }

    // If there were no changes made, we return a message indicating that
    if output == "" {
        return "No changes made to memory protection.", nil
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
		output += fmt.Sprintf("Terminating process with PID: %d...\n", pid)

		cmd := exec.Command("taskkill", "/PID", fmt.Sprint(pid), "/F")
		processOutput, err := cmd.CombinedOutput()
		if err != nil {
			output += fmt.Sprintf("Error terminating process: %s\n", err)
			return output, err
		}

		output += fmt.Sprintf("Process terminated successfully: %s\n", string(processOutput))
	}
	
	if output == "" {
        return "No kill function did not run ruccessfully", nil
    }

	return output, nil
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// Function to set the priority of a process
func SetProcessPriority(proc string, priority uint32) (string, error) {
	var output string
    pids, err := FindPidByNamePowerShell(proc)
    if err != nil {
        return "error finding process: ", err
    }
    
    for _, hpid := range pids {
		output += fmt.Sprintf("Setting process priority to Pprocess: %s...\n", proc)
        pid := uint32(hpid)		
        handle, err := windows.OpenProcess(PROCESS_ALL_ACCESS, false, pid)
        if err != nil {
            output += fmt.Sprintf("Error terminating process: %s\n", err)
            return output, err
        }
        defer windows.CloseHandle(handle)
    
        err = windows.SetPriorityClass(handle, priority)
        if err != nil {
            return "error setting priority class: ", err
        } 
        output += fmt.Sprintf("Priority setted successfully for %s\n", proc)
    }
	return output, nil
}


// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

//suspend te thread of the process
func SuspendProcess(proc string) (string, error) {
    pids, err := FindPidByNamePowerShell(proc) 
    if err != nil {
        return "error finding process: ", err
    }

    for _, pid := range pids {
        hpid := uint32(pid)
        handle, err := windows.OpenThread(THREAD_SUSPEND_RESUME, false, hpid)
        if err != nil {
            continue
            // Try the next PID if one fails
        }

        retVal, _, err := procSuspendThread.Call(uintptr(handle))
        if err != nil {
            windows.CloseHandle(handle) // Close handle upon error
            continue // Try the next PID
        }

        if retVal == 0xFFFFFFFF { // Check for failure
            windows.CloseHandle(handle) // Ensure we clean up
            continue // Try the next PID
        }

        windows.CloseHandle(handle) // Close handle after suspension
        return string(retVal), nil
    }

    return "", nil // None suspended
}


// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// Close the handle
func closeHandle(handle syscall.Handle) {
	procCloseHandle.Call(uintptr(handle))
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// Function to resume a suspended process
func ResumeProcess(proc string) (string, error){
    pids, err := FindPidByNamePowerShell(proc) 
    if err != nil {
        return "error finding process: ", err
    }

    for _, pid := range pids {
        hpid := uint32(pid)
        handle, err := windows.OpenThread(THREAD_SUSPEND_RESUME, false, hpid)
        if err != nil {
            continue
            // Try the next PID if one fails
        }

        retVal, _, err := procSuspendThread.Call(uintptr(handle))
        if err != nil {
            windows.CloseHandle(handle) // Close handle upon error
            continue // Try the next PID
        }

        if retVal == 0xFFFFFFFF { // Check for failure
            windows.CloseHandle(handle) // Ensure we clean up
            continue // Try the next PID
        }

        windows.CloseHandle(handle) // Close handle after suspension
        return string(retVal), nil
    }

    return "", nil // None suspended
}



// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// Function to read memory from a specific process
func ReadMemory(proc string, address uintptr, size int) ([]string, error) {
    var output []string
    output = append(output, "Listing all processes...")

    pids, err := FindPidByNamePowerShell(proc)
    if err != nil {
        return output, fmt.Errorf("failed to find PIDs: %w", err)
    }

    for _, pid := range pids {
        hProcess, err := windows.OpenProcess(PROCESS_ALL_ACCESS, false, uint32(pid))
        if err != nil {
            output = append(output, fmt.Sprintf("failed to open process with PID %d: %v", pid, err))
            continue
        }

        // Close the process handle properly 
        defer windows.CloseHandle(hProcess)

        dataBytes := make([]byte, size)
        
        err = windows.ReadProcessMemory(hProcess, address, &dataBytes[0], uintptr(len(dataBytes)), nil)
        if err != nil {
            output = append(output, fmt.Sprintf("failed to read memory from PID %d: %v", pid, err))
            continue
        }

        // Add the read data to output in a readable format (for example, as a hex string)
        output = append(output, fmt.Sprintf("Read memory from PID %d: %x", pid, dataBytes))
    }

    return output, nil
}




// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// Function to write data to a specific memory address of a process
func WriteMemory(proc string, address int, data string) (string, error){
    pids, err := FindPidByNamePowerShell(proc)
    if err != nil {
        return "", err
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
            return "", err
        }
        return string(dataBytes), nil
    }
    return "", nil
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

    var connections []string
    for _, line := range lines {
        if line != "" {
            connections = append(connections, line) // Collect non-empty lines
        }
    }

    if len(connections) == 0 {
        return "No active connections found.", nil
    }

    return strings.Join(connections, "<br />"), nil // Join lines with HTML line breaks for formatting
}







// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func DisplayHelp() {
	fmt.Println("\033[36mThis is a tool for process analysis, is suggested to use the 'generic' args as first one... ")
	fmt.Println("Use help for display this massage...")
	fmt.Println("\033[36mUsage: ProcHandle <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  list                   \033[37mList all running processes on the system.")
	fmt.Println("  info <proc_name>             \033[37mRetrieve detailed information for a specific process by its PID.")
	fmt.Println("  kill <proc_name>        \033[37mTerminate a process by its PID.")
	fmt.Println("  set-priority <process_name> <priority> \033[37mSet the priority for a process. Priority can be one of: low, normal, high, realtime.")
	fmt.Println("  suspend <proc_name>          \033[37mSuspend a process by its PID.")
	fmt.Println("  cpnnection          \033[37mRetrievs a list of all current connection.")
	fmt.Println("  resume <proc_name>           \033[37mResume a suspended process by its PID.")
	fmt.Println("  read-memory <process_name> <address> <size> \033[37mRead memory at a specific address of a process.")
	fmt.Println("  write-memory <process_name><address> <data> \033[37mWrite data to a specific memory address of a process.")
	fmt.Println("  protect <process_name> <lpAddress> <dwSize> <flNewProtect> \033[37mChange the type of permits of a specific memory region which belongs to the process given by name.")
	
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
				lastOutput = "<p style='color:red;'>Error: " + err.Error() + "</p>" // Display the error message
			} else {
				lastOutput = "<p><strong></strong> Retrieved active connections:<br />" + result + "</p>"
			}
		
            
		case "generic":
			output, err := generic()
			if err != nil {
				fmt.Println("Error:", err)
				lastOutput = "<p style='color:red;'>Error: " + err.Error() + "</p>" // Display the error message
				renderForm(w)
			}
			lastOutput = "<p><strong></strong> Retrieved active PROCESSES:<br />" + output + "</p>"

			
        
		case "list":
			conn, err := ListInfoProcesses()
			if err != nil {
				fmt.Println("Error:", err)
				renderForm(w)
				return
			}
		
		
			lastOutput = "<p><strong></strong> Listed all processes:</p><ul>"
			for _, proc := range conn {
				lastOutput += "<li>" + proc + "</li>"
			}
			lastOutput += "</ul>"
		
		
		case "info":
			if len(args) < 1 {
				lastOutput = "<p style='color:red;'>Error: Missing process name argument</p>"
				renderForm(w)
				return
			}
			result, err := GetProcessInfo(args[0])
			if err != nil {
				lastOutput = fmt.Sprintf("<p style='color:red;'>Error: %v</p>", err)
			} else {
				lastOutput = fmt.Sprintf("<p><strong></strong><br />%s</p>", result) 
			}
		
		

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
		
			// Call the VirtualProtection function and handle the output
			output, err := VirtualProtection(processName, uintptr(lpAddress), uintptr(dwSize), flNewProtect)
			if err != nil {
				lastOutput = fmt.Sprintf("<p style='color:red;'>Error: %v</p>", err)
			} else {
				lastOutput = "<p><strong>Output:</strong> " + output + "</p>" // Display the output from the VirtualProtection function
			}
		
			

        case "kill":
            if len(args) < 1 {
                lastOutput = "<p style='color:red;'>Error: Missing process name</p>"
                renderForm(w)
                return
            }
            TerminateProcess(args[0])
            lastOutput = fmt.Sprintf("<p><strong></strong> Process %s terminated</p>", args[0])

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

            if _, err := SetProcessPriority(processName, priority); err != nil {
                lastOutput = fmt.Sprintf("<p style='color:red;'>Error setting process priority: %v</p>", err)
                renderForm(w)
                return
            }
            lastOutput = fmt.Sprintf("<p><strong></strong> Process %s priority set to %s</p>", processName, priorityStr)

        case "suspend":
            if len(args) < 1 {
                lastOutput = "<p style='color:red;'>Error: Missing process name</p>"
                renderForm(w)
                return
            }
            _, err := SuspendProcess(args[0])
            if err != nil {
                lastOutput = fmt.Sprintf("<p style='color:red;'>Error: %v</p>", err)
            }
            lastOutput = fmt.Sprintf("<p><strong></strong> Process %s suspended </p>", args[0])

        case "resume":
            if len(args) < 1 {
                lastOutput = "<p style='color:red;'>Error: Missing process name</p>"
                renderForm(w)
                return
            }
            ResumeProcess(args[0])
            lastOutput = fmt.Sprintf("<p><strong></strong> Process %s resumed</p>", args[0])

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
            data, err := ReadMemory(processName, uintptr(address), size)
            if err != nil {
                lastOutput = fmt.Sprintf("<p style='color:red;'>Error Calling function: %v</p>", err)
                renderForm(w)
                return
            }
            for _, single_data := range data{
                lastOutput = fmt.Sprintf("<p><strong></strong> Read memory from process %s, data: %s</p>", processName, single_data)
            }
            

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
            lastOutput = fmt.Sprintf("<p><strong></strong> Wrote memory to process %s</p>", processName)

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
	            backdrop-filter: blur(5px); /* Background blur effect */
	        }
	        .container {
	            background-color: rgba(255, 255, 255, 0.99); /* More opaque background */
	            padding: 40px; /* Increased padding */
	            border-radius: 12px;
	            box-shadow: 0 4px 30px rgba(0, 0, 0, 0.4); /* Increased shadow effect */
	            width: 500px; /* Increased width */
	            text-align: center;
	            position: relative;
	        }
	        h2 {
	            font-size: 24px; /* Larger font size for the header */
	            margin-bottom: 20px; /* Space below the header */
	        }
	        input[type="text"] {
	            width: 100%;
	            padding: 15px; /* Increased padding */
	            margin-bottom: 15px;
	            border: 1px solid #ccc;
	            border-radius: 4px;
	            font-size: 18px; /* Larger font size for input */
	        }
	        button {
	            padding: 15px 25px; /* Increased padding */
	            background-color: #007BFF;
	            color: #fff;
	            border: none;
	            border-radius: 4px;
	            cursor: pointer;
	            font-size: 18px; /* Larger font size for button */
	        }
	        button:hover {
	            background-color: #0056b3;
	        }
	        .output {
	            margin-top: 20px;
	            text-align: left;
	            background-color: #f9f9f9;
	            padding: 20px; /* Increased padding */
	            border-radius: 4px;
	            border: 1px solid #ccc;
	            max-height: 400px; /* Increased height for output */
	            overflow-y: auto;
	            font-size: 16px; /* Larger font size for output */
	        }
	        .footer {
	            position: absolute;
	            bottom: 10px; /* Positioning at the bottom */
	            right: 10px; /* Positioning at the right */
	            font-size: 14px; /* Font size for footer link */
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
	        <div class="footer">
	            <a href="https://github.com/riccio8/Offensive-defensive-tolls/blob/main/src/processes/go/ui/ProcHandleUI.go" target="_blank">GitHub</a>
	        </div>
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
