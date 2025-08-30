package image

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	util "github.com/troppes/portable-container-engine/internal/util"
)

func TestDownload(t *testing.T) {
	tests := []struct {
		name        string
		imageName   string
		wantErr     bool
		errContains string
	}{
		{
			name:      "Valid image name",
			imageName: "alpine:latest",
			wantErr:   false,
		},
		{
			name:        "Invalid image name",
			imageName:   "not/a/valid:image",
			wantErr:     true,
			errContains: "UNAUTHORIZED",
		},
		{
			name:        "Empty image name",
			imageName:   "",
			wantErr:     true,
			errContains: "could not parse reference",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref, img, err := download(tt.imageName)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				} else if tt.errContains != "" && !util.Contains(err.Error(), tt.errContains) {
					t.Errorf("error message %q does not contain %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if ref == nil {
				t.Error("expected reference to not be nil")
			}

			if img == nil {
				t.Error("expected image to not be nil")
			}

			if !util.Contains(ref.String(), tt.imageName) {
				t.Errorf("reference %v does not contain input image name %v", ref.String(), tt.imageName)
			}
		})
	}
}

func TestRetrieveImage(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "pce-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name        string
		imageName   string
		extract     bool
		wantErr     bool
		errContains string
		checkFile   func(t *testing.T, path string)
	}{
		{
			name:      "Download alpine as tar",
			imageName: "alpine:latest",
			extract:   false,
			wantErr:   false,
			checkFile: func(t *testing.T, path string) {
				// Check if file exists and ends with .tar
				if !strings.HasSuffix(path, ".tar") {
					t.Errorf("expected tar file, got: %s", path)
				}
				if _, err := os.Stat(path); os.IsNotExist(err) {
					t.Errorf("file %s does not exist", path)
				}
			},
		},
		{
			name:      "Download and Extract alpine image",
			imageName: "alpine:latest",
			extract:   true,
			wantErr:   false,
			checkFile: func(t *testing.T, path string) {
				if _, err := os.Stat(filepath.Join(path, "bin")); os.IsNotExist(err) {
					t.Errorf("extracted image missing /bin directory")
				}
				if _, err := os.Stat(filepath.Join(path, "etc")); os.IsNotExist(err) {
					t.Errorf("extracted image missing /etc directory")
				}
			},
		},
		{
			name:        "Invalid image name",
			imageName:   "not/a/valid:image",
			extract:     false,
			wantErr:     true,
			errContains: "UNAUTHORIZED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, _, err := RetrieveImage(tt.imageName, tt.extract, tmpDir)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				} else if tt.errContains != "" && !util.Contains(err.Error(), tt.errContains) {
					t.Errorf("error message %q does not contain %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if path == "" {
				t.Error("expected path to not be empty")
				return
			}

			if tt.checkFile != nil {
				tt.checkFile(t, path)
			}
		})
	}
}
