package backlight

import "log"

type BacklightDriverLogger struct {
	delegate BacklightDriver
}

func (bdl BacklightDriverLogger) Close() error {
	log.Printf("[backlight][BacklightDriver][Close] Closing.")
	return bdl.delegate.Close()
}

func (bdl BacklightDriverLogger) On() {
	log.Printf("[backlight][BacklightDriver][On] Turning backlight on.")
	bdl.delegate.On()
}

func (bdl BacklightDriverLogger) Off() {
	log.Printf("[backlight][BacklightDriver][Off] Turning backlight off.")
	bdl.delegate.Off()
}

func (bdl BacklightDriverLogger) SetLevel(level int) {
	log.Printf("[backlight][BacklightDriver][Off] Setting backlight level to %v", level)
	bdl.delegate.Off()
}
