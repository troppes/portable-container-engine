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
			if err := root.Mkdir(target, 0755); err != nil && !os.IsExist(err) {
				fmt.Printf("Warning: Failed to create directory %s: %s\n", target, err.Error())
				continue
			}
		case tar.TypeReg:
			if dir := filepath.Dir(target); dir != "." {
				if err := root.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
					fmt.Printf("Warning: Failed to create parent directory for %s: %s\n", target, err.Error())
					continue
				}
			}

			outFile, err := root.Create(target)
			if err != nil {
				fmt.Printf("Warning: Failed to create file %s: %s\n", target, err.Error())
				continue
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				fmt.Printf("Warning: Failed to copy content to %s: %s\n", target, err.Error())
				continue
			}

			err = outFile.Chmod(0755)
			if err != nil {
				fmt.Printf("Warning: Failed to set permissions on %s: %s\n", target, err.Error())
			}
			outFile.Close()

		case tar.TypeSymlink:
			// Remove existing file/symlink if it exists
			if _, err := root.Lstat(header.Name); err == nil {
				if err := root.Remove(header.Name); err != nil {
					fmt.Printf("Warning: Failed to remove existing file %s: %s\n", header.Name, err.Error())
					continue
				}
			}

			if err := root.Symlink(header.Linkname, header.Name); err != nil {
				fmt.Printf("Warning: Failed to create symlink %s -> %s: %s\n", header.Name, header.Linkname, err.Error())
				continue
			}

		case tar.TypeLink:
			if err := root.Link(header.Linkname, target); err != nil {
				fmt.Printf("Warning: Failed to create hard link %s -> %s: %s\n", target, header.Linkname, err.Error())
				continue
			}
		case tar.TypeFifo:
			fmt.Printf("Warning: Skipping FIFO: %s\n", header.Name)
		case tar.TypeChar:
			fmt.Printf("Warning: Skipping TypeChar: %s\n", header.Name)
		case tar.TypeBlock:
			fmt.Printf("Warning: Skipping TypeBlock: %s\n", header.Name)
		default:
			fmt.Printf("Warning: Skipping unknown file type %s in %s\n", string(header.Typeflag), header.Name)
		}
	}

	return nil
}
