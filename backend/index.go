package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// serveIndex serves the index.html file while injecting a JS config in the head
func serveIndex(c *gin.Context, frontendFiles string, pathPrefix string) {
	indexPath := filepath.Join(frontendFiles, "index.html")

	data, err := os.ReadFile(indexPath)
	if err != nil {
		log.Panicf("Index.html unreadable")
	}
	html := string(data)


	// injecting Mikochi config
	configScript := fmt.Sprintf(
		"<script>window.__MIKOCHI_CONFIG__={pathPrefix:%q};</script></head>",
		pathPrefix,
	)
	html = strings.Replace(
		html,
		"</head>",
		configScript,
		1,
	)

	// rewriting js/css assets path
	assetsPath := fmt.Sprintf("%sassets", pathPrefix)
	html = strings.Replace(
		html,
		"/assets",
		assetsPath,
		-1,
	)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

// serveIndexMiddleware will call custom serveIndex logic before static.ServeRoot
func serveIndexMiddleware(frontendFiles string, pathPrefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := normalizePath(c.Request.URL.Path)

		isRoot := path == pathPrefix
		isIndex := path == pathPrefix+"index.html"

		if isRoot || isIndex {
			serveIndex(c, frontendFiles, pathPrefix)
			c.Abort()
			return
		}

		c.Next()
	}
}

func normalizePath(path string) string {
	if path == "" {
		path = "/"
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	return path
}
