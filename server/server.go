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

// Package server provides the ability to forward Nagios Notifications
//to an Asterisk phone sytem.
package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/spf13/viper"

	"github.com/axamon/sauron2/reperibili"
	"github.com/gorilla/mux"
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

//Dettagli non usato al momento
type Dettagli struct {
	Info  string `json:"info,omitempty"`
	State string `json:"state,omitempty"`
}

var people []Notifica

//CreatePersonEndpoint non usata al momento
func CreatePersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	fmt.Println(params)
	var errore Notifica
	err := json.NewDecoder(req.Body).Decode(&errore)
	if err != nil {
		fmt.Println(err)
	}
	errore.Hostname = params["hostname"]
	errore.Service = params["servizio"]
	errore.Messaggio = params["messaggio"]
	json.NewEncoder(w).Encode(errore)
	fmt.Println(errore)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

//GetReper recupera il reperibile attuale per la piattaforma
func GetReper(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Accel-Expires", "0")
	vars := mux.Vars(r)

	piattaforma := vars["piatta"]

	reperibile, err := reperibili.GetReperibile(piattaforma)
	if err != nil {
		respondWithError(w, http.StatusNoContent, err.Error())
		return
	}

	result := fmt.Sprintf("Il reperibile per %s è: %s. Cell: %s", piattaforma, reperibile.Cognome, reperibile.Cellulare)

	respondWithJSON(w, http.StatusFound, result)
	return
}

//SetReper inserisce reperibilità
func SetReper(w http.ResponseWriter, r *http.Request) {
	/* var p reperibili.Contatto
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close() */

	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Accel-Expires", "0")
	r.ParseForm()

	nome := r.PostFormValue("nome")
	cognome := r.PostFormValue("cognome")
	cellulare := r.PostFormValue("cellulare")
	piattaforma := r.PostFormValue("piattaforma")

	oggi := time.Now().Format("20060102")

	err := reperibili.AddRuota(nome, cognome, cellulare, piattaforma, oggi, "gruppo6")
	if err != nil {
		fmt.Println("errorone", err.Error(), cellulare)
		return
	}
	fmt.Println("inserito reperibile: ", nome, cognome, cellulare)

	/* err := reperibili.AddRuota(p.Nome, p.Cognome, p.Cellulare, "CDN", "20180101", "gruppo6")
	if err != nil {
		fmt.Println("errorone")
	}
	*/
	/* 	fmt.Println(p)

	respondWithJSON(w, http.StatusCreated, p) */
	return
}

//CreateNotifica riceve gli alerts
func CreateNotifica(w http.ResponseWriter, r *http.Request) {

	//Crea p come tipo Notifica con i suoi structs
	var p Notifica
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	//fmt.Println(p) //debug

	hostname, err := url.QueryUnescape(p.Hostname)

	service, err := url.QueryUnescape(p.Service)

	piattaforma, err := url.QueryUnescape(p.Piattaforma)

	reperibile, err := url.QueryUnescape(p.Reperibile)

	cellulare, err := url.QueryUnescape(p.Cellulare)

	messaggio, err := url.QueryUnescape(p.Messaggio)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	result := fmt.Sprintf("Ok. campi ricevuti: Hostname: %s, Service: %s, Piattaforma: %s, Reperibile: %s, Cellulare: %s, Messaggio: %s", hostname, service, piattaforma, reperibile, cellulare, messaggio)

	respondWithJSON(w, http.StatusCreated, result)

	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), result)

	if ok := verificacampo(reperibile); ok == true {
		fmt.Println("reperibile è un cell: ", reperibile)
	}

	scheletro :=
		`Channel: SIP/999` + reperibile + `@10.31.18.26
MaxRetries: 5 
RetryTime: 300 
WaitTime: 60 
Context: nagios-notify 
Extension: s 
Archive: Yes 
Set: CONTACT_NAME="Gringo" 
Set: PLAT_NAME="` + piattaforma + `" 
Set: NOT_TYPE="PROBLEM" 
Set: HOST_ALIAS="` + hostname + `" 
Set: SERVICE_NAME="` + service + `" 
Set: STATUS="Critico" 
Set: NOT_HEAD_MSG="è stato riscontrato un problema" 
Set: SRV_MSG="il server ` + hostname + ` è in critical a causa di ` + service + `"`

	CreateCall(scheletro)

	return
}

//CreateCall crea il file .call che serve ad Asterisk
func CreateCall(notifica string) {
	path := viper.GetString("CallPath")
	file, err := os.Create(path + "exampleTest.call") // Truncates if file already exists, be careful!
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer file.Close() // Make sure to close the file when you're done

	len, err := file.WriteString(notifica)
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
	fmt.Printf("\nLength: %d bytes", len)
	fmt.Printf("\nFile Name: %s", file.Name())
}

func verificacampo(value string) (ok bool) {
	test := regexp.MustCompile(`^.*[0-9]{10}$`)
	return test.MatchString(value)

}
