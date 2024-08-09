package cache

import (
	"sync"
	"time"
)

var (
	Cache       []string
	CacheMutex  sync.Mutex
	CacheCond          = sync.NewCond(&CacheMutex)
	Key         []byte = []byte("zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz")
	FilesFound         = 0
	FilesError         = 0
	GlobalStart        = false
)

func MonitorCache(done chan struct{}) {
	for {
		time.Sleep(1 * time.Second)
		CacheMutex.Lock()
		if len(Cache) == 0 {
			close(done)
			CacheMutex.Unlock()
			return
		}
		CacheMutex.Unlock()
	}
}
