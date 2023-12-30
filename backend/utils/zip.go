package utils

import (
	"archive/zip"
	"io"
	"os"

	"github.com/pkg/errors"
)

func AddFileToZip(writer *zip.Writer, path string, zipPath string) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return errors.Wrap(err, "failed to get file info")
	}

	if fileInfo.IsDir() {
		return errors.New("file is a directory")
	}

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return errors.Wrap(err, "failed to create header")
	}

	header.Method = zip.Deflate
	header.Name = zipPath

	fileWriter, err := writer.CreateHeader(header)
	if err != nil {
		return errors.Wrap(err, "failed to create file writer")
	}

	_, err = io.Copy(fileWriter, file)
	return errors.Wrap(err, "failed to copy file")
}
