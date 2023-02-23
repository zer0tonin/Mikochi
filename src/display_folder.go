package main

import (
	"html/template"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

var folderTemplate *template.Template

func displayFolder(w http.ResponseWriter, path string) {
	pathInDataDir := getPathInDataDir(path)
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
