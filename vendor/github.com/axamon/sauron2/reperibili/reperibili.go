package reperibili

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"database/sql"

	"github.com/axamon/cripta"
	"github.com/axamon/sauron2/sms"
	//serve per gestire i db sqlite
	_ "github.com/mattn/go-sqlite3"
)

//Contatto info puculiari del reperibile
type Contatto struct {
	Nome      string
	Cognome   string
	Cellulare string
}

//Ruota struct con informazioni reperibili
type Ruota struct {
	ID          int
	Nome        string
	Cognome     string
	Cellulare   string
	Piattaforma string
	Giorno      string
	Gruppo      string
}

var filecsv = flag.String("f", "reperibilita.csv", "Percorso del file csv per la reperibilità")
var piattaforma = flag.String("p", "CDN", "La piattaforma di cui desideri ricavare il reperibile")

var contatti []Ruota

const (
	//il file dove creare il DB
	dbfile = "reperibili.db"

	//Crea la tabella dei reperibili
	createruota = `
	CREATE TABLE IF NOT EXISTS ruota (
		id	integer PRIMARY KEY AUTOINCREMENT,
		nome	varchar ( 255 ),
		cognome	varchar ( 255 ),
		cellulare	varchar ( 255 ),
		piattaforma varchar ( 255 ),
		giorno      varchar ( 255 ),
		gruppo      varchar ( 255 )
	);`
)

//Male non fa
var db *sql.DB

//Opendb Opens the DB on the file system
func Opendb(pathtofile string) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", pathtofile)
	creadb, err := db.Prepare(createruota)
	_, err = creadb.Exec()
	/* if err != nil {
		//fmt.Println(err.Error())
		return nil, fmt.Errorf("Problema ad aprire il DB %s", err.Error())
	}
	*/
	return
}

//AddRuota Aggiunge un reperibile al DB
func AddRuota(nome, cognome, cellulare, piattaforma, giorno, gruppo string) (err error) {
	if ok := sms.Verificacellulare(cellulare); ok != true {
		return fmt.Errorf("Cellulare inserito non nel formato +39(10)cifre")
	}

	db, err := Opendb(dbfile)
	defer db.Close()

	addreperiperbilita, err := db.Prepare("INSERT INTO ruota (nome, cognome, cellulare, piattaforma, giorno, gruppo) VALUES (?, ?,?,?,?,?)")

	_, err = addreperiperbilita.Exec(nome, cognome, cellulare, piattaforma, giorno, gruppo)

	return
}

//IDRuota restituisce id della riga inserita per ultima
func IDRuota(giorno, piattaforma string) (id int, err error) {
	db, err := Opendb(dbfile)
	defer db.Close()

	cercaid, err := db.Prepare("select id from ruota where giorno = ? and piattaforma = ? order by id desc limit 1")

	row := cercaid.QueryRow(giorno, piattaforma)

	row.Scan(&id)

	return
}

//GetReperibile resistuisce le info del reperibile
func GetReperibile(piattaforma string) (Reperibile Contatto, err error) {
	//t è il timestamp di adesso
	var t = time.Now()

	//limite delle 7. Fino alle 7 del mattino seguente il reperibile che viene visualizzato è quello del giorno prima
	var limite7 = time.Date(t.Year(), t.Month(), t.Day(), 7, 0, 0, 0, t.Location())

	var ieri = time.Now().Add(-24 * time.Hour).Format("20060102")
	var oggi = time.Now().Format("20060102")
	//var domani = time.Now().Add(24 * time.Hour).Format("20060102")

	db, err := Opendb(dbfile)
	defer db.Close()

	var giorno string

	switch {
	//Se è prima delle 7 di mattino restituisce il Reperibile del giorno prima
	case t.Before(limite7):
		giorno = ieri
	default:
		giorno = oggi
	}
	if piattaforma == "CDN" {
		giorno = oggi
	}
	getrep, err := db.Prepare("select nome, cognome, cellulare from ruota where giorno= ? and piattaforma=? order by id desc limit 1")

	row := getrep.QueryRow(giorno, piattaforma)

	err = row.Scan(&Reperibile.Nome, &Reperibile.Cognome, &Reperibile.Cellulare)

	if len(Reperibile.Cellulare) == 0 {
		err = fmt.Errorf("Nessun reperibile settato")
	}

	return

}

func recuperavariabile(variabile string) (result string, err error) {
	if result, ok := os.LookupEnv(variabile); ok && len(result) != 0 {
		return result, nil
	}
	return "", fmt.Errorf("la variabile %s non esiste o è vuota", variabile)
}

func recuperavariabilecifrata(variabile, passwd string) (result string, err error) {
	if str, ok := os.LookupEnv(variabile); ok && len(str) != 0 {
		result = cripta.Decifra(str, passwd)
		return result, nil
	}
	return "", fmt.Errorf("la variabile %s non esiste o è vuota", variabile)
}

//Chiamareperibile contatta il reperibile in turno
func Chiamareperibile(TO, NOME, COGNOME string) (sid string, err error) {

	twilionumber, err := recuperavariabile("TWILIONUMBER")

	//Recupera l'accountsid di Twilio dallla variabile d'ambiente
	accountSid, err := recuperavariabile("TWILIOACCOUNTSID")

	//Recupera il token supersegreto dalla variabile d'ambiente
	authToken, err := recuperavariabile("TWILIOAUTHTOKEN")

	//Questa è la url di Twilio per le chiamate vocali
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Calls.json"

	v := url.Values{}
	v.Set("To", TO)
	v.Set("From", twilionumber)

	//Questa è la url in cui si possono aggiungere i campi da far pronunciare a Twilio
	//EH5cef42aa1454fc2326780c8f08c6d568 è l'identificativo del twiml da richiamare
	v.Set("Url", "https://handler.twilio.com/twiml/EH5cef42aa1454fc2326780c8f08c6d568?NOME="+NOME+"&COGNOME="+COGNOME)
	rb := *strings.NewReader(v.Encode())

	//Crea il client http
	client := &http.Client{}

	req, err := http.NewRequest("POST", urlStr, &rb)

	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)

	// make request
	var data map[string]interface{}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {

		bodyBytes, errb := ioutil.ReadAll(resp.Body)
		if errb != nil {
			fmt.Fprintln(os.Stdout, errb.Error())
		}
		err := json.Unmarshal(bodyBytes, &data)
		if err != nil {
			return "", err
		}
	}
	defer resp.Body.Close()
	//fmt.Println(data) //debug

	//se la mappa contiene un valore per sid lo ritorna
	if val, ok := data["sid"]; ok {
		sid = val.(string)
	}

	return

}

/* func main() {
	AddRuota("Alberto", "IERI", "+393357291533", "CDN", "20180610", "GRUPPO6")
	AddRuota("Alberto", "Oggi", "+393357291533", "CDN", "20180611", "GRUPPO6")
	AddRuota("Alberto", "IERI", "+393357291533", "APS", "20180610", "GRUPPO6")
	AddRuota("Alberto", "Oggi", "+393357291533", "APS", "20180611", "GRUPPO6")

	rep, err := GetReperibile("CDN")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("CDN", rep.Cognome)
	rep, _ = GetReperibile("APS")
	fmt.Println("APS", rep.Cognome)

	sid, _ := Chiamareperibile(rep.Cellulare, rep.Nome, rep.Cognome)

	fmt.Print(sid)

} */
