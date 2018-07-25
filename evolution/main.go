package main

import (
	"database/sql"
	"fmt"
	"time"
)

type azioni interface {
	chiamaReperibile() error
	archivia() error
}

type notifica struct {
	ruotaReperibilità string
	gruppo            string
}

func (n notifica) chiamaReperibile() error {
	
}

func (n notifica) archivia() error {
	timestamp := time.Now()
	database, _ := sql.Open("sqlite3", "./sarumann.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS notifiche (id INTEGER PRIMARY KEY, server TEXT, servizio TEXT, piattaforma TEXT, reperibile TEXT, cellulare TEXT, messaggio TEXT, timestamp INT)")
	statement.Exec()
	statement, err := database.Prepare("INSERT INTO notifiche (server, servizio, piattaforma, reperibile, cellulare, messaggio, timestamp) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = statement.Exec(server, servizio, piattaforma, reperibile, cellulare, messaggio, timestamp.Unix())
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}
