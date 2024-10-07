package main

import (
	"fmt"
	"syscall"
	"unsafe"
)


// GUID that identify the etw provider interface
var providerGUID = syscall.GUID{
	Data1: 0xB1A5F2C5, Data2: 0x1B6D, Data3: 0x4D08,
	Data4: [8]byte{0xAB, 0xC1, 0xA0, 0xC9, 0x8D, 0x23, 0x46, 0x29},
}


var providerHandle uintptr


type EVENT_DESCRIPTOR struct {
	Id      uint16
	Version uint8
	Channel uint8
	Level   uint8
	Opcode  uint8
	Task    uint16
	Keyword uint64
}

const (
	JobObjectLimitInformation          = 9

	JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE          = 0x00002000 // Kill processes when the job is closed
	JOB_OBJECT_LIMIT_SUSPEND_RESUME             = 0x00000004 // Allow processes to be suspended and resumed
	JOB_OBJECT_LIMIT_SET_INFORMATION            = 0x00000010 // Allow setting information for the job
	JOB_OBJECT_LIMIT_JOB_MEMORY                 = 0x00000200 // Limit the memory usage of the job
	JOB_OBJECT_LIMIT_ACTIVE_PROCESS             = 0x00000008 // Maximum number of simultaneously active processes
	JOB_OBJECT_LIMIT_BREAKAWAY_OK               = 0x00000800 // Prevent processes from creating child processes outside the job
	JOB_OBJECT_LIMIT_DIE_ON_UNHANDLED_EXCEPTION = 0x00000400 // Forces unhandled exceptions to terminate the job
	JOB_OBJECT_LIMIT_PROCESS_MEMORY             = 0x00000100 // Limit the memory usage of each process
)

var (
	kernel32                 = syscall.NewLazyDLL("kernel32.dll")
	CreateJobObjectA         = kernel32.NewProc("CreateJobObjectA")
	SetInformationJobObject  = kernel32.NewProc("SetInformationJobObject")
	CreateProcess            = kernel32.NewProc("CreateProcessW")
	AssignProcessToJobObject = kernel32.NewProc("AssignProcessToJobObject")
	CloseHandle              = kernel32.NewProc("CloseHandle")
	Advapi32                 = syscall.NewLazyDLL("advapi32.dll")
	procEventRegister        = Advapi32.NewProc("EventRegister")
	procEventUnregister      = Advapi32.NewProc("EventUnregister")
	procEventWrite           = Advapi32.NewProc("EventWrite")
)
	

const (
	CREATE_NEW_CONSOLE = 0x00000010
)

type JobObjectExtendedLimitInformation struct {
	BasicLimitInformation JobObjectBasicLimitInformation
	IoInfo                IO_COUNTERS
	ProcessMemoryLimit    uintptr
	JobMemoryLimit        uintptr
	PeakProcessMemoryUsed uintptr
	PeakJobMemoryUsed     uintptr
}

type JobObjectBasicLimitInformation struct {
	PerProcessUserTimeLimit int64
	PerJobUserTimeLimit     int64
	LimitFlags              uint32
	MinimumWorkingSetSize   uintptr
	MaximumWorkingSetSize   uintptr
	ActiveProcessLimit      uint32
	Affinity                uintptr
	PriorityClass           uint32
	SchedulingClass         uint32
}

type IO_COUNTERS struct {
	ReadOperationCount  uint64
	WriteOperationCount uint64
	OtherOperationCount uint64
	ReadTransferCount   uint64
	WriteTransferCount  uint64
	OtherTransferCount  uint64
}

type Path string

// Global variable for executable path
var exePathGlobal Path

// SetFilePath sets the executable path
func SetFilePath(exePath string) error {
	exePathGlobal = Path(exePath)
	return nil
}

// StartExeInJob starts the process in the given job
func StartExeInJob(jobHandle syscall.Handle) error {
	var startupInfo syscall.StartupInfo
	var processInfo syscall.ProcessInformation

	startupInfo.Cb = uint32(unsafe.Sizeof(startupInfo))
	startupInfo.Flags = syscall.STARTF_USESTDHANDLES | syscall.STARTF_USESHOWWINDOW
	startupInfo.ShowWindow = syscall.SW_SHOW

	utf16ExePath, err := syscall.UTF16PtrFromString(string(exePathGlobal))
	if err != nil {
		return err
	}

	ret, _, err := CreateProcess.Call(
		uintptr(unsafe.Pointer(utf16ExePath)),
		0,
		0,
		0,
		0,
		CREATE_NEW_CONSOLE,
		0,
		0,
		uintptr(unsafe.Pointer(&startupInfo)),
		uintptr(unsafe.Pointer(&processInfo)),
	)

	if ret == 0 {
		return fmt.Errorf("CreateProcess failed: %v", err)
	}

	// Assign the process to the job object
	if retValue, _, err := AssignProcessToJobObject.Call(uintptr(jobHandle), uintptr(processInfo.Process)); retValue == 0 {
		return err
	}

	CloseHandle.Call(uintptr(processInfo.Process))
	CloseHandle.Call(uintptr(processInfo.Thread))

	return nil
}

// CreateSandboxJob creates a job object with specific limits
func CreateSandboxJob() (syscall.Handle, error) {
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

	limitInfo := JobObjectExtendedLimitInformation{
		BasicLimitInformation: JobObjectBasicLimitInformation{
			LimitFlags: uint32(jobLimit),
		},
	}

	ret, _, err := SetInformationJobObject.Call(
		jobHandle,
		uintptr(JobObjectLimitInformation),
		uintptr(unsafe.Pointer(&limitInfo)),
		uintptr(unsafe.Sizeof(limitInfo)),
	)

	if ret == 0 {
		lastErr := syscall.GetLastError()
		return 0, fmt.Errorf("SetInformationJobObject failed with error code %d: %s", lastErr, err)
	}

	return syscall.Handle(jobHandle), nil
}


func registerProvider() error {
	r, _, err := procEventRegister.Call(
		uintptr(unsafe.Pointer(&providerGUID)),
		0, 0,
		uintptr(unsafe.Pointer(&providerHandle)),
	)
	if r != 0 {
		return err
	}
	return nil
}

func unregisterProvider() error {
	r, _, err := procEventUnregister.Call(providerHandle)
	if r != 0 {
		return err
	}
	return nil
}


func writeEvent(eventDescriptor *EVENT_DESCRIPTOR, message string) error {
	utf16Message, err := syscall.UTF16PtrFromString(message)
	if err != nil {
		return err
	}

	r, _, syscallErr := procEventWrite.Call(
		providerHandle,
		uintptr(unsafe.Pointer(eventDescriptor)),
		1,
		uintptr(unsafe.Pointer(utf16Message)),
	)
	if r != 0 {
		return syscallErr
	}
	return nil
}


// Help function to demonstrate the usage of SetFilePath, CreateSandboxJob, and StartExeInJob
func Help() {
	fmt.Println("=== Help Function Demonstration ===")
	fmt.Println("This example will show how to set the executable path, create a sandbox job, and run the executable within the sandbox.\n")

	// Step 1: Set the executable path
	fmt.Println("Step 1: Setting the executable path")
	fmt.Println("Code: err := SetFilePath(\"C:\\Path\\To\\YourExecutable.exe\")")
	err := SetFilePath("C:\\Path\\To\\YourExecutable.exe")
	if err != nil {
		fmt.Printf("Error setting executable path: %v\n", err)
		return
	}
	fmt.Println("Executable path set successfully.\n")

	// Step 2: Create a job object
	fmt.Println("Step 2: Creating a sandbox job object")
	fmt.Println("Code: jobHandle, err := CreateSandboxJob()")
	jobHandle, err := CreateSandboxJob()
	if err != nil {
		fmt.Printf("Error creating job object: %v\n", err)
		return
	}
	fmt.Println("Sandbox job object created successfully.")

	// Ensure the job handle is properly closed when no longer needed
	defer func() {
		fmt.Println("\nClosing the job handle")
		fmt.Println("Code: CloseHandle.Call(uintptr(jobHandle))")
		CloseHandle.Call(uintptr(jobHandle))
		fmt.Println("Job handle closed.")
	}()

	// Step 3: Start the executable process in the job object
	fmt.Println("\nStep 3: Starting the process in the sandbox job")
	fmt.Println("Code: err = StartExeInJob(jobHandle)")
	err = StartExeInJob(jobHandle)
	if err != nil {
		fmt.Printf("Error starting process in job object: %v\n", err)
		return
	}
	fmt.Println("Process successfully started in the sandbox.\n")

	fmt.Println("=== End of Help Demonstration ===")
}

func main() {
	return
}


// compile using ```go build -o fullblock.dll -buildmode=c-shared .\FullBloccking.go```
