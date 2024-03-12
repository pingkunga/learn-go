package http

import (
	"encoding/json"
	"io"
	"net/http"
)

type Response struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
	Info string `json:"info"`
}

func MaketHttpCall(url string) (*Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := &Response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return response, err
}
