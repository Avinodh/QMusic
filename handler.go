package main

/************** Functions that handle all routes/ Defines our API *************/

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"github.com/zmb3/spotify"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"log"
)

var (
	store       = sessions.NewCookieStore([]byte(os.Getenv("COOKIE_STORE_AUTHENTICATION_KEY")), []byte(os.Getenv("COOKIE_STORE_ENCRYPTION_KEY")))
	redirectURL = os.Getenv("REDIRECT_URL")
	auth        = spotify.NewAuthenticator(redirectURL, spotify.ScopeUserReadPrivate, spotify.ScopePlaylistModifyPublic)
	db          *sql.DB
	dberr       error
)

func init() {
	/**************** DATABASE CONNECTION ********************/
	db, dberr = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if dberr != nil {
		panic(dberr)
	}

	p, _ := db.Prepare("INSERT INTO party VALUES ($1, $2, $3, $4, $5, $6, $7)")
	_, e := p.Exec("1","2","3","4","5","6","7")
	if e != nil {
		panic(e)
	}
	/**********************************************************/
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

	log.Printf("%s", respBody)
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

	pc = new(Party_Controller)
	pc.PartyHostUserId = Spotify_User_Object.Id

	Spotify_User_Object.DisplayPic = Spotify_User_Object.ProfilePic[0].Url
	/******************************************/

	/*body, _ := ioutil.ReadFile("www/dashboard.html")
	fmt.Fprint(rw, string(body))*/


	t, err := template.ParseFiles("www/dashboard.html")
	checkErr(err)

	err = t.Execute(rw, Spotify_User_Object)
	checkErr(err)
}

func RenderDashboard(rw http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("www/dashboard.html")
	checkErr(err)

	err = t.Execute(rw, Spotify_User_Object)
	checkErr(err)
}

func GetCurrentPlaylist(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprint(rw, pc.PlaylistId)
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

	createPlaylistUrl := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", Spotify_User_Object.Id)
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


	// pc = TheMasterController.AddPartyController(r.Form["secret-code"][0])
	pc.CreateParty(r, playlist.Id)

	/*body, _ := ioutil.ReadFile("www/search.html")
	fmt.Fprint(rw, string(body))*/
	http.Redirect(rw, r, "/search", http.StatusSeeOther)
}

func SearchSong(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var query = r.FormValue("searchsong")

	getTrackUrl := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track", url.QueryEscape(query))

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

func FindRecommendedSongs(rw http.ResponseWriter, r *http.Request) {
	getTrackUrl := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists/%s/tracks", pc.PartyHostUserId, pc.PlaylistId)
	httpClient := &http.Client{}
	authHeader := fmt.Sprintf("Bearer %s", Spotify_Auth_Object.AccessToken)
	req, _ := http.NewRequest("GET", getTrackUrl, nil)

	req.Header.Set("Authorization", authHeader)
	res, _ := httpClient.Do(req)
	resBody, _ := ioutil.ReadAll(res.Body)

	var result ViewTracks
	err := json.Unmarshal([]byte(resBody), &result)
	if err != nil {
		panic(err)
	}

	// Need at least one song
	songId, _ := json.Marshal(result.Items[0].TrackItem.Id)

	getRecommendedUrl := fmt.Sprintf("https://api.spotify.com/v1/recommendations?seed_tracks=%s&market=US", songId[1:len(songId)-1])

	httpClient = &http.Client{}
	req, _ = http.NewRequest("GET", getRecommendedUrl, nil)

	req.Header.Set("Authorization", authHeader)
	res, _ = httpClient.Do(req)
	resBody, _ = ioutil.ReadAll(res.Body)
	log.Printf("%s\n%s", songId, getRecommendedUrl)

	var recommendedResult ViewRecommendedTracks
	err = json.Unmarshal([]byte(resBody), &recommendedResult)
	if err != nil {
		panic(err)
	}

	jsonResult, _ := json.Marshal(recommendedResult)
	log.Printf("\n\n\n%s\n%s\n%s\n%s\n\n\n", jsonResult, resBody, req, recommendedResult)
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
	addtoPlaylistUrl := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists/%s/tracks?uris=spotify:track:%s", pc.PartyHostUserId, pc.PlaylistId, trackId)
	req, _ := http.NewRequest("POST", addtoPlaylistUrl, nil)
	req.Header.Set("Authorization", authHeader)
	_, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(rw, "Successfully added track to playlist.")
}

func GetHostParties(rw http.ResponseWriter, r *http.Request) {
	var (
		active_time    string
		secret_code    string
		party_name     string
		party_location string
		playlist_id    string
	)

	p, err := db.Prepare("SELECT active_time, secret_code, party_name, party_location, playlist_id from party WHERE host_id=$1")
	checkErr(err)

	rows, err := p.Query(Spotify_User_Object.Id)
	checkErr(err)

	var hostParties = HostParties{}

	for rows.Next() {
		err := rows.Scan(&active_time, &secret_code, &party_name, &party_location, &playlist_id)
		checkErr(err)

		hostParty := HostParty{ActiveTime: active_time, SecretCode: secret_code, PartyName: party_name, PartyLocation: party_location, PlaylistId: playlist_id}
		hostParties = append(hostParties, hostParty)
	}

	err = rows.Err()
	checkErr(err)
	defer p.Close()
	defer rows.Close()

	if err := json.NewEncoder(rw).Encode(hostParties); err != nil {
		panic(err)
	}
}

func FindParties(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var query = r.FormValue("usercode")

	var (
		active_time    string
		secret_code    string
		party_name     string
		party_location string
		playlist_id    string
	)

	p, err := db.Prepare("SELECT active_time, secret_code, party_name, party_location, playlist_id from party WHERE secret_code=$1")
	checkErr(err)

	rows, err := p.Query(query)
	checkErr(err)

	var hostParties = HostParties{}

	for rows.Next() {
		err := rows.Scan(&active_time, &secret_code, &party_name, &party_location, &playlist_id)
		checkErr(err)

		hostParty := HostParty{ActiveTime: active_time, SecretCode: secret_code, PartyName: party_name, PartyLocation: party_location, PlaylistId: playlist_id}
		hostParties = append(hostParties, hostParty)
	}

	pc.PlaylistId = playlist_id

	err = rows.Err()
	checkErr(err)
	defer p.Close()
	defer rows.Close()


	//	if err := json.NewEncoder(rw).Encode(hostParties); err != nil {
		//	panic(err)
	//}

	if hostParties[0].SecretCode != query {
		fmt.Fprint(rw,"Code didn't match")
	} else {
		http.Redirect(rw, r, "/search", http.StatusSeeOther)
	}

}


func ViewPlaylist(rw http.ResponseWriter, r *http.Request) {
	getTrackUrl := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists/%s/tracks", pc.PartyHostUserId, pc.PlaylistId)
	httpClient := &http.Client{}
	authHeader := fmt.Sprintf("Bearer %s", Spotify_Auth_Object.AccessToken)
	req, _ := http.NewRequest("GET", getTrackUrl, nil)

	req.Header.Set("Authorization", authHeader)
	res, _ := httpClient.Do(req)
	resBody, _ := ioutil.ReadAll(res.Body)

	var result ViewTracks
	err := json.Unmarshal([]byte(resBody), &result)
	if err != nil {
		panic(err)
	}

	/*t, err := template.ParseFiles("www/templates/playlist.html")
	if err != nil {
		fmt.Fprint(rw, err)
		return
	}
	t.Execute(rw, result)*/
	//jsonResult, _ := json.Marshal(result)
	//fmt.Fprint(rw, string(jsonResult))
	//fmt.Fprint(rw, result)
	jsonResult, _ := json.Marshal(result.Items)
	fmt.Fprint(rw, string(jsonResult))
}

func RenderPlaylist(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var playlistId = r.FormValue("playlist_id")
	pc.PlaylistId = playlistId
	var (
		hostParty HostParty
		party_name string
		party_location string
		active_time string
	)
	row := db.QueryRow("SELECT active_time, party_name, party_location from party WHERE playlist_id=$1", pc.PlaylistId)
	err := row.Scan(&active_time, &party_name, &party_location)

	hostParty = HostParty{ActiveTime: active_time, PartyName: party_name, PartyLocation: party_location}

	t, err := template.ParseFiles("www/playlist.html")
	checkErr(err)

	err = t.Execute(rw, hostParty)
	checkErr(err)
}

func RemoveTrack(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var trackId = r.FormValue("trackId")

	removeFromPlaylistUrl := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists/%s/tracks", pc.PartyHostUserId, pc.PlaylistId)
	params := fmt.Sprintf("{\"tracks\":[{\"uri\":\"spotify:track:%s\"}]}", trackId)
	reader := strings.NewReader(params)
	req, _ := http.NewRequest("DELETE", removeFromPlaylistUrl, reader)

	// DELETE request to add track to playlist
	httpClient := &http.Client{}
	authHeader := fmt.Sprintf("Bearer %s", Spotify_Auth_Object.AccessToken)
	req.Header.Set("Authorization", authHeader)
	_, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(rw, "Successfully Deleted Song!")
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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
