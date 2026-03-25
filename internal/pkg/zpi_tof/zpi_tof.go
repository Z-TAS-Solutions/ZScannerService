package zpi_trigger

import (
	"log"

	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-vl53l0x"
	"periph.io/x/host/v3"
)

type ZToF struct {
	sensor       *vl53l0x.Vl53l0x
	i2c_bus      *i2c.I2C
	lastDistance uint16
}

func NewZToF() (*ZToF, error) {
	if _, err := host.Init(); err != nil {
		return nil, err
	}

	i2c_bus, err := i2c.NewI2C(0x29, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer i2c_bus.Close()

	sensor := vl53l0x.NewVl53l0x()

	err = sensor.Reset(i2c_bus)
	if err != nil {
		log.Fatal(err)
	}

	err = sensor.Init(i2c_bus)
	if err != nil {
		log.Fatal(err)
	}
	return &ZToF{sensor: sensor, i2c_bus: i2c_bus, lastDistance: 9999}, nil
}

func (t *ZToF) Read() (uint16, error) {
	distance, err := t.sensor.ReadRangeSingleMillimeters(t.i2c_bus)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Measured range = %v mm", distance)

	t.lastDistance = uint16(distance)

	return t.lastDistance, nil
}
