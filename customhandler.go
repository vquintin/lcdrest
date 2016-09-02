package lcdrest

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type adder struct {
	rm *randomMessage
}

func (a *adder) apply(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	message := r.FormValue("message")
	a.rm.Add(key, message)
}

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func getRoutes(rm *randomMessage) []route {
	a := &adder{rm}
	return []route{
		route{
			Name:        "Add",
			Method:      "POST",
			Pattern:     "/",
			HandlerFunc: a.apply,
		},
	}
}

func NewCustomHandler(writer io.Writer) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	rm := NewRandomMessage(writer)
	routes := getRoutes(rm)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}
