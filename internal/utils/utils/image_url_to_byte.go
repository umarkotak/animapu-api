package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"

	"github.com/sirupsen/logrus"
)

func ImageUrlToJpegByte(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Error(fmt.Errorf("failed to get image: status code %d", resp.StatusCode))
		return nil, fmt.Errorf("failed to get image: status code %d", resp.StatusCode)
	}

	// Decode the image
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// Encode the image into the buffer as PNG
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 100})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return buf.Bytes(), nil
}
