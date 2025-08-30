package image

import (
	"archive/tar"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractImage(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "pce-extract-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name      string
		tarFiles  []tarFile
		wantErr   bool
		checkFunc func(t *testing.T, dir string)
	}{
		{
			name: "Basic files and directories",
			tarFiles: []tarFile{
				{
					name:     "testdir",
					typeflag: tar.TypeDir,
					mode:     0755,
				},
				{
					name:     "testdir/file1.txt",
					typeflag: tar.TypeReg,
					content:  []byte("test content"),
					mode:     0644,
				},
			},
			checkFunc: func(t *testing.T, dir string) {
				// Check directory was created
				dirInfo, err := os.Stat(filepath.Join(dir, "testdir"))
				if err != nil {
					t.Errorf("failed to stat directory: %v", err)
					return
				}
				if !dirInfo.IsDir() {
					t.Error("testdir is not a directory")
				}

				// Check file was created with correct content
				content, err := os.ReadFile(filepath.Join(dir, "testdir/file1.txt"))
				if err != nil {
					t.Errorf("failed to read file: %v", err)
					return
				}
				if string(content) != "test content" {
					t.Errorf("file content mismatch, got %q, want %q", string(content), "test content")
				}
			},
		},
		{
			name: "Symlink handling",
			tarFiles: []tarFile{
				{
					name:     "target.txt",
					typeflag: tar.TypeReg,
					content:  []byte("target content"),
					mode:     0644,
				},
				{
					name:     "link.txt",
					typeflag: tar.TypeSymlink,
					linkname: "target.txt",
				},
			},
			checkFunc: func(t *testing.T, dir string) {
				// Check symlink points to correct file
				link := filepath.Join(dir, "link.txt")
				target, err := os.Readlink(link)
				if err != nil {
					t.Errorf("failed to read symlink: %v", err)
					return
				}
				if target != "target.txt" {
					t.Errorf("symlink points to %q, want %q", target, "target.txt")
				}

				// Check content through symlink
				content, err := os.ReadFile(link)
				if err != nil {
					t.Errorf("failed to read through symlink: %v", err)
					return
				}
				if string(content) != "target content" {
					t.Errorf("content through symlink mismatch, got %q, want %q", string(content), "target content")
				}
			},
		},
		{
			name: "Path traversal attempt",
			tarFiles: []tarFile{
				{
					name:     "../outside.txt",
					typeflag: tar.TypeReg,
					content:  []byte("should not exist"),
					mode:     0644,
				},
			},
			wantErr: true,
		},
		{
			name: "Deep directory structure",
			tarFiles: []tarFile{
				{
					name:     "dir1",
					typeflag: tar.TypeDir,
					mode:     0755,
				},
				{
					name:     "dir1/dir2",
					typeflag: tar.TypeDir,
					mode:     0755,
				},
				{
					name:     "dir1/dir2/dir3",
					typeflag: tar.TypeDir,
					mode:     0755,
				},
				{
					name:     "dir1/dir2/dir3/file.txt",
					typeflag: tar.TypeReg,
					content:  []byte("deep file"),
					mode:     0644,
				},
			},
			checkFunc: func(t *testing.T, dir string) {
				path := filepath.Join(dir, "dir1/dir2/dir3/file.txt")
				content, err := os.ReadFile(path)
				if err != nil {
					t.Errorf("failed to read deep file: %v", err)
					return
				}
				if string(content) != "deep file" {
					t.Errorf("deep file content mismatch, got %q, want %q", string(content), "deep file")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test directory for this specific test
			testDir := filepath.Join(tmpDir, tt.name)
			if err := os.MkdirAll(testDir, 0755); err != nil {
				t.Fatalf("failed to create test directory: %v", err)
			}

			var buf bytes.Buffer
			tarWriter := tar.NewWriter(&buf)

			// Add files to tar archive
			for _, tf := range tt.tarFiles {
				if err := writeTarEntry(tarWriter, tf); err != nil {
					t.Fatalf("failed to write tar entry: %v", err)
				}
			}
			tarWriter.Close()

			err := ExtractImage(&buf, testDir)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Run custom checks
			if tt.checkFunc != nil {
				tt.checkFunc(t, testDir)
			}
		})
	}
}

type tarFile struct {
	name     string
	typeflag byte
	mode     int64
	content  []byte
	linkname string
}

func writeTarEntry(tw *tar.Writer, tf tarFile) error {
	header := &tar.Header{
		Name:     tf.name,
		Mode:     tf.mode,
		Size:     int64(len(tf.content)),
		Typeflag: tf.typeflag,
		Linkname: tf.linkname,
	}

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	if tf.typeflag == tar.TypeReg {
		if _, err := tw.Write(tf.content); err != nil {
			return err
		}
	}

	return nil
}
