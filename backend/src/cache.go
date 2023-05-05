package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/rjeczalik/notify"
)

// fileCache is a global cache of FileInfos from the watched data directory
var fileCache map[string]fs.FileInfo
var fileCacheMutex sync.Mutex

// resets the cache
func resetCache() {
	fileCache = map[string]fs.FileInfo{}
	cacheFolder(fileCache, "/")
}

// recursively initalizes the cache
// we do not use filepath.Walk to avoid a mess with relative / absolute paths
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

		fileCacheMutex.Lock()
		cache[relativePath] = fileInfo
		fileCacheMutex.Unlock()

		if dirEntry.IsDir() {
			cacheFolder(cache, relativePath+"/")
		}
	}
}

// refreshes the cache on data dir changes
func watchDataDir() {
	c := make(chan notify.EventInfo, 1)

	// watcg the move, create, remove, rename events
	if err := notify.Watch(dataDir + "/...", c, notify.InMovedTo, notify.Create, notify.Remove, notify.Rename); err != nil {
		panic(err)
	}
	defer notify.Stop(c)

	for ei := range c {
		log.Printf("Event %d on file %s", ei.Event(), ei.Path())
		// kinda brutal solution, could be improved
		resetCache()
	}
}
