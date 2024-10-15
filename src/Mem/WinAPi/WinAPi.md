# THIS FILE CONTAINS A LOT OF INTRESTING WIN APIS THAT COULD USE FOR CONTROL THE DEVICE (almost for reverse engeeniring):

## 1. File and Process Manipulation APIs

These APIs allow interacting with files, directories, and processes on Windows.

- **CreateFile** – Opens files or devices.
- **ReadFile** – Reads data from an open file.
- **WriteFile** – Writes data to an open file.
- **CloseHandle** – Closes file, process, or other handles.
- **CreateProcess** – Starts a new process.
- **TerminateProcess** – Ends a specific process.
- **OpenProcess** – Opens an existing process for manipulation.
- **GetExitCodeProcess** – Gets the exit code of a process.
- **VirtualAllocEx** – Allocates memory in another process's space.
- **WriteProcessMemory** – Writes memory in another process.

## 2. Networking and Communication APIs

These APIs allow opening network ports, managing connections, and exchanging data.

- **Socket** – Creates a network socket.
- **Bind** – Binds a socket to a local address.
- **Listen** – Puts a socket in listening mode for connections.
- **Accept** – Accepts incoming connections.
- **Connect** – Connects a socket to a remote server.
- **Send** – Sends data through a socket.
- **Recv** – Receives data through a socket.
- **WSAStartup** – Initializes the Winsock socket library.
- **GetHostByName** – Resolves a host name to an IP address.
- **GetAddrInfo** – Resolves names and IP addresses (modern alternative to `GetHostByName`).
- **HttpSendRequest** – Sends HTTP requests (via WinINet).
- **InternetOpen** – Opens an Internet connection.
- **InternetConnect** – Connects to a remote server over HTTP, FTP, or other protocols.
- **InternetReadFile** – Reads data from an Internet connection.

## 3. Memory Manipulation and Security APIs

These APIs provide access to and manipulation of system or other process memory.

- **OpenProcessToken** – Opens an access token for a process.
- **LookupPrivilegeValue** – Looks up privileges associated with a token.
- **AdjustTokenPrivileges** – Adjusts the privileges of an access token.
- **GetModuleHandle** – Gets the handle of a loaded module (DLL).
- **GetProcAddress** – Gets the address of an exported function from a DLL.
- **VirtualProtect** – Changes memory protection permissions.
- **VirtualQueryEx** – Gets information about a process's memory.
- **HeapAlloc** – Allocates memory on a heap.
- **HeapFree** – Frees memory from a heap.
- **RtlMoveMemory** – Copies data in memory (similar to `memcpy`).
- **CryptAcquireContext** – Creates a handle for a cryptographic context.
- **CryptEncrypt** – Encrypts data using a cryptographic context.

## 4. Portable Executable (PE) File Manipulation APIs

These APIs allow inspecting or manipulating PE binary files (executables, DLLs).

- **ImageNtHeader** – Gets the NT header of a PE file.
- **ImageDirectoryEntryToData** – Accesses data in a specific directory of a PE file.
- **CheckSumMappedFile** – Computes or verifies the checksum of a PE file.
- **MapViewOfFile** – Maps a file into memory for faster access.
- **UnmapViewOfFile** – Unmaps a file from memory.
- **GetFileSize** – Gets the size of a file.

## 5. Registry and System Service APIs

These APIs allow interacting with the Windows registry or managing services.

- **RegOpenKeyEx** – Opens a registry key.
- **RegQueryValueEx** – Retrieves data associated with a registry value.
- **RegSetValueEx** – Sets a value in the registry.
- **OpenSCManager** – Opens the Service Control Manager (for managing services).
- **CreateService** – Creates a new system service.
- **StartService** – Starts a service.
- **ControlService** – Sends control commands to a service.
- **DeleteService** – Deletes a service from the system.

## 6. Command and Shell Execution APIs

These APIs allow executing commands and interacting with the Windows shell.

- **ShellExecute** – Executes an application or opens a document.
- **WinExec** – Executes a specified application.
- **CreateProcessWithLogonW** – Executes a process with specific credentials.
- **SystemParametersInfo** – Retrieves or sets system information.


### 1. Module, Library, and Driver Management APIs

These APIs allow loading and managing DLLs, modules, and low-level drivers.

- **LoadLibrary** – Loads a DLL into a process.
- **FreeLibrary** – Unloads a loaded DLL.
- **GetModuleFileName** – Gets the file name of a module.
- **LoadLibraryEx** – Loads a library with extended options.
- **EnumDeviceDrivers** – Lists device drivers in the system.
- **GetDeviceDriverBaseName** – Gets the file name of a device driver.
- **GetDeviceDriverFileName** – Gets the full path of a device driver.
- **QueryDosDevice** – Gets the DOS device mappings.
- **NtLoadDriver** – Loads a driver into the kernel (native function, requires elevated privileges).
- **NtUnloadDriver** – Unloads a driver from the kernel (native function).

### 2. Thread and Synchronization APIs

These APIs allow managing threads, semaphores, and synchronization between processes or threads.

- **CreateThread** – Creates a new thread in a process.
- **SuspendThread** – Suspends a thread.
- **ResumeThread** – Resumes a suspended thread.
- **Sleep** – Puts a thread to sleep for a specified time.
- **WaitForSingleObject** – Waits for a single object (like a thread or process) to complete.
- **WaitForMultipleObjects** – Waits for multiple objects to complete.
- **SetEvent** – Sets an event state for thread synchronization.
- **ResetEvent** – Resets a synchronized event state.
- **CreateMutex** – Creates a mutex for synchronizing access to shared resources.
- **ReleaseMutex** – Releases a mutex.
- **CreateSemaphore** – Creates a semaphore for resource synchronization.
- **ReleaseSemaphore** – Releases a semaphore.
- **EnterCriticalSection** – Enters a critical section (for thread-safe access).
- **LeaveCriticalSection** – Exits a critical section.
- **InterlockedIncrement** – Atomically increments a variable (thread-safe).
- **InterlockedDecrement** – Atomically decrements a variable (thread-safe).

### 3. System and Resource Management APIs

These APIs provide access to system information, resource management, and statistics.

- **GetSystemInfo** – Retrieves system information (CPU, architecture, etc.).
- **GlobalMemoryStatusEx** – Gets information about system memory.
- **GetTickCount** – Gets the system uptime in milliseconds.
- **GetSystemTime** – Retrieves the system time (in UTC).
- **SetSystemTime** – Sets the system time (requires admin privileges).
- **GetCurrentProcessId** – Gets the ID of the current process.
- **GetCurrentThreadId** – Gets the ID of the current thread.
- **GetProcessAffinityMask** – Gets a process's CPU core affinity mask.
- **SetProcessAffinityMask** – Sets a process's CPU core affinity mask.
- **GetPerformanceInfo** – Retrieves system resource usage statistics.

### 4. GUI and Window Management APIs

These APIs allow interacting with the Windows GUI, manipulating windows and controls.

- **FindWindow** – Finds an existing window by title or class.
- **ShowWindow** – Changes the display state of a window (minimize, maximize, etc.).
- **SetWindowText** – Sets a window's title text.
- **SendMessage** – Sends a message to a window or control.
- **PostMessage** – Sends a message to a window or control asynchronously.
- **EnumWindows** – Lists all open windows in the system.
- **GetWindowThreadProcessId** – Gets the process and thread ID associated with a window.
- **GetForegroundWindow** – Retrieves the currently active foreground window.
- **SetForegroundWindow** – Sets a window as the active foreground window.
- **MoveWindow** – Changes a window's position and size.
- **GetWindowRect** – Retrieves the coordinates and size of a window.
- **ClipCursor** – Restricts the mouse cursor's movement within a window.

### 5. Hardware and Peripheral Access APIs

These APIs allow interacting with hardware devices, like the keyboard and mouse.

- **GetAsyncKeyState** – Checks the state of a key on the keyboard.
- **GetKeyState** – Gets the state of a key (pressed or not).
- **mouse_event** – Simulates mouse events (movement, clicks, scroll).
- **keybd_event** – Simulates keyboard events (key press and release).
- **MapVirtualKey** – Converts a virtual key code to a scan code and vice versa.
- **RegisterHotKey** – Registers a key combination as a system-wide hotkey.
- **UnregisterHotKey** – Deregisters a system-wide hotkey.

# MORE IN THIS PDF:
[winapi_hack en.pdf](https://github.com/user-attachments/files/17379229/winapi_hack.en.pdf)
