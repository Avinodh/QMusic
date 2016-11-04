package main

/************** Functions that handle all routes/ Defines our API *************/

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	_ "github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/zmb3/spotify"
	"io/ioutil"
	"net/http"
	_"net/url"
	"os"
	"strings"
)

var (
	store       = sessions.NewCookieStore([]byte(os.Getenv("COOKIE_STORE_AUTHENTICATION_KEY")), []byte(os.Getenv("COOKIE_STORE_ENCRYPTION_KEY")))
	redirectURL = os.Getenv("REDIRECT_URL")
	auth        = spotify.NewAuthenticator(redirectURL, spotify.ScopeUserReadPrivate, spotify.ScopePlaylistModifyPublic)
)

func init() {
	store.Options = &sessions.Options{
		Path: "/",
		//MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
	}
}

func Index(rw http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadFile("www/index.html")
	fmt.Fprint(rw, string(body))
}

func AuthorizeSpotify(rw http.ResponseWriter, r *http.Request) {
	// Generating Session Token
	session, _ := store.Get(r, "UserAuthStore")

	session_token, _ := GenerateRandomString(32)

	// Setting Session Token
	session.Values["session_token"] = &session_token
	err := sessions.Save(r, rw)
	if err != nil {
		fmt.Fprint(rw, "Session not saved! ", err)
	}

	url := auth.AuthURL(session_token)

	fmt.Fprint(rw, url)
}

func Dashboard(rw http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "UserAuthStore")
	//session_token := fmt.Sprint(session.Values["session_token"])

	// use the same state string here that you used to generate the URL
	r.ParseForm()
	code := r.FormValue("code")

	authTokenUrl := fmt.Sprintf("https://accounts.spotify.com/api/token")

	// POST request to fetch Access Token
	params := fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s", code, redirectURL, os.Getenv("SPOTIFY_ID"), os.Getenv("SPOTIFY_SECRET"))
	reader := strings.NewReader(params)
	req, err := http.NewRequest("POST", authTokenUrl, reader)
	if err != nil {
		// handle err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := http.DefaultClient.Do(req)

	respBody, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(respBody, &Spotify_Auth_Object)

	session.Values["spotify_access_token"] = &Spotify_Auth_Object.AccessToken

	/******* Fetch user's Spotify ID *********/
	authHeader := fmt.Sprintf("Bearer %s", Spotify_Auth_Object.AccessToken)
	getTrackUrl := fmt.Sprintf("https://api.spotify.com/v1/me")

	httpClient := &http.Client{}
	req, _ = http.NewRequest("GET", getTrackUrl, nil)
	req.Header.Set("Authorization", authHeader)

	res, _ := httpClient.Do(req)
	resBody, _ := ioutil.ReadAll(res.Body)

	// Spotify_User_Object now contains User ID, Display Name, and Profile Picture URL
	err = json.Unmarshal([]byte(resBody), &Spotify_User_Object)
	if err != nil {
		panic(err)
	}
	/******************************************/

	body, _ := ioutil.ReadFile("www/dashboard.html")
	fmt.Fprint(rw, string(body))
}

func RenderSearch(rw http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadFile("www/search.html")
	fmt.Fprint(rw, string(body))
}

func CreatePartyController(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	playlistName := r.FormValue("user") // Actually the party/playlist name

	/********* Create a new Spotify Playlist ********/
	authHeader := fmt.Sprintf("Bearer %s", Spotify_Auth_Object.AccessToken)
	/*data := url.Values{}
	data.Set("name", playlistName)*/
	httpClient := &http.Client{}

	createPlaylistUrl := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists",Spotify_User_Object.Id)
	params := fmt.Sprintf("{\"name\":\"%s\"}", playlistName)
	reader := strings.NewReader(params)

	req, _ := http.NewRequest("POST", createPlaylistUrl, reader)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")
	res, _ := httpClient.Do(req)
	resBody, _ := ioutil.ReadAll(res.Body)

	/***** Fetch newly created Playlist details *****/

	var playlist Playlist
	err := json.Unmarshal([]byte(resBody), &playlist)
	if err != nil {
		panic(err)
	}

	/************************************************/

	pc = TheMasterController.AddPartyController(r.Form["secret-code"][0])
	pc.CreateParty(r, playlist.Id)
	
	/*body, _ := ioutil.ReadFile("www/search.html")
	fmt.Fprint(rw, string(body))*/
	 http.Redirect(rw, r, "/search", http.StatusSeeOther)
}

func SearchSong(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var query = r.FormValue("searchsong")

	getTrackUrl := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track", query)

	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", getTrackUrl, nil)
	res, _ := httpClient.Do(req)
	resBody, _ := ioutil.ReadAll(res.Body)

	//var songList SongList
	var result SearchResult
	err := json.Unmarshal([]byte(resBody), &result)
	if err != nil {
		panic(err)
	}

	jsonResult, _ := json.Marshal(result.Tracks.Items)
	fmt.Fprint(rw, string(jsonResult))
}

func AddSongToPlaylist(rw http.ResponseWriter, r *http.Request) {
	if pc.PlaylistId == "" {
		fmt.Fprint(rw, "No Playlist created.")
		return
	}
	r.ParseForm()
	var trackId = r.FormValue("trackId")
	// POST request to add track to playlist
	httpClient := &http.Client{}
	authHeader := fmt.Sprintf("Bearer %s", Spotify_Auth_Object.AccessToken)
	addtoPlaylistUrl := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists/%s/tracks?uris=spotify:track:%s",pc.PartyHostUserId, pc.PlaylistId, trackId)
	req, _ := http.NewRequest("POST", addtoPlaylistUrl, nil)
	req.Header.Set("Authorization", authHeader)
	_, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(rw, "Successfully added track to playlist.")
}

/************** BEGIN SECTION: HELPER FUNCTIONS *************/

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

/************** END SECTION: HELPER FUNCTIONS *************/

/*Fetch Playlists API endpoint */
/*httpClient := &http.Client{}

  getPlaylistsUrl := fmt.Sprint("https://api.spotify.com/v1/users/12124324757/playlists")
  authHeader := fmt.Sprintf("Bearer %s", "BQCdp2D2ocHOyYPdP-1CDlaAsH6uKn3TMhoCPpJFpwWPufArd_01VBjXWoevKmuqqoala4yrFbWAW15VbGM4uA0yHhzmTf0_IJVz-VA5iVO-cKxwF4sHGr23osmxGg7_E9MBCPUwQbmEmNrDeDwBnrWO_nd8p3tz")

  req, _ := http.NewRequest("GET", getPlaylistsUrl, nil)
  req.Header.Set("Authorization", authHeader)
  res, _ := httpClient.Do(req)
  resBody, _ := ioutil.ReadAll(res.Body)

  fmt.Fprint(rw, string(resBody))*/
