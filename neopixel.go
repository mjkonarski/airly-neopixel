package main

import (
	"fmt"
	"log"

	"github.com/tarm/serial"
)

type NeopixelColor struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

type Neopixel struct {
	SerialConfig *serial.Config
	Serial       *serial.Port
}

func newNeopixel(port string, baud int) *Neopixel {

	serialConfig := &serial.Config{
		Name: port,
		Baud: baud,
	}

	return &Neopixel{
		SerialConfig: serialConfig,
	}
}

func (neopixel *Neopixel) open() {
	serial, err := serial.OpenPort(neopixel.SerialConfig)
	if err != nil {
		log.Fatal(err)
	}

	neopixel.Serial = serial
}

func (neopixel *Neopixel) close() {
	neopixel.Serial.Close()
}

func (neopixel *Neopixel) setColor(pixelID int, color NeopixelColor) {
	command := fmt.Sprintf("%d %d %d %d\n",
		pixelID, color.Red, color.Green, color.Blue)

	_, err := neopixel.Serial.Write([]byte(command))
	if err != nil {
		log.Fatal(err)
	}
}
