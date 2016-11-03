package main

import (
	"net/http"
)

type Spotify_User struct {
	Id          string `json:"uid,Number"`
	DisplayName string `json:"displayName,Number"`
	ProfilePic  string `json:"images"`
}

type Spotify_Auth struct {
	AccessToken    string `json:"access_token"`
	TokenType     string `json:"token_type"`
	ExpiresIn    int   `json:"expires_in"`
    RefreshToken string `json:"refresh_token"`
	Scope string `json:"scope"`

type Song struct {
  SongName string
  Artist string
  Album string
  Artwork string
  SongURI string
}

type Party struct {
	IsActive bool
	PartyHost string
	Location string
	SecretCode string
	ActiveTime string
	MyController *Party_Controller
}

type Party_Controller struct {
	Active *Party
	AuthToken string
	RefreshToken string
  PartyHostUserId string
  PlaylistId string
}

type Master_Controller struct {
	PartyControllers map[string]*Party_Controller
}

func (pc *Party_Controller) CreateParty(r* http.Request) bool {
    var new_party *Party = new(Party)
    new_party.IsActive = true
    new_party.PartyHost = r.Form["user"][0]
    new_party.Location = r.Form["location"][0]
    new_party.SecretCode = r.Form["secret-code"][0]
    new_party.ActiveTime = r.Form["active-time"][0]
    pc.Active = new_party
    pc.AuthToken = Spotify_Auth_Object.AccessToken
    pc.RefreshToken = Spotify_Auth_Object.RefreshToken
    return true
}

func (mc *Master_Controller) AddPartyController(id string) *Party_Controller {
	mc.PartyControllers[id] = new(Party_Controller)
	return mc.PartyControllers[id]
}

func InitializeController() *Master_Controller{
	var mc *Master_Controller = new(Master_Controller)
	mc.PartyControllers = make(map[string]*Party_Controller)
	return mc
}

var TheMasterController = InitializeController()
var Spotify_Auth_Object Spotify_Auth

