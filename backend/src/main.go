package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

// Config (read from env)
var dataDir string
var jwtSecret string
var username string
var password string


func main() {
	dataDir = os.Getenv("data_dir")
	jwtSecret = os.Getenv("jwt_secret")
	username = os.Getenv("username")
	password = os.Getenv("password")

	fmt.Println("Caching " + dataDir)
	resetCache()

	go watchDataDir()

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/browse/*path", checkJWT, browseFolder)
	r.GET("/stream/*path", checkJWT, streamFile)
	r.GET("/refresh", checkJWT, refresh)
	r.POST("/login", login)

	host := os.Getenv("host")
	fmt.Println("Listening on " + host)
	r.Run(host)
}
