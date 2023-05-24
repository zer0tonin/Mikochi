package main

import (
	"log"
	"net/http"
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

	// in production builds, this route serves the frontend files
	// in the dev environment, this is handled by the frontend container
	r.NoRoute(gin.WrapH(http.FileServer(gin.Dir("./static", false))))

	api := r.Group("/api")

	// business logic
	api.GET("/browse/*path", checkJWT, browseFolder)
	api.GET("/stream/*path", checkSingleUseJWT, streamFile)
	api.PUT("/move/*path", checkJWT, move)
	api.DELETE("/delete/*path", checkJWT, delete)
	api.PUT("/upload/*path", checkJWT, upload)

	// authentication
	api.GET("/refresh", checkJWT, refresh)
	api.GET("/single-use", checkJWT, singleUse)
	api.POST("/login", login)

	// k8s ready/live check
	r.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	host := os.Getenv("host")
	log.Print("Listening on " + host)
	r.Run(host)
}
