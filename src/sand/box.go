package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	kernel32                 = syscall.NewLazyDLL("Kernel32.dll")
	advapi32                 = syscall.NewLazyDLL("advapi32.dll")
	procEventRegister        = advapi32.NewProc("EventRegister")
	CreateJobObjectA         = kernel32.NewProc("CreateJobObjectA")
	SetInformationJobObject  = kernel32.NewProc("SetInformationJobObject")
	CreateProcess            = kernel32.NewProc("CreateProcessW") // Use CreateProcessW for Unicode support
	CloseHandle              = kernel32.NewProc("CloseHandle")
	AssignProcessToJobObject = kernel32.NewProc("AssignProcessToJobObject")
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
	IoInfo                syscall.IoCounters
	ProcessMemoryLimit    uintptr
	JobMemoryLimit        uintptr
}
const (
	// Job object flags
	JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE          = 0x00002000 // Kill processes when the job is closed
	JOB_OBJECT_LIMIT_SUSPEND_RESUME             = 0x00000004 // Allow processes to be suspended and resumed
	JOB_OBJECT_LIMIT_SET_INFORMATION            = 0x00000010 // Allow setting information for the job
	JOB_OBJECT_LIMIT_JOB_MEMORY                 = 0x00000200 // Limit the memory usage of the job
	JOB_OBJECT_LIMIT_ACTIVE_PROCESS             = 0x00000008 // Maximum number of simultaneously active processes
	JOB_OBJECT_LIMIT_AFFINITY                   = 0x00000010 // All processes use the same processor affinity
	JOB_OBJECT_LIMIT_BREAKAWAY_OK               = 0x00000800 // Prevent processes from creating child processes outside the job
	JOB_OBJECT_LIMIT_DIE_ON_UNHANDLED_EXCEPTION = 0x00000400 // Forces unhandled exceptions to terminate the job
	JOB_OBJECT_LIMIT_JOB_TIME                   = 0x00000004 // User-mode execution time limit for the job
	JOB_OBJECT_LIMIT_PRESERVE_JOB_TIME          = 0x00000040 // Preserve previously set job time limits
	JOB_OBJECT_LIMIT_PRIORITY_CLASS             = 0x00000020 // All processes use the same priority class
	JOB_OBJECT_LIMIT_PROCESS_MEMORY             = 0x00000100 // Limit the memory usage of each process
	JOB_OBJECT_LIMIT_PROCESS_TIME               = 0x00000002 // User-mode execution time limit for each process
	JOB_OBJECT_LIMIT_SCHEDULING_CLASS           = 0x00000080 // All processes use the same scheduling class
	JOB_OBJECT_LIMIT_SILENT_BREAKAWAY_OK        = 0x00001000 // Allows breakaway of child processes
	JOB_OBJECT_LIMIT_SUBSET_AFFINITY            = 0x00004000 // Allows processes to use a subset of the processor affinity
	JOB_OBJECT_LIMIT_WORKINGSET                 = 0x00000001 // Limit the working set size of each process

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
		lastErr := syscall.GetLastError() // Get the last error code
		return 0, fmt.Errorf("CreateJobObjectA failed with error code %d: %s", lastErr, err)
	}

	// Set job limits 
	jobLimit := JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE |
		JOB_OBJECT_LIMIT_SUSPEND_RESUME |
		JOB_OBJECT_LIMIT_SET_INFORMATION |
		JOB_OBJECT_LIMIT_JOB_MEMORY |
		JOB_OBJECT_LIMIT_ACTIVE_PROCESS |
		JOB_OBJECT_LIMIT_AFFINITY |
		JOB_OBJECT_LIMIT_BREAKAWAY_OK | // Prevent processes from creating child processes outside the job
		JOB_OBJECT_LIMIT_DIE_ON_UNHANDLED_EXCEPTION |
		JOB_OBJECT_LIMIT_JOB_TIME |
		JOB_OBJECT_LIMIT_PRESERVE_JOB_TIME |
		JOB_OBJECT_LIMIT_PRIORITY_CLASS |
		JOB_OBJECT_LIMIT_PROCESS_MEMORY |
		JOB_OBJECT_LIMIT_PROCESS_TIME |
		JOB_OBJECT_LIMIT_SCHEDULING_CLASS |
		JOB_OBJECT_LIMIT_SILENT_BREAKAWAY_OK |
		JOB_OBJECT_LIMIT_SUBSET_AFFINITY |
		JOB_OBJECT_LIMIT_WORKINGSET

	// Create a structure for job limits
	limitInfo := JobObjectBasicLimitInformation{
		LimitFlags: uint32(jobLimit),
	}

	// Call SetInformationJobObject to set the job limits
	ret, _, err := SetInformationJobObject.Call(
		jobHandle,                                   // Handle to the job object
		uintptr(JobObjectLimitInformation),          // Information class
		uintptr(unsafe.Pointer(&limitInfo)),         // Pointer to the limit information
		uintptr(unsafe.Sizeof(limitInfo)),           // Size of the limit information
	)

	if ret == 0 {
		CloseHandle.Call(jobHandle)
		lastErr := syscall.GetLastError() // Get the last error code
		return 0, fmt.Errorf("SetInformationJobObject failed with error code %d: %s", lastErr, err)
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

func main() {
	jobHandle, err := CreateSandboxJob()
	if err != nil {
		fmt.Printf("Error during the creation of the job object: %v\n", err)
		return
	}
	defer CloseHandle.Call(uintptr(jobHandle))

	exePath := "C:\\Users\\ricci\\aaaaaCODE\\secure\\go\\sand\\box.exe" // path of the executable
	if err := StartProcessInJob(jobHandle, exePath); err != nil {
		fmt.Printf("Error starting the process: %v\n", err)
		return
	}

	fmt.Println("Binary is in the sandbox...")
}
