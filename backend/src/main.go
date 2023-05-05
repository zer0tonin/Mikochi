package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

// Config (read from env)
var dataDir string
var jwtSecret []byte 
var username string
var password string

// Whitelist of single-use JWTs (for streams)
var tokenWhitelist map[string]bool

func main() {
	dataDir = os.Getenv("data_dir")
	jwtSecret = []byte(os.Getenv("jwt_secret"))
	username = os.Getenv("username")
	password = os.Getenv("password")

	tokenWhitelist = map[string]bool{}

	fmt.Println("Caching " + dataDir)
	resetCache()

	go watchDataDir()

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/browse/*path", checkJWT, browseFolder)
	r.GET("/stream/*path", checkSingleUseJWT, streamFile)
	r.GET("/refresh", checkJWT, refresh)
	r.GET("/single-use", checkJWT, singleUse)
	r.POST("/login", login)

	host := os.Getenv("host")
	fmt.Println("Listening on " + host)
	r.Run(host)
}
