// Copyright 2009 Thomas Jager <mail@jager.no> All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Driver for ADAM-4000 series I/O Modules from Advantech

package main

import (
	"errors"
	"fmt"
)

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
