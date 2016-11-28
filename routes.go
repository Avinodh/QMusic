package main

/**************** Declares all routes/API ***************/

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"CreatePartyController",
		"POST",
		"/createparty",
		CreatePartyController,
	},

	Route {
		"FindParties",
		"POST",
		"/findparties",
		FindParties,
		},

	Route{
		"AuthorizeSpotify",
		"GET",
		"/authspotify",
		AuthorizeSpotify,
	},

	Route{
		"Dashboard",
		"GET",
		"/dashboard",
		Dashboard,
	},

	Route{
		"RenderSearch",
		"GET",
		"/search",
		RenderSearch,
	},

	Route{
		"SearchSong",
		"GET",
		"/searchsong",
		SearchSong,
	},

	Route{
		"AddSongToPlaylist",
		"POST",
		"/addsong",
		AddSongToPlaylist,
	},

	Route{
		"ViewPlaylist",
		"GET",
		"/viewplaylist",
		ViewPlaylist,
	},

	Route{
		"FindRecommendedSongs",
		"GET",
		"/findrecommendedsongs",
		FindRecommendedSongs,
	},

	Route{
		"RenderPlaylist",
		"GET",
		"/playlist",
		RenderPlaylist,
	},
	Route{
		"GetHostParties",
		"GET",
		"/gethostparties",
		GetHostParties,
	},
	Route{
		"RenderDashboard",
		"GET",
		"/renderdashboard",
		RenderDashboard,
	},
	Route{
		"GetCurrentPlaylist",
		"GET",
		"/currentplaylist",
		GetCurrentPlaylist,
	},
	Route{
		"RemoveTrack",
		"POST",
		"/removetrack",
		RemoveTrack,
	},
}

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {

		/*********** LOGGER CODE *************/
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		/*************************************/

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler) //Analogous to Handler(route.handlerFunc)
	}

	s := http.StripPrefix("/", http.FileServer(http.Dir("./www")))
	router.PathPrefix("/").Handler(s)

	return router
}
