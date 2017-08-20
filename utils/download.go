package utils

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"

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

	if response.StatusCode >= http.StatusBadRequest {
		b, _ := ioutil.ReadAll(response.Body)
		return errors.Errorf("Http Error: %s", string(b))
	}

	// download to output file
	_, err = io.Copy(output, response.Body)
	if err != nil {
		return errors.Wrap(err, "Download error")
	}

	return nil
}
