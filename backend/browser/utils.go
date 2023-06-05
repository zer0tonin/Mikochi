package browser

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/lithammer/fuzzysearch/fuzzy"
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
	return filepath.Join(viper.GetString("DATA_DIR") + cleanedPath)
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
		fmt.Println("here")
		return false // file not in dir
	}
	// we want search to be case insensitive
	rel = strings.ToLower(rel)
	search = strings.ToLower(search)
	fmt.Println(fuzzy.RankMatch(search, rel))
	return fuzzy.RankMatch(search, rel) > 0 && !strings.HasPrefix(rel, "../")
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
