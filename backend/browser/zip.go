package browser

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ZipCompressor struct{}

func NewZipCompressor() *ZipCompressor {
	return &ZipCompressor{}
}

func (z *ZipCompressor) CompressDirectory(dir string, w io.Writer) error {
	zipWriter := zip.NewWriter(w)

	root, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return err
	}

	if err := z.writeDirectoryToZip(root, root, zipWriter); err != nil {
		_ = zipWriter.Close()
		return err
	}

	return zipWriter.Close()
}

func (z *ZipCompressor) writeDirectoryToZip(directory string, root string, zipWriter *zip.Writer) error {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		currentPath := filepath.Join(directory, entry.Name())

		if entry.IsDir() {
			if err := z.writeDirectoryToZip(currentPath, root, zipWriter); err != nil {
				return err
			}
			continue
		}

		fileInfo, err := entry.Info()
		if err != nil {
			return err
		}

		if fileInfo.Mode()&os.ModeSymlink != 0 {
			resolvedPath, err := filepath.EvalSymlinks(currentPath)
			if err != nil {
				return err
			}

			resolvedInfo, err := os.Stat(resolvedPath)
			if err != nil {
				return err
			}

			if resolvedInfo.IsDir() {
				if err := z.writeDirectoryToZip(resolvedPath, root, zipWriter); err != nil {
					return err
				}
				continue
			}

			currentPath = resolvedPath
			fileInfo = resolvedInfo
		}

		if !fileInfo.Mode().IsRegular() {
			continue
		}

		if err := z.writeFileToZip(currentPath, root, zipWriter, fileInfo); err != nil {
			return err
		}
	}

	return nil
}

func (z *ZipCompressor) writeFileToZip(path string, root string, zipWriter *zip.Writer, fileInfo os.FileInfo) error {
	relPath, err := filepath.Rel(root, path)
	if err != nil {
		return err
	}

	if relPath == "." || relPath == ".." || strings.HasPrefix(relPath, ".."+string(filepath.Separator)) {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}

	header.Name = filepath.ToSlash(relPath)
	header.Method = zip.Deflate

	zipFileWriter, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(zipFileWriter, file)
	return err
}

func (z *ZipCompressor) GetFileExtension() string {
	return ".zip"
}
