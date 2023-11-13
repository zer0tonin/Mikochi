package browser

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/spf13/viper"
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

// searchInDir returns the results of a search query inside a given directory
func searchInDir(dir, search string) []FileDescription {
	children := []string{}
	for file := range fileCache {
		if filepath.HasPrefix(file, dir) {
			children = append(children, file)
		}
	}

	matches := fuzzy.RankFindNormalizedFold(search, children)
	sort.Sort(matches)

	results := []FileDescription{}
	for _, match := range matches {
		fileInfo := fileInfoToFileDescription(fileCache[match.Target], match.Target)
		results = append(results, fileInfo)
	}
	return results
}

func isDir(filepath string) (bool, error) {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}

func sendTarGz() error {
	return nil
}
