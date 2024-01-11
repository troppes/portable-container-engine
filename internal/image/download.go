package image

import (
	"log"
	"os"
	"regexp"

	name "github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	remote "github.com/google/go-containerregistry/pkg/v1/remote"
	tarball "github.com/google/go-containerregistry/pkg/v1/tarball"
)

func DownloadImage(imageName string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	ref, img, err := download(imageName)

	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`([^/]*):`) // match everything after the last slash and before the tag
	matches := re.FindStringSubmatch(ref.String())

	if len(matches) != 2 {
		panic("Image name not correctly found")
	}

	savePath := dir + "/" + matches[1] + ".tar"

	err = tarball.WriteToFile(savePath, ref, img)
	if err != nil {
		return "", err
	}

	return savePath, nil
}

func DownloadAndExtract(imageName string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	ref, img, err := download(imageName)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`.*/(.+):`)
	matches := re.FindStringSubmatch(ref.String())

	if len(matches) != 2 {
		panic("Image name not found")
	}

	savePath := dir + "/" + matches[1]

	layers, err := img.Layers()
	if err != nil {
		log.Fatal(err)
	}

	for _, layer := range layers {
		r, err := layer.Uncompressed()
		if err != nil {
			log.Fatal(err)
		}

		err = ExtractImage(r, savePath)
		if err != nil {
			log.Fatal(err)
		}
	}

	return savePath, nil
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
