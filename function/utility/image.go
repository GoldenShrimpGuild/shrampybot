package utility

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Image struct {
	Data     []byte
	Length   int
	Url      string
	MimeType string
	AltText  string
	Width    uint16
	Height   uint16
}

// Creates a new Image by substituting width/height in a twitch url string
// and fetching the bytes.
func NewFromThumbnailURL(url string, width uint16, height uint16, altText string) (*Image, error) {
	image := &Image{
		Url:     strings.Replace(url, "{width}x{height}", fmt.Sprintf("%vx%v", width, height), 1),
		Width:   width,
		Height:  height,
		AltText: altText,
	}

	resp, err := http.Get(image.Url)
	if err != nil || resp.StatusCode != 200 {
		return &Image{}, err
	}

	image.MimeType = resp.Header.Get("Content-type")
	image.Length, err = resp.Body.Read(image.Data)
	if err != nil {
		log.Printf("Could not retrieve stream preview image: %v\n", url)
		return &Image{}, err
	}

	return image, nil
}

func (i *Image) GetReader() *bytes.Reader {
	return bytes.NewReader(i.Data)
}
