package image

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ExtractImage(r io.Reader, dest string) error {
	tarReader := tar.NewReader(r)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("ExtractTar: Next() failed: %s", err.Error())
		}

		target := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("ExtractTar: Mkdir() failed: %s", err.Error())
			}
		case tar.TypeReg:
			outFile, err := os.Create(target)
			if err != nil {
				return fmt.Errorf("ExtractTar: Create() failed: %s", err.Error())
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("ExtractTar: Copy() failed: %s", err.Error())
			}
			outFile.Close()
			err = os.Chmod(target, 0755)
			if err != nil {
				return fmt.Errorf("ExtractTar: Chmod() failed: %s", err.Error())
			}
		case tar.TypeSymlink:
			err := os.Symlink(header.Linkname, target)
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
