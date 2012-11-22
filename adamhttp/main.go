// Copyright 2012 Thomas Jager <mail@jager.no> All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Driver for ADAM-4000 series I/O Modules from Advantech

package main

import (
        "go-adam4000"
	"fmt"
        "net"
	"net/http"
        "log"
        "strconv"
        "bufio"
        "html/template"
)

var conn net.Conn

func DetailServer(w http.ResponseWriter, req *http.Request) {
    s1, err := template.ParseFiles("templates/header.tmpl", "templates/footer.tmpl", "templates/detail.tmpl")
    if err != nil {
        fmt.Fprintf(w, "Error: %s\n", err)
    }
    req.ParseForm()
    if strunit, ok := req.Form["unit"]; ok {
        if unit, err := strconv.Atoi(strunit[0]); err == nil {
	    adam := adam4000.NewADAM4000(byte(unit), bufio.NewReader(conn), bufio.NewWriter(conn))
            adam.Retries = 1
            adam.GetConfig()
            adam.GetVersion()
            adam.GetName()
            adam.GetAllValue()
            s1.ExecuteTemplate(w, "detail", adam)
        } else {
            fmt.Fprintf(w, "Malformed Unit Def (Not integer)");
        }
    } else {
            fmt.Fprintf(w, "Malformed Unit Def (Missing)");
    }
}

func main() {
        var err error
	conn, err = net.Dial("tcp", "192.168.0.60:5301")
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
        http.HandleFunc("/detail", DetailServer)
        err = http.ListenAndServe(":8084", nil)
        if err != nil {
            log.Fatal("Listen and server: ", err)
        }
}
