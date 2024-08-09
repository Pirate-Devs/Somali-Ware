package structs

type IPInfo struct {
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	CountryCode string  `json:"countryCode"`
}

type SystemInfo struct {
	Username     string `json:"hostname"`
	HWID         string `json:"hwid"`
	Date         string `json:"date"`
	TimeZone     string `json:"timezone"`
	LatLong      string `json:"location"`
	FilesFound   string `json:"files_found"`
	CountryCode  string `json:"country_code"`
	EncryptedKey string `json:"encrypted_key"`
}
