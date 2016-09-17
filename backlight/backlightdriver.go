package backlight

import (
	"time"

	"github.com/vquintin/i2c"
)

const (
	MAXBACKLIGHT = 100
	PWMPERIOD    = time.Second / 100
)

type BacklightDriver struct {
	levelC chan<- int
	onC    chan<- bool
}

func NewBacklightDriver(lcd *i2c.Lcd) BacklightDriver {
	levelC := make(chan int)
	onC := make(chan bool)
	go backlightRoutine(lcd, levelC, onC)
	return BacklightDriver{
		levelC: levelC,
		onC:    onC,
	}
}

func (bd BacklightDriver) Close() error {
	//TODO
	return nil
}

func (bd BacklightDriver) On() {
	bd.onC <- true
}

func (bd BacklightDriver) Off() {
	bd.onC <- false
}

func (bd BacklightDriver) SetLevel(level int) {
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
				lcd.BacklightOff()
			} else {
				lcd.BacklightOff()
			}
		case level = <-levelC:
		case on = <-onC:
		}
	}
}
