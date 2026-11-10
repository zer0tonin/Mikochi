package browser

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

type TarGzCompressor struct {}

func NewTarGzCompressor() *TarGzCompressor {
	return &TarGzCompressor{}
}

func (t *TarGzCompressor) CompressDirectory(filepath string, w io.Writer) error {
	gzipWriter := gzip.NewWriter(w)
	defer gzipWriter.Close()
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	t.writeDirectoryToTarGz(filepath, tarWriter, filepath)

	return nil
}

func (t *TarGzCompressor) writeDirectoryToTarGz(directory string, tarWriter *tar.Writer, subPath string) error {
	files, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	for _, file := range files {
		currentPath := filepath.Join(directory, file.Name())
		if file.IsDir() {
			err := t.writeDirectoryToTarGz(currentPath, tarWriter, subPath)
			if err != nil {
				return err
			}
		} else {
			fileInfo, err := file.Info()
			if err != nil {
				return err
			}
			err = t.writeFileToTarGz(currentPath, tarWriter, fileInfo, subPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *TarGzCompressor) writeFileToTarGz(path string, tarWriter *tar.Writer, fileInfo os.FileInfo, subPath string) error {
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
