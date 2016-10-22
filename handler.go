package main

/************** Functions that handle all routes/ Defines our API *************/

import (
    "fmt"
    //"github.com/gorilla/mux"
    "io/ioutil"
    "net/http"
)

func init() {
}

func Index(rw http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("www/index.html")
    fmt.Fprint(rw, string(body))
}
