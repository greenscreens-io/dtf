//go:build linux
// +build linux

package winapi

func ShellExecuteAndWait(hwnd uintptr, lpOperation, lpFile, lpParameters, lpDirectory string, nShowCmd int) error {
	return nil
}

func ShellExecuteNowait(hwnd uintptr, lpOperation, lpFile, lpParameters, lpDirectory string, nShowCmd int) error {
	return nil
}

func ShellExecuteEx(pExecInfo uintptr) error {
	return nil
}
