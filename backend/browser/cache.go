package browser

import (
	"io/fs"
	"iter"
	"log"
	"maps"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rjeczalik/notify"
)

type FileCache struct {
	cache         map[string]fs.FileInfo
	mutex         sync.Mutex
	dataDir       string
	pathConverter *PathConverter
}

func NewFileCache(dataDir string, pathConverter *PathConverter) *FileCache {
	return &FileCache{
		cache:         map[string]fs.FileInfo{},
		mutex:         sync.Mutex{},
		dataDir:       dataDir,
		pathConverter: pathConverter,
	}
}

// resets the cache
func (f *FileCache) Reset() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	log.Printf("Caching %s", f.dataDir)
	defer func() {
		if r := recover(); r != nil {
			log.Print("Failed to refresh cache")
		}
	}()

	newFileCache := map[string]fs.FileInfo{}
	f.cacheFolder(newFileCache, "/")

	// just doing this at once should avoid excessive lock/unlock
	f.cache = newFileCache
	log.Print("Refreshed cache")
}

// recursively initalizes the cache
// we do not use filepath.Walk to avoid a mess with relative / absolute paths
func (f *FileCache) cacheFolder(cache map[string]fs.FileInfo, path string) {
	rootPath := f.pathConverter.GetAbsolutePath(path)
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
func (f *FileCache) WatchDataDir() {
	c := make(chan notify.EventInfo, 1)

	// watch the create, remove, rename events on the data dir and sub directories
	if err := notify.Watch(f.dataDir+"/...", c, notify.Create, notify.Remove, notify.Rename); err != nil {
		panic(err)
	}
	defer notify.Stop(c)

	for ei := range c {
		log.Printf("Event %d on file %s", ei.Event(), ei.Path())
		// kinda brutal solution, could be improved
		f.Reset()
	}
}

func (f *FileCache) Iterate() iter.Seq2[string, fs.FileInfo] {
	return maps.All(f.cache)
}

func (f *FileCache) Get(key string) (fs.FileInfo, bool) {
	res, ok := f.cache[key]
	return res, ok
}
