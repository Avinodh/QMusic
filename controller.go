package main

/*
import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func (pc *Party_Controller) addToPlaylist(s Song) {
	//maintain an internal queue of Songs?
	apiUrl := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists/%s/tracks/", pc.PartyHostUserId, pc.PlaylistId)
	data := url.Values{}
	data.Set("uris", s.SongURI)

	u, _ := url.ParseRequestURI(apiUrl)
	urlStr := fmt.Sprintf("%v", u) // "https://api.com/user/"

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
	r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := client.Do(r)
	fmt.Println(resp.Status)
}*/
