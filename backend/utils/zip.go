package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

func AddFileToZip(writer *zip.Writer, path string, zipPath string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("file is a directory")
	}

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return fmt.Errorf("failed to create header: %w", err)
	}

	header.Method = zip.Deflate
	header.Name = zipPath

	fileWriter, err := writer.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("failed to create file writer: %w", err)
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}
	return nil
}
