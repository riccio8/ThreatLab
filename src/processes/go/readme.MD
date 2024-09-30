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
ProcHandle.exe info <pid>
```

#### Arguments
- `<pid>` (integer): The Process ID of the process you want to query.

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

### 4. `set-priority`

#### Description
Sets the priority class for a given process.

#### Usage
```bash
ProcHandle.exe set-priority <pid> <priority>
```

#### Arguments
- `<pid>` (integer): The Process ID of the process whose priority you want to change.
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
ProcHandle.exe suspend <pid>
```

#### Arguments
- `<pid>` (integer): The Process ID of the process you want to suspend.

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
ProcHandle.exe resume <pid>
```

#### Arguments
- `<pid>` (integer): The Process ID of the process you want to resume.

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

### 8. `write-memory`

#### Description
Writes data to a specific memory address of a process.

#### Usage
```bash
ProcHandle.exe write-memory <pid> <address> <data>
```

#### Arguments
- `<pid>` (integer): The Process ID of the target process.
- `<address>` (hexadecimal): The memory address to write to.
- `<data>` (string): The data to write at the specified address.

#### Returns
- Confirmation message indicating if the data was successfully written.

#### Exceptions
- **InvalidPID**: If the provided PID does not exist.
- **InvalidAddress**: If the provided address is not valid.
- **PermissionDenied**: If the tool lacks permissions to write to the memory of the process.

## Example Usage

Here are some examples of how to use the tool:

1. List all running processes:
   ```bash
   ProcHandle.exe list
   ```

2. Get information about a specific process:
   ```bash
   ProcHandle.exe info 1234
   ```

3. Terminate a process:
   ```bash
   ProcHandle.exe terminate 1234
   ```

4. Change the priority of a process:
   ```bash
   ProcHandle.exe set-priority 1234 high
   ```

5. Suspend a process:
   ```bash
   ProcHandle.exe suspend 1234
   ```

6. Resume a suspended process:
   ```bash
   ProcHandle.exe resume 1234
   ```

7. Read memory from a process:
   ```bash
   ProcHandle.exe read-memory 1234 0x00400000 256
   ```

8. Write data to a process's memory:
   ```bash
   ProcHandle.exe write-memory 1234 0x00400000 "Hello World"
   ```

## FINALLY

This tool is a powerful command-line utility for managing processes in a Windows environment. It is essential to run it with appropriate permissions to access or modify processes successfully. Use the `help` command to see all available commands and their descriptions.
```
