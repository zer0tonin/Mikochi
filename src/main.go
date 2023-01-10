package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var folderTemplate *template.Template

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Failed to load config")
	}

	folderTemplate = template.Must(template.ParseFiles("./templates/folder.html"))
}

func getPathInDataDir(path string) string {
	cleanedPath := filepath.Clean(path)
	return filepath.Join(viper.GetString("dataDir") + cleanedPath)
}

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

func streamFile(w http.ResponseWriter, r *http.Request) {
	path := getPathInDataDir(r.URL.Path[7:])
	http.ServeFile(w, r, path)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

			http.Error(w, "Not Found", http.StatusNotFound)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	fmt.Println("Listening on " + viper.GetString("host"))
	http.ListenAndServe(viper.GetString("host"), nil)
}
