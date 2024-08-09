package taskbar

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	procFindWindow   = user32.NewProc("FindWindowW")
	procShowWindow   = user32.NewProc("ShowWindow")
	procSetWindowPos = user32.NewProc("SetWindowPos")
)

const (
	SW_HIDE        = 0
	SW_SHOW        = 5
	HWND_BOTTOM    = 1
	SWP_NOMOVE     = 0x0002
	SWP_NOSIZE     = 0x0001
	SWP_NOACTIVATE = 0x0010
	SWP_HIDEWINDOW = 0x0080
	SWP_SHOWWINDOW = 0x0040
)

func HideTaskBar() error {
	out, err := syscall.UTF16PtrFromString("Shell_TrayWnd")
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	taskbarHandle, _, _ := procFindWindow.Call(
		uintptr(unsafe.Pointer(out)),
		uintptr(unsafe.Pointer(nil)),
	)

	procShowWindow.Call(taskbarHandle, SW_HIDE)
	procSetWindowPos.Call(taskbarHandle, HWND_BOTTOM, 0, 0, 0, 0, SWP_HIDEWINDOW|SWP_NOMOVE|SWP_NOSIZE|SWP_NOACTIVATE)

	return nil
}

func ShowTaskBar() error {
	out, err := syscall.UTF16PtrFromString("Shell_TrayWnd")
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	taskbarHandle, _, _ := procFindWindow.Call(
		uintptr(unsafe.Pointer(out)),
		uintptr(unsafe.Pointer(nil)),
	)

	procShowWindow.Call(taskbarHandle, SW_SHOW)
	procSetWindowPos.Call(taskbarHandle, HWND_BOTTOM, 0, 0, 0, 0, SWP_SHOWWINDOW|SWP_NOMOVE|SWP_NOSIZE|SWP_NOACTIVATE)

	return nil

}
