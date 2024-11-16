package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// Archiver 接口定义了解压方法
type Archiver interface {
	Extract(targetPath string, path string, maxConcurrency int) error
}

// Zip 结构体实现了 Archiver 接口
type Zip struct{}

func (z *Zip) Extract(targetPath string, path string, maxConcurrency int) error {
	zipReader, err := zip.OpenReader(path)
	if err != nil {
		return fmt.Errorf("failed to read zip file: %w", err)
	}
	defer zipReader.Close()

	// Calculate total size of all files in the ZIP archive
	var totalSize int64
	for _, file := range zipReader.File {
		totalSize += file.FileInfo().Size()
	}

	bar := progressbar.DefaultBytes(
		totalSize,
		fmt.Sprintf("Extract %s\n", filepath.Base(path)),
	)
	// Semaphore for limiting concurrency
	sem := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	for _, file := range zipReader.File {
		wg.Add(1)
		sem <- struct{}{} // acquire a token

		go func(file *zip.File) {
			defer wg.Done()
			defer func() { <-sem }() // release the token when done

			// 获取每个文件的目标路径
			targetFilePath := filepath.Join(targetPath, file.Name)
			if file.FileInfo().IsDir() {
				// 如果是目录，则创建目录
				if err := os.MkdirAll(targetFilePath, file.Mode()); err != nil {
					fmt.Printf("failed to create directory: %v\n", err)
				}
			} else {
				if err := os.MkdirAll(filepath.Dir(targetFilePath), os.ModePerm); err != nil {
					fmt.Printf("failed to create directory: %v\n", err)
				}
				// 解压文件
				if err := extractZipFile(file, targetFilePath, bar); err != nil {
					fmt.Printf("failed to extract zip file: %v\n", err)
				}
			}
		}(file)
	}

	wg.Wait() // wait for all goroutines to finish
	return nil
}

func extractZipFile(file *zip.File, targetPath string, bar *progressbar.ProgressBar) error {
	inFile, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open zip file inside: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	// Copy the content to the destination file
	n, err := io.Copy(outFile, inFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Update progress after extraction
	return bar.Add64(n)
}

// TarGz 结构体实现了 Archiver 接口
type TarGz struct{}

func (t *TarGz) Extract(targetPath string, path string, maxConcurrency int) error {
	reader, err := os.Open(path)
	if err != nil {
		return err
	}
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	// Calculate total size of all files in the TAR archive
	var totalSize int64
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar file: %w", err)
		}
		totalSize += header.Size
	}

	bar := progressbar.DefaultBytes(
		totalSize,
		fmt.Sprintf("Extract %s\n", filepath.Base(path)),
	)
	// Semaphore for limiting concurrency
	sem := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar file: %w", err)
		}

		wg.Add(1)
		sem <- struct{}{} // acquire a token

		go func(header *tar.Header) {
			defer wg.Done()
			defer func() { <-sem }() // release the token when done

			targetFilePath := filepath.Join(targetPath, header.Name)
			switch header.Typeflag {
			case tar.TypeDir:
				if err := os.MkdirAll(targetFilePath, os.ModePerm); err != nil {
					fmt.Printf("failed to create directory: %v\n", err)
				}
			case tar.TypeReg:
				if err := extractTarFile(tarReader, targetFilePath, header, bar); err != nil {
					fmt.Printf("failed to extract tar file: %v\n", err)
				}
			}
		}(header)
	}

	wg.Wait() // wait for all goroutines to finish
	return nil
}

func extractTarFile(tarReader io.Reader, targetPath string, _ *tar.Header, bar *progressbar.ProgressBar) error {
	outFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	n, err := io.Copy(outFile, tarReader)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	return bar.Add64(n)
}

// 工厂函数，根据压缩类型返回对应的 Archiver 实现
func NewArchiver(archiveType string) (Archiver, error) {
	switch archiveType {
	case "zip":
		return &Zip{}, nil
	case "tar.gz":
		return &TarGz{}, nil
	default:
		return nil, fmt.Errorf("unsupported archive type: %s", archiveType)
	}
}
