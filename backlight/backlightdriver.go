package backlight

import (
	"io"
	"time"

	"github.com/vquintin/i2c"
)

const (
	MAXBACKLIGHT = 100
	PWMPERIOD    = time.Second / 100
)

type BacklightDriver interface {
	On()
	Off()
	SetLevel(level int)
	io.Closer
}

type backlightDriver struct {
	levelC chan<- int
	onC    chan<- bool
}

func NewBacklightDriver(lcd *i2c.Lcd) BacklightDriver {
	levelC := make(chan int)
	onC := make(chan bool)
	go backlightRoutine(lcd, levelC, onC)
	return backlightDriver{
		levelC: levelC,
		onC:    onC,
	}
}

func (bd backlightDriver) Close() error {
	//TODO
	return nil
}

func (bd backlightDriver) On() {
	bd.onC <- true
}

func (bd backlightDriver) Off() {
	bd.onC <- false
}

func (bd backlightDriver) SetLevel(level int) {
	if level < 0 {
		bd.levelC <- 0
	} else if level > MAXBACKLIGHT {
		bd.levelC <- MAXBACKLIGHT
	} else {
		bd.levelC <- level
	}
}

func backlightRoutine(lcd *i2c.Lcd, levelC <-chan int, onC <-chan bool) {
	pwm := time.NewTicker(PWMPERIOD)
	on := true
	level := MAXBACKLIGHT
	for {
		select {
		case <-pwm.C:
			if on {
				lcd.BacklightOn()
				time.Sleep(time.Duration(level) * PWMPERIOD / MAXBACKLIGHT)
			}
			lcd.BacklightOff()
		case level = <-levelC:
		case on = <-onC:
		}
	}
}
