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

//Notifica sono le info che si ricevono dai nagios che vengono
//elaborate per creare le chiamate automatiche
type Notifica struct {
	//Time        time.Time `json:"timestamp,omitempty"`
	Hostname    string `json:"hostname,omitempty"`
	Service     string `json:"servizio,omitempty"`
	Piattaforma string `json:"piattaforma,omitempty"`
	Reperibile  string `json:"reperibile,omitempty"`
	Cellulare   string `json:"cellulare,omitempty"`
	Messaggio   string `json:"messaggio,omitempty"`
}

//Dettagli non usato al momento ma servirà a gestire le
//risposte del centralino virtuale asterisk
type Dettagli struct {
	Info  string `json:"info,omitempty"`
	State string `json:"state,omitempty"`
}

var people []Notifica

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
//passata come argomento
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

//SetReper inserisce reperibilità in un archivio condiviso
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

//CreateNotifica riceve gli alerts dei nagios e li utilizza per
//allertare telefonicamente il reperibile in turno
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
		log.Fatal(err.Error())
		return
	}

	result := fmt.Sprintf("Ok. campi ricevuti: Hostname: %s, Service: %s, Piattaforma: %s, Reperibile: %s, Cellulare: %s, Messaggio: %s", hostname, service, piattaforma, reperibile, cellulare, messaggio)

	respondWithJSON(w, http.StatusCreated, result)
	log.Println("ok")

	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), result)

	CreateCall(hostname, service, piattaforma, reperibile, cellulare, messaggio)

	return
}

//CreateCall crea il file .call che serve ad Asterisk per contattare il reperibile
func CreateCall(hostname, service, piattaforma, reperibile, cellulare, messaggio string) (err error) {

	//Trasforma il campo passato in una stringa di 10 numeri
	cell, err := verificaCell(reperibile)
	if err != nil {
		log.Printf("Cellulare non gestibile: %s\n", err.Error())
		return
	}

	scheletro :=
		`Channel: SIP/999` + cell + `@10.31.18.26
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
Set: SRV_MSG="sul server ` + hostname + ` il servizio ` + service + ` è in critical` + messaggio + `"`

	//dove salavare i file in maniera che asterisk li possa scaricare
	//nel nostro caso equivale a dove nginx tiene i contenuti statici del webserver
	//le informazioni sono nel file nascosto .sarumann.yaml che l'utente deve avere
	//nella propria $HOME
	path := viper.GetString("CallPath")
	file, err := os.Create(path + "exampleTest.call") // Truncates if file already exists, be careful!
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer file.Close() // Make sure to close the file when you're done

	_, err = file.WriteString(scheletro)
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
	//fmt.Printf("\nLength: %d bytes", len)
	fmt.Printf("\nFile Name: %s\n", file.Name())
	return
}

//verificaCell verifica che il cell sia una stringa di 10 cifre
func verificaCell(value string) (cell string, err error) {

	//se value ha meno di 10 cifre non è buono
	if len(value) < 10 {
		err := fmt.Errorf("Cellulare con poche cifre: %v", len(value))
		log.Println(err.Error())
		return "", err
	}
	//cell10cifre prende gli ultimi 10 caratteri del value
	cell10cifre := string(value[len(value)-10:])

	//test verifica che il valore sia composto da esattamente 10 cifre
	test := regexp.MustCompile(`^[0-9]{10}$`)
	switch {
	case test.MatchString(cell10cifre) == true:
		cell = cell10cifre
	default:
		cell = ""
		err = fmt.Errorf("Il cellulare non è corretto")
		log.Println(err.Error())
		return "", err
	}

	return

}
