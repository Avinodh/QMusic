package main


type party struct {
	is_active bool
	location string
	my_controller *party_controller
}

type party_controller struct {
	id string
	active *party
	auth_token string

}

type master_controller struct {
	party_controllers map[string]*party_controller
}

func (mc *master_controller) Add_party_controller() bool {
	var new_pc *party_controller = new(party_controller)
	new_pc.id = "rand"
	new_pc.active = new(party)
	new_pc.active.is_active = true
	new_pc.active.location = "xyz"
	new_pc.active.my_controller = new_pc
	new_pc.auth_token = "abcdef"
	mc.party_controllers[new_pc.id] = new_pc 
	return true
}