package main

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"fmt"
)

var fileCache map[string]os.DirEntry
var searchTemplate *template.Template

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

func searchFile(w http.ResponseWriter, r *http.Request) {
	result := map[string]os.DirEntry{}
	path := filepath.Clean(r.URL.Path[8:])
	for fileName, dirEntry := range fileCache {
		fmt.Println(dirEntry.Name())
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
