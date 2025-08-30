package image

import (
	"path/filepath"
	"regexp"

	name "github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	remote "github.com/google/go-containerregistry/pkg/v1/remote"
	tarball "github.com/google/go-containerregistry/pkg/v1/tarball"
)

func RetrieveImage(imageName string, extract bool, basePath string) (string, *v1.ConfigFile, error) {
	// Use the provided basePath instead of current working directory
	dir := basePath

	ref, img, err := download(imageName)
	if err != nil {
		return "", nil, err
	}

	re := regexp.MustCompile(`(?:.+/)?([^:@]+)(?::.+)?`)
	matches := re.FindStringSubmatch(ref.String())
	if len(matches) != 2 {
		panic("Image name not correctly found")
	}

	baseName := matches[1]

	if !extract {
		// Save image tarball
		savePath := filepath.Join(dir, baseName+".tar")
		err = tarball.WriteToFile(savePath, ref, img)
		if err != nil {
			return "", nil, err
		}
		return savePath, nil, nil
	}

	// Extract image layers
	savePath := filepath.Join(dir, baseName)
	layers, err := img.Layers()
	if err != nil {
		return "", nil, err
	}

	for _, layer := range layers {
		r, err := layer.Uncompressed()
		if err != nil {
			return "", nil, err
		}

		err = ExtractImage(r, savePath)
		if err != nil {
			return "", nil, err
		}
	}

	configFile, err := img.ConfigFile()
	if err != nil {
		return "", nil, err
	}

	return savePath, configFile, nil
}

func download(imageName string) (name.Reference, v1.Image, error) {
	ref, err := name.ParseReference(imageName)
	if err != nil {
		return nil, nil, err
	}

	desc, err := remote.Get(ref)
	if err != nil {
		return nil, nil, err
	}

	img, err := desc.Image()
	if err != nil {
		return nil, nil, err
	}

	return ref, img, nil
}
