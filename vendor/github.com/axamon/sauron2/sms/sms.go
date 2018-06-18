package sms

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	raven "github.com/getsentry/raven-go"
)

func init() {
	raven.SetDSN("https://3c8659c6cced4338a494519ca736de01:a9831014cc5540468e380a7bb5343108@sentry.io/1222500")
}

//Verificacellulare si assicura che il cellulare inserito sia nel formato corretto
func Verificacellulare(CELLULARE string) (ok bool) {

	re := regexp.MustCompile(`^\+39[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]$`)
	return re.MatchString(CELLULARE)

}

//Quanto sono compliant con almeno alcune delle 12 best practices di GO!
//https://talks.golang.org/2013/bestpractices.slide#1
func recuperavariabile(variabile string) (result string, err error) {
	if result, ok := os.LookupEnv(variabile); ok && len(result) != 0 {
		return result, nil
	}
	err = fmt.Errorf("la variabile %s non esiste o è vuota", variabile)
	fmt.Fprintln(os.Stderr, err.Error())
	raven.CaptureError(err, nil)
	return "", err
}

//Inviasms invia sms via Twilio
func Inviasms(to, body string) (result string, err error) {

	//Recupera il numero di Twilio dallla variabile d'ambiente
	TWILIONUMBER, err := recuperavariabile("TWILIONUMBER")

	//Recupera TWILIOACCOUNTSID  dallla variabile d'ambiente
	TWILIOACCOUNTSID, err := recuperavariabile("TWILIOTWILIOACCOUNTSID")

	//Recupera il token supersegreto dalla variabile d'ambiente
	TWILIOAUTHTOKEN, err := recuperavariabile("TWILIOTWILIOAUTHTOKEN")

	//TODO vedere se riesce a prendere anche le variabili da ambiente windows...
	//...ma anche no! :)

	//Crea la URL necessaria per richiamare la funzionalità degli SMS di Twilio
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + TWILIOACCOUNTSID + "/Messages.json"

	//Valorizza i campi per l'invio del SMS
	v := url.Values{}
	v.Set("To", to)             //Esempio: "+393357291532"
	v.Set("From", TWILIONUMBER) //Esempio "+17372041296"
	v.Set("Body", body)

	//impacchettiamo tutte le variabile insieme
	rb := *strings.NewReader(v.Encode())

	//Creiamo un client http
	client := &http.Client{}

	//Creiamo la http request da inviare dopo
	req, err := http.NewRequest("POST", urlStr, &rb)

	//Utiliziamo l'autenticazione basic
	req.SetBasicAuth(TWILIOACCOUNTSID, TWILIOAUTHTOKEN)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Inviamo la request e salviamo la http response
	resp, err := client.Do(req)

	if err != nil {
		err = fmt.Errorf("Errore nella ricesione response: %s", err.Error())
		fmt.Fprintln(os.Stderr, err.Error())
		raven.CaptureError(err, nil)
		//return "", err
	}

	result = resp.Status

	//controlliamo che ha da dire la response
	//Restituisce codice e significato, se ricevi 201 CREATED allora è ok.
	//fmt.Println(resp.Status)

	raven.CaptureMessage(resp.Status, nil)

	return
}
