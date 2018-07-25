package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//Notifica è lo stuct in arrivo dai nagios
type Notifica struct {
	//Time        time.Time `json:"timestamp,omitempty"`
	Hostname    string `json:"hostname,omitempty"`
	Service     string `json:"servizio,omitempty"`
	Piattaforma string `json:"piattaforma,omitempty"`
	Reperibile  string `json:"reperibile,omitempty"`
	Cellulare   string `json:"cellulare,omitempty"`
	Messaggio   string `json:"messaggio,omitempty"`
}

//LogNotifica archivia la notifica in un Database
func LogNotifica(n Notifica) (err error) {

	//recupera il timestamp di adesso
	timestamp := time.Now()

	//apre il database
	database, _ := sql.Open("sqlite3", "./sarumann.db")

	//chiudi Database una volta fatto
	defer database.Close()

	//prepara la creazione della tabella notifiche se non esite
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS notifiche (id INTEGER PRIMARY KEY, server TEXT, servizio TEXT, piattaforma TEXT, reperibile TEXT, cellulare TEXT, messaggio TEXT, timestamp INT)")
	//esegue la creazione della tabella notifiche se non esiste già nel database
	statement.Exec()

	//prepara l'inserimenti della notifica
	statement, err = database.Prepare("INSERT INTO notifiche(server, servizio, piattaforma, reperibile, cellulare, messaggio, timestamp) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
	}
	//esegue l'inserimento della notifica passata come argomento della funzione
	_, err = statement.Exec(n.Hostname, n.Service, n.Piattaforma, n.Reperibile, n.Cellulare, n.Messaggio, timestamp.Unix())
	if err != nil {
		fmt.Println(err.Error())
	}
	return

}

//AntiStorm evita che il reperibile riceva troppe chiamate
func AntiStorm(piattaforma string) (err error) {
	database, _ := sql.Open("sqlite3", "./sarumann.db")
	defer database.Close()
	row := database.QueryRow("SELECT timestamp FROM notifiche where piattaforma = ? order by timestamp desc limit 1", piattaforma)
	var last string
	row.Scan(&last)
	fmt.Println(last) //debug
	lastint, err := strconv.Atoi(last)
	if err != nil {
		fmt.Println("errore")
	}

	oraepoch := time.Now().Unix()
	//Se non sono passati tot secondi dall'ultima notifica allora esce
	tot := 10 ///1800 secondi uguale mezz'ora
	if lastint+(tot) > int(oraepoch) {
		err = fmt.Errorf("Troppe chiamate al reperibile per %s, è permessa una sola chiamata ogni %d secondi", piattaforma, tot)
		return err
	}
	return nil
}

func main() {

	n := Notifica{
		Hostname:    "srv1",
		Service:     "www",
		Piattaforma: "CDN",
		Reperibile:  "+393333333333",
		Cellulare:   "",
		Messaggio:   "",
	}

	err := AntiStorm(n.Piattaforma)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	LogNotifica(n)

}
