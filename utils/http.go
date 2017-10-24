package utils

import (
	"encoding/json"
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

// Send http POST request.
// Note that json encoding errors are ignored. We expect a valid payload
// based on https://www.zupzup.org/io-pipe-go/
func HttpPostJson(url string, payload interface{}) (*http.Response, error) {
	pr, pw := io.Pipe()

	go func() {
		// close the writer, so the reader knows there's no more data
		defer pw.Close()

		// write json data to the PipeReader through the PipeWriter
		json.NewEncoder(pw).Encode(&payload)
	}()

	// JSON from the PipeWriter lands in the PipeReader
	// ...and we send it off...
	var resp *http.Response
	resp, err := http.Post(url, "application/json", pr)
	if err != nil {
		return nil, errors.Wrap(err, "http.Post")
	}

	return resp, nil
}
