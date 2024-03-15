package main

import (
	"log"
	"net/http"

	"github.com/zer0tonin/mikochi/auth"
	"github.com/zer0tonin/mikochi/browser"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetDefault("DATA_DIR", "/data")
	viper.SetDefault("JWT_SECRET", auth.GenerateRandomSecret())
	viper.SetDefault("USERNAME", "root")
	viper.SetDefault("PASSWORD", "pass")
	viper.SetDefault("HOST", "0.0.0.0:8080")
	viper.SetDefault("ENV", "production")
	viper.AutomaticEnv()

	browser.ResetCache()

	go browser.WatchDataDir()

	if viper.GetString("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()
	r.Use(gin.Recovery())

	// in production builds, this route serves the frontend files
	// in the dev environment, this is handled by the frontend container
	r.Use(static.ServeRoot("/", "./static"))
	r.NoRoute(func (c *gin.Context) {
		// we let the client-side routing take over
		c.File("./static/index.html")
	})

	api := r.Group("/api")

	// business logic
	api.GET("/browse/*path", auth.CheckJWT, browser.BrowseFolder)
	api.GET("/stream/*path", auth.CheckSingleUseJWT, browser.StreamFile)
	api.PUT("/move/*path", auth.CheckJWT, browser.Move)
	api.DELETE("/delete/*path", auth.CheckJWT, browser.Delete)
	api.PUT("/upload/*path", auth.CheckJWT, browser.Upload)
	api.PUT("/mkdir/*path", auth.CheckJWT, browser.Mkdir)

	// authentication
	api.GET("/refresh", auth.CheckJWT, auth.Refresh)
	api.GET("/single-use", auth.CheckJWT, auth.SingleUse)
	api.POST("/login", auth.Login)

	// k8s ready/live check
	r.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	host := viper.GetString("HOST")
	log.Print("Listening on " + host)

	var err error
	if viper.IsSet("CERT_CA") && viper.IsSet("CERT_KEY") {
		err = r.RunTLS(host, viper.GetString("CERT_CA"), viper.GetString("CERT_KEY"))
	} else {
		err = r.Run(host)
	}

	if err != nil {
		log.Panicf("Failed to launch web server: %s", err.Error())
	}
}
