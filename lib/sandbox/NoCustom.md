# fullblock.dll

## Overview

The `fullblock.dll` provides functions for managing Windows job objects and logging events using Event Tracing for Windows (ETW). It allows you to create a sandboxed environment for running processes with specified resource limits.

### Key Functions

1. **SetFilePath**: Sets the path to the executable to run.
2. **CreateSandboxJob**: Creates a job object with specific resource limits.
3. **StartExeInJob**: Starts the executable in the created job object (so before use the ```CreateSandboxJob``` function).
4. **registerProvider**: Registers an ETW provider.
5. **unregisterProvider**: Unregisters the ETW provider.
6. **writeEvent**: Writes an event to the ETW provider.

### Function Documentation

#### 1. `SetFilePath(exePath Path) error`

- **Description**: Sets the global executable path for the process.
- **Parameters**: 
  - `exePath (string)`: The path of the executable, u pass the string it will be converte din the local type path.
- **Returns**: 
  - `error`: An error if setting the path fails.

**Go Usage Example**:
```go
err := SetFilePath("C:\\Path\\To\\YourExecutable.exe")
if err != nil {
    // Handle error
}
```

**C++ Usage Snippet**:
```cpp
typedef int (*SetFilePathFunc)(const wchar_t*);
SetFilePathFunc SetFilePath = (SetFilePathFunc)GetProcAddress(hInstLibrary, "SetFilePath");
int result = SetFilePath(L"C:\\Path\\To\\YourExecutable.exe");
```

#### 2. `CreateSandboxJob() (syscall.Handle, error)`

- **Description**: Creates a job object to impose resource limits.
- **Returns**: 
  - `syscall.Handle`: Handle to the created job object.
  - `error`: An error if creation fails.

**Go Usage Example**:
```go
jobHandle, err := CreateSandboxJob()
if err != nil {
    // Handle error
}
defer CloseHandle.Call(uintptr(jobHandle)) // Ensure handle is closed
```

**C++ Usage Snippet**:
```cpp
typedef HANDLE (*CreateSandboxJobFunc)();
CreateSandboxJobFunc CreateSandboxJob = (CreateSandboxJobFunc)GetProcAddress(hInstLibrary, "CreateSandboxJob");
HANDLE jobHandle = CreateSandboxJob();
```

#### 3. `StartExeInJob(jobHandle syscall.Handle) error`

- **Description**: Starts the executable set by `SetFilePath()` in the given job object.
- **Parameters**: 
  - `jobHandle (syscall.Handle)`: Handle to the job object.
- **Returns**: 
  - `error`: An error if starting the process fails.

**Go Usage Example**:
```go
err = StartExeInJob(jobHandle)
if err != nil {
    // Handle error
}
```

**C++ Usage Snippet**:
```cpp
typedef int (*StartExeInJobFunc)(HANDLE);
StartExeInJobFunc StartExeInJob = (StartExeInJobFunc)GetProcAddress(hInstLibrary, "StartExeInJob");
int result = StartExeInJob(jobHandle);
```

#### 4. `registerProvider() error`

- **Description**: Registers the ETW provider.
- **Returns**: 
  - `error`: An error if registration fails.

**Go Usage Example**:
```go
err := registerProvider()
if err != nil {
    // Handle error
}
```

**C++ Usage Snippet**:
```cpp
typedef int (*RegisterProviderFunc)();
RegisterProviderFunc registerProvider = (RegisterProviderFunc)GetProcAddress(hInstLibrary, "registerProvider");
int result = registerProvider();
```

#### 5. `unregisterProvider() error`

- **Description**: Unregisters the ETW provider.
- **Returns**: 
  - `error`: An error if unregistration fails.

**Go Usage Example**:
```go
err := unregisterProvider()
if err != nil {
    // Handle error
}
```

**C++ Usage Snippet**:
```cpp
typedef int (*UnregisterProviderFunc)();
UnregisterProviderFunc unregisterProvider = (UnregisterProviderFunc)GetProcAddress(hInstLibrary, "unregisterProvider");
int result = unregisterProvider();
```

#### 6. `writeEvent(eventDescriptor *EVENT_DESCRIPTOR, message string) error`

- **Description**: Writes an event to the ETW provider.
- **Parameters**: 
  - `eventDescriptor`: Pointer to the event descriptor.
  - `message`: The message to log.
- **Returns**: 
  - `error`: An error if writing the event fails.

**Go Usage Example**:
```go
eventDesc := &EVENT_DESCRIPTOR{Id: 1, Level: 2}
err = writeEvent(eventDesc, "Event message")
if err != nil {
    // Handle error
}
```

**C++ Usage Snippet**:
```cpp
typedef int (*WriteEventFunc)(EVENT_DESCRIPTOR*, const wchar_t*);
WriteEventFunc writeEvent = (WriteEventFunc)GetProcAddress(hInstLibrary, "writeEvent");
EVENT_DESCRIPTOR eventDesc = {1, 2, 0, 0, 0, 0, 0}; // Set appropriate values
int result = writeEvent(&eventDesc, L"Event message");
```

### Build Instructions

1. Ensure you have Go installed.
2. Compile the Go file to create the DLL:
   ```bash
   go build -o fullblock.dll -buildmode=c-shared your_go_file.go
   ```
## CURRENT JOB LIMITS:
		JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE 
		JOB_OBJECT_LIMIT_SUSPEND_RESUME 
		JOB_OBJECT_LIMIT_SET_INFORMATION
		JOB_OBJECT_LIMIT_JOB_MEMORY 
		JOB_OBJECT_LIMIT_ACTIVE_PROCESS 
		JOB_OBJECT_LIMIT_BREAKAWAY_OK 
		JOB_OBJECT_LIMIT_DIE_ON_UNHANDLED_EXCEPTION 
		JOB_OBJECT_LIMIT_PROCESS_MEMORY
