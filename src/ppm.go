package main

import (
	"device/arm"
	"machine"
)

// --- Confgurable ------------------------------------------------------------
var ppmPin = machine.D3

// --- Implementation ---------------------------------------------------------

var ppmChanNumber = -1

var ppmFrameLen uint32 = 22500 * sysCyclesPerMicrosecond
var ppmOffLen uint32 = 300 * sysCyclesPerMicrosecond

var sysCyclesPerMicrosecond = machine.CPUFrequency() / 1_000_000

func ppmInit() {
	ppmPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ppmPin.Low()
	arm.SetupSystemTimer(sysCyclesPerMicrosecond)
}

// --- Interrupt Handler ------------------------------------------------------

//export SysTick_Handler
func timer_isr() {
	// separator
	if ppmPin.Get() {
		ppmPin.Low()
		ppmChanNumber++
		if ppmChanNumber > 7 {
			ppmChanNumber = -1
		}
		arm.SetupSystemTimer(ppmOffLen)
		return
	}
	// regular channel
	if ppmChanNumber != -1 {
		ppmPin.High()
		arm.SetupSystemTimer(uint32(channels[ppmChanNumber])*sysCyclesPerMicrosecond - ppmOffLen)
		return
	}
	// padding
	ppmPin.High()
	sum := uint16(0)
	for _, value := range channels {
		sum += value
	}
	arm.SetupSystemTimer(ppmFrameLen - uint32(sum)*sysCyclesPerMicrosecond - 8*ppmOffLen)
}
