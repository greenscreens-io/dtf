//go:build windows
// +build windows

package winapi

import (
	"syscall"
)

const (
	_SEE_MASK_DEFAULT            = 0x00000000
	_SEE_MASK_CLASSNAME          = 0x00000001
	_SEE_MASK_CLASSKEY           = 0x00000003
	_SEE_MASK_IDLIST             = 0x00000004
	_SEE_MASK_INVOKEIDLIST       = 0x0000000C
	_SEE_MASK_ICON               = 0x00000010
	_SEE_MASK_HOTKEY             = 0x00000020
	_SEE_MASK_NOCLOSEPROCESS     = 0x00000040
	_SEE_MASK_CONNECTNETDRV      = 0x00000080
	_SEE_MASK_NOASYNC            = 0x00000100
	_SEE_MASK_FLAG_DDEWAIT       = 0x00000100
	_SEE_MASK_DOENVSUBST         = 0x00000200
	_SEE_MASK_FLAG_NO_UI         = 0x00000400
	_SEE_MASK_UNICODE            = 0x00004000
	_SEE_MASK_NO_CONSOLE         = 0x00008000
	_SEE_MASK_ASYNCOK            = 0x00100000
	_SEE_MASK_NOQUERYCLASSSTORE  = 0x01000000
	_SEE_MASK_HMONITOR           = 0x00200000
	_SEE_MASK_NOZONECHECKS       = 0x00800000
	_SEE_MASK_WAITFORINPUTIDLE   = 0x02000000
	_SEE_MASK_FLAG_LOG_USAGE     = 0x04000000
	_SEE_MASK_FLAG_HINST_IS_SITE = 0x08000000
)

const (
	_ERROR_BAD_FORMAT = 11
)

const (
	_SE_ERR_FNF             = 2
	_SE_ERR_PNF             = 3
	_SE_ERR_ACCESSDENIED    = 5
	_SE_ERR_OOM             = 8
	_SE_ERR_DLLNOTFOUND     = 32
	_SE_ERR_SHARE           = 26
	_SE_ERR_ASSOCINCOMPLETE = 27
	_SE_ERR_DDETIMEOUT      = 28
	_SE_ERR_DDEFAIL         = 29
	_SE_ERR_DDEBUSY         = 30
	_SE_ERR_NOASSOC         = 31
)

type (
	dword     uint32
	hinstance syscall.Handle
	hkey      syscall.Handle
	hwnd      syscall.Handle
	ulong     uint32
	lpctstr   uintptr
	lpvoid    uintptr
)

// SHELLEXECUTEINFO struct
type _SHELLEXECUTEINFO struct {
	cbSize         dword
	fMask          ulong
	hwnd           hwnd
	lpVerb         lpctstr
	lpFile         lpctstr
	lpParameters   lpctstr
	lpDirectory    lpctstr
	nShow          int
	hInstApp       hinstance
	lpIDList       lpvoid
	lpClass        lpctstr
	hkeyClass      hkey
	dwHotKey       dword
	hIconOrMonitor syscall.Handle
	hProcess       syscall.Handle
}
