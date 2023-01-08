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

	return router
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"AuthPost",
		strings.ToUpper("Post"),
		"/auth",
		AuthPost,
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
