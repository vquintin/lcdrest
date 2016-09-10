package lcdrest

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

type customHandler struct {
	rm       *randomMessage
	delegate http.Handler
}

func (ch customHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ch.delegate.ServeHTTP(w, r)
}

func (ch customHandler) put(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	rawBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	values, err := url.ParseQuery(string(rawBody))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("[DEBUG][lcdRest][customHandler][put] %v", values)
	messageSlice, ok := values["message"]
	if !ok || len(messageSlice) != 1 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		message := messageSlice[0]
		_, created := ch.rm.Put(key, message)
		if created {
			w.WriteHeader(http.StatusCreated)
		}
	}
}

func (ch customHandler) get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	message, found := ch.rm.Get(key)
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Printf("[ERROR][lcdRest][customHandler][get] %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else if !found {
		w.WriteHeader(http.StatusNotFound)
	}
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

func NewCustomHandler(rm *randomMessage) http.Handler {
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
