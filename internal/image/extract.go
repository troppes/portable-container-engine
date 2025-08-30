package image

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ExtractImage(r io.Reader, dest string) error {
	tarReader := tar.NewReader(r)

	dest, err := filepath.Abs(dest)
	if err != nil {
		return fmt.Errorf("ExtractTar: failed to get absolute path: %s", err.Error())
	}

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("ExtractTar: Next() failed: %s", err.Error())
		}

		target, err := validatePath(dest, header.Name)
		if err != nil {
			return fmt.Errorf("ExtractTar: invalid path %s: %s", header.Name, err.Error())
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("ExtractTar: Mkdir() failed: %s", err.Error())
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("ExtractTar: MkdirAll() failed: %s", err.Error())
			}

			outFile, err := os.Create(target)
			if err != nil {
				return fmt.Errorf("ExtractTar: Create() failed: %s", err.Error())
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("ExtractTar: Copy() failed: %s", err.Error())
			}
			outFile.Close()
			err = os.Chmod(target, 0755)
			if err != nil {
				return fmt.Errorf("ExtractTar: Chmod() failed: %s", err.Error())
			}
		case tar.TypeSymlink:
			linkTarget, err := validateSymlinkTarget(dest, target, header.Linkname)
			if err != nil {
				return fmt.Errorf("ExtractTar: invalid symlink target %s: %s", header.Linkname, err.Error())
			}

			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("ExtractTar: MkdirAll() failed: %s", err.Error())
			}

			err = os.Symlink(linkTarget, target)
			if err != nil {
				return fmt.Errorf("ExtractTar: Symlink() failed: %s", err.Error())
			}
		default:
			return fmt.Errorf(
				"ExtractTar: unknown type: %s in %s",
				string(header.Typeflag),
				header.Name)
		}
	}

	return nil
}

// validatePath ensures the target path is within the destination directory
func validatePath(dest, path string) (string, error) {
	// Clean the path to resolve any .. or . elements
	cleanPath := filepath.Clean(path)

	// Check for absolute paths or paths starting with ..
	if filepath.IsAbs(cleanPath) || strings.HasPrefix(cleanPath, "..") {
		return "", fmt.Errorf("path attempts to escape destination directory")
	}

	// Join with destination and get absolute path
	target := filepath.Join(dest, cleanPath)
	target, err := filepath.Abs(target)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %s", err.Error())
	}

	// Ensure the target is within the destination directory
	if !strings.HasPrefix(target, dest+string(filepath.Separator)) && target != dest {
		return "", fmt.Errorf("path attempts to escape destination directory")
	}

	return target, nil
}

// validateSymlinkTarget ensures the symlink target is safe
func validateSymlinkTarget(dest, symlinkPath, linkTarget string) (string, error) {
	// If absolute path, reject it
	if filepath.IsAbs(linkTarget) {
		return "", fmt.Errorf("absolute symlink targets are not allowed")
	}

	// Clean the link target
	cleanTarget := filepath.Clean(linkTarget)

	// Calculate what the symlink would resolve to
	symlinkDir := filepath.Dir(symlinkPath)
	resolvedTarget := filepath.Join(symlinkDir, cleanTarget)
	resolvedTarget, err := filepath.Abs(resolvedTarget)
	if err != nil {
		return "", fmt.Errorf("failed to resolve symlink target: %s", err.Error())
	}

	// Fail if the resolved target is outside the destination directory
	if !strings.HasPrefix(resolvedTarget, dest+string(filepath.Separator)) && resolvedTarget != dest {
		return "", fmt.Errorf("symlink target attempts to escape destination directory")
	}

	return cleanTarget, nil
}
