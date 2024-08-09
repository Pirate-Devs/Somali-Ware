package background

import (
	"Somali-Ware/settings"
	"encoding/base64"
	"io"
	"os"
	"syscall"
	"unsafe"
)

var (
	user32                   = syscall.NewLazyDLL("user32.dll")
	procSystemParametersInfo = user32.NewProc("SystemParametersInfoW")
)

const (
	SPI_SETDESKWALLPAPER = 0x0014
	SPIF_UPDATEINIFILE   = 0x01
	SPIF_SENDCHANGE      = 0x02
	SPI_GETDESKWALLPAPER = 0x0073
	MAX_PATH             = 260
)

func SetWallpaper() error {
	SaveWallpaper()
	wallpaper_path := "C:\\Users\\Public\\wallpaper.jpg"
	decoded, _ := base64.StdEncoding.DecodeString(settings.Wallpaper_Image)
	err := os.WriteFile("C:\\Users\\Public\\wallpaper.jpg", decoded, 0644)
	if err != nil {
		return err
	}

	lpFile, err := syscall.UTF16PtrFromString(wallpaper_path)
	if err != nil {
		return err
	}

	ret, _, err := procSystemParametersInfo.Call(
		SPI_SETDESKWALLPAPER,
		0,
		uintptr(unsafe.Pointer(lpFile)),
		SPIF_UPDATEINIFILE|SPIF_SENDCHANGE,
	)
	if ret == 0 {
		return err
	}
	return nil
}

func SaveWallpaper() error {
	var wallpaperPath [MAX_PATH]uint16
	_, _, err := syscall.NewLazyDLL("user32.dll").NewProc("SystemParametersInfoW").Call(
		SPI_GETDESKWALLPAPER,
		uintptr(MAX_PATH),
		uintptr(unsafe.Pointer(&wallpaperPath)),
		0,
	)
	if err != syscall.Errno(0) {
		return err
	}
	file_path := syscall.UTF16ToString(wallpaperPath[:])

	sourceFile, err := os.Open(file_path)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create("C:\\Users\\Public\\wallpaper_original.jpg")
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func ResetWallpaper() error {
	wallpaper_path := "C:\\Users\\Public\\wallpaper_original.jpg"

	lpFile, err := syscall.UTF16PtrFromString(wallpaper_path)
	if err != nil {
		return err
	}

	ret, _, err := procSystemParametersInfo.Call(
		SPI_SETDESKWALLPAPER,
		0,
		uintptr(unsafe.Pointer(lpFile)),
		SPIF_UPDATEINIFILE|SPIF_SENDCHANGE,
	)
	if ret == 0 {
		return err
	}
	return nil
}
