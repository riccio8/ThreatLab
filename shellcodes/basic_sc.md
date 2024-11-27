### **1. File Deletion**

#### **File Deletion**

**Windows Functions Used:**

- `DeleteFile` - Deletes the specified file from the filesystem.

---

### **2. Creating a New File**

#### **Creating a New File**

**Windows Functions Used:**

- `CreateFile` - Creates a new file with the desired name and attributes.
- `WriteFile` - Writes data into the created file.

---

### **3. Writing Data to a File**

#### **Writing Data to a File**

**Windows Functions Used:**

- `CreateFile` - Creates or opens a file for writing.
- `WriteFile` - Writes the specified data into the file.

---

### **4. Renaming a File**

#### **Renaming a File**

**Windows Functions Used:**

- `MoveFile` - Renames a file or moves it to a different directory.

---

### **5. Hiding a File**

#### **Hiding a File**

**Windows Functions Used:**

- `SetFileAttributes` - Sets file attributes to make the file hidden by using the `FILE_ATTRIBUTE_HIDDEN` flag.

---

### **6. Creating a Hidden Folder**

#### **Creating a Hidden Folder**

**Windows Functions Used:**

- `CreateDirectory` - Creates a directory for storing files.
- `SetFileAttributes` - Hides the directory using the `FILE_ATTRIBUTE_HIDDEN` attribute.

---

### **7. Opening a Process**

#### **Opening a Process for Manipulation**

**Windows Functions Used:**

- `OpenProcess` - Opens an existing process with specified access rights for manipulation.

---

### **8. Writing to Another Process's Memory**

#### **Writing to Another Process's Memory**

**Windows Functions Used:**

- `OpenProcess` - Opens the target process with required access rights.
- `VirtualAllocEx` - Allocates memory in the target process.
- `WriteProcessMemory` - Writes the data (shellcode or payload) into the allocated memory.

---

### **9. Terminating a Process**

#### **Terminating a Process**

**Windows Functions Used:**

- `OpenProcess` - Opens the process that needs to be terminated.
- `TerminateProcess` - Terminates the specified process.

---

### **10. Changing File Permissions**

#### **Changing File Permissions**

**Windows Functions Used:**

- `SetFileAttributes` - Changes file permissions, such as setting it to `READONLY`.

---

### **11. Creating a Thread**

#### **Creating a Thread in the Current Process**

**Windows Functions Used:**

- `CreateThread` - Creates a new thread in the current process to execute specific code or shellcode.

---

### **12. Redirecting Input/Output**

#### **Redirecting Input/Output of a Process**

**Windows Functions Used:**

- `CreatePipe` - Creates a pipe for redirecting the standard input and output of a process.
- `SetStdHandle` - Redirects the standard input/output to the created pipe.

---

### **13. Finding Files**

#### **Finding Files in a Directory**

**Windows Functions Used:**

- `FindFirstFile` - Starts the search for files in a directory.
- `FindNextFile` - Continues searching for additional files matching the pattern.
- `FindClose` - Closes the search handle when done.

---

### **14. Modifying File Attributes**

#### **Modifying File Attributes**

**Windows Functions Used:**

- `SetFileAttributes` - Modifies the attributes of the specified file, such as making it hidden or read-only.

---

### **15. Message Box for User Interaction**

#### **Displaying a Message Box**

**Windows Functions Used:**

- `MessageBox` - Displays a message box with text to the user, useful for alerting or tricking the user into interaction.

---

### **16. Network Communication via Sockets**

#### **Network Communication Using Sockets**

**Windows Functions Used:**

- `socket` - Creates a socket for network communication.
- `connect` - Establishes a connection to a remote server.
- `send` - Sends data to the remote server.
- `recv` - Receives data from the remote server.

---

### **17. Executing a Command**

#### **Executing a Command via Shell**

**Windows Functions Used:**

- `WinExec` - Executes a command in a new process, typically used to run shell commands.
- `CreateProcess` - Creates a process and executes a command or shell script.

---

### **18. File Copying**

#### **Copying a File to Another Location**

**Windows Functions Used:**

- `CopyFile` - Copies a file from a source path to a destination path.

---

### **19. Hiding Process Window**

#### **Hiding the Window of a Running Process**

**Windows Functions Used:**

- `ShowWindow` - Changes the window state of the process to hide its window.
- `FindWindow` - Finds the window by its title or class name to hide it.

---

### **20. Changing System Time**

#### **Changing the System Time**

**Windows Functions Used:**

- `SetSystemTime` - Sets the system time to a new value.
- `GetSystemTime` - Retrieves the current system time for comparison.

---

### **21. Launching a Payload from Memory**

#### **Launching a Payload from Memory**

**Windows Functions Used:**

- `VirtualAlloc` - Allocates memory in the current process for payload storage.
- `WriteProcessMemory` - Writes the payload code into the allocated memory.
- `CreateRemoteThread` - Executes the payload by creating a new thread in the process.

---

### **22. Changing Console Window Title**

#### **Changing the Title of the Console Window**

**Windows Functions Used:**

- `SetConsoleTitle` - Changes the console window title to a desired string.

---

### **23. Running Processes Silently**

#### **Running Processes Without a Visible Window**

**Windows Functions Used:**

- `CreateProcess` - Creates a process with the `CREATE_NO_WINDOW` flag to prevent the window from showing.

---

### **24. Setting Environment Variables**

#### **Setting Environment Variables for Persistence**

**Windows Functions Used:**

- `SetEnvironmentVariable` - Sets an environment variable that can be used to store configuration or paths.

---

### **25. Storing Data in the Registry**

#### **Storing Data in the Windows Registry for Persistence**

**Windows Functions Used:**

- `RegOpenKeyEx` - Opens a registry key for writing.
- `RegSetValueEx` - Sets the value for a registry key to store data (e.g., executable path for persistence).
- `RegCloseKey` - Closes the registry key after manipulation.

---
