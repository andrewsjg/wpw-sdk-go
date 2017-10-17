// +build !windows

package apputils

// Check if process is alive.

import (
	"os"
	"syscall"
)

func IsProcAlive(pid int) (bool, error) {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false, err
	} else {
		err := process.Signal(syscall.Signal(0))
		if err == nil {
			return true, nil
		}
		return false, err
	}
}
