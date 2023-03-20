package main

import (
	"fmt"
	"html/template"
	"net/http"

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

	resetCache()
}

func main() {
	http.HandleFunc("/", routes)
	go watchDataDir()
	fmt.Println("Listening on " + viper.GetString("host"))
	http.ListenAndServe(viper.GetString("host"), nil)
}
