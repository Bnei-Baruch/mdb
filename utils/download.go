package utils

import (
	"os"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

func DownloadUrl(url string, path string) error {
	// create output file
	output, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "Create output file")
	}
	defer output.Close()

	// do http call
	response, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "Http error")
	}
	defer response.Body.Close()

	// download to output file
	_, err = io.Copy(output, response.Body)
	if err != nil {
		return errors.Wrap(err, "Download error")
	}

	return nil
}
