package tmdb

import (
	"io/ioutil"
	"net/http"
)

type HttpGetter func(string) (*http.Response, error)

func getBody(url string, getter HttpGetter) ([]byte, error) {
	resp, err := getter(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
