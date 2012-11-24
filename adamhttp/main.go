// Copyright 2012 Thomas Jager <mail@jager.no> All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Driver for ADAM-4000 series I/O Modules from Advantech

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go-adam4000"
	"html/template"
	"log"
	"net"
	"net/http"
	"strconv"
)

type TemplateInventory struct {
	Units map[int]*adam4000.ADAM4000

	Scanning         bool
	Scanning_Address int
}

var scanning bool
var scanning_address int

var units map[int]*adam4000.ADAM4000

var conn net.Conn

func JsonServer(w http.ResponseWriter, req *http.Request) {
	enc := json.NewEncoder(w)
	req.ParseForm()
	if strunit, ok := req.Form["unit"]; ok {
		if unit, err := strconv.Atoi(strunit[0]); err == nil {
			if _, ok = units[unit]; !ok {
				units[unit] = adam4000.NewADAM4000(byte(unit), bufio.NewReader(conn), bufio.NewWriter(conn))
			}
			units[unit].Retries = 1
			values, err := units[unit].GetAllValue()
			if err != nil {
				enc.Encode(map[string]interface{}{"error": true, "errorstr": "Unit not found"})
				return
			}
			enc.Encode(map[string]interface{}{"values": values, "error": false})
		} else {
			enc.Encode(map[string]interface{}{"error": true, "errorstr": "Malformed Unit Def (Not integer)"})
		}
	} else {
		enc.Encode(map[string]interface{}{"error": true, "errorstr": "Missing Unit Def"})
	}
}

func DetailServer(w http.ResponseWriter, req *http.Request) {
	s1, err := template.ParseFiles("templates/header.tmpl", "templates/footer.tmpl", "templates/detail.tmpl")
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
	}
	req.ParseForm()
	if strunit, ok := req.Form["unit"]; ok {
		if unit, err := strconv.Atoi(strunit[0]); err == nil {
			if _, ok = units[unit]; !ok {
				units[unit] = adam4000.NewADAM4000(byte(unit), bufio.NewReader(conn), bufio.NewWriter(conn))
			}
			if _, ok := req.Form["setconfig"]; ok {
				if setaddress_str, ok := req.Form["setaddress"]; ok {
					if setaddress, err := strconv.Atoi(setaddress_str[0]); err == nil && setaddress <= 255 {
						units[unit].Address = byte(setaddress)
					}
				}
				if setinputrange_str, ok := req.Form["setinputrange"]; ok {
					if setinputrange, err := strconv.Atoi(setinputrange_str[0]); err == nil && setinputrange <= 255 {
						units[unit].InputRange = adam4000.InputRangeCode(setinputrange)
					}
				}
				units[unit].SetConfig()
			}
			units[unit].Retries = 1
			err := units[unit].GetConfig()
			if err != nil {
				fmt.Fprintf(w, "unit at address %d not found.", unit)
				return
			}
			units[unit].GetVersion()
			units[unit].GetName()
			units[unit].GetAllValue()
			s1.ExecuteTemplate(w, "detail", units[unit])
		} else {
			fmt.Fprintf(w, "Malformed Unit Def (Not integer)")
		}
	} else {
		fmt.Fprintf(w, "Malformed Unit Def (Missing)")
	}
}

func OverviewServer(w http.ResponseWriter, req *http.Request) {
	s1, err := template.ParseFiles("templates/header.tmpl", "templates/footer.tmpl", "templates/overview.tmpl")
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}
	req.ParseForm()
	if _, ok := req.Form["scan"]; ok && !scanning {
		scanning = true
		go ADAMScanner()
	}
	if _, ok := req.Form["stopscan"]; ok {
		scanning = false
	}
	d1 := &TemplateInventory{units, scanning, scanning_address}
	s1.ExecuteTemplate(w, "overview", d1)
}

func ADAMScanner() {
	for i := 0; i <= 255 && scanning; i++ {
		scanning_address = i
		adam := adam4000.NewADAM4000(byte(i), bufio.NewReader(conn), bufio.NewWriter(conn))
		err := adam.GetConfig()
		if err != nil {
			continue
		}
		adam.GetVersion()
		adam.GetName()
		adam.GetAllValue()
		units[i] = adam
	}
	scanning = false
	scanning_address = 0
}

func main() {
	var err error
	scanning = false
	units = make(map[int]*adam4000.ADAM4000)
	conn, err = net.Dial("tcp", "192.168.0.60:5301")
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	http.HandleFunc("/", OverviewServer)
	http.HandleFunc("/detail", DetailServer)
	http.HandleFunc("/json/data", JsonServer)
	err = http.ListenAndServe(":8084", nil)
	if err != nil {
		log.Fatal("Listen and server: ", err)
	}
}
