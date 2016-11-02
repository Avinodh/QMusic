package main

import (
	"net/http"
)

type Spotify_User struct {
	Id          string `json:"uid,Number"`
	DisplayName string `json:"displayName,Number"`
	ProfilePic  string `json:"images"`
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
	Id string
	Active *Party
	AuthToken string
	RefreshToken string

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