// Copyright (C) 2018 by Alberto Bregliano <alberto.bregliano@pm.me>
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package client provides the ability to test sarumann server
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

//Notifica sono le info che si ricevono dai nagios
type Notifica struct {
	//Time        time.Time `json:"timestamp,omitempty"`
	Hostname    string `json:"hostname,omitempty"`
	Service     string `json:"servizio,omitempty"`
	Piattaforma string `json:"piattaforma,omitempty"`
	Reperibile  string `json:"reperibile,omitempty"`
	Cellulare   string `json:"cellulare,omitempty"`
	Messaggio   string `json:"messaggio,omitempty"`
}

//SendPost invia via POST le variabili ricevute
//endpoint è la url dove si vuole inviare il POST
func SendPost(endpoint, hostname, service, piatta, rep, cell, msg string) (err error) {

	p := &Notifica{
		Hostname:    hostname,
		Service:     service,
		Piattaforma: piatta,
		Reperibile:  rep,
		Cellulare:   cell,
		Messaggio:   msg,
	}

	//Crea un nuovo buffer
	b := new(bytes.Buffer)

	//Encoda p in b
	json.NewEncoder(b).Encode(p)

	//invia a endpoint b con le informazioni di p encodade in json via POST
	res, err := http.Post(endpoint, "application/json; charset=utf-8", b)
	if err != nil {
		err = fmt.Errorf("Problema creazione post %s", err.Error())
		log.Fatal(err.Error())
	}

	//serve per vedere il body della response
	io.Copy(os.Stdout, res.Body)

	return
}
