package browser

import "path/filepath"

type PathConverter struct {
	dataDir string
}

func NewPathConverter(dataDir string) *PathConverter {
	return &PathConverter{
		dataDir: dataDir,
	}
}

func (p *PathConverter) GetAbsolutePath(path string) string {
	cleanedPath := filepath.Clean(path)
	if cleanedPath == "." || cleanedPath == ".." {
		return p.dataDir
	}
	return filepath.Join(p.dataDir + cleanedPath)
}
