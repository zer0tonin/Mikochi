package main

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func getAbsolutePath(path string) string {
	cleanedPath := filepath.Clean(path)
	return filepath.Join(dataDir + cleanedPath)
}

// GET /search
func searchFile(c *gin.Context) {
	path := filepath.Clean(c.Param("path"))

	results := []string{}
	for fileName, _ := range fileCache {
		if strings.Contains(filepath.Clean(fileName), path) {
			results = append(results, fileName)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"dirEntries": results,
		"isRoot":     false,
	})
}

// GET /stream
func streamFile(c *gin.Context) {
	pathInDataDir := getAbsolutePath(c.Param("path"))
	c.File(pathInDataDir)
}

// GET /browse
func browseFolder(c *gin.Context) {
	path := filepath.Clean(c.Param("path"))

	results := []string{}
	for fileName, _ := range fileCache {
		if strings.HasPrefix(filepath.Clean(fileName), path) {
			results = append(results, fileName)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"dirEntries": results,
		"isRoot":     path == dataDir,
	})
}
