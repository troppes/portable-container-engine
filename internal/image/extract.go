package image

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/safeopen"
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

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(filepath.Join(dest, filepath.Dir(header.Name)), 0755); err != nil {
				return fmt.Errorf("ExtractTar: Mkdir() failed: %s", err.Error())
			}
			if err := os.MkdirAll(filepath.Join(dest, header.Name), 0755); err != nil {
				return fmt.Errorf("ExtractTar: Mkdir() failed: %s", err.Error())
			}

		case tar.TypeReg:
			outFile, err := safeopen.CreateBeneath(dest, header.Name)
			if err != nil {
				return fmt.Errorf("ExtractTar: CreateBeneath() failed: %s", err.Error())
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("ExtractTar: Copy() failed: %s", err.Error())
			}

			err = outFile.Chmod(0755)
			outFile.Close()
			if err != nil {
				return fmt.Errorf("ExtractTar: Chmod() failed: %s", err.Error())
			}

		case tar.TypeSymlink:
			target := filepath.Join(dest, header.Name)
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
