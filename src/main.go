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

func readDataDir(path string) (result []string, err error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return
	}

	for _, file := range files {
		result = append(result, file.Name())
	}
	return result, nil
}

func displayFolder(w http.ResponseWriter, path string) {
	files, err := readDataDir(getPathInDataDir(path))
	if err != nil {
		// Most likely explanation for error is that the directory doesn't exist
		http.Error(w, "Directory not found", http.StatusNotFound)
		return
	}

	err = folderTemplate.Execute(
		w,
		map[string]interface{}{
			"files": files,
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
				http.Redirect(w, r, "/browse", 301)
				return
			}

			if r.URL.Path == "/browse" {
				displayFolder(w, "/")
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
