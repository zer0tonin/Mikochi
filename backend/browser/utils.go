package browser

import (
	"io/fs"
	"os"
	"path/filepath"
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

// fileInDir returns if a file is contained in a directory (and not a sub-directory)
func fileInDir(file, dir string) bool {
	return filepath.Dir(file) == dir
}

func isDir(filepath string) (bool, error) {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}
