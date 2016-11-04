package main

import (
    "fmt"
    "io"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

var (
    server   *httptest.Server
    reader   io.Reader
    usersUrl string
)

func init() {
	//Creating new server with the user handlers
    server = httptest.NewServer(NewRouter())
    //Grab the address for the API endpoint
    usersUrl = fmt.Sprintf("%s/authspotify", server.URL) 
}

func TestAuthSpotify(t *testing.T) {
	reader = strings.NewReader("")

	request, err := http.NewRequest("GET", usersUrl, reader)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}


