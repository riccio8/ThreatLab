/*
 * Copyright 2023-2024 Riccardo Adami. All rights reserved.
 * License: https://github.com/riccio8/ThreatLab/blob/main/LICENSE
 */


package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

// for access privileges
const (
	PAGE_NOACCESS          = 0x00000001
	PAGE_READONLY          = 0x00000002
	PAGE_READWRITE         = 0x00000004
	PAGE_WRITECOPY         = 0x00000008
	PAGE_EXECUTE           = 0x00000010
	PAGE_EXECUTE_READ      = 0x00000020
	PAGE_EXECUTE_READWRITE = 0x00000040
	PAGE_EXECUTE_WRITECOPY = 0x00000080
	PAGE_GUARD             = 0x00000100
	PAGE_NOCACHE           = 0x00000200
	PAGE_WRITECOMBINE      = 0x00000400
	PAGE_TARGETS_INVALID   = 0x40000000
	PAGE_TARGETS_NO_UPDATE = 0x40000000
)

var (
	kernel32              = syscall.NewLazyDLL("kernel32.dll")
	procAllocConsole      = kernel32.NewProc("AllocConsole")
	procSetConsoleTitleW  = kernel32.NewProc("SetConsoleTitleW")
	isDebuggerPresent     = kernel32.NewProc("IsDebuggerPresent")
	RemoteDebuggerPresent = kernel32.NewProc("CheckRemoteDebuggerPresent")
	getTickCount          = kernel32.NewProc("GetTickCount")
	GetTickCount64        = kernel32.NewProc("GetTickCount64")
	debugBreak            = kernel32.NewProc("DebugBreak")
	VirtualProtect        = kernel32.NewProc("VirtualProtect")
	GetLastError          = kernel32.NewProc("GetLastError")

	ntdll                   = syscall.NewLazyDLL("ntdll.dll")
	dbgUiRemoteBreakin      = ntdll.NewProc("DbgUiRemoteBreakin")
	DebugActiveProcess      = ntdll.NewProc("NtDebugActiveProcess")
	QueryInformationProcess = ntdll.NewProc("NtQueryInformationProcess")
	SetInformationProcess   = ntdll.NewProc("NtSetInformationProcess")
	QueryObject             = ntdll.NewProc("NtQueryObject")
	rtlAdjustPrivilege      = ntdll.NewProc("RtlAdjustPrivilege")
	RaiseHardError          = ntdll.NewProc("NtRaiseHardError")
)

// ANSI color codes
const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Reset  = "\033[0m"
)

func getLastError() error {
	ret, _, _ := GetLastError.Call()
	if ret != 0 {
		return fmt.Errorf("GetLastError: %d", ret)
	}
	return nil
}

func FindPidByNamePowerShell(processName string) ([]int, error) {
	powershellCommand := fmt.Sprintf(`Get-Process -Name %s | Where-Object { $_.Parent -eq $null } | Select-Object -ExpandProperty Id`, processName)

	cmd := exec.Command("powershell", "-Command", powershellCommand)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%serror executing PowerShell command: %v%s", Red, err, Reset)
	}

	output := strings.TrimSpace(out.String())
	if output == "" {
		return nil, fmt.Errorf("%sno process found with the name: %s%s", Red, processName, Reset)
	}

	lines := strings.Split(output, "\n")
	var pids []int
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if pid, err := strconv.Atoi(line); err == nil {
			pids = append(pids, pid)
		}
	}

	if len(pids) == 0 {
		return nil, fmt.Errorf("%s[ERROR] no process found with the name: %s%s", Red, processName, Reset)
	}

	return pids, nil
}

func isDebuggerActive() bool {
	result, _, err := isDebuggerPresent.Call()

	if err != nil {
		log.Print(err)
		return false
	}

	if result != 0 {
		log.Println("Debugger is active")
		return true

		log.Println("No debugger detected")
	}
	return result != 0
}

func checkRemoteDebuggerPresent() bool {
	fmt.Println("\n[+] Hooked CheckRemoteDebuggerPresent")
	ret, _, _ := RemoteDebuggerPresent.Call()
	return ret != 0
}

func GetTickCount() bool {
	fmt.Println("\n[+] Hooked GetTickCount")
	ret, _, _ := getTickCount.Call()
	return ret != 0
}

func DbgBreak() bool {
	fmt.Println("\n[+] Attempting to trigger DebugBreak")
	if isDebuggerActive() {
		fmt.Println("[!] Debugger detected, but skipping DebugBreak to avoid crash.")
		ret, _, _ := debugBreak.Call()
		return ret != 0
	}
	return false
}

func GetTicketCount64() bool {
	fmt.Println("\n[+] Hooked GetTickCount64")
	ret, _, _ := GetTickCount64.Call()
	return ret != 0
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func DbgUiRemoteBreakin() bool {
	fmt.Println("\n[+] Hooked DbgUiRemoteBreakin")

	resultChan := make(chan bool)
	go func() {
		ret, _, _ := dbgUiRemoteBreakin.Call()
		resultChan <- ret != 0
	}()

	select {
	case result := <-resultChan:
	
		return result
	case <-time.After(10 * time.Second):
	
		fmt.Println("[-] Timeout: moving to the next function.")
		return false
	}
}

func debugActiveProcess() bool {
	
	resultChan := make(chan bool)
	go func() {
		ret, _, _ := dbgUiRemoteBreakin.Call()
		resultChan <- ret != 0
	}()

	select {
	case result := <-resultChan:
	
		return result
	case <-time.After(10 * time.Second):
	
		fmt.Println("[-] Timeout: moving to the next function.")
		return false
	}
}

func queryInformationProcess() bool {
	fmt.Println("\n[+] Hooked QueryInformationProcess")
	ret, _, _ := QueryInformationProcess.Call()
	return ret != 0
}

func setInformationProcess() bool {
	fmt.Println("\n[+] Hooked SetInformationProcess")
	ret, _, _ := SetInformationProcess.Call()
	return ret != 0
}

func queryObject() bool {
	fmt.Println("\n[+] Hooked QueryObject")
	ret, _, _ := QueryObject.Call()
	return ret != 0
}

func AdjustPrivilege() bool {
	fmt.Println("\n[+] Hooked AdjustPrivilege")
	ret, _, _ := rtlAdjustPrivilege.Call()
	return ret != 0
}

func raiseHardError() bool {
	fmt.Println("\n[+] Hooked raiseHardError")

	var response uint32
	errorStatus := uintptr(0xC0000001) // or could be also 0xC000021A – STATUS_SYSTEM_PROCESS_TERMINATED, which typically raises a critical system error. That's 0xC0000001 – STATUS_UNSUCCESSFUL, a generic error code.
	numberOfParameters := uintptr(0)
	unicodeStringParameterMask := uintptr(0)
	parameters := uintptr(unsafe.Pointer(nil))
	validResponseOptions := uintptr(6)

	ret, _, _ := RaiseHardError.Call(
		errorStatus,
		numberOfParameters,
		unicodeStringParameterMask,
		parameters,
		validResponseOptions,
		uintptr(unsafe.Pointer(&response)),
	)

	if ret == 0 {
		fmt.Printf("%s[INFO] NtRaiseHardError executed successfully, Response: %d%s\n", Green, response, Reset)
	} else {
		fmt.Printf("%s[ERROR] NtRaiseHardError failed, return code: %x%s\n", Red, ret, Reset)
	}

	return ret == 0
}

func patch() {

	patch := []byte{0x90, 0x90, 0x90, 0x90, 0x90, 0x90}
	address := uintptr(isDebuggerPresent.Addr())

	fmt.Printf("[INFO] IsDebuggerPresent address: %x\n", address)

	var oldProtect uintptr

	ret, _, err := VirtualProtect.Call(address, uintptr(len(patch)), PAGE_EXECUTE_READWRITE, uintptr(unsafe.Pointer(&oldProtect)))
	if ret == 0 {
		fmt.Println("[ERROR] VirtualProtect failed:", err)
		fmt.Println(getLastError())
		return
	}

	copy((*[6]byte)(unsafe.Pointer(address))[:], patch)

	ret, _, err = VirtualProtect.Call(address, uintptr(len(patch)), oldProtect, uintptr(unsafe.Pointer(&oldProtect)))
	if ret == 0 {
		fmt.Println("[ERROR] Failed to restore memory protection:", err)
		fmt.Println(getLastError())
		return
	}

	fmt.Println("[INFO] IsDebuggerPresent successfully patched at:", address)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: antidbg <process name>")
		return
	}

	fmt.Println("Allocating console...")
	procAllocConsole.Call()

	title := "AntiDebuggi"
	titleUTF16 := syscall.StringToUTF16(title)
	procSetConsoleTitleW.Call(uintptr(unsafe.Pointer(&titleUTF16[0])))

	pids, err := FindPidByNamePowerShell(os.Args[1])
	if err != nil {
		fmt.Printf("%sError finding process: %v%s", Red, err, Reset)
		return
	}

	fmt.Printf("%s[INFO] Found Process IDs: %v%s\n", Green, pids, Reset)

	for _, pid := range pids {
		fmt.Println(Blue + "[X] Attempting to patch process with PID:\t " + fmt.Sprint(pid) + Reset)
		active := isDebuggerActive()
		if active {
			fmt.Println("[+] Hooked IsDebuggerPresent")
			patch()
		} else {
			fmt.Println("[-] IsDebuggerPresent function not found")
		}

		if checkRemoteDebuggerPresent() {
			patch()
		} else {
			fmt.Println(Red + "[!] Debugger not detected! Skipping patch for this control." + Reset)
		}

		if GetTickCount() {
			patch()
		} else {
			fmt.Println(Red + "[!] Debugger not detected! Skipping patch for this control." + Reset)
		}

		if GetTicketCount64() {
			patch()
		} else {
			fmt.Println(Red + "[!] Debugger not detected! Skipping patch for this control." + Reset)
		}

		if DbgUiRemoteBreakin() {
			patch()
		} else {
			fmt.Println(Red + "[!] Debugger not detected! Skipping patch for this control." + Reset)
		}

		if debugActiveProcess() {
			patch()
		} else {
			fmt.Println(Red + "[!] Debugger not detected! Skipping patch for this control." + Reset)
		}

		if queryInformationProcess() {
			patch()
		} else {
			fmt.Println(Red + "[!] Debugger not detected! Skipping patch for this control." + Reset)
		}

		if setInformationProcess() {
			patch()
		} else {
			fmt.Println(Red + "[!] Debugger not detected! Skipping patch for this control." + Reset)
		}

		if queryObject() {
			patch()
		} else {
			fmt.Println(Red + "[!] Debugger not detected! Skipping patch for this control." + Reset)
		}

		if AdjustPrivilege() {
			patch()
		} else {
			fmt.Println(Red + "[!] Debugger not detected! Skipping patch for this control." + Reset)
		}

		if raiseHardError() {
			patch()
		} else {
			fmt.Println(Red + "[!] Debugger not detected! Skipping patch for this control." + Reset)
		}
	}

	fmt.Println("\nFinished executing.")
}
