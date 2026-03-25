package zpi_indicator

import (
	"fmt"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

type ZLED struct {
	redPin, greenPin, bluePin    gpio.PinIO
	lastRed, lastGreen, lastBlue uint32
}

func NewLED(redPin, greenPin, bluePin int) *ZLED {
	host.Init()
	pin14 := gpioreg.ByName("GPIO14")
	pin15 := gpioreg.ByName("GPIO15")
	pin18 := gpioreg.ByName("GPIO18")

	return &ZLED{
		redPin:   pin14,
		greenPin: pin15,
		bluePin:  pin18,
	}
}

func (l *ZLED) Set(r, g, b uint32) error {
	if r > 0 {
		l.redPin.Out(gpio.High)
	} else {
		l.redPin.Out(gpio.Low)
	}
	if g > 0 {
		l.greenPin.Out(gpio.High)
	} else {
		l.greenPin.Out(gpio.Low)
	}
	if b > 0 {
		l.bluePin.Out(gpio.High)
	} else {
		l.bluePin.Out(gpio.Low)
	}

	l.lastRed, l.lastGreen, l.lastBlue = r, g, b
	fmt.Printf("R:%d G:%d B:%d\n", r, g, b)
	return nil
}

func (l *ZLED) Get() (uint32, uint32, uint32) {
	return l.lastRed, l.lastGreen, l.lastBlue
}
