package backlight

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vquintin/i2c"
)

type BacklightHandler struct {
	delegate http.Handler
	bd       BacklightDriver
}

func (bh BacklightHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bh.delegate.ServeHTTP(w, r)
}

func (bh BacklightHandler) Close() error {
	return bh.bd.Close()
}

func (bh BacklightHandler) on(w http.ResponseWriter, r *http.Request) {
	bh.bd.On()
}

func (bh BacklightHandler) off(w http.ResponseWriter, r *http.Request) {
	bh.bd.Off()
}

func (bh BacklightHandler) setLevel(w http.ResponseWriter, r *http.Request) {
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
	levelSlice, ok := values["level"]
	if !ok || len(levelSlice) != 1 {
		w.WriteHeader(http.StatusBadRequest)
	} else if level, err := strconv.Atoi(levelSlice[0]); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		bh.bd.SetLevel(level)
	}
}

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func getRoutes(bh BacklightHandler) []route {
	return []route{
		route{
			Name:        "On",
			Method:      "POST",
			Pattern:     "/on",
			HandlerFunc: bh.on,
		},
		route{
			Name:        "Off",
			Method:      "POST",
			Pattern:     "/off",
			HandlerFunc: bh.on,
		},
		route{
			Name:        "Level",
			Method:      "PUT",
			Pattern:     "/level",
			HandlerFunc: bh.setLevel,
		},
	}
}

func NewBacklightHandler(lcd *i2c.Lcd) BacklightHandler {
	router := mux.NewRouter().StrictSlash(true)
	bd := NewBacklightDriver(lcd)
	msh := BacklightHandler{bd: bd}
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
