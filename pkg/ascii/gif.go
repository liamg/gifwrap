package ascii

import (
	"crypto/tls"
	"image/gif"
	"io"
	"net/http"
	"os"
	"time"
)

func FromURL(url string, skipTLSVerify bool) (*Renderer, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipTLSVerify},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 30,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return FromReader(resp.Body)
}

func FromFile(path string) (*Renderer, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return FromReader(file)
}

func FromReader(reader io.Reader) (*Renderer, error) {
	img, err := gif.DecodeAll(reader)
	if err != nil {
		return nil, err
	}

	return &Renderer{
		image: img,
	}, nil
}
