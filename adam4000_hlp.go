// Copyright 2012 Thomas Jager <mail@jager.no> All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Driver for ADAM-4000 series I/O Modules from Advantech

package adam4000

import (
	"errors"
	"fmt"
	"time"
)

func (a *ADAM4000) startReader() {
	a.readChan = make(chan []byte)
	a.errorChan = make(chan error)
	for {
		str, err := a.rc.ReadBytes('\r')
		if err != nil {
			a.errorChan <- err
			return
		}
		fmt.Printf("--> %s\n", str)
		a.readChan <- str
	}
}

func (a *ADAM4000) comResF(format string, va ...interface{}) ([]byte, error) {
	buf := fmt.Sprintf(format, va...)
	fmt.Printf("<-- %s\n", buf)
	_, err := fmt.Fprint(a.wc, buf)
	if err != nil {
		return nil, err
	}
	a.wc.Flush()
	var str []byte
	retry := a.Retries
	for {
		select {
		case err = <-a.errorChan:
			return nil, err
		case str = <-a.readChan:
			return str, nil
		case <-time.After(a.Timeout):
			if retry <= 0 {
				return nil, errors.New("No reply from module")
			}
			retry--
			fmt.Printf("<-- (Retry) %s\n", buf)
			_, err := fmt.Fprint(a.wc, buf)
			if err != nil {
				return nil, err
			}
			a.wc.Flush()

		}
	}
	return nil, nil
}
