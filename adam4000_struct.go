// Copyright 2009 Thomas Jager <mail@jager.no> All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Driver for ADAM-4000 series I/O Modules from Advantech

package main

import (
	"bufio"
	"fmt"
)

type DataFormatCode byte

const (
	EngineeringUnits = iota
	PercentofFSR
	TwosComplement
	Ohms
)

func (df DataFormatCode) String() string {
	switch df {
	case EngineeringUnits:
		return "Engineering Units"
	case PercentofFSR:
		return "% of FSR"
	case TwosComplement:
		return "Twos Complement"
	case Ohms:
		return "Ohms"
	}
	return "Undefined"
}

type InputRangeCode byte

const (
	c15mV InputRangeCode = iota
	c50mV
	c100mV
	c500mV
	c1V
	c2_5V
	_
	c4_20mA
	c10V
	c5V
	_
	_
	_
	c20mA
	_
	cJTc
	cKTc
	cTTc
	cETc
	cRTc
	cSTc
	cBTc
)

func (br InputRangeCode) String() string {
	switch br {
	case c15mV:
		return "+/- 15mV"
	case c50mV:
		return "+/- 50mV"
	case c100mV:
		return "+/- 100mV"
	case c500mV:
		return "+/- 500mV"
	case c1V:
		return "+/- 1V"
	case c2_5V:
		return "+/- 2.5V"
	case c4_20mA:
		return "4~20mA"
	case c10V:
		return "+/- 10V"
	case c5V:
		return "+/- 5V"
	case c20mA:
		return "+/- 20mA"
	case cJTc:
		return "Type-J TC"
	case cKTc:
		return "Type-K TC"
	case cTTc:
		return "Type-T TC"
	case cETc:
		return "Type-E TC"
	case cRTc:
		return "Type-R TC"
	case cSTc:
		return "Type-S TC"
	case cBTc:
		return "Type-B TC"
	}
	return "Undefined"
}

type BaudRateCode byte

const (
	c1200bps BaudRateCode = (iota + 0x03)
	c2400bps
	c4800bps
	c9600bps
	c19200bps
	c38400bps
)

func (br BaudRateCode) String() string {
	switch br {
	case c1200bps:
		return "1200"
	case c2400bps:
		return "2400"
	case c4800bps:
		return "4800"
	case c9600bps:
		return "9600"
	case c19200bps:
		return "19200"
	case c38400bps:
		return "38400"
	}
	return "Undefined"
}

type ADAM4000 struct {
	Address    byte // from config
	address    byte
	InputRange InputRangeCode
	BaudRate   BaudRateCode
	DataFormat DataFormatCode

	Name, Version string

	Value []float64

	Integration_time bool //true = 50Hz, false = 60Hz
	Checksum         bool //true = Checksum enabled

	rc *bufio.Reader
	wc *bufio.Writer

	readChan  chan []byte
	errorChan chan error

	Retries int
}

func (a ADAM4000) String() string {
	return fmt.Sprintf("Address = %d, InputRange = %s, BaudRate = %s, Name = %s, Version = %s, Integration = %t, Checksum = %t, Data Format = %s", a.Address, a.InputRange, a.BaudRate, a.Name, a.Version, a.Integration_time, a.Checksum, a.DataFormat)
}
