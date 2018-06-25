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
	"net/url"
	"os"
	"regexp"
)

const (
	//Credenziali impostate su NGINX del server
	username = "sarumann"
	password = "pippo"
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

func verificaCell(value string) (err error) {
	test := regexp.MustCompile(`^.*[0-9]{10}$`)
	switch {
	case test.MatchString(value) == true:
		err = nil
	default:
		errore := fmt.Sprintf("Il cellulare non è corretto, %s", value)
		err = fmt.Errorf(errore)
		//log.Fatal(err.Error())
	}
	return
}

func verificaMsg(value string) (err error) {
	switch {
	case len(value) < 50:
		err = nil
	default:
		err = fmt.Errorf("Messaggio eccede 50 caratteri")
		//log.Fatal(err.Error())
	}
	return
}

func verificaEndpoint(value string) (err error) {

	_, err = url.ParseRequestURI(value)

	return
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

	err = verificaCell(p.Cellulare)

	err = verificaMsg(p.Messaggio)

	err = verificaEndpoint(endpoint)

	//Crea un nuovo buffer
	b := new(bytes.Buffer)

	//Encoda la notifica p nel buffer b
	err = json.NewEncoder(b).Encode(p)

	//Prepara il client http
	client := &http.Client{}

	//Imposta la request
	req, err := http.NewRequest("POST", endpoint, b)

	//Aggiunge username e password impostate sul server web di arrivo
	req.SetBasicAuth(username, password)

	//Avvia il client e riceve la response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	//bodyText, err := ioutil.ReadAll(resp.Body)
	io.Copy(os.Stdout, resp.Body)

	return
}
