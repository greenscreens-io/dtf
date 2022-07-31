//go:build darwin
// +build darwin

package winapi

// _ShellExecuteAndWait is version of ShellExecuteEx which want process
func ShellExecuteAndWait(hwnd uintptr, lpOperation, lpFile, lpParameters, lpDirectory string, nShowCmd int) error {
	return nil
}

// _ShellExecuteNoWait is version of ShellExecuteEx which don't want process
func ShellExecuteNowait(hwnd uintptr, lpOperation, lpFile, lpParameters, lpDirectory string, nShowCmd int) error {
	return nil
}

// ShellExecuteEx is Windows API
func ShellExecuteEx(pExecInfo uintptr) error {
	return nil
}
