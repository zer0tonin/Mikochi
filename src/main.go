package main

import (
	"fmt"
	"html/template"
	"os"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Failed to load config")
	}

	folderTemplate = template.Must(template.ParseFiles("./templates/folder.html"))
	searchTemplate = template.Must(template.ParseFiles("./templates/search.html"))

	fileCache = map[string]os.DirEntry{}
	cacheFolder(fileCache, "/")
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

			if strings.HasPrefix(r.URL.Path, "/search/") {
				searchFile(w, r)
			}

			http.Error(w, "Not Found", http.StatusNotFound)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	fmt.Println("Listening on " + viper.GetString("host"))
	http.ListenAndServe(viper.GetString("host"), nil)
}
