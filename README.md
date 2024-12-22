---

# Offensive & Defensive Tools Repository

![Golang](https://img.shields.io/badge/Golang-Tools-00ADD8?style=flat&logo=go) [![CodeQL Advanced](https://github.com/riccio8/Offensive-defensive-tools/actions/workflows/codeql.yml/badge.svg?branch=main)](https://github.com/riccio8/Offensive-defensive-tools/actions/workflows/codeql.yml) ![Build Status](https://img.shields.io/badge/build-passing-brightgreen) ![Version](https://img.shields.io/badge/version-1.0.0-blue)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

## **Overview**

This repository includes a collection of tools, resources and libraries designed for offensive and defensive security operations, covering areas such as process and memory analysis, network security, vulnerability detection, and exploit development. It provides a versatile set of resources for researchers, cybersecurity professionals, and developers.

## **Directory Structure**

- **`lib/`**: Core libraries and DLLs for various security-related tasks.  
  - Includes C++, Python, Go, and Assembly-based libraries for memory, network, and sandboxing functionalities.
  
- **`shellcodes/`**: Pre-built shellcode examples for penetration testing and exploit development.  
  - Contains scripts for reverse shells, directory hiding, and more.

- **`src/`**: Main source files for various projects.  
  - **`Mem/`**: Tools for memory analysis, hex viewers, ROP chain detection, and system calls management.  
  - **`net/`**: Network-related utilities focusing on DOS/DDOS testing and network resilience.  
  - **`processes/`**: Process monitoring and manipulation tools, including privilege escalation scripts and anti-debugging mechanisms.
  -  **`resources/`**: That's a directory containing documents, videos and some other resources than may help u understanding deeply some args, for both beginners and advanced
  - **`rules/`**: YARA rules for malware detection, ransomware identification, keylogging detection, and other threat signatures.

## **Key Features**

- **Process & Memory Analysis**:  
  Identify and analyze running processes, heap memory, and vulnerabilities in real-time.

- **Network Security Tools**:  
  Utilities designed to simulate and analyze various network threats like DOS and DDOS attacks.

- **Sandboxing**:  
  Isolated environments for malware testing and secure code execution.

- **Cross-Platform Capabilities**:  
  Primarily focused on Windows, with some support for Linux environments.

## **Technologies Used**

- **Assembly**: For low-level system manipulation.  
- **C++**: High-performance process and memory handling tools.  
- **Golang**: Lightweight and efficient network and process utilities.  
- **Python**: Scripting for process automation and vulnerability detection.  
- **YARA**: Advanced threat detection through rule-based signatures.

## **Disclaimer**

⚠️ **Disclaimer**: All tools in this repository are intended for **educational and legal security research purposes only**. Misuse of these tools for malicious purposes is strictly prohibited. The author takes no responsibility for any illegal use.

---

*Note: Continuous improvement is ongoing, especially for assembly-related tools!*
