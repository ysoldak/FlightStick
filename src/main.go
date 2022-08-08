package main

import (
	"machine"
	"time"
)

var led = machine.LED_GREEN

var controlPins = []machine.Pin{
	machine.D4, // calibration: record pot middle positions
}

var potPins = []machine.ADC{
	{Pin: machine.A0},
	{Pin: machine.A1},
	{Pin: machine.A2},
	{Pin: machine.A3},
}

var channels = []uint16{1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500}

type AxisCalibration struct {
	min uint16
	mid uint16
	max uint16
}

var calibration = []*AxisCalibration{}

var controlDefaults = []bool{}

// ----------------------------------------------------------------------------

func main() {

	ledInit()
	ppmInit()
	controlInit()
	potInit()

	go listenControls()
	go updateChannels()

	for {
		if led.Get() {
			led.Low()
		} else {
			led.High()
		}
		time.Sleep(time.Second)
	}

}

func ledInit() {
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
}

func controlInit() {
	for _, pin := range controlPins {
		pin.Configure(machine.PinConfig{Mode: machine.PinInput})
		controlDefaults = append(controlDefaults, pin.Get())
	}
}

func potInit() {
	for _, pin := range potPins {
		pin.Configure(machine.ADCConfig{})
		calibration = append(calibration, &AxisCalibration{min: 10000, mid: pin.Get(), max: 60000})
	}
}

// ----------------------------------------------------------------------------

func listenControls() {
	for {
		// calibrate
		if controlPins[0].Get() != controlDefaults[0] {
			for i, c := range calibration {
				c.mid = potPins[i].Get()
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// ----------------------------------------------------------------------------

func updateChannels() {
	for {
		for i, pot := range potPins {
			channels[i] = potToChannel(i, pot)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func potToChannel(i int, pot machine.ADC) uint16 {
	value := pot.Get()
	if calibration[i].min > value {
		calibration[i].min = value
	}
	if calibration[i].max < value {
		calibration[i].max = value
	}
	if value < calibration[i].mid {
		potRange := calibration[i].mid - calibration[i].min
		return 988 + uint16(float64(value-calibration[i].min)/float64(potRange)*512)
	}
	if value > calibration[i].mid {
		potRange := calibration[i].max - calibration[i].mid
		return 1500 + uint16(float64(value-calibration[i].mid)/float64(potRange)*512)
	}
	return 1500
}
