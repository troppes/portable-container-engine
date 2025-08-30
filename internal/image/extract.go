package image

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/safeopen"
)

func ExtractImage(r io.Reader, dest string) error {
	tarReader := tar.NewReader(r)

	absDest, err := filepath.Abs(dest)
	if err != nil {
		return fmt.Errorf("ExtractTar: failed to get absolute path of destination: %w", err)
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("ExtractTar: Next() failed: %w", err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			dirPath := filepath.Join(absDest, header.Name)
			if err := os.MkdirAll(dirPath, 0o755); err != nil {
				return fmt.Errorf("ExtractTar: MkdirAll() failed: %w", err)
			}

		case tar.TypeReg:
			outFile, err := safeopen.CreateBeneath(absDest, header.Name)
			if err != nil {
				return fmt.Errorf("ExtractTar: CreateBeneath() failed: %w", err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("ExtractTar: Copy() failed: %w", err)
			}

			if err := outFile.Chmod(0o755); err != nil {
				outFile.Close()
				return fmt.Errorf("ExtractTar: Chmod() failed: %w", err)
			}

			outFile.Close()

		case tar.TypeSymlink:
			symlinkLocation := filepath.Join(absDest, header.Name)
			linkTarget := filepath.Clean(header.Linkname)

			resolvedTarget := filepath.Join(filepath.Dir(symlinkLocation), linkTarget)
			absResolvedTarget, err := filepath.Abs(resolvedTarget)
			if err != nil {
				return fmt.Errorf("ExtractTar: failed to resolve symlink target: %w", err)
			}

			destSeparator := absDest + string(filepath.Separator)
			if !strings.HasPrefix(absResolvedTarget, destSeparator) && absResolvedTarget != absDest {
				return fmt.Errorf("ExtractTar: symlink target escapes destination directory: %s -> %s",
					header.Name, header.Linkname)
			}

			if err := os.Symlink(header.Linkname, symlinkLocation); err != nil {
				return fmt.Errorf("ExtractTar: Symlink() failed: %w", err)
			}

		default:
			return fmt.Errorf("ExtractTar: unsupported file type: %c in %s",
				header.Typeflag, header.Name)
		}
	}

	return nil
}
