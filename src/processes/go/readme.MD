---

# Process Management CLI Tool - Documentation

## Overview

This tool provides a comprehensive set of commands to manage and interact with processes on a Windows system using golang. Below is a detailed description of each available function, including usage examples, arguments, return values, and exceptions.

## Functions

### 0. help

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

### 1. `list`

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

### 2. `info`

#### Description
Retrieves detailed information about a specific process by its PID.

#### Usage
```bash
ProcHandle.exe info <proc_name>
```

#### Arguments
- `<process_name>` (string): The Process name that you want to query.

#### Returns
- Process details, including memory usage, executable path, start time, and priority class.

#### Exceptions
- **InvalidPID**: If the provided PID does not exist or is not accessible.
- **PermissionDenied**: If the tool lacks the permissions to access the process details.

### 3. `terminate`

#### Description
Terminates a process given its PID.

#### Usage
```bash
ProcHandle.exe terminate <proc_name>
```

#### Arguments
- `<proc_name>` (string): The Process name that you want to terminate.

#### Returns
- Confirmation message indicating if the process was terminated successfully.

#### Exceptions
- **PermissionDenied**: If the tool does not have permissions to terminate the process.
- **InvalidPID**: If the provided PID does not exist.
- **AccessDenied**: If the process cannot be terminated (e.g., a critical system process).

### 4. `set-priority`

#### Description
Sets the priority class for a given process.

#### Usage
```bash
ProcHandle.exe set-priority <proc_name> <priority>
```

#### Arguments
- `<Proc_name>` (string): The Process name whose priority you want to change.
- `<priority>` (string): The priority class. Possible values are: `low`, `normal`, `high`, `realtime`.

#### Returns
- Confirmation message indicating if the priority was successfully changed.

#### Exceptions
- **InvalidPID**: If the provided PID does not exist.
- **InvalidPriority**: If the specified priority class is not valid.
- **PermissionDenied**: If the tool lacks the permissions to modify the priority of the process.

### 5. `suspend`

#### Description
Suspends all threads of a given process.

#### Usage
```bash
ProcHandle.exe suspend <proc_name>
```

#### Arguments
- `<proc_name>` (string): The Process name that u want to suspend.

#### Returns
- Confirmation message indicating if the process was successfully suspended.

#### Exceptions
- **InvalidPID**: If the provided PID does not exist.
- **PermissionDenied**: If the tool does not have permissions to suspend the process.

### 6. `resume`

#### Description
Resumes a suspended process.

#### Usage
```bash
ProcHandle.exe resume <proc_name>
```

#### Arguments
- `<Proc_name>` (string): The Process name that you want to resume.

#### Returns
- Confirmation message indicating if the process was successfully resumed.

#### Exceptions
- **InvalidPID**: If the provided PID does not exist.
- **ProcessNotSuspended**: If the specified process is not in a suspended state.

### 7. `read-memory`

#### Description
Reads a specific memory address of a given process.

#### Usage
```bash
ProcHandle.exe read-memory <proc_name> <address> <size>
```

#### Arguments
- `<proc_name>` (string): The Process name of the target process.
- `<address>` (hexadecimal): The memory address to read.
- `<size>` (integer): The number of bytes to read.

#### Returns
- The content of the memory at the specified address.

#### Exceptions
- **InvalidPID**: If the provided PID does not exist.
- **InvalidAddress**: If the provided address is not valid.
- **PermissionDenied**: If the tool lacks permissions to read the memory of the process.

### 8. `write-memory`

#### Description
Writes data to a specific memory address of a process.

#### Usage
```bash
ProcHandle.exe write-memory <proc_name> <address> <data>
```

#### Arguments
- `<proc_name>` (string): The process name of the target process.
- `<address>` (hexadecimal): The memory address to write to.
- `<data>` (string): The data to write at the specified address.

#### Returns
- Confirmation message indicating if the data was successfully written.

#### Exceptions
- **InvalidPID**: If the provided PID does not exist.
- **InvalidAddress**: If the provided address is not valid.
- **PermissionDenied**: If the tool lacks permissions to write to the memory of the process.

### 9. `protect`

#### Description
Protect address of the specified process with a specific permission

#### Usage
```bash
protect <process_name> <lpAddress> <dwSize> <flNewProtect>
```

#### Arguments
- `<process_name>` (string): The name of the target process whose memory permissions you want to modify.
- `<lpAddress>` (hexadecimal): The memory address of the region that will be protected.
- `<dwSize>` (int): The size (in bytes) of the memory region that you want to protect.
- `<flNewProtect>` (string or integer): The new protection level for the specified memory region. This parameter should be provided as a valid memory protection constant, such as `PAGE_READONLY`, `PAGE_READWRITE`, `PAGE_EXECUTE`, etc.

#### Returns
- A confirmation message indicating if the memory protection was successfully changed. If successful, it shows that the specified address was modified with the new permissions.

#### Exceptions
- **InvalidPID**: Raised when the process ID for the given process name does not exist.
- **InvalidAddress**: Raised when the provided address is not valid within the memory space of the process.
- **PermissionDenied**: Raised if the tool lacks permissions to modify the memory of the target process. This usually occurs if the script is not run with sufficient privileges (e.g., administrator rights).



## Example Usage

Here are some examples of how to use the tool:

1. List all running processes:
   ```bash
   ProcHandle.exe list
   ```

2. Get information about a specific process:
   ```bash
   ProcHandle.exe info code
   ```

3. Terminate a process:
   ```bash
   ProcHandle.exe terminate code
   ```

4. Change the priority of a process:
   ```bash
   ProcHandle.exe set-priority code high
   ```

5. Suspend a process:
   ```bash
   ProcHandle.exe suspend code
   ```

6. Resume a suspended process:
   ```bash
   ProcHandle.exe resume code
   ```

7. Read memory from a process:
   ```bash
   ProcHandle.exe read-memory code 0x00400000 256
   ```

8. Write data to a process's memory:
   ```bash
   ProcHandle.exe write-memory code 0x00400000 "Hello World"
   ```
   
9. Protect an address region:
   ```bash
   ProcHandle.exe protect code 0x7ffdf000 4096 PAGE_READWRITE
   ```

## FINALLY

This tool is a powerful command-line utility for managing processes in a Windows environment. It is essential to run it with appropriate permissions to access or modify processes successfully. Use the `help` command to see all available commands and their descriptions.
