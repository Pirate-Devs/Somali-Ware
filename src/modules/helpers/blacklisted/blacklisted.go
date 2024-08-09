package blacklisted

import "os"

var Blacklisted_file_types = []string{
	"exe",
	"com",
	"dll",
	"cpl",
	"somalia",
	"ini",
	"cfg",
	"reg",
	"sys",
	"drv",
}

var Blacklisted_file = "C:\\Users\\" + os.Getenv("USERNAME") + "\\Desktop\\DO_NOT_DELETE.html"
