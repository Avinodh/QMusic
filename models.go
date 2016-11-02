package main

type Spotify_User struct {
	Id          string `json:"uid,Number"`
	DisplayName string `json:"displayName,Number"`
	ProfilePic  string `json:"images"`
}

type Party struct {
	IsActive bool
	Location string
	MyController *Party_Controller
}

type Party_Controller struct {
	Id string
	Active *Party
	AuthToken string

}

type Master_Controller struct {
	PartyControllers map[string]*Party_Controller
}

func (mc *master_controller) Add_party_controller() bool {
	var new_pc *Party_controller = new(Party_Controller)
	new_pc.Id = "rand"
	new_pc.Active = new(Party)
	new_pc.Active.IsActive = true
	new_pc.Active.Location = "xyz"
	new_pc.Active.MyController = new_pc
	new_pc.AuthToken = "abcdef"
	mc.party_controllers[new_pc.Id] = new_pc 
	return true
}