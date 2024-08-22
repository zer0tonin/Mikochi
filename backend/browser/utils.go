package browser

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/fs"
	"os"
	"path/filepath"

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

func isDir(filepath string) (bool, error) {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}

func writeTarGz(filepath string, w io.Writer) error {
	gzipWriter := gzip.NewWriter(w)
	defer gzipWriter.Close()
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	writeDirectoryToTarGz(filepath, tarWriter, filepath)

	return nil
}

func writeDirectoryToTarGz(directory string, tarWriter *tar.Writer, subPath string) error {
	files, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	for _, file := range files {
		currentPath := filepath.Join(directory, file.Name())
		if file.IsDir() {
			err := writeDirectoryToTarGz(currentPath, tarWriter, subPath)
			if err != nil {
				return err
			}
		} else {
			fileInfo, err := file.Info()
			if err != nil {
				return err
			}
			err = writeFileToTarGz(currentPath, tarWriter, fileInfo, subPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func writeFileToTarGz(path string, tarWriter *tar.Writer, fileInfo os.FileInfo, subPath string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}

	subPath, err = filepath.EvalSymlinks(subPath)
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(fileInfo, path)
	if err != nil {
		return err
	}
	header.Name = path[len(subPath):]

	err = tarWriter.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(tarWriter, file)
	return err
}
