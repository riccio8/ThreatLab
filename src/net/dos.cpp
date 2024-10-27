#ifndef UNICODE
#define UNICODE 1
#endif

#include <iostream>
#include <windows.h>
#include <winsock2.h>
#include <ws2tcpip.h>
#include <stdio.h>

#pragma comment(lib, "ws2_32.lib")  // Required for linking with Ws2_32.lib

// #define SERVER "localhost"  
// #define BUFFER_SIZE 1024  
#define PORT 8080
#define SOCK_STREAM 1
#define AF_INET 2
#define IPPROTO_TCP 6

#define WIN32_LEAN_AND_MEAN


int main()
{

    WORD wVersionRequested;
    WSADATA wsaData;
    int err;

/* Use the MAKEWORD(lowbyte, highbyte) macro declared in Windef.h */
    wVersionRequested = MAKEWORD(2, 2);

    err = WSAStartup(wVersionRequested, &wsaData);
    if (err != 0) {
        /* Tell the user that we could not find a usable */
        /* Winsock DLL.                                  */
        std::cerr << "WSAStartup failed with error: \n" << err << std::endl;
        return 1;
    }

    if (LOBYTE(wsaData.wVersion) != 2 || HIBYTE(wsaData.wVersion) != 2) {
        /* Tell the user that we could not find a usable */
        /* WinSock DLL.                                  */
        std::cerr << "Could not find a usable version of Winsock.dll" << std::endl;
        WSACleanup();
        return 1;
    }
    else
        std::cout << "The Winsock 2.2 dll was found okay" << std::endl;
        
        
        
    int iResult = 0;

//    int i = 1;

    SOCKET sock = INVALID_SOCKET;
    int iFamily = AF_UNSPEC;
    int iType = 0;
    int iProtocol = 0;

// The Winsock DLL is acceptable. Proceed to use it. 
    
    sock = socket(iFamily, iType, iProtocol);
    if (sock == INVALID_SOCKET){
        std::cerr << L"socket function failed with error \n" <<  WSAGetLastError() << std::endl;
    }
    else {
        std::cout << L"socket function succeeded\n" << std::endl;
    }
    
    
// SHUTTING DOWN CONNECTION AND CLOSING SOCKET
    
    iResult = shutdown(sock, SD_SEND);
    if (iResult == SOCKET_ERROR) {
        std::cerr << L"shutdown failed with error \n" << WSAGetLastError() << std::endl;
        closesocket(sock);
        WSACleanup();
        return 1;
    }

    closesocket(sock);
    
    std::cout << L"shutdown and closesocket function passed\n" << std::endl;
    
    
        
/*Then call WSACleanup when done using the Winsock dll */
    
    WSACleanup();
    
    return 0;

}

/*  THIS IS JUST THE FIRST PART, THE INITIAL SETUP, I'LL ADD THREAD AND SO ON, JUST WAIT */

// Compile using: g++ -Wall -o sock .\dos.cpp -lws2_32
