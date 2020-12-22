package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// DownloadTextFile downloads a file from a given URL or return an error otherwise
func DownloadTextFile(URL string) (string, error) {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", fmt.Errorf("Received non 200 response code: %d", response.StatusCode)
	}
	contents := new(bytes.Buffer)
	_, err = io.Copy(contents, response.Body)
	if err != nil {
		return "", err
	}
	return contents.String(), nil
}
