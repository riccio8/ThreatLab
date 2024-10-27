#ifndef UNICODE
#define UNICODE 1
#endif

#include <winsock2.h>  // Include winsock2.h before windows.h
#include <windows.h>
#include <ws2tcpip.h>
#include <iostream>
#include <stdio.h>

#define PORT 8080
#define WIN32_LEAN_AND_MEAN

int main() {
    WORD wVersionRequested = MAKEWORD(2, 2);
    WSADATA wsaData;
    int err;

    err = WSAStartup(wVersionRequested, &wsaData);
    if (err != 0) {
        std::cerr << "WSAStartup failed with error: " << err << std::endl;
        return 1;
    }

    if (LOBYTE(wsaData.wVersion) != 2 || HIBYTE(wsaData.wVersion) != 2) {
        std::cerr << "Could not find a usable version of Winsock.dll" << std::endl;
        WSACleanup();
        return 1;
    } else {
        std::cout << "The Winsock 2.2 dll was found okay" << std::endl;
    }

    SOCKET sock = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
    if (sock == INVALID_SOCKET) {
        std::cerr << "Socket function failed with error: " << WSAGetLastError() << std::endl;
        WSACleanup();
        return 1;
    } else {
        std::cout << "Socket function succeeded" << std::endl;
    }
    
    sockaddr_in serverAddr;
    serverAddr.sin_family = AF_INET;
    serverAddr.sin_port = htons(80); // Port 80 
    inet_pton(AF_INET, "127.0.0.1", &serverAddr.sin_addr); 
    
    
    int connectResult = connect(sock, (sockaddr*)&serverAddr, sizeof(serverAddr));
    if (connectResult == SOCKET_ERROR) {
        std::cerr << "Error during connection: " << WSAGetLastError() << std::endl;
        std::cerr << GetLastError() << std::endl;
        closesocket(sock);
        WSACleanup();
        return 1;
    }else{
        std::cout << "Connected" << std::endl;
    }
    
    const char* message = "FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF";  
    int sendResult = send(sock, message, strlen(message), 0);
    if (sendResult == SOCKET_ERROR) {
        std::cerr << "Error sending data: \n" << WSAGetLastError() << std::endl;
        closesocket(sock);
        WSACleanup();
        return 1;
    }else{
        std::cout << "Data sent successfully" << std::endl;
    }


    
    // Shutting down connection and closing socket
    int iResult = shutdown(sock, SD_SEND);
    if (iResult == SOCKET_ERROR) {
        std::cerr << "Shutdown failed with error: " << WSAGetLastError() << std::endl;
        closesocket(sock);
        WSACleanup();
        return 1;
    } else {
        std::cout << "No shutdown requested" << std::endl;
    }

    closesocket(sock);
    std::cout << "Shutdown and closesocket function passed" << std::endl;

    WSACleanup();
    return 0;
}


/*  THIS IS JUST THE FIRST PART, THE INITIAL SETUP, I'LL ADD THREAD AND SO ON, JUST WAIT */

// Compile using: g++ -Wall -o sock .\dos.cpp -lws2_32
