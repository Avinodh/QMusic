
package main

/************** Functions that handle all routes/ Defines our API *************/

import (
    "fmt"
    //"github.com/gorilla/mux"
    "crypto/rand"
    "encoding/base64"
    "github.com/gorilla/sessions"
    "github.com/zmb3/spotify"
    "io/ioutil"
    "net/http"
    "os"
)

var (
    store       = sessions.NewCookieStore([]byte(os.Getenv("COOKIE_STORE_AUTHENTICATION_KEY")), []byte(os.Getenv("COOKIE_STORE_ENCRYPTION_KEY")))
    redirectURL = os.Getenv("REDIRECT_URL")
    auth        = spotify.NewAuthenticator(redirectURL, spotify.ScopeUserReadPrivate)
)

func init() {
    store.Options = &sessions.Options{
        Path:     "/",
        MaxAge:   3600 * 8, // 8 hours
        HttpOnly: true,
    }
}

func Index(rw http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("www/index.html")
    fmt.Fprint(rw, string(body))
}

func AuthorizeSpotify(rw http.ResponseWriter, r *http.Request) {
    // Setting GET request query parameters
    /*v := url.Values{}
      v.Set("client_id", client_spotify_id)
      v.Set("response_type", "code")
      v.Set("redirect_uri", "http://qmusicapp.herokuapp.com/dashboard")

      queryString := fmt.Sprintf("https://accounts.spotify.com/authorize/?%s", v.Encode())

      resp, _ := http.Get(queryString)
      fmt.Fprint(rw, resp)*/

    // Generating Session Token
    session, _ := store.Get(r, "UserAuthStore")

    session_token, _ := GenerateRandomString(32)

    // Setting Session Token
    session.Values["session_token"] = &session_token
    err := sessions.Save(r, rw)
    if err != nil {
        fmt.Fprint(rw, "Session not saved! ", err)
    }

    // get the user to this URL - how you do that is up to you
    // you should specify a unique state string to identify the session
    url := auth.AuthURL(session_token)

    fmt.Fprint(rw, url)
}

func Dashboard(rw http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "UserAuthStore")
    session_token := fmt.Sprint(session.Values["session_token"])

    // use the same state string here that you used to generate the URL
    token, err := auth.Token(session_token, r)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusNotFound)
        return
    }
    // create a client using the specified token
    client := auth.NewClient(token)

    client_info, _ := client.CurrentUser()

   // Extract user data from spotify client and populate Spotify User model
    /*user := Spotify_User{
        Id:          client_info.User.ID,
        DisplayName: client_info.User.DisplayName,
        ProfilePic:  client_info.User.Images[2].URL,
    }*/

    client_info = client_info

    body, _ := ioutil.ReadFile("www/dashboard.html")
    fmt.Fprint(rw, string(body))
    // the client can now be used to make authenticated requests
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

/************** END SECTION: HELPER FUNCTIONS *************/
