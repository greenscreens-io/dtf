//go:build windows
// +build windows

package winapi

import (
	"errors"
	"syscall"
	"unsafe"
)

//Set windows console title
func SetConsoleTitle(title string) (int, error) {

	if len(title) == 0 {
		return 0, errors.New("title required")
	}

	handle, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return 0, err
	}
	defer syscall.FreeLibrary(handle)

	proc, err := syscall.GetProcAddress(handle, "SetConsoleTitleW")
	if err != nil {
		return 0, err
	}

	lpFp, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return 0, err
	}

	r, _, err := syscall.Syscall(proc, 1, uintptr(unsafe.Pointer(lpFp)), 0, 0)
	return int(r), err
}
