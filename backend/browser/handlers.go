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
	fileCache *FileCache
	pathConverter *PathConverter
}

func NewBrowserHandlers(fileCache *FileCache, pathConverter *PathConverter) *BrowserHandlers {
	return &BrowserHandlers{
		fileCache: fileCache,
		pathConverter: pathConverter,
	}
}

// GET /stream
// StreamFile streams the content of the requested file
func (b *BrowserHandlers) StreamFile(c *gin.Context) {
	path := c.Param("path")
	pathInDataDir := b.pathConverter.GetAbsolutePath(path)

	dir, err := isDir(pathInDataDir)
	if err != nil {
		log.Printf("Err: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "Failed to stat file",
		})
		return
	}

	if dir {
		c.Header(
			"Content-Disposition",
			fmt.Sprintf("attachment; filename=%s.tar.gz", path[1:len(path)-1]),
		)
		writeTarGz(pathInDataDir, c.Writer)
	} else {
		c.File(pathInDataDir)
	}
}

// GET /browse
// BrowseFolder returns the content of a directory and/or search results
func (b *BrowserHandlers) BrowseFolder(c *gin.Context) {
	path := filepath.Clean(c.Param("path"))
	search := c.Query("search")

	var results []FileDescription
	if search == "" {
		results = b.browseDir(path)
	} else {
		search = filepath.Clean(search)
		results = b.searchInDir(path, search)
	}

	c.JSON(http.StatusOK, gin.H{
		"fileInfos": results,
		"isRoot":    path == "/",
	})
}

// browseDir returns the content of a directory
func (b *BrowserHandlers) browseDir(dir string) []FileDescription {
	results := []FileDescription{}
	for file, fileInfo := range b.fileCache.Iterate() {
		if fileInDir(file, dir) {
			results = append(results, fileInfoToFileDescription(fileInfo, file))
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
		fileInfo, ok := b.fileCache.Get(match.Target)
		if !ok {
			continue
		}
		fileDescription := fileInfoToFileDescription(fileInfo, match.Target)
		results = append(results, fileDescription)
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

	path := b.pathConverter.GetAbsolutePath(c.Param("path"))
	newPath := b.pathConverter.GetAbsolutePath(command.NewPath)

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
}

// DELETE /delete
// Delete deletes a file from the filesystem
func (b *BrowserHandlers) Delete(c *gin.Context) {
	// we make a synchronous reset cache to avoid querying /browse on stale data
	defer b.fileCache.Reset()

	path := b.pathConverter.GetAbsolutePath(c.Param("path"))

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

	pathInDataDir := b.pathConverter.GetAbsolutePath(c.Param("path"))
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
func (b *BrowserHandlers) Mkdir(c *gin.Context) {
	// we make a synchronous reset cache to avoid querying /browse on stale data
	defer b.fileCache.Reset()

	pathInDataDir := b.pathConverter.GetAbsolutePath(c.Param("path"))
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
