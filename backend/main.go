package main

import (
	"log"
	"net/http"

	"github.com/zer0tonin/mikochi/auth"
	"github.com/zer0tonin/mikochi/browser"

	"github.com/gin-contrib/static"
	"github.com/gin-contrib/gzip"
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
	viper.SetDefault("NO_AUTH", "false")
	viper.SetDefault("GZIP", "false")
	viper.AutomaticEnv()

	browser.ResetCache()

	authMiddleware := auth.NewAuthMiddleware(viper.GetString("NO_AUTH") != "true", viper.GetString("JWT_SECRET"))
	authHandlers := auth.NewAuthHandlers(authMiddleware, viper.GetString("USERNAME"), viper.GetString("PASSWORD"))

	go browser.WatchDataDir()

	if viper.GetString("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()
	r.Use(gin.Recovery())
	if viper.GetBool("GZIP") {
		r.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	// in production builds, this route serves the frontend files
	// in the dev environment, this is handled by the frontend container
	r.Use(static.ServeRoot("/", "./static"))
	r.NoRoute(func(c *gin.Context) {
		// we let the client-side routing take over
		c.File("./static/index.html")
	})

	api := r.Group("/api")

	// business logic
	api.GET("/browse/*path", authMiddleware.CheckAuth, browser.BrowseFolder)
	api.GET("/stream/*path", authMiddleware.CheckStreamAuth, browser.StreamFile)
	api.PUT("/move/*path", authMiddleware.CheckAuth, browser.Move)
	api.DELETE("/delete/*path", authMiddleware.CheckAuth, browser.Delete)
	api.PUT("/upload/*path", authMiddleware.CheckAuth, browser.Upload)
	api.PUT("/mkdir/*path", authMiddleware.CheckAuth, browser.Mkdir)

	// authentication
	api.GET("/refresh", authMiddleware.CheckAuth, authHandlers.Refresh)
	api.GET("/single-use", authMiddleware.CheckAuth, authHandlers.SingleUse)
	api.POST("/login", authHandlers.Login)

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
