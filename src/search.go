package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var fileCache map[string]os.DirEntry
var searchTemplate *template.Template

func resetCache() {
	fileCache = map[string]os.DirEntry{}
	cacheFolder(fileCache, "/")
}

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

func searchFile(w http.ResponseWriter, r *http.Request) {
	result := map[string]os.DirEntry{}
	path := filepath.Clean(r.URL.Path[8:])
	for fileName, dirEntry := range fileCache {
		if strings.Contains(filepath.Clean(fileName), path) {
			result[fileName] = dirEntry
		}
	}

	err := searchTemplate.Execute(
		w,
		map[string]interface{}{
			"result": result,
		},
	)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
