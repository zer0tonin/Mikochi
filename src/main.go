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

	fmt.Println(viper.GetString("dataDir"))
	files, err := os.ReadDir(viper.GetString("dataDir") + path)
	if err != nil {
		return
	}

	for _, file := range files {
		fmt.Println(file.Name())
		result = append(result, file.Name())
	}
	return result, nil
}

func displayFolder(w http.ResponseWriter, path string) error {
	files, err := readDataDir(path)
	if err != nil {
		return err
	}

	folderTemplate.Execute(
		w,
		map[string]interface{}{
			"files": files,
		},
	)
	return nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Path == "/" {
				displayFolder(w, "")
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	fmt.Println("Listening on " + viper.GetString("host"))
	http.ListenAndServe(viper.GetString("host"), nil)
}
