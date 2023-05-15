package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// Config (read from env)
var dataDir string
var jwtSecret []byte
var username string
var password string

func main() {
	dataDir = os.Getenv("data_dir")
	jwtSecret = []byte(os.Getenv("jwt_secret"))
	username = os.Getenv("username")
	password = os.Getenv("password")

	tokenWhitelist = map[string]string{}

	log.Print("Caching " + dataDir)
	resetCache()

	go watchDataDir()

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/browse/*path", checkJWT, browseFolder)
	r.GET("/stream/*path", checkSingleUseJWT, streamFile)
	r.PUT("/move/*path", checkJWT, move)

	r.GET("/refresh", checkJWT, refresh)
	r.GET("/single-use", checkJWT, singleUse)
	r.POST("/login", login)

	r.GET("/ready", func (c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	host := os.Getenv("host")
	log.Print("Listening on " + host)
	r.Run(host)
}
