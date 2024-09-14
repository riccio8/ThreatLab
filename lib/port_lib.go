package port_lib

import (
	"log"
	"net"
	"strconv"
	"time"

)

func ScanPort(protocol, hostname string, port int) bool {
	address := hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeOut(protocol, address, 60*time.Second)

	if err != nil {
		log.Fatal("[ERROR] Port is close... ", err)
		return false 
	}

	defer conn.Close()

	return true
}
