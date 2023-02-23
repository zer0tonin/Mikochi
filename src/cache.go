package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var fileCache map[string]os.DirEntry

// resets the cache
func resetCache() {
	fileCache = map[string]os.DirEntry{}
	cacheFolder(fileCache, "/")
}

// recursively initalizes the cache
func cacheFolder(cache map[string]os.DirEntry, path string) {
	dirEntries, err := os.ReadDir(getAbsolutePath(path))
	if err != nil {
		panic(err)
	}

	for _, dirEntry := range dirEntries {
		relativePath := filepath.Clean(path + dirEntry.Name())
		cache[relativePath] = dirEntry
		if dirEntry.IsDir() {
			cacheFolder(cache, relativePath + "/")
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
				if (event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename)) {
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

    err = watcher.Add(viper.GetString("dataDir"))
    if err != nil {
		panic(err)
    }

	// Block goroutine
    <-make(chan struct{})
}

