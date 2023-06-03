package main

import (
	"log"
	"net/http"

	"github.com/zer0tonin/mikochi/auth"
	"github.com/zer0tonin/mikochi/browser"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetDefault("DATA_DIR", "/data")
	viper.SetDefault("JWT_SECRET", auth.GenerateRandomSecret())
	viper.SetDefault("USERNAME", "root")
	viper.SetDefault("PASSWORD", "pass")
	viper.SetDefault("HOST", "0.0.0.0:8080")
	viper.AutomaticEnv()

	browser.ResetCache()

	go browser.WatchDataDir()

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// in production builds, this route serves the frontend files
	// in the dev environment, this is handled by the frontend container
	r.NoRoute(gin.WrapH(http.FileServer(gin.Dir("./static", false))))

	api := r.Group("/api")

	// business logic
	api.GET("/browse/*path", auth.CheckJWT, browser.BrowseFolder)
	api.GET("/stream/*path", auth.CheckSingleUseJWT, browser.StreamFile)
	api.PUT("/move/*path", auth.CheckJWT, browser.Move)
	api.DELETE("/delete/*path", auth.CheckJWT, browser.Delete)
	api.PUT("/upload/*path", auth.CheckJWT, browser.Upload)

	// authentication
	api.GET("/refresh", auth.CheckJWT, auth.Refresh)
	api.GET("/single-use", auth.CheckJWT, auth.SingleUse)
	api.POST("/login", auth.Login)

	// k8s ready/live check
	r.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	host := viper.GetString("HOST")
	log.Print("Listening on " + host)
	r.Run(host)
}
