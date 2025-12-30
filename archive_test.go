package main

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"
)

func TestTarGz_Extract(t *testing.T) {
	// 1. Create a temporary tar.gz file
	tempDir := t.TempDir()
	tarFile := filepath.Join(tempDir, "test.tar.gz")
	
	f, err := os.Create(tarFile)
	if err != nil {
		t.Fatal(err)
	}
	
	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)
	
	files := map[string]string{
		"file1.txt": "content1",
		"dir/file2.txt": "content2",
	}
	
	for name, content := range files {
		hdr := &tar.Header{
			Name: name,
			Mode: 0600,
			Size: int64(len(content)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	
	tw.Close()
	gw.Close()
	f.Close()
	
	// 2. Extract it
	extractDir := filepath.Join(tempDir, "extracted")
	archiver := &TarGz{}
	err = archiver.Extract(extractDir, tarFile, 1)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}
	
	// 3. Verify contents
	for name, content := range files {
		path := filepath.Join(extractDir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("Failed to read extracted file %s: %v", name, err)
			continue
		}
		if string(data) != content {
			t.Errorf("File %s content mismatch. Want %s, got %s", name, content, string(data))
		}
	}
}
