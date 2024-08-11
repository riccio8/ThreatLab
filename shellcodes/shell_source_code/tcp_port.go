package main

import (
 "os/exec"
 "syscall"
)

// wil work only for windows, if u want to do the same things on linux u have to delete the sys call

func main() {
 cmd := exec.Command("cmd.exe", "/c", "netsh advfirewall firewall add rule name=\"Apri Porta 8080\" dir=in action=allow protocol=TCP localport=8080")
 cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
 cmd.Run()
}

// it can be run only in administration mode
