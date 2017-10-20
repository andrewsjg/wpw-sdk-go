// +build windows

package apputils

// Check if process is alive.

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	// defined by the Win32 API
	th32cs_snapprocess uintptr = 0x2
)

var (
	kernel                   = syscall.MustLoadDLL("kernel32.dll")
	CreateToolhelp32Snapshot = kernel.MustFindProc("CreateToolhelp32Snapshot")
	Process32First           = kernel.MustFindProc("Process32FirstW")
	Process32Next            = kernel.MustFindProc("Process32NextW")
)

// ProcessEntry32 structure defined by the Win32 API
type processEntry32 struct {
	dwSize              uint32
	cntUsage            uint32
	th32ProcessID       uint32
	th32DefaultHeapID   int
	th32ModuleID        uint32
	cntThreads          uint32
	th32ParentProcessID uint32
	pcPriClassBase      int32
	dwFlags             uint32
	szExeFile           [syscall.MAX_PATH]uint16
}

func getProcessEntry(pid int) (pe *processEntry32, err error) {
	snapshot, _, e1 := CreateToolhelp32Snapshot.Call(th32cs_snapprocess, uintptr(0))
	if snapshot == uintptr(syscall.InvalidHandle) {
		err = fmt.Errorf("CreateToolhelp32Snapshot: %v", e1)
		return
	}
	defer syscall.CloseHandle(syscall.Handle(snapshot))

	var processEntry processEntry32
	processEntry.dwSize = uint32(unsafe.Sizeof(processEntry))
	ok, _, e1 := Process32First.Call(snapshot, uintptr(unsafe.Pointer(&processEntry)))
	if ok == 0 {
		err = fmt.Errorf("Process32First: %v", e1)
		return
	}

	for {
		if processEntry.th32ProcessID == uint32(pid) {
			pe = &processEntry
			return
		}

		ok, _, e1 = Process32Next.Call(snapshot, uintptr(unsafe.Pointer(&processEntry)))
		if ok == 0 {
			err = fmt.Errorf("Process32Next: %v", e1)
			return
		}
	}
}

func isProcAlive_windows(pid int) (bool, error) {
	_, err := getProcessEntry(pid)
	if err != nil {
		return false, err
	}
	return true, nil
}

func IsProcAlive(pid int) (bool, error) {
	return isProcAlive_windows(pid)
}
