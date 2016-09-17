package backlight

import "log"

type BacklightDriverLogger struct {
	Delegate BacklightDriver
}

func (bdl BacklightDriverLogger) Close() error {
	log.Printf("[backlight][BacklightDriver][Close] Closing.")
	return bdl.Delegate.Close()
}

func (bdl BacklightDriverLogger) On() {
	log.Printf("[backlight][BacklightDriver][On] Turning backlight on.")
	bdl.Delegate.On()
}

func (bdl BacklightDriverLogger) Off() {
	log.Printf("[backlight][BacklightDriver][Off] Turning backlight off.")
	bdl.Delegate.Off()
}

func (bdl BacklightDriverLogger) SetLevel(level int) {
	log.Printf("[backlight][BacklightDriver][SetLevel] Setting backlight level to %v", level)
	bdl.Delegate.SetLevel(level)
}
