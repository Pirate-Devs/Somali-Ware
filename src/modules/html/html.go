package html

import (
	"Somali-Ware/settings"
	"encoding/base64"
)

func GetHTML() string {
	decoded, err := base64.StdEncoding.DecodeString(settings.HTML_CODE)
	if err != nil {
		decoded = []byte(settings.HTML_CODE)
	}
	return string(decoded)
}
