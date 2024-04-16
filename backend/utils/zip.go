package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func ExtractZip(zipPath string, dst string) error {
	archive, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Wrap in a function to ensure defer is called before the next iteration
		err = func() error {
			dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer dstFile.Close()

			fileInArchive, err := f.Open()
			if err != nil {
				return fmt.Errorf("failed to open file in archive: %w", err)
			}
			defer fileInArchive.Close()

			if _, err := io.Copy(dstFile, fileInArchive); err != nil {
				return fmt.Errorf("failed to copy file: %w", err)
			}
			return nil
		}()

		if err != nil {
			return err
		}
	}
	return nil
}
