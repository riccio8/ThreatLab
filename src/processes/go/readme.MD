---

# Process Management CLI Tool - Documentation

## Overview

This tool provides a comprehensive set of commands to manage and interact with processes on a Windows system using golang. Below is a detailed description of each available function, including usage examples, arguments, return values, and exceptions.

## Functions

### 0. `help`

#### Description
Print the help message

#### Usage
```bash
ProcHandle.exe help
```

#### Arguments
- None

#### Returns
- A list of the args

#### Exceptions
- **None** (i hope)

### 1. `ListProcesses`

#### Description
Lists all the running processes on the system.

#### Usage
```bash
ProcHandle.exe list
```

#### Arguments
- None

#### Returns
- A list of running processes, including their PID (Process ID) and process name.

#### Exceptions
- **PermissionDenied**: If the tool doesn't have sufficient permissions to access the list of processes.
- **SnapshotFailed**: If there is an error creating the snapshot of processes.

### 2. `GetProcessInfo`

#### Description
Retrieves detailed information about a specific process by its PID.

#### Usage
```bash
ProcHandle.exe info <pid>
```

#### Arguments
- `<pid>` (integer): The Process ID of the process you want to query.

#### Returns
- Process details, including memory usage, executable path, start time, and priority class.

#### Exceptions
- **InvalidPID**: If the provided PID does not exist or is not accessible.
- **PermissionDenied**: If the tool lacks the permissions to access the process details.

### 3. `TerminateProcess`

#### Description
Terminates a process given its PID.

#### Usage
```bash
ProcHandle.exe terminate <pid>
```

#### Arguments
- `<pid>` (integer): The Process ID of the process you want to terminate.

#### Returns
- Confirmation message indicating if the process was terminated successfully.

#### Exceptions
- **PermissionDenied**: If the tool does not have permissions to terminate the process.
- **InvalidPID**: If the provided PID does not exist.
- **AccessDenied**: If the process cannot be terminated (e.g., a critical system process).

### 4. `SetProcessPriority`

#### Description
Sets the priority class for a given process.

#### Usage
```bash
ProcHandle.exe set-priority <pid> <priority>
```

#### Arguments
- `<pid>` (integer): The Process ID of the process whose priority you want to change.
- `<priority>` (string): The priority class. Possible values are: `low`, `below-normal`, `normal`, `above-normal`, `high`, `realtime`.

#### Returns
- Confirmation message indicating if the priority was successfully changed.

#### Exceptions
- **InvalidPID**: If the provided PID does not exist.
- **InvalidPriority**: If the specified priority class is not valid.
- **PermissionDenied**: If the tool lacks the permissions to modify the priority of the process.

### 5. `SuspendProcess`

#### Description
Suspends all threads of a given process.

#### Usage
```bash
ProcHandle.exe suspend <pid>
```

#### Arguments
- `<pid>` (integer): The Process ID of the process you want to suspend.

#### Returns
- Confirmation message indicating if the process was successfully suspended.

#### Exceptions
- **InvalidPID**: If the provided PID does not exist.
- **PermissionDenied**: If the tool does not have permissions to suspend the process.

### 6. `ResumeProcess`

#### Description
Resumes a suspended process.

#### Usage
```bash
ProcHandle.exe resume <pid>
```

#### Arguments
- `<pid>` (integer): The Process ID of the process you want to resume.

#### Returns
- Confirmation message indicating if the process was successfully resumed.

#### Exceptions
- **InvalidPID**: If the provided PID does not exist.
- **ProcessNotSuspended**: If the specified process is not in a suspended state.

### 7. `ReadMemory`

#### Description
Reads a specific memory address of a given process.

#### Usage
```bash
ProcHandle.exe read-memory <pid> <address> <size>
```

#### Arguments
- `<pid>` (integer): The Process ID of the target process.
- `<address>` (hexadecimal): The memory address to read.
- `<size>` (integer): The number of bytes to read.

#### Returns
- The content of the memory at the specified address.

#### Exceptions
- **InvalidPID**: If the provided PID does not exist.
- **InvalidAddress**: If the provided address is not valid.
- **PermissionDenied**: If the tool lacks permissions to read the memory of the process.

### 8. `WriteMemory`

#### Description
Writes data to a specific memory address of a given process.

#### Usage
```bash
ProcHandle.exe write-memory <pid> <address> <data>
```

#### Arguments
- `<pid>` (integer): The Process ID of the target process.
- `<address>` (hexadecimal): The memory address to write to.
- `<data>` (hexadecimal string): The data to write at the specified address.

#### Returns
- Confirmation message indicating if the memory was written successfully.

#### Exceptions
- **InvalidPID**: If the provided PID does not exist.
- **InvalidAddress**: If the provided address is not valid.
- **PermissionDenied**: If the tool lacks permissions to write to the memory of the process.

## Error Handling

All functions provide meaningful error messages to help users diagnose issues quickly. Common error codes include:

- **PermissionDenied**: Indicates a lack of privileges to perform the requested action. You may need to run the tool with elevated privileges (Administrator mode).
- **InvalidPID**: The specified Process ID does not correspond to an active process. Ensure the PID is correct.
- **InvalidAddress**: The memory address provided is either out of range or cannot be accessed.
- **AccessDenied**: Some system processes or processes owned by other users may not be accessible.

## License

This tool is provided under the MIT License. See `LICENSE` for more details.

## Contributing

Contributions are welcome! Please read `CONTRIBUTING.md` for information on how to contribute to this project.

## Contact

For further information or support, feel free to open an issue or contact the maintainer at riccioadami@gmail.com.

---
