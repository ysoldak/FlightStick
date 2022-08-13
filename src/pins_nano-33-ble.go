//go:build nano_33_ble
// +build nano_33_ble

package main

import "machine"

var pinAdc = []machine.ADC{
	{Pin: machine.A0},
	{Pin: machine.A1},
	{Pin: machine.A2},
	{Pin: machine.A3},
	{Pin: machine.A4},
	{Pin: machine.A5},
	{Pin: machine.A6},
	{Pin: machine.A7},
}
