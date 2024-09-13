package lib

import (
	"net"
	"strconv"
	"fmt"
	"time"
)
type ScanResult struct{
	Port int
	State string
}

fmt.Println("[INFO]Library in progress...")
func scanPort(protocol, hostname string, port int) ScanResult{
	result = ScanResult{Port: port}
	address := hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout(protocol, address, 60*time.Second)

	if err != nil{
		result.State = "Cloesd"
		return result
	}

	defer conn.Close()
	result.State = "Open"
	return result

	fmt.Println(Port, "\n", port, "\n", result)
}

func InintialScan(hostname string) []ScanResult {
	var results []ScanResult

	for i := 0, i <= 1024; 1++{
		results = append (results, Scanport("tcp", hostname, i))
	}

	return results

}
