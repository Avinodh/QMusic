package main

// TODO: separate package for testing

import (
	"fmt"
	//"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var (
	server *httptest.Server
	//reader   io.Reader
)

func init() {
	// Create server for testing
	server = httptest.NewServer(NewRouter())
}

// func TestMain(m *testing.M) {
// 	fmt.Println("Testing initiated")
// 	os.Exit(m.Run())
// }

func TestAuthSpotify(t *testing.T) {
	fmt.Println("Starting TestAuthSpotify")
	//t.log("Starting TestAuthSpotify")

	authSpotifyUrl := fmt.Sprintf("%s/authspotify", server.URL)

	reader := strings.NewReader("")

	req, err := http.NewRequest("GET", authSpotifyUrl, reader)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}

func TestCreatePartyController(t *testing.T) {

	fmt.Println("Starting TestCreatePartyController")

	createPartyUrl := fmt.Sprintf("%s/dashboard/createparty", server.URL)

	form := url.Values{}
	form.Add("user", "my_user")
	form.Add("location", "my_location")
	form.Add("secret-code", "my_code")
	form.Add("active-time", "my_time")

	reader := strings.NewReader(form.Encode())

	req, err := http.NewRequest("POST", createPartyUrl, reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}

}

func TestSearchSong(t *testing.T) {

	fmt.Println("Starting TestSearchSong")

	searchSongUrl := fmt.Sprintf("%s/searchsong", server.URL)

	form := url.Values{}
	form.Add("searchsong", "halo")

	reader := strings.NewReader(form.Encode())

	req, err := http.NewRequest("GET", searchSongUrl, reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}
