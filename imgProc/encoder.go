package imgProc

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
)

func GetImage(imageName string) image.Image {
	desiredImage := fmt.Sprintf("./imgProc/stock/%v", imageName)

	bytesRes, err := os.ReadFile(desiredImage)
	if err != nil {
		panic(err)
	}

	imageFormat := http.DetectContentType(bytesRes)
	// ImageBase64 := base64.StdEncoding.EncodeToString(bytesRes)

	var imageInf image.Image

	switch imageFormat {
	case "image/jpeg":
		imageInf, err = jpeg.Decode(bytes.NewReader(bytesRes))
		if err != nil {
			panic(err)
		}
	case "image/png":
		imageInf, err = png.Decode(bytes.NewReader(bytesRes))
		if err != nil {
			panic(err)
		}
	}

	return imageInf
}
