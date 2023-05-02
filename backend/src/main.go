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
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/browse/*path", checkJWT, browseFolder)
	r.GET("/stream/*path", checkJWT, streamFile)
	r.POST("/login", login)

	host := os.Getenv("host")
	fmt.Println("Listening on " + host)
	r.Run(host)
}
