package browser

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type BrowserHandlers struct {
	fileCache     *FileCache
	pathConverter *PathConverter
	compressor    Compressor
}

type Compressor interface {
	CompressDirectory(string, io.Writer) error
	GetFileExtension() string
}

func NewBrowserHandlers(fileCache *FileCache, pathConverter *PathConverter, compressor Compressor) *BrowserHandlers {
	return &BrowserHandlers{
		fileCache:     fileCache,
		pathConverter: pathConverter,
		compressor:    compressor,
	}
}

// GET /stream
// StreamFile streams the content of the requested file
func (b *BrowserHandlers) StreamFile(c *gin.Context) {
	relativePath := c.Param("path")
	if strings.HasSuffix(relativePath, "/") {
		// path do not have a trailing / in the cache
		relativePath = relativePath[:len(relativePath)-1]
	}

	absolutePath, err := b.pathConverter.GetAbsolutePath(relativePath)
	if err != nil {
		log.Printf("Err: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "Invalid path",
		})
		return
	}

	fileDescription, ok := b.fileCache.Get(relativePath)
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"err": "File not found",
		})
		return
	}

	if fileDescription.IsDir {
		c.Header(
			"Content-Disposition",
			fmt.Sprintf("attachment; filename=%s%s", relativePath[1:len(relativePath)], b.compressor.GetFileExtension()),
		)
		err = b.compressor.CompressDirectory(absolutePath, c.Writer)
		if err != nil {
			log.Printf("Err: %s", err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"err": "Failed to compress directory",
			})
			return
		}
	} else {
		c.File(absolutePath)
	}
}

// GET /browse
// BrowseFolder returns the content of a directory and/or search results
func (b *BrowserHandlers) BrowseFolder(c *gin.Context) {
	relativePath := filepath.Clean(c.Param("path"))
	search := c.Query("search")

	var results []FileDescription
	if search == "" {
		results = b.browseDir(relativePath)
	} else {
		search = filepath.Clean(search)
		results = b.searchInDir(relativePath, search)
	}

	c.JSON(http.StatusOK, gin.H{
		"fileInfos": results,
		"isRoot":    relativePath == "/",
	})
}

// browseDir returns the content of a directory
func (b *BrowserHandlers) browseDir(dir string) []FileDescription {
	results := []FileDescription{}
	for _, fileDescription := range b.fileCache.Iterate() {
		if fileDescription.IsInDir(dir) {
			results = append(results, fileDescription)
		}
	}
	return results
}

// searchInDir returns the results of a search query inside a given directory
func (b *BrowserHandlers) searchInDir(dir, search string) []FileDescription {
	children := []string{}
	for file := range b.fileCache.Iterate() {
		if strings.HasPrefix(file, dir) {
			children = append(children, file)
		}
	}

	matches := fuzzy.RankFindNormalizedFold(search, children)
	sort.Sort(matches)

	results := []FileDescription{}
	for _, match := range matches {
		if fileDescription, ok := b.fileCache.Get(match.Target); ok {
			results = append(results, fileDescription)
		}
	}
	return results
}

// PUT /move
// Move is used to move a file or change its name
func (b *BrowserHandlers) Move(c *gin.Context) {
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

	// we make a synchronous reset cache to avoid querying /browse on stale data
	defer b.fileCache.Reset()

	absolutePath, err := b.pathConverter.GetAbsolutePath(c.Param("path"))
	if err != nil {
		log.Printf("Err: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "Invalid path",
		})
		return
	}

	newAbsolutePath, err := b.pathConverter.GetAbsolutePath(command.NewPath)
	if err != nil {
		log.Printf("Err: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "Invalid path",
		})
		return
	}

	err = os.Rename(absolutePath, newAbsolutePath)
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
}

// DELETE /delete
// Delete deletes a file from the filesystem
func (b *BrowserHandlers) Delete(c *gin.Context) {
	// we make a synchronous reset cache to avoid querying /browse on stale data
	defer b.fileCache.Reset()

	absolutePath, err := b.pathConverter.GetAbsolutePath(c.Param("path"))
	if err != nil {
		log.Printf("Err: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "Invalid path",
		})
		return
	}

	err = os.RemoveAll(absolutePath)
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
}

// PUT /upload
// Upload writes data received in multi-part to the disk
func (b *BrowserHandlers) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Err: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "Couldn't read file from request",
		})
		return
	}

	// we make a synchronous reset cache to avoid querying /browse on stale data
	defer b.fileCache.Reset()

	absolutePath, err := b.pathConverter.GetAbsolutePath(c.Param("path"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "Invalid path",
		})
		return
	}

	dst, err := os.Create(absolutePath)
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
func (b *BrowserHandlers) Mkdir(c *gin.Context) {
	// we make a synchronous reset cache to avoid querying /browse on stale data
	defer b.fileCache.Reset()

	absolutePath, err := b.pathConverter.GetAbsolutePath(c.Param("path"))
	if err != nil {
		log.Printf("Err: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "Invalid path",
		})
		return
	}

	err = os.Mkdir(absolutePath, 0755)
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
