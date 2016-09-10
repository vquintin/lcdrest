package lcdrest

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type customHandler struct {
	rm       randomMessage
	delegate http.Handler
}

func (ch customHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ch.delegate.ServeHTTP(w, r)
}

func (ch customHandler) put(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	log.Printf("[DEBUG][lcdRest][customHandler][put] %v", r)
	ch.rm.Put(key, "")
}

func (ch customHandler) get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	message, _ := ch.rm.Get(key)
	_, err := w.Write([]byte(message))
	log.Printf("[ERROR][lcdRest][customHandler][get] %v", err)
}

func (ch customHandler) delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	ch.rm.Delete(key)
}

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func getRoutes(ch customHandler) []route {
	return []route{
		route{
			Name:        "Put",
			Method:      "PUT",
			Pattern:     "/{key}",
			HandlerFunc: ch.put,
		},
		route{
			Name:        "Get",
			Method:      "GET",
			Pattern:     "/{key}",
			HandlerFunc: ch.get,
		},
		route{
			Name:        "Delete",
			Method:      "DELETE",
			Pattern:     "/{key}",
			HandlerFunc: ch.delete,
		},
	}
}

func NewCustomHandler(rm randomMessage) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	ch := customHandler{rm: rm}
	routes := getRoutes(ch)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	ch.delegate = router
	return ch
}
