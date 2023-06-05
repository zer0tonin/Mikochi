package browser

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// GET /stream
// StreamFile streams the content of the requested file
func StreamFile(c *gin.Context) {
	pathInDataDir := getAbsolutePath(c.Param("path"))
	c.File(pathInDataDir)
}

// GET /browse
// BrowseFolder returns the content of a directory and/or search results
func BrowseFolder(c *gin.Context) {
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
// Move is used to move a file or change its name
func Move(c *gin.Context) {
	var command struct {
		NewPath string `json:"newPath"`
	}
	err := c.BindJSON(&command)
	if err != nil {
		log.Printf("Err: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "Couldn't deserialize command",
		})
		return
	}

	path := getAbsolutePath(c.Param("path"))
	newPath := getAbsolutePath(command.NewPath)

	err = os.Rename(path, newPath)
	if err != nil {
		log.Printf("Err: %s", err.Error())
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
// Delete deletes a file from the filesystem
func Delete(c *gin.Context) {
	path := getAbsolutePath(c.Param("path"))

	err := os.RemoveAll(path)
	if err != nil {
		log.Printf("Err: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "Couldn't remove file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

	// cache refresh should be triggered automatically
}

// PUT /upload
// Upload writes data received in multi-part to the disk
func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Err: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "Couldn't read file from request",
		})
		return
	}

	pathInDataDir := getAbsolutePath(c.Param("path"))
	dst, err := os.Create(pathInDataDir)
	if err != nil {
		log.Printf("Err: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "Couldn't create destination file",
		})
		return
	}
	defer dst.Close()

	src, err := file.Open()
	if err != nil {
		log.Printf("Err: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "Failed to open file from request",
		})
		return
	}
	defer src.Close()

	// Copy the file to the destination
	_, err = io.Copy(dst, src)
	if err != nil {
		log.Printf("Err: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "Failed to write file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// PUT /mkdir
// Mkdir creates a new (empty) directory
func Mkdir(c *gin.Context) {
	pathInDataDir := getAbsolutePath(c.Param("path"))
	err := os.Mkdir(pathInDataDir, 0755)
	if err != nil {
		log.Printf("Err: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "Failed to write file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
