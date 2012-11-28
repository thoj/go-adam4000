// Copyright 2012 Thomas Jager <mail@jager.no> All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Driver for ADAM-4000 series I/O Modules from Advantech

package adam4000

import (
	"bufio"
	"fmt"
	"time"
)

type DataFormatCode byte

const (
	DataFormatEngUnits = iota
	DataFormatPercentFSR
	DataFormatTwosComplement
	DataFormatOhms
)

func (df DataFormatCode) String() string {
	switch df {
	case DataFormatEngUnits:
		return "Engineering Units"
	case DataFormatPercentFSR:
		return "% of FSR"
	case DataFormatTwosComplement:
		return "Twos Complement"
	case DataFormatOhms:
		return "Ohms"
	}
	return "Undefined"
}

type InputRangeCode byte

const (
	Range15mV InputRangeCode = iota
	Range50mV
	Range100mV
	Range500mV
	Range1V
	Range2_5V
	_
	Range4_20mA
	Range10V
	Range5V
	_
	_
	_
	Range20mA
	_
	RangeJTc
	RangeKTc
	RangeTTc
	RangeETc
	RangeRTc
	RangeSTc
	RangeBTc
)

func (br InputRangeCode) String() string {
	switch br {
	case Range15mV:
		return "+/- 15mV"
	case Range50mV:
		return "+/- 50mV"
	case Range100mV:
		return "+/- 100mV"
	case Range500mV:
		return "+/- 500mV"
	case Range1V:
		return "+/- 1V"
	case Range2_5V:
		return "+/- 2.5V"
	case Range4_20mA:
		return "4~20mA"
	case Range10V:
		return "+/- 10V"
	case Range5V:
		return "+/- 5V"
	case Range20mA:
		return "+/- 20mA"
	case RangeJTc:
		return "Type-J TC"
	case RangeKTc:
		return "Type-K TC"
	case RangeTTc:
		return "Type-T TC"
	case RangeETc:
		return "Type-E TC"
	case RangeRTc:
		return "Type-R TC"
	case RangeSTc:
		return "Type-S TC"
	case RangeBTc:
		return "Type-B TC"
	}
	return "Undefined"
}

type BaudRateCode byte

const (
	BaudRate1200bps BaudRateCode = (iota + 0x03)
	BaudRate2400bps
	BaudRate4800bps
	BaudRate9600bps
	BaudRate19200bps
	BaudRate38400bps
)

func (br BaudRateCode) String() string {
	switch br {
	case BaudRate1200bps:
		return "1200"
	case BaudRate2400bps:
		return "2400"
	case BaudRate4800bps:
		return "4800"
	case BaudRate9600bps:
		return "9600"
	case BaudRate19200bps:
		return "19200"
	case BaudRate38400bps:
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
	Timeout time.Duration
}

func (a ADAM4000) String() string {
	return fmt.Sprintf("Address = %d, InputRange = %s, BaudRate = %s, Name = %s, Version = %s, Integration = %t, Checksum = %t, Data Format = %s", a.Address, a.InputRange, a.BaudRate, a.Name, a.Version, a.Integration_time, a.Checksum, a.DataFormat)
}
