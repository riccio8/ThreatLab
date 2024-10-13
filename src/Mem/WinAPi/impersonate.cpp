#include <windows.h>
#include <iostream>
#include <sddl.h>

void ImpersonateUser() {
    HANDLE hToken = NULL;

    if (!OpenProcessToken(GetCurrentProcess(), TOKEN_DUPLICATE | TOKEN_QUERY, &hToken)) {
        std::cerr << "OpenProcessToken failed. Error: " << GetLastError() << std::endl;
        return;
    }

    HANDLE hImpersonatedToken = NULL;
    if (!DuplicateToken(hToken, SecurityImpersonation, &hImpersonatedToken)) {
        std::cerr << "DuplicateToken failed. Error: " << GetLastError() << std::endl;
        CloseHandle(hToken);
        return;
    }

    if (!ImpersonateLoggedOnUser(hImpersonatedToken)) {
        std::cerr << "ImpersonateLoggedOnUser failed. Error: " << GetLastError() << std::endl;
        CloseHandle(hImpersonatedToken);
        CloseHandle(hToken);
        return;
    }

    char userName[256];
    DWORD userNameSize = sizeof(userName);
    if (GetUserNameA(userName, &userNameSize)) {
        std::cout << "Impersonated user: " << userName << std::endl;
    } else {
        std::cerr << "GetUserNameA failed. Error: " << GetLastError() << std::endl;
    }

    RevertToSelf();

    CloseHandle(hImpersonatedToken);
    CloseHandle(hToken);
}

int main() {
    ImpersonateUser();
    return 0;
}

// compile: g++ -o impersonate impersonate.cpp
