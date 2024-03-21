package http

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"
)

// Mock httpserver with real server
var get = http.Get

func Tes1tMarkHttpWithRealServer(t *testing.T) {
	get = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{"id": 1, "name": "PingkungA", "age": 33, "info": "dotnet dev/ blogger"}`)),
		}, nil
	}

	want := &Response{
		ID:   1,
		Name: "PingkungA",
		Age:  33,
		Info: "dotnet dev/ blogger",
	}

	t.Run("Success Real Server response", func(t *testing.T) {

		resp, err := MaketHttpCall("http://localhost:8080")
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(resp, want) {
			t.Errorf("expected (%v), got (%v)", want, resp)
		}

		if !errors.Is(err, nil) {
			t.Errorf("expected (%v), got (%v)", nil, err)
		}
	})
}
