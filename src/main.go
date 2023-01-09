package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
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

func readDataDir(path string) (result []string, err error) {
	if strings.Contains(path, "..") {
		return result, fmt.Errorf("Invalid path")
	}

	// FIXME: filepath.Join
	files, err := os.ReadDir(viper.GetString("dataDir") + path)
	if err != nil {
		return
	}

	for _, file := range files {
		result = append(result, file.Name())
	}
	return result, nil
}

func displayFolder(w http.ResponseWriter, path string) {
	files, err := readDataDir(path)
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
	path := r.URL.Path[7:]
	// FIXME WTF???
	if strings.Contains(path, "..") {
		fmt.Println("Here")
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	http.ServeFile(w, r, viper.GetString("dataDir") + path)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// TODO: two path /browse and /stream
			// / should redirect to /browse
			if r.URL.Path == "/" {
				http.Redirect(w, r, "/browse", 301)
				return
			}

			if r.URL.Path == "/browse" {
				displayFolder(w, "")
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
