package main

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func fileMatchesSearch(file, dir, search string) bool {
	rel, err := filepath.Rel(dir, file)
	if err != nil {
		return false // file not in dir
	}
	return strings.Contains(rel, search) && !strings.HasPrefix(rel, "../")
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

// POST /login
func login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := c.BindJSON(credentials)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "Couldn't deserialize credentials",
		})
		return
	}

	if credentials.Username != os.Getenv("USERNAME") || credentials.Password != os.Getenv("PASSWORD") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"err": "Invalid credentials",
		})
		return
	}

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 730)),
		Issuer: "Mikochi",
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "Failed to generate authentication token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": signedToken,
	})
}
