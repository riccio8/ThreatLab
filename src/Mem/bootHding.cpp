/*
 * Copyright 2023-2024 Riccardo Adami. All rights reserved.
 * License: https://github.com/riccio8/ThreatLab/blob/main/LICENSE
 */


#include <iostream>
#include <windows.h>

void hideBootFolder() {
    // Hide the boot folder
    system("attrib +h +s +r C:\\boot");         // To unhide: attrib -h -s -r C:\boot
    std::cout << "'C:\\boot' folder hidden.\n";
}

void denyAccessToBootFolder() {
    // Deny access to all users
    system("icacls C:\\boot /deny Everyone:F");        // Reset permissions: icacls C:\boot /reset
    std::cout << "Access to 'C:\\boot' folder denied for everyone.\n";
}

void modifyBootConfiguration() {
    // Modify the system boot path
    system("bcdedit /set {current} path \\Windows\\System32\\winload.exe");
    std::cout << "Boot path modified.\n";
}

int main() {
    int choice;
    std::cout << "Choose an operation:\n";
    std::cout << "1. Hide 'C:\\boot' folder\n";
    std::cout << "2. Deny access to 'C:\\boot' folder\n";
    std::cout << "3. Modify boot configuration\n";
    std::cout << "Choice: ";
    std::cin >> choice;

    switch (choice) {
        case 1:
            hideBootFolder();
            break;
        case 2:
            denyAccessToBootFolder();
            break;
        case 3:
            modifyBootConfiguration();
            break;
        default:
            std::cout << "Invalid choice.\n";
            break;
    }

    return 0;
}
