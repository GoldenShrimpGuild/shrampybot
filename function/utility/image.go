package utility

import (
	"bytes"
	"fmt"
	"io"
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
	Width    int
	Height   int
}

// Creates a new Image by substituting width/height in a twitch url string
// and fetching the bytes.
func NewFromThumbnailURL(url string, width int, height int, altText string) (*Image, error) {
	image := &Image{
		Url:     strings.Replace(url, "{width}x{height}", fmt.Sprintf("%vx%v", width, height), 1),
		Width:   width,
		Height:  height,
		AltText: altText,
	}

	resp, err := http.Get(image.Url)
	if err != nil || resp.StatusCode > 299 {
		return &Image{}, err
	}

	image.MimeType = resp.Header.Get("Content-type")

	image.Data, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Could not retrieve stream preview image: %v\n", url)
		resp.Body.Close()
		return &Image{}, err
	}
	resp.Body.Close()
	return image, nil
}

func (i *Image) GetReader() *bytes.Reader {
	return bytes.NewReader(i.Data)
}
