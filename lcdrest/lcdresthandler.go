package lcdrest

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/vquintin/i2c"
	"github.com/vquintin/lcdrest/backlight"
	"github.com/vquintin/lcdrest/messagestore"
	"github.com/vquintin/lcdrest/randommessage"
)

type LcdRestHandler struct {
	bh       backlight.BacklightHandler
	msh      messagestore.MessageStoreHandler
	rm       randommessage.RandomMessage
	delegate http.Handler
}

func (lrh LcdRestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lrh.delegate.ServeHTTP(w, r)
}

func (lrh LcdRestHandler) Close() error {
	lrh.bh.Close()
	lrh.msh.Close()
	lrh.rm.Close()
	return nil
}

type route struct {
	Prefix      string
	HandlerFunc http.HandlerFunc
}

func getRoutes(lrh LcdRestHandler) []route {
	return []route{
		route{
			Prefix:      "/backlight",
			HandlerFunc: lrh.bh.ServeHTTP,
		},
		route{
			Prefix:      "/messages",
			HandlerFunc: lrh.msh.ServeHTTP,
		},
	}
}

func NewLcdRestHandler(messageStore messagestore.MessageStore, lcd *i2c.Lcd, duration time.Duration) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	lrh := LcdRestHandler{
		bh:       backlight.NewBacklightHandler(lcd),
		msh:      messagestore.NewMessageStoreHandler(messageStore),
		rm:       randommessage.NewRandomMessage(messageStore, &lcdWrapper{lcd}, duration),
		delegate: router,
	}
	routes := getRoutes(lrh)
	for _, route := range routes {
		router.PathPrefix(route.Prefix).
			Handler(route.HandlerFunc)
	}
	return lrh
}
