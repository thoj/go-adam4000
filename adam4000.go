// Copyright 2009 Thomas Jager <mail@jager.no> All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Driver for ADAM-4000 series I/O Modules from Advantech

package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

func NewADAM4000(addr byte, rc *bufio.Reader, wc *bufio.Writer) *ADAM4000 {
	var a ADAM4000
	a.address = addr
	a.rc = rc
	a.wc = wc
	a.Value = make([]float64, 8)
	a.Retries = 3
	go a.startReader()
	return &a
}

func (a *ADAM4000) GetName() (string, error) {
	resp, err := a.comResF("$%02XM\r", a.address)
	if err != nil {
		return "", err
	}
	a.Name = strings.Trim(string(resp[3:]), "\r ")
	return a.Name, nil
}

func (a *ADAM4000) GetAllValue() ([]float64, error) {
	resp, err := a.comResF("#%02X\r", a.address)
	if err != nil {
		return nil, err
	}
	values := string(resp[1:])
	fmt.Printf("%d\n", len(values))
	if len(values) == 57 {
		a.Value[0], err = strconv.ParseFloat(values[0:7], 64)
		a.Value[1], err = strconv.ParseFloat(values[7:14], 64)
		a.Value[2], err = strconv.ParseFloat(values[14:21], 64)
		a.Value[3], err = strconv.ParseFloat(values[21:28], 64)
		a.Value[4], err = strconv.ParseFloat(values[28:35], 64)
		a.Value[5], err = strconv.ParseFloat(values[35:42], 64)
		a.Value[6], err = strconv.ParseFloat(values[42:49], 64)
		a.Value[7], err = strconv.ParseFloat(values[49:56], 64)
	} else {
		intvals := make([]int64, 8)
		intvals[0], err = strconv.ParseInt(values[0:4], 16, 64)
		intvals[1], err = strconv.ParseInt(values[4:8], 16, 64)
		intvals[2], err = strconv.ParseInt(values[8:12], 16, 64)
		intvals[3], err = strconv.ParseInt(values[12:16], 16, 64)
		intvals[4], err = strconv.ParseInt(values[16:20], 16, 64)
		intvals[5], err = strconv.ParseInt(values[20:24], 16, 64)
		intvals[6], err = strconv.ParseInt(values[24:28], 16, 64)
		intvals[7], err = strconv.ParseInt(values[28:32], 16, 64)
		a.Value[0] = float64(intvals[0])
		a.Value[1] = float64(intvals[1])
		a.Value[2] = float64(intvals[2])
		a.Value[3] = float64(intvals[3])
		a.Value[4] = float64(intvals[4])
		a.Value[5] = float64(intvals[5])
		a.Value[6] = float64(intvals[6])
		a.Value[7] = float64(intvals[7])
	}
	return a.Value, err
}

func (a *ADAM4000) GetChannelValue(n int) (float64, error) {
	resp, err := a.comResF("#%02X%d\r", a.address, n)
	if err != nil {
		return float64(0), err
	}
	values := string(resp[1:])
	if len(values) == 7 {
		a.Value[n], err = strconv.ParseFloat(values[0:7], 64)
	} else {
		intval, _ := strconv.ParseInt(values[0:4], 16, 64)
		a.Value[n] = float64(intval)
	}
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
		return 0, err
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
	fmt.Printf("%X\n", data)
	a.Integration_time = data[0]&byte(1<<7) > 0
	a.Checksum = data[0]&byte(1<<6) > 0
	a.DataFormat = DataFormatCode(data[0] & byte(2))
	if a.Address != a.address {
		fmt.Printf("Warning: Configured address (%d) differs from connected address (%d), in init mode?\n", a.Address, a.address)
	}
	return nil
}

func (a *ADAM4000) SetConfig() error {
	data := byte(a.DataFormat)
	if a.Integration_time {
		data |= byte(1 << 7)
	}
	if a.Checksum {
		data |= byte(1 << 6)
	}

	_, err := a.comResF("%%%02X%02XFF%02X%02X\r", a.address, a.Address, byte(a.BaudRate), data)
	return err
}
