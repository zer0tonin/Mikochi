package browser

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

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
