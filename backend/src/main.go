package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

// The watched directectory
var dataDir string

func main() {
	dataDir = os.Getenv("data_dir")

	fmt.Println("Caching " + dataDir)
	resetCache()

	go watchDataDir()

	r := gin.Default()
	r.GET("/browse/*path", browseFolder)
	r.GET("/stream/:path", streamFile)
	r.GET("/search/:path", searchFile)
	host := os.Getenv("host")
	fmt.Println("Listening on " + host)
	r.Run(host)
}
