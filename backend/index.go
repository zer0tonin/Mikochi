package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// serveIndex serves the index.html file while injecting a JS config in the head
func serveIndex(c *gin.Context, frontendFiles string, pathPrefix string) bool {
	indexPath := filepath.Join(frontendFiles, "index.html")

	data, err := os.ReadFile(indexPath)
	if err != nil {
		return false
	}

	configScript := fmt.Sprintf(
		"<script>window.__MIKOCHI_CONFIG__={pathPrefix:%q};</script>",
		pathPrefix,
	)

	html := string(data)
	html = strings.Replace(
		html,
		"<!-- mikochi-config -->",
		configScript,
		1,
	)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)

	return true
}

// serveIndexMiddleware will call custom serveIndex logic before static.ServeRoot
func serveIndexMiddleware(frontendFiles string, pathPrefix string) gin.HandlerFunc {
    return func(c *gin.Context) {
        path := normalizePath(c.Request.URL.Path)
        isRoot := path == pathPrefix
        isIndex := path == pathPrefix + "index.html"

        if isRoot || isIndex {
            if serveIndex(c, frontendFiles, pathPrefix) {
                c.Abort()
                return
            }
        }

        c.Next()
    }
}
