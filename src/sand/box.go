package main

import (
	"fmt"
	"syscall"
	"unsafe"
)


// IMPORT REOURCES: https://learn.microsoft.com/en-us/windows/win32/api/winnt/ns-winnt-jobobject_limit_violation_information
//                  https://learn.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-createjobobjecta


var (
	kernel32                  = syscall.NewLazyDLL("Kernel32.dll")
	CreateJobObjectA          = kernel32.NewProc("CreateJobObjectA")
	SetInformationJobObject   = kernel32.NewProc("SetInformationJobObject")
	CreateProcess             = kernel32.NewProc("CreateProcessW") // Use CreateProcessW for Unicode support
	CloseHandle               = kernel32.NewProc("CloseHandle")
	AssignProcessToJobObject  = kernel32.NewProc("AssignProcessToJobObject")
	QueryInformationJobObject = kernel32.NewProc("QueryInformationJobObject")
	GetProcessId              = kernel32.NewProc("GetProcessId")
)

type JobObjectBasicLimitInformation struct {
	PerProcessUserTimeLimit syscall.Filetime
	PerJobUserTimeLimit     syscall.Filetime
	LimitFlags              uint32
	MinimumWorkingSetSize   uintptr
	MaximumWorkingSetSize   uintptr
	ActiveProcessCount      uint32
	TotalProcessCount       uint32
	InterfaceProcessCount   uint32
	Reserved1               uint32
	Reserved2               [3]uint32
}

type JobObjectExtendedLimitInformation struct {
	BasicLimitInformation JobObjectBasicLimitInformation
	ProcessMemoryLimit    uintptr
	JobMemoryLimit        uintptr
}

const (
	JobObjectBasicProcessIdList = 5 // Information class for basic process ID list
)

const (
	// Job object flags
	JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE          = 0x00002000 // Kill processes when the job is closed
	JOB_OBJECT_LIMIT_SUSPEND_RESUME             = 0x00000004 // Allow processes to be suspended and resumed
	JOB_OBJECT_LIMIT_SET_INFORMATION            = 0x00000010 // Allow setting information for the job
	JOB_OBJECT_LIMIT_JOB_MEMORY                 = 0x00000200 // Limit the memory usage of the job
	JOB_OBJECT_LIMIT_ACTIVE_PROCESS             = 0x00000008 // Maximum number of simultaneously active processes
	JOB_OBJECT_LIMIT_BREAKAWAY_OK               = 0x00000800 // Prevent processes from creating child processes outside the job
	JOB_OBJECT_LIMIT_DIE_ON_UNHANDLED_EXCEPTION = 0x00000400 // Forces unhandled exceptions to terminate the job
	JOB_OBJECT_LIMIT_PROCESS_MEMORY             = 0x00000100 // Limit the memory usage of each process

	// Job information classes
	JobObjectLimitInformation = 9 // The class used to set job limits

	// Create process flags
	CREATE_NEW_CONSOLE = 0x00000010 // Create a new console for the process
)

// CreateSandboxJob creates a job object for sandboxing processes
func CreateSandboxJob() (syscall.Handle, error) {
	// Create a job object
	jobHandle, _, err := CreateJobObjectA.Call(0, 0)
	if jobHandle == 0 {
		lastErr := syscall.GetLastError()
		return 0, fmt.Errorf("CreateJobObjectA failed with error code %d: %s", lastErr, err)
	}

	// Set job limits
	jobLimit := JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE |
		JOB_OBJECT_LIMIT_SUSPEND_RESUME |
		JOB_OBJECT_LIMIT_SET_INFORMATION |
		JOB_OBJECT_LIMIT_JOB_MEMORY |
		JOB_OBJECT_LIMIT_ACTIVE_PROCESS |
		JOB_OBJECT_LIMIT_BREAKAWAY_OK |
		JOB_OBJECT_LIMIT_DIE_ON_UNHANDLED_EXCEPTION |
		JOB_OBJECT_LIMIT_PROCESS_MEMORY

	// Create a structure for job limits
	limitInfo := JobObjectExtendedLimitInformation{
		BasicLimitInformation: JobObjectBasicLimitInformation{
			LimitFlags: uint32(jobLimit),
		},
	}

	// Call SetInformationJobObject to set the job limits
	ret, _, err := SetInformationJobObject.Call(
		jobHandle,                           // Handle to the job object
		uintptr(JobObjectLimitInformation),  // Information class
		uintptr(unsafe.Pointer(&limitInfo)), // Pointer to the limit information
		uintptr(unsafe.Sizeof(limitInfo)),   // Size of the limit information
	)

	if err != nil {

		// CloseHandle.Call(jobHandle)
		// lastErr := syscall.GetLastError()
		// fmt.Println(lastErr, err)
		fmt.Println(ret)
		// return 0, fmt.Errorf("SetInformationJobObject failed with error code %d: %s", lastErr, err)

	}

	// Return the handle of the job object
	return syscall.Handle(jobHandle), nil
}

func StartProcessInJob(jobHandle syscall.Handle, exePath string) error {
	var startupInfo syscall.StartupInfo
	var processInfo syscall.ProcessInformation

	startupInfo.Cb = uint32(unsafe.Sizeof(startupInfo))
	startupInfo.Flags = syscall.STARTF_USESTDHANDLES | syscall.STARTF_USESHOWWINDOW
	startupInfo.ShowWindow = syscall.SW_SHOW

	utf16ExePath, err := syscall.UTF16PtrFromString(exePath)
	if err != nil {
		return err
	}

	ret, _, err := CreateProcess.Call(
		uintptr(unsafe.Pointer(utf16ExePath)), // Pointer to the executable file path
		0,                                     // Command line arguments (none in this case)
		0,                                     // Security attributes for the process (none)
		0,                                     // Security attributes for the thread (none)
		0,                                     // Do not inherit handles
		CREATE_NEW_CONSOLE,                    // Create a new console for the process
		0,                                     // Environment block (none)
		0,                                     // Current directory (none)
		uintptr(unsafe.Pointer(&startupInfo)), // Pointer to startup information
		uintptr(unsafe.Pointer(&processInfo)), // Pointer to process information
	)

	if ret == 0 {
		return err
	}

	// Add the process to the JOB OBJECT
	if retValue, _, err := AssignProcessToJobObject.Call(uintptr(jobHandle), uintptr(processInfo.Process)); retValue == 0 {
		return err
	}
	CloseHandle.Call(uintptr(processInfo.Process))
	CloseHandle.Call(uintptr(processInfo.Thread))

	return nil
}

func IsProcessInJob(processHandle syscall.Handle, jobHandle syscall.Handle) (bool, error) {
	var processIdList [1024]uint32 // Buffer for the process ID list
	var size uint32

	ret, _, err := QueryInformationJobObject.Call(
		uintptr(jobHandle),
		uintptr(JobObjectBasicProcessIdList),
		uintptr(unsafe.Pointer(&processIdList[0])), // Pointer to the buffer
		uintptr(len(processIdList)*4),              // Size of the buffer
		uintptr(unsafe.Pointer(&size)),             // Pointer to the size variable
	)

	if ret == 0 {
		return false, err
	}

	// Calculate how many process IDs we received.
	numProcesses := size / 4 // Each process ID is 4 bytes (uint32)

	// Check if our process ID is in the list.
	var targetProcessId uint32
	handleInfo, _, err := GetProcessId.Call(uintptr(processHandle))
	if err == nil {
		targetProcessId = uint32(handleInfo)
	}

	for i := uint32(0); i < numProcesses; i++ {
		if processIdList[i] == targetProcessId {
			return true, nil
		}
	}

	return false, nil
}

func main() {
	jobHandle, err := CreateSandboxJob()
	if err != nil {
		fmt.Printf("Error during the creation of the job object: %v\n", err)
		return
	}
	defer CloseHandle.Call(uintptr(jobHandle))

	exePath := "C:/Users/ricci/aaaaaCODE/secure/go/sand/test.exe" // path of the executable
	if err := StartProcessInJob(jobHandle, exePath); err != nil {
		fmt.Printf("Error starting the process: %v\n", err)
		return
	}

	fmt.Println("Binary is in the sandbox...")

	// Check if the process is in the job
	isInJob, err := IsProcessInJob(syscall.Handle(jobHandle), jobHandle)
	if err != nil {
		fmt.Printf("Error checking if the process is in the job: %v\n", err)
		return
	}
	if isInJob {
		fmt.Println("The process is running in the job.")
	} else {
		fmt.Println("The process is not running in the job.")
	}
}
