---

### Reverse Shell in Assembly (64-bit)

This is a **64-bit reverse shell** implemented in **x86-64 assembly**, designed to create a connection to a remote server, redirect input/output to a socket, and execute `cmd.exe` on the target system. The code utilizes Windows API functions to set up the networking environment and spawn a process that connects back to the attacker machine.

### Overview of the Components

The reverse shell works by establishing a TCP connection to a remote server and creating a new process (`cmd.exe`) that communicates over the socket. This is useful for penetration testing and exploitation purposes, as it allows the attacker to execute arbitrary commands on the target machine from a remote location.

#### Key APIs Used:
1. **`LoadLibraryA`**: Dynamically loads the `ws2_32.dll` library required for networking.
2. **`WSAStartup`**: Initializes the Windows Sockets API (Winsock) for TCP/IP communication.
3. **`socket`**: Creates a socket used for communication (specifically, a TCP socket).
4. **`connect`**: Establishes a connection to a specified IP and port.
5. **`dup2`**: Duplicates the socket to `stdin`, `stdout`, and `stderr`, so all input/output goes through the socket.
6. **`CreateProcessA`**: Spawns a new process (`cmd.exe`), which will have its input/output redirected to the socket.

### How It Works:

1. **Loading the DLL**: The reverse shell loads `ws2_32.dll` using `LoadLibraryA`, which is necessary for socket communication.
2. **Setting Up Winsock**: It initializes Winsock with version `2.2` using `WSAStartup`, allowing it to create network connections.
3. **Socket Creation**: A socket is created with the `socket` function (TCP over IPv4).
4. **Connection**: The shell connects to a remote IP (default is `127.0.0.1`) on a specified port (`4444` by default). You can modify these values to target any machine or port you want.
5. **Redirection of Input/Output**: The shell redirects `stdin`, `stdout`, and `stderr` to the socket using `dup2`. This way, all interactions with `cmd.exe` will go through the network connection.
6. **Executing `cmd.exe`**: Finally, `cmd.exe` is spawned with `CreateProcessA`, and the shell is ready to accept commands over the network.

### Compilation Instructions:

1. **Assemble the code**:
   ```bash
   nasm -f win64 -o rev .\reverseShell.asm 
   ```



### Notes:
- The code assumes the IP address is set to `127.0.0.1` and the port is set to `4444`. If you want to connect to a different server, simply modify the `add rbx, 0x7F000001` and `mov bx, 0x5C11` lines to the desired IP and port.
- The shell redirects input and output to the socket using `dup2`, allowing you to interact with the remote machine through `cmd.exe`.
- This shell is designed to work on **Windows** systems, utilizing Windows Sockets API (`ws2_32.dll`).
