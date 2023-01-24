package handler

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	static := http.StripPrefix("/static/", http.FileServer(http.Dir("./resources/static/")))    // file server
	scripts := http.StripPrefix("/scripts/", http.FileServer(http.Dir("./resources/scripts/"))) // file server
	router.PathPrefix("/static/").Handler(static)
	router.PathPrefix("/scripts/").Handler(scripts)

	return router
}

var routes = Routes{
	Route{
		"Index",
		strings.ToUpper("Get"),
		"/",
		Index,
	},

	Route{
		"Main",
		strings.ToUpper("Get"),
		"/main",
		Main,
	},

	Route{
		"AuthGet",
		strings.ToUpper("Get"),
		"/auth",
		AuthGet,
	},

	Route{
		"AuthPost",
		strings.ToUpper("Post"),
		"/auth",
		AuthPost,
	},

	Route{
		"LogoutGet",
		strings.ToUpper("Get"), // да, лучше было бы сделать POST
		"/logout",
		LogoutGet,
	},

	Route{
		"AddAdvertGet",
		strings.ToUpper("Get"),
		"/advert",
		AddAdvertGet,
	},

	Route{
		"AddAdvertPost",
		strings.ToUpper("Post"),
		"/advert",
		AddAdvertPost,
	},

	Route{
		"AddCarPost",
		strings.ToUpper("Post"),
		"/car",
		AddCarPost,
	},

	Route{
		"OrderIdDelete",
		strings.ToUpper("Delete"),
		"/order/{id}",
		OrderIdDelete,
	},

	Route{
		"OrderIdGet",
		strings.ToUpper("Get"),
		"/order/{id}",
		OrderIdGet,
	},

	Route{
		"RegistrationGet",
		strings.ToUpper("Get"),
		"/registration",
		RegistrationGet,
	},

	Route{
		"RegistrationPost",
		strings.ToUpper("Post"),
		"/registration",
		RegistrationPost,
	},
}
