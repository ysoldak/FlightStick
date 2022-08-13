package main

import (
	"machine"
	"time"

	"github.com/ysoldak/FlightStick/src/trainer"
)

var Version string

var channels = []uint16{1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500}

// amount of microseconds around 1500 result in 1500 returned
// this removes jitter around mid stick position and helps with INAV launch mode
// TODO calibrate for jitter automatically on mid-stick calibration
const channelTolerance = 0x18

var t trainer.Trainer

// ----------------------------------------------------------------------------

func init() {
	initLeds()
	initPins()
	initAdc()
}

func main() {

	// Trainer (Bluetooth or PPM)
	if !pinSelectPPM.Get() { // Low means connected to GND => PPM output requested
		t = trainer.NewPPM(pinOutputPPM) // PPM wire
	} else {
		t = trainer.NewPara()
	}
	t.Configure()
	go t.Run()

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

// ----------------------------------------------------------------------------

func listenControls() {
	for {
		if pinCalibrate.Get() != pinCalibrateDefault {
			calibrateAdc()
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// ----------------------------------------------------------------------------

func updateChannels() {
	for {
		for i, pot := range pinAdc {
			channels[i] = adcToChannel(i, pot)
			if channels[i] < 988+channelTolerance {
				channels[i] = 988
			}
			if 1500-channelTolerance < channels[i] && channels[i] < 1500+channelTolerance {
				channels[i] = 1500
			}
			if 2012-channelTolerance < channels[i] {
				channels[i] = 2012
			}
		}
		for i, value := range channels {
			t.SetChannel(i, value)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func adcToChannel(i int, pot machine.ADC) uint16 {
	value := pot.Get()
	if adcCalibration[i].min > value {
		adcCalibration[i].min = value
	}
	if adcCalibration[i].max < value {
		adcCalibration[i].max = value
	}
	if value < adcCalibration[i].mid {
		potRange := adcCalibration[i].mid - adcCalibration[i].min
		return 988 + uint16(float64(value-adcCalibration[i].min)/float64(potRange)*512)
	}
	if value > adcCalibration[i].mid {
		potRange := adcCalibration[i].max - adcCalibration[i].mid
		return 1500 + uint16(float64(value-adcCalibration[i].mid)/float64(potRange)*512)
	}
	return 1500
}
