package client

import (
	"bytes"
	"encoding/json"
	"io"
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

//SendPost posts data to the chosed server endpoint
func SendPost(endpoint, hostname, service, piatta, rep, cell, msg string) (err error) {

	p := &Notifica{
		Hostname:    hostname,
		Service:     service,
		Piattaforma: piatta,
		Reperibile:  rep,
		Cellulare:   cell,
		Messaggio:   msg,
	}

	b := new(bytes.Buffer)

	json.NewEncoder(b).Encode(p)

	res, err := http.Post(endpoint, "application/json; charset=utf-8", b)
	io.Copy(os.Stdout, res.Body)

	return
}
