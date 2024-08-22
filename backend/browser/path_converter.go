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

// TODO: test
func (p *PathConverter) GetAbsolutePath(path string) string {
	cleanedPath := filepath.Clean(path)
	return filepath.Join(p.dataDir + cleanedPath)
}