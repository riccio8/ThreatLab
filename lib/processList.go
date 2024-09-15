package processList

import (
	"fmt"
	"time"

	"github.com/Microsoft/go-winio/pkg/process"
)


func ListHandle() ([]uint32, error) {
	fmt.Println("[X] Ricerca dei processi in corso...")
	time.Sleep(2 * time.Second)


	procs, err := process.EnumProcesses()
	if err != nil {
		return nil, err
	}

	return procs, nil
}
