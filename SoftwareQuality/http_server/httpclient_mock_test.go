package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func HandlerMyAPI(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"id": 1, "name": "PingkungA", "age": 33, "info": "dotnet dev/ blogger"}`))
}

func TestMakeHttp(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(HandlerMyAPI))

	want := &Response{
		ID:   1,
		Name: "PingkungA",
		Age:  33,
		Info: "dotnet dev/ blogger",
	}

	t.Run("Sccess Server response", func(t *testing.T) {
		defer server.Close()

		resp, err := MaketHttpCall(server.URL)

		if !reflect.DeepEqual(resp, want) {
			t.Errorf("expected (%v), got (%v)", want, resp)
		}

		if !errors.Is(err, nil) {
			t.Errorf("expected (%v), got (%v)", nil, err)
		}
	})
}
