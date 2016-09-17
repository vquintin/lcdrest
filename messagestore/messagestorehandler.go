package messagestore

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

type MessageStoreHandler struct {
	ms       MessageStore
	delegate http.Handler
}

func (ch MessageStoreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ch.delegate.ServeHTTP(w, r)
}

func (rm MessageStoreHandler) Close() error {
	return nil
}

func (ch MessageStoreHandler) put(w http.ResponseWriter, r *http.Request) {
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
	messageSlice, ok := values["message"]
	if !ok || len(messageSlice) != 1 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		message := messageSlice[0]
		_, created := ch.ms.Put(key, message)
		if created {
			w.WriteHeader(http.StatusCreated)
		}
	}
}

func (ch MessageStoreHandler) get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	message, found := ch.ms.Get(key)
	if found {
		_, err := w.Write([]byte(message))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (ch MessageStoreHandler) delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	_, found := ch.ms.Delete(key)
	if !found {
		w.WriteHeader(http.StatusNotFound)
	}
}

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func getRoutes(ch MessageStoreHandler) []route {
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

func NewMessageStoreHandler(ms MessageStore) MessageStoreHandler {
	router := mux.NewRouter().StrictSlash(true)
	msh := MessageStoreHandler{ms: ms}
	routes := getRoutes(msh)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	msh.delegate = router
	return msh
}
