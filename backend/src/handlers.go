package main

import (
	"io"
	"io/fs"
	"net/http"
	"os"
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
		Path:  path,
	}
}

func getAbsolutePath(path string) string {
	cleanedPath := filepath.Clean(path)
	return filepath.Join(dataDir + cleanedPath)
}

// GET /stream
// streamFile streams the content of the requested file
func streamFile(c *gin.Context) {
	pathInDataDir := getAbsolutePath(c.Param("path"))
	c.File(pathInDataDir)
}

// fileInDir returns if a file is contained in a directory (and not a sub-directory)
func fileInDir(file, dir string) bool {
	return filepath.Dir(file) == dir
}

// browseDir returns the content of a directory
func browseDir(dir string) []FileDescription {
	results := []FileDescription{}
	for file, fileInfo := range fileCache {
		if fileInDir(file, dir) {
			results = append(results, fileInfoToFileDescription(fileInfo, fileInfo.Name()))
		}
	}
	return results
}

// fileMatchesSearch checks that a file is contained in a directory (or its subdiretories) and matches the search query
func fileMatchesSearch(file, dir, search string) bool {
	rel, err := filepath.Rel(dir, file)
	if err != nil {
		return false // file not in dir
	}
	// we want search to be case insensitive
	rel = strings.ToLower(rel)
	search = strings.ToLower(search)
	return strings.Contains(rel, search) && !strings.HasPrefix(rel, "../")
}

// searchInDir returns the results of a search query inside a given directory
func searchInDir(dir, search string) []FileDescription {
	results := []FileDescription{}
	for file, fileInfo := range fileCache {
		if fileMatchesSearch(file, dir, search) {
			path, _ := strings.CutPrefix(file, dir+"/")
			results = append(results, fileInfoToFileDescription(fileInfo, path))
		}
	}
	return results
}

// GET /browse
// browseFolder returns the content of a directory and/or search results
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

// PUT /move
// move is used to move a file or change its name
func move(c *gin.Context) {
	var command struct {
		NewPath string `json:"newPath"`
	}
	err := c.BindJSON(&command)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "Couldn't deserialize command",
		})
		return
	}

	path := getAbsolutePath(c.Param("path"))
	newPath := getAbsolutePath(command.NewPath)

	err = os.Rename(path, newPath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "Couldn't move file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

	// cache refresh should be triggered automatically
}

// DELETE /delete
// delete deletes a file from the filesystem
func delete(c *gin.Context) {
	path := getAbsolutePath(c.Param("path"))

	err := os.Remove(path)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "Couldn't move file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

	// cache refresh should be triggered automatically
}

// PUT /upload
// writes data received in multi-part to the disk
func upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Couldn't read file from request",
		})
		return
	}

	pathInDataDir := getAbsolutePath(c.Param("path"))
	dst, err := os.Create(pathInDataDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "Couldn't create destination file",
		})
		return
	}
	defer dst.Close()

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "Failed to open file from request",
		})
		return
	}
	defer src.Close()

	// Copy the file to the destination
	_, err = io.Copy(dst, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "Failed to write file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
