package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vquintin/i2c"
	"github.com/vquintin/lcdrest"
)

const (
	rs uint8 = iota
	rw
	en
	backlight
	d4
	d5
	d6
	d7
)

func main() {
	var bus = flag.Int("device", 1, "i2c bus for lcd screen (e.g. 1 for /dev/i2c-1)")
	var address = flag.Uint("address", 0x3f, "i2c address for lcd screen")
	var port = flag.Uint("port", 8080, "Server port")
	flag.Parse()
	i2cConnector, err := i2c.New(uint8(*address), *bus)
	if err != nil {
		log.Fatal("Can't open i2c bus")
	}
	lcd, err := i2c.NewLcd(i2cConnector, en, rw, rs, d4, d5, d6, d7, backlight)
	if err != nil {
		log.Fatal("Cant' open lcd screen")
	}
	lcd.BacklightOn()
	rm := lcdrest.NewRandomMessage(&lcdWrapper{lcd}, 10*time.Second)
	handler := lcdrest.NewCustomHandler(rm)
	serverAddress := fmt.Sprintf(":%v", *port)
	log.Fatal(http.ListenAndServe(serverAddress, handler))
}
