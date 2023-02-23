package main

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var fileCache map[string]os.DirEntry
var searchTemplate *template.Template

func searchFile(w http.ResponseWriter, r *http.Request) {
	result := []os.DirEntry{}
	path := filepath.Clean(r.URL.Path[7:])
	for fileName, dirEntry := range fileCache {
		if strings.Contains(filepath.Clean(fileName), path) {
			result = append(result, dirEntry)
		}
	}

	err := searchTemplate.Execute(
		w,
		map[string]interface{}{
			"dirEntries": result,
		},
	)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

