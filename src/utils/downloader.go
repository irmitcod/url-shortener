package utils

import (
	"url-shortener/src/utils/rest_error"
	"io"
	"log"
	"net/http"
)

func DownloadFile(URL string) ([]byte, rest_error.RestErr) {
	// Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return nil, rest_error.NewBadRequestError("url not valid")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, rest_error.NewNotFoundError("urls not found")
	}

	file, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, rest_error.NewBadRequestError("Error while from response body")
	}

	return file, nil
}

type Result struct {
	Result []byte
	err    rest_error.RestErr
}

func DownloadFilePool(wid int, inputUrl <-chan string, results chan<- Result) {
	for url := range inputUrl {
		log.Printf("Worker ID: %d\n", wid)
		// Get the response bytes from the url
		response, err := http.Get(url)
		var res Result
		if err != nil {
			res = Result{
				Result: nil,
				err:    rest_error.NewBadRequestError("url not valid"),
			}
		} else if response.StatusCode != http.StatusOK {
			res = Result{
				Result: nil,
				err:    rest_error.NewNotFoundError("urls not found"),
			}
		} else {
			file, err := io.ReadAll(response.Body)
			if err != nil {
				res = Result{
					Result: nil,
					err:    rest_error.NewNotFoundError("urls not found"),
				}
			} else {
				res = Result{
					Result: file,
					err:    nil,
				}
			}
		}

		results <- res
	}

}
