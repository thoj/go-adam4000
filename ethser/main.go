// Copyright 2012 Thomas Jager <mail@jager.no> All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Driver for ADAM-4000 series I/O Modules from Advantech

package main

import (
	"bufio"
        "go-adam4000"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "192.168.0.60:5301")
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	adam := adam4000.NewADAM4000(0, bufio.NewReader(conn), bufio.NewWriter(conn))
	adam.Retries = 1
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
	err = adam.SetChannelRange(4, adam4000.RangeKTc)
	if err != nil {
		fmt.Printf("%s", err)
	}
	rangec, err := adam.GetChannelRange(4)
	fmt.Printf("%s\n", rangec)
	if err != nil {
		fmt.Printf("%s", err)
	}
	_, err = adam.GetAllValue()
	if err != nil {
		fmt.Printf("%s", err)
	}
	val, err := adam.GetChannelValue(3)
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s\n", adam)
	fmt.Printf("%X: %s, %s, %s, %s %v %f\n", adam.Address, adam.BaudRate, adam.InputRange, adam.Version, adam.Name, adam.Value, val)
	adam.Address = 7
	adam.DataFormat = adam4000.DataFormatEngUnits
	adam.SetConfig()
	conn.Close()
}
