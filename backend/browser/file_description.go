package browser

import (
	"io/fs"
	"path/filepath"
)

// FileDescription is a serializable FileInfo with path
type FileDescription struct {
	IsDir bool   `json:"isDir"`
	Size  int64  `json:"size"`
	Path  string `json:"path"`
}

func fileInfoToFileDescription(fileInfo fs.FileInfo, relativePath string) FileDescription {
	return FileDescription{
		IsDir: fileInfo.IsDir(),
		Size:  fileInfo.Size(),
		Path:  relativePath,
	}
}

func (f *FileDescription) IsInDir(dir string) bool {
	return filepath.Dir(f.Path) == dir
}
