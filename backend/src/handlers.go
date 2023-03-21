package main

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// FileDescription is just a serializable FileInfo
type FileDescription struct {
	Name string `json:"name"`
	IsDir bool `json:"isDir"`
	Size int64 `json:"size"`
}

func getAbsolutePath(path string) string {
	cleanedPath := filepath.Clean(path)
	return filepath.Join(dataDir + cleanedPath)
}

func fileInfoToFileDescription(fileInfo fs.FileInfo) FileDescription {
	return FileDescription{
		Name: fileInfo.Name(),
		IsDir: fileInfo.IsDir(),
		Size: fileInfo.Size(),
	}
}

// GET /search
func searchFile(c *gin.Context) {
	path := filepath.Clean(c.Param("path"))

	results := []FileDescription{}
	for fileName, fileInfo := range fileCache {
		if strings.Contains(filepath.Clean(fileName), path) {
			results = append(results, fileInfoToFileDescription(fileInfo))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"fileInfos": results,
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

	results := []FileDescription{}
	for fileName, fileInfo := range fileCache {
		if strings.HasPrefix(filepath.Clean(fileName), path) {
			results = append(results, fileInfoToFileDescription(fileInfo))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"fileInfos": results,
		"isRoot":     path == dataDir,
	})
}
