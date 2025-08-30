package image

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ExtractImage(r io.Reader, dest string) error {
	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("ExtractTar: MkdirAll() failed: %s", err.Error())
	}

	root, err := os.OpenRoot(dest)
	if err != nil {
		return fmt.Errorf("ExtractTar: OpenRoot() failed: %s", err.Error())
	}

	defer root.Close()

	tarReader := tar.NewReader(r)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("ExtractTar: Next() failed: %s", err.Error())
		}

		target := filepath.Clean(header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := root.Mkdir(target, 0755); err != nil {
				return fmt.Errorf("ExtractTar: Mkdir() failed: %s", err.Error())
			}
		case tar.TypeReg:
			if dir := filepath.Dir(target); dir != "." {
				if err := root.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
					return fmt.Errorf("ExtractTar: Mkdir() failed: %s", err.Error())
				}
			}

			outFile, err := root.Create(target)
			if err != nil {
				return fmt.Errorf("ExtractTar: Create() failed: %s", err.Error())
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("ExtractTar: Copy() failed: %s", err.Error())
			}
			outFile.Close()

		case tar.TypeSymlink:
			if err := root.Symlink(header.Linkname, header.Name); err != nil {
				return fmt.Errorf("ExtractTar: Symlink() failed: %w", err)
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
