package browser

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

type PathConverter struct {
	dataDir string
}

func NewPathConverter(dataDir string) *PathConverter {
	return &PathConverter{
		dataDir: dataDir,
	}
}

func (p *PathConverter) GetAbsolutePath(path string) (string, error) {
	cleanedPath := filepath.Clean(path)
	if cleanedPath == "" || cleanedPath == "." || cleanedPath == ".." {
		return p.dataDir, nil
	}

	res := filepath.Join(p.dataDir + cleanedPath)

	if rel, err := filepath.Rel(p.dataDir, res); err != nil || strings.HasPrefix(rel, "..") {
		log.Printf("Possible directory escape attempt: %s", path)
		return p.dataDir, fmt.Errorf("Invalid target path: %s (%w)", path, err)
	}

	return res, nil
}
