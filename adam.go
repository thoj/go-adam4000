// Copyright 2009 Thomas Jager <mail@jager.no> All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Driver for ADAM-4000 series I/O Modules from Advantech

package main

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type DataFormatCode byte

const (
	EngineeringUnits = 0
	PercentofFSR
	TwosComplement
	Ohms
)

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

	Name, Version string

	Value []float64

	integration_time, checksum bool

	rc *bufio.Reader
	wc *bufio.Writer
}

func main() {
	conn, err := net.Dial("tcp", "192.168.0.60:5301")
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	adam := NewADAM4000(0, bufio.NewReader(conn), bufio.NewWriter(conn))
	err = adam.GetConfig()
	if err != nil {
		fmt.Printf("%s", err)
	}
	_, err = adam.GetVersion()
	if err != nil {
		fmt.Printf("%s", err)
	}
	_, err = adam.GetName()
	if err != nil {
		fmt.Printf("%s", err)
	}
        err = adam.SetChannelRange(4, cKTc)
	if err != nil {
		fmt.Printf("%s", err)
	}
        rangec, err := adam.GetChannelRange(4)
        fmt.Printf("%s\n", rangec)
	if err != nil {
		fmt.Printf("%s", err)
	}
	_, err = adam.ReadAll()
	if err != nil {
		fmt.Printf("%s", err)
	}
	val, err := adam.ReadChannel(3)
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%X: %s, %s, %s, %s %v %f\n", adam.Address, adam.BaudRate, adam.InputRange, adam.Version, adam.Name, adam.Value, val)
	conn.Close()
}

func NewADAM4000(addr byte, rc *bufio.Reader, wc *bufio.Writer) *ADAM4000 {
	return &ADAM4000{address: addr, rc: rc, wc: wc, Value: make([]float64, 8)}
}

//Data Format 
//Bit   1               2               3   4   5   6    7  8
//      Integration     Checksum        n/a N/A n/A N/A  dATA
//      Time            Status                           fORMAT

//AANNTTCCFF\r
// AA = Address
// NN = New Address
// TT = Input Range (4015, 4019 = 00)
// CC = Baud Rate
// FF = Data Format (8 bits) (00: Engineering Units, 01: % of FSR, 10: Two's compliment, 11: Ohms)

func (a *ADAM4000) comResF(format string, va ...interface{}) ([]byte, error) {
        buf := fmt.Sprintf(format, va...)
	fmt.Printf("<-- %s\n", buf)
	_, err := fmt.Fprint(a.wc, buf)
	if err != nil {
		return nil, err
	}
	a.wc.Flush()
	str, err := a.rc.ReadBytes('\r')
	fmt.Printf("--> %s\n", str)
	if err != nil {
		return nil, err
	}
	if str[0] == '?' {
		return nil, errors.New("Module returned invalid command .")
	}
	return str, nil
}

func (a *ADAM4000) GetName() (string, error) {
	resp, err := a.comResF("$%02XM\r", a.address)
	if err != nil {
		return "", err
	}
	a.Name = strings.Trim(string(resp[3:]), "\r ")
	return a.Name, nil
}

func (a *ADAM4000) ReadAll() ([]float64, error) {
	resp, err := a.comResF("#%02X\r", a.address)
	if err != nil {
		return nil, err
	}
	values := string(resp[1:])
	a.Value[0], err = strconv.ParseFloat(values[0:7], 64)
	a.Value[1], err = strconv.ParseFloat(values[7:14], 64)
	a.Value[2], err = strconv.ParseFloat(values[14:21], 64)
	a.Value[3], err = strconv.ParseFloat(values[21:28], 64)
	a.Value[4], err = strconv.ParseFloat(values[28:35], 64)
	a.Value[5], err = strconv.ParseFloat(values[35:42], 64)
	a.Value[6], err = strconv.ParseFloat(values[42:49], 64)
	a.Value[7], err = strconv.ParseFloat(values[49:56], 64)
	return a.Value, err
}

func (a *ADAM4000) ReadChannel(n int) (float64, error) {
	resp, err := a.comResF("#%02X%d\r", a.address, n)
	if err != nil {
		return float64(0), err
	}
	values := string(resp[1:])
	a.Value[n], err = strconv.ParseFloat(values[0:7], 64)
	return a.Value[n], err
}

func (a *ADAM4000) GetVersion() (string, error) {
	resp, err := a.comResF("$%02XF\r", a.address)
	if err != nil {
		return "", err
	}
	a.Version = strings.Trim(string(resp[3:]), "\r ")
	return a.Version, nil
}

func (a *ADAM4000) SetChannelRange(channel int, rangec InputRangeCode) error {
	_, err := a.comResF("$%02X7C%dR%02X\r", a.address, channel, byte(rangec))
	if err != nil {
		return err
	}
	return nil
}

func (a *ADAM4000) GetChannelRange(channel int) (InputRangeCode, error) {
	resp, err := a.comResF("$%02X8C%d\r", a.address, channel)
	if err != nil {
		return 0,err
	}
        rangec := make([]byte, 1)
	hex.Decode(rangec, resp[6:8])
	return InputRangeCode(rangec[0]), nil
}

func (a *ADAM4000) SyncronizeRead() error {
	//Stub
        return nil
}

func (a *ADAM4000) SyncronizedValue() ([]float64, error) {
	//Stub
        return nil, nil
}

func (a *ADAM4000) GetConfig() error {
	resp, err := a.comResF("$%02X2\r", a.address)
	if err != nil {
		return err
	}

	addr := make([]byte, 1)
	typecode := make([]byte, 1)
	baud := make([]byte, 1)
	data := make([]byte, 1)

	hex.Decode(addr, resp[1:3])
	hex.Decode(typecode, resp[3:5])
	hex.Decode(baud, resp[5:7])
	hex.Decode(data, resp[7:9])

	a.Address = addr[0]
	a.InputRange = InputRangeCode(typecode[0])
	a.BaudRate = BaudRateCode(baud[0])
	return nil
}
