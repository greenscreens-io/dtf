//go:build windows
// +build windows

package winapi

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	modshell32         = syscall.NewLazyDLL("shell32.dll")
	procShellExecuteEx = modshell32.NewProc("ShellExecuteExW")
)

// _ShellExecuteAndWait is version of ShellExecuteEx which want process
func ShellExecuteAndWait(hwnd hwnd, lpOperation, lpFile, lpParameters, lpDirectory string, nShowCmd int) error {
	var lpctstrVerb, lpctstrParameters, lpctstrDirectory, lpctstrFile lpctstr

	lpOp, err := syscall.UTF16PtrFromString(lpOperation)
	if err != nil {
		return err
	}
	if len(lpOperation) != 0 {
		lpctstrVerb = lpctstr(unsafe.Pointer(lpOp))
	}

	lpPp, err := syscall.UTF16PtrFromString(lpParameters)
	if err != nil {
		return err
	}
	if len(lpParameters) != 0 {
		lpctstrParameters = lpctstr(unsafe.Pointer(lpPp))
	}

	lpDp, err := syscall.UTF16PtrFromString(lpDirectory)
	if err != nil {
		return err
	}
	if len(lpDirectory) != 0 {
		lpctstrDirectory = lpctstr(unsafe.Pointer(lpDp))
	}

	lpFp, err := syscall.UTF16PtrFromString(lpFile)
	if err != nil {
		return err
	}
	if len(lpFile) != 0 {
		lpctstrFile = lpctstr(unsafe.Pointer(lpFp))
	}

	i := &_SHELLEXECUTEINFO{
		fMask:        _SEE_MASK_NOCLOSEPROCESS,
		hwnd:         hwnd,
		lpVerb:       lpctstrVerb,
		lpFile:       lpctstrFile,
		lpParameters: lpctstrParameters,
		lpDirectory:  lpctstrDirectory,
		nShow:        nShowCmd,
	}
	i.cbSize = dword(unsafe.Sizeof(*i))
	return ShellExecuteEx(i)
}

// _ShellExecuteNoWait is version of ShellExecuteEx which don't want process
func ShellExecuteNowait(hwnd hwnd, lpOperation, lpFile, lpParameters, lpDirectory string, nShowCmd int) error {
	var lpctstrVerb, lpctstrParameters, lpctstrDirectory, lpctstrFile lpctstr

	lpOp, err := syscall.UTF16PtrFromString(lpOperation)
	if err != nil {
		return err
	}
	if len(lpOperation) != 0 {
		lpctstrVerb = lpctstr(unsafe.Pointer(lpOp))
	}

	lpPp, err := syscall.UTF16PtrFromString(lpParameters)
	if err != nil {
		return err
	}
	if len(lpParameters) != 0 {
		lpctstrParameters = lpctstr(unsafe.Pointer(lpPp))
	}

	lpDp, err := syscall.UTF16PtrFromString(lpDirectory)
	if err != nil {
		return err
	}
	if len(lpDirectory) != 0 {
		lpctstrDirectory = lpctstr(unsafe.Pointer(lpDp))
	}

	lpFp, err := syscall.UTF16PtrFromString(lpFile)
	if err != nil {
		return err
	}
	if len(lpFile) != 0 {
		lpctstrFile = lpctstr(unsafe.Pointer(lpFp))
	}

	i := &_SHELLEXECUTEINFO{
		fMask:        _SEE_MASK_DEFAULT,
		hwnd:         hwnd,
		lpVerb:       lpctstrVerb,
		lpFile:       lpctstrFile,
		lpParameters: lpctstrParameters,
		lpDirectory:  lpctstrDirectory,
		nShow:        nShowCmd,
	}
	i.cbSize = dword(unsafe.Sizeof(*i))
	return ShellExecuteEx(i)
}

// ShellExecuteEx is Windows API
func ShellExecuteEx(pExecInfo *_SHELLEXECUTEINFO) error {
	ret, _, _ := procShellExecuteEx.Call(uintptr(unsafe.Pointer(pExecInfo)))
	if ret == 1 && pExecInfo.fMask&_SEE_MASK_NOCLOSEPROCESS != 0 {
		s, e := syscall.WaitForSingleObject(syscall.Handle(pExecInfo.hProcess), syscall.INFINITE)
		switch s {
		case syscall.WAIT_OBJECT_0:
			break
		case syscall.WAIT_FAILED:
			return os.NewSyscallError("WaitForSingleObject", e)
		default:
			return errors.New("unexpected result from WaitForSingleObject")
		}
	}
	errorMsg := ""
	if pExecInfo.hInstApp != 0 && pExecInfo.hInstApp <= 32 {
		switch int(pExecInfo.hInstApp) {
		case _SE_ERR_FNF:
			errorMsg = "The specified file was not found"
		case _SE_ERR_PNF:
			errorMsg = "The specified path was not found"
		case _ERROR_BAD_FORMAT:
			errorMsg = "The .exe file is invalid (non-Win32 .exe or error in .exe image)"
		case _SE_ERR_ACCESSDENIED:
			errorMsg = "The operating system denied access to the specified file"
		case _SE_ERR_ASSOCINCOMPLETE:
			errorMsg = "The file name association is incomplete or invalid"
		case _SE_ERR_DDEBUSY:
			errorMsg = "The DDE transaction could not be completed because other DDE transactions were being processed"
		case _SE_ERR_DDEFAIL:
			errorMsg = "The DDE transaction failed"
		case _SE_ERR_DDETIMEOUT:
			errorMsg = "The DDE transaction could not be completed because the request timed out"
		case _SE_ERR_DLLNOTFOUND:
			errorMsg = "The specified DLL was not found"
		case _SE_ERR_NOASSOC:
			errorMsg = "There is no application associated with the given file name extension"
		case _SE_ERR_OOM:
			errorMsg = "There was not enough memory to complete the operation"
		case _SE_ERR_SHARE:
			errorMsg = "A sharing violation occurred"
		default:
			errorMsg = fmt.Sprintf("Unknown error occurred with error code %v", pExecInfo.hInstApp)
		}
	} else {
		return nil
	}
	return errors.New(errorMsg)
}
