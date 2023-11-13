package browser

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rjeczalik/notify"
	"github.com/spf13/viper"
)

// fileCache is a global cache of FileInfos from the watched data directory
var fileCache = map[string]fs.FileInfo{}
var fileCacheMutex = sync.Mutex{}

// resets the cache
func ResetCache() {
	log.Printf("Caching %s", viper.GetString("DATA_DIR"))
	defer func() {
		if r := recover(); r != nil {
			log.Print("Failed to refresh cache")
		}
	}()

	newFileCache := map[string]fs.FileInfo{}
	cacheFolder(newFileCache, "/")

	// just doing this at once should avoid excessive lock/unlock
	fileCacheMutex.Lock()
	fileCache = newFileCache
	fileCacheMutex.Unlock()
	log.Print("Refreshed cached")
}

// recursively initalizes the cache
// we do not use filepath.Walk to avoid a mess with relative / absolute paths
func cacheFolder(cache map[string]fs.FileInfo, path string) {
	rootPath := getAbsolutePath(path)
	err := filepath.WalkDir(
		rootPath,
		func(path string, dirEntry fs.DirEntry, err error) error {
			relativePath := strings.TrimPrefix(path, rootPath)
			if relativePath == "" {
				return nil
			}
			fileInfo, err := dirEntry.Info()
			if err != nil {
				return err
			}

			cache[relativePath] = fileInfo
			return nil
		},
	)
	if err != nil {
		log.Panicf("Error while refreshing cache: %s", err.Error())
	}
}

// refreshes the cache on data dir changes
func WatchDataDir() {
	c := make(chan notify.EventInfo, 1)

	// watcg the create, remove, rename events on the data dir and sub directories
	if err := notify.Watch(viper.GetString("DATA_DIR")+"/...", c, notify.Create, notify.Remove, notify.Rename); err != nil {
		panic(err)
	}
	defer notify.Stop(c)

	for ei := range c {
		log.Printf("Event %d on file %s", ei.Event(), ei.Path())
		// kinda brutal solution, could be improved
		ResetCache()
	}
}
