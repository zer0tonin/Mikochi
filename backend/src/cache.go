package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// fileCache is a global cache of FileInfos from the watched data directory
var fileCache map[string]fs.FileInfo

// resets the cache
func resetCache() {
	fileCache = map[string]fs.FileInfo{}
	cacheFolder(fileCache, "/")
}

// recursively initalizes the cache
// TODO: lock
func cacheFolder(cache map[string]fs.FileInfo, path string) {
	dirEntries, err := os.ReadDir(getAbsolutePath(path))
	if err != nil {
		panic(err)
	}

	for _, dirEntry := range dirEntries {
		relativePath := filepath.Clean(path + dirEntry.Name())
		fileInfo, err := dirEntry.Info()
		if err != nil {
			panic(err)
		}
		cache[relativePath] = fileInfo
		if dirEntry.IsDir() {
			cacheFolder(cache, relativePath+"/")
		}
	}
}

// refreshes the cache on data dir changes
func watchDataDir() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
					fmt.Println("event:", event)
					resetCache()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(dataDir)
	if err != nil {
		panic(err)
	}

	// Block goroutine
	<-make(chan struct{})
}
