package main

import "machine"

var led = machine.LED_GREEN

var (
	pinSelectPPM = machine.D8
	pinOutputPPM = machine.D10
	pinCalibrate = machine.D4
)

var pinCalibrateDefault bool

const (
	pinAdcDefMin = 0x2000
	pinAdcDefMid = 0x8000
	pinAdcDefMax = 0xE000
	pinAdcCalLim = 0x3000
)

type AdcCalibration struct {
	min uint16
	mid uint16
	max uint16
}

var adcCalibration = []*AdcCalibration{}

// ----------------------------------------------------------------------------

func initLeds() {
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
}

func initPins() {
	pinSelectPPM.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	pinOutputPPM.Configure(machine.PinConfig{Mode: machine.PinOutput})
	pinCalibrate.Configure(machine.PinConfig{Mode: machine.PinInput})
	pinCalibrateDefault = pinCalibrate.Get()
}

func initAdc() {
	for _, pin := range pinAdc {
		pin.Configure(machine.ADCConfig{})
		value := pin.Get()
		min := uint16(pinAdcDefMin)
		mid := uint16(pinAdcDefMid)
		max := uint16(pinAdcDefMax)
		if value < pinAdcDefMin+pinAdcCalLim {
			min = value
		}
		if pinAdcDefMid-pinAdcCalLim < value && value < pinAdcDefMid+pinAdcCalLim {
			mid = value
		}
		adcCalibration = append(adcCalibration, &AdcCalibration{min: min, mid: mid, max: max})
	}
}

func calibrateAdc() {
	for i, c := range adcCalibration {
		value := pinAdc[i].Get()
		if value < pinAdcDefMin*2 {
			c.min = value
		}
		if pinAdcDefMid-pinAdcCalLim < value && value < pinAdcDefMid+pinAdcCalLim {
			c.mid = value
		}
	}
}
