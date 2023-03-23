package main

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// FileDescription is a serializable FileInfo with path
type FileDescription struct {
	IsDir bool   `json:"isDir"`
	Size  int64  `json:"size"`
	Path  string `json:"path"`
}

func fileInfoToFileDescription(fileInfo fs.FileInfo, path string) FileDescription {
	return FileDescription{
		IsDir: fileInfo.IsDir(),
		Size:  fileInfo.Size(),
		Path: path,
	}
}

func getAbsolutePath(path string) string {
	cleanedPath := filepath.Clean(path)
	return filepath.Join(dataDir + cleanedPath)
}

// GET /stream
func streamFile(c *gin.Context) {
	pathInDataDir := getAbsolutePath(c.Param("path"))
	c.File(pathInDataDir)
}

// fileInDir returns if a file is contained in a directory (and not a sub-directory)
func fileInDir(file, dir string) bool {
	return filepath.Dir(file) == dir
}

func browseDir(dir string) []FileDescription {
	results := []FileDescription{}
	for file, fileInfo := range fileCache {
		if fileInDir(file, dir) {
			results = append(results, fileInfoToFileDescription(fileInfo, fileInfo.Name()))
		}
	}
	return results
}

// TODO test
func fileMatchesSearch(file, dir, search string) bool {
	return strings.Contains(file, search) && strings.HasPrefix(file, dir)
}

func searchInDir(dir, search string) []FileDescription {
	results := []FileDescription{}
	for file, fileInfo := range fileCache {
		if fileMatchesSearch(file, dir, search) {
			path, _ := strings.CutPrefix(file, dir + "/")
			results = append(results, fileInfoToFileDescription(fileInfo, path))
		}
	}
	return results
}

// GET /browse
func browseFolder(c *gin.Context) {
	path := filepath.Clean(c.Param("path"))

	search := c.Query("search")
	var results []FileDescription
	if search == "" {
		results = browseDir(path)
	} else {
		search = filepath.Clean(search)
		results = searchInDir(path, search)
	}

	c.JSON(http.StatusOK, gin.H{
		"fileInfos": results,
		"isRoot":    path == "/",
	})
}
