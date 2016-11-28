package main

import (
	"net/http"
)

type Spotify_User struct {
	Id          string  `json:"id"`
	DisplayName string  `json:"display_name"`
	ProfilePic  []Image `json:"images"`
	DisplayPic	string `json:"display_pic"`
}

type Image struct {
	Url string `json:"url"`
}

type Spotify_Auth struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type SearchResult struct {
	Tracks SearchTracks `json:"tracks"`
}

type SearchTracks struct {
	Items []Track `json:"items"`
}

type Track struct {
	Id        string   `json:"id"`
	TrackName string   `json:"name"`
	Artists   []Artist `json:"artists"`
}

type ViewTracks struct {
	Items []ViewTrack `json:"items"`
}

type ViewTrack struct {
	TrackItem Track `json:"track"`
}

type ViewRecommendedTracks struct {
	Items []Track `json:"tracks"`
}

type Playlist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Artist struct {
	AristName string `json:"name"`
}

type HostParty struct {
	IsActive      bool   `json:"is_active"`
	ActiveTime    string `json:"active_time"`
	SecretCode    string `json:"secret_code"`
	PartyName     string `json:"party_name"`
	PartyLocation string `json:"party_location"`
	PlaylistId    string `json:"playlist_id"`
	MyController  *Party_Controller
}

type Party_Controller struct {
	PartyHostUserId string
	PlaylistId      string
}

/* type Master_Controller struct {
	PartyControllers map[string]*Party_Controller
}*/

func (pc *Party_Controller) CreateParty(r *http.Request, playlistId string) bool {
	var new_party *HostParty = new(HostParty)
	new_party.IsActive = true
	new_party.PartyName = r.Form["user"][0]
	new_party.PartyLocation = r.Form["location"][0]
	new_party.SecretCode = r.Form["secret-code"][0]
	new_party.ActiveTime = r.Form["active-time"][0]
	new_party.PlaylistId = playlistId

	pc.PartyHostUserId = Spotify_User_Object.Id
	pc.PlaylistId = playlistId

	/**** Saving to Database ****/
	p, _ := db.Prepare("INSERT INTO partycontroller VALUES ($1, $2, $3, $4, $5)")
	_, e := p.Exec(new_party.SecretCode, Spotify_Auth_Object.AccessToken, Spotify_Auth_Object.RefreshToken, pc.PartyHostUserId, pc.PlaylistId)
	if e != nil {
		panic(e)
	}

	p, _ = db.Prepare("INSERT INTO party VALUES ($1, $2, $3, $4, $5, $6, $7)")
	_, e = p.Exec(new_party.IsActive, new_party.ActiveTime, new_party.SecretCode, pc.PartyHostUserId, new_party.PartyName, new_party.PartyLocation, pc.PlaylistId)
	if e != nil {
		panic(e)
	}
	/****************************/
	return true
}

/*func (mc *Master_Controller) AddPartyController(id string) *Party_Controller {
	mc.PartyControllers[id] = new(Party_Controller)
	return mc.PartyControllers[id]
}

func InitializeController() *Master_Controller {
	mc = new(Master_Controller)
	mc.PartyControllers = make(map[string]*Party_Controller)
	return mc
}*/

//var TheMasterController = InitializeController()
var Spotify_Auth_Object Spotify_Auth
var Spotify_User_Object Spotify_User
//var mc *Master_Controller
var pc *Party_Controller

type HostParties []HostParty
