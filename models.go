package main

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
}


var Spotify_Auth_Object Spotify_Auth 