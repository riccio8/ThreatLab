package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func scanPort(protocol, hostname str, port int) bool{
	address := hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout(protocol, address, 60*time.Second())

	if err != nil{
		return false
	}

	defer conn.Close()
	return true
}

func main(){
	fmt.Println("Port scanning")
	open := scanPort("tcp", "localhost", 80)
	fmt.Println("Port status: %t\n", open)
}
