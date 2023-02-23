package main

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var searchTemplate *template.Template
var folderTemplate *template.Template

// GET /search
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

func getAbsolutePath(path string) string {
	cleanedPath := filepath.Clean(path)
	return filepath.Join(viper.GetString("dataDir") + cleanedPath)
}

// GET /stream
func streamFile(w http.ResponseWriter, r *http.Request) {
	path := getAbsolutePath(r.URL.Path[7:])
	http.ServeFile(w, r, path)
}

// GET /browse
func displayFolder(w http.ResponseWriter, path string) {
	pathInDataDir := getAbsolutePath(path)
	dirEntries, err := os.ReadDir(pathInDataDir)
	if err != nil {
		// Most likely explanation for error is that the directory doesn't exist
		http.Error(w, "Directory not found", http.StatusNotFound)
		return
	}

	err = folderTemplate.Execute(
		w,
		map[string]interface{}{
			"dirEntries": dirEntries,
			"isRoot": pathInDataDir == viper.GetString("dataDir"),
		},
	)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func routes(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
	    if r.URL.Path == "/" {
		    http.Redirect(w, r, "/browse/", 301)
		    return
	    }

	    // if the browse page is not suffixed by "/", we append one to simplify navigation
	    if strings.HasPrefix(r.URL.Path, "/browse") && r.URL.Path[len(r.URL.Path)-1] != '/' {
		    http.Redirect(w, r, r.URL.Path + "/", 301)
		    return
	    }

	    if strings.HasPrefix(r.URL.Path, "/browse/") {
		    displayFolder(w, r.URL.Path[7:])
		    return
	    }

	    if strings.HasPrefix(r.URL.Path, "/stream/") {
		    streamFile(w, r)
		    return
	    }

	    if strings.HasPrefix(r.URL.Path, "/search/") {
		    searchFile(w, r)
		    return
	    }

	    http.Error(w, "Not Found", http.StatusNotFound)
    default:
	    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}
