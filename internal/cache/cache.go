package cache

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"sync"
	"time"
)

type Cache struct {
	accessFlag   sync.RWMutex
	data         map[string][]byte
	path         string
	lastLoadTime time.Time
	ttl          time.Duration
}

func New(ttl time.Duration) *Cache {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = os.Getenv("HOME")
	}

	return &Cache{
		data:         make(map[string][]byte),
		path:         path.Join(homeDir, ".config/ghelper/cache"),
		lastLoadTime: time.Now(),
		ttl:          ttl,
	}
}

func (c *Cache) Load() {
	fileStat, _ := os.Stat(c.path)
	isCacheFileExist := fileStat != nil
	if isCacheFileExist && time.Since(fileStat.ModTime()) < c.ttl {
		file, err := os.OpenFile(c.path, os.O_RDWR, 0755)
		if err != nil {
			panic(fmt.Errorf("can't file file: %s", err))
		}
		defer func() {
			_ = file.Close()
		}()

		rw := bufio.NewReadWriter(bufio.NewReader(file), bufio.NewWriter(file))
		data, err := io.ReadAll(rw)
		if err != nil {
			panic(fmt.Errorf("can't read data: %s", err))
		}
		err = json.Unmarshal(data, &c.data)
		if err != nil {
			fmt.Println("corrupted Cache. will be stored new one")
		}
	}
}

func (c *Cache) Store() {
	file, err := os.OpenFile(c.path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(fmt.Errorf("can't open file: %s", err))
	}
	defer func() {
		_ = file.Close()
	}()
	data, err := json.Marshal(c.data)
	if err != nil {
		panic(fmt.Errorf("can't marshal data: %s", err))
	}
	_, err = file.Write(data)
	if err != nil {
		panic(fmt.Errorf("can't write data: %s", err))
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.accessFlag.RLock()
	defer c.accessFlag.RUnlock()

	data, ok := c.data[key]

	return data, ok
}

func (c *Cache) Set(key string, value []byte) {
	c.accessFlag.Lock()
	defer c.accessFlag.Unlock()

	c.data[key] = value
}
