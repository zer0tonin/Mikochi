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
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
	Size  int64  `json:"size"`
}

func fileInfoToFileDescription(fileInfo fs.FileInfo) FileDescription {
	return FileDescription{
		Name:  fileInfo.Name(),
		IsDir: fileInfo.IsDir(),
		Size:  fileInfo.Size(),
	}
}

func getAbsolutePath(path string) string {
	cleanedPath := filepath.Clean(path)
	return filepath.Join(dataDir + cleanedPath)
}

// GET /search
func searchFile(c *gin.Context) {
	path := filepath.Clean(c.Param("path"))

	results := []FileDescription{}
	for file, fileInfo := range fileCache {
		if strings.Contains(filepath.Clean(file), path) {
			results = append(results, fileInfoToFileDescription(fileInfo))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"fileInfos": results,
		"isRoot":    false,
	})
}

// GET /stream
func streamFile(c *gin.Context) {
	pathInDataDir := getAbsolutePath(c.Param("path"))
	c.File(pathInDataDir)
}

// fileInDir returns if a file is contained in a directory (and not a sub-directory)
func fileInDir(file, dir string) bool {
	suffix, found := strings.CutPrefix(file, dir)
	if !found {
		return false
	}
	return !strings.ContainsRune(suffix, '/')
}

// suffixPathWithSlash adds a trailing / at the end of a path if it is not already present
func suffixPathWithSlash(dir string) string {
	if strings.HasSuffix(dir, "/") {
		return dir
	}
	return dir + "/"
}

// GET /browse
func browseFolder(c *gin.Context) {
	path := filepath.Clean(c.Param("path"))
	path = suffixPathWithSlash(path)

	results := []FileDescription{}
	for file, fileInfo := range fileCache {
		if fileInDir(filepath.Clean(file), path) {
			results = append(results, fileInfoToFileDescription(fileInfo))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"fileInfos": results,
		"isRoot":    path == "/",
	})
}
