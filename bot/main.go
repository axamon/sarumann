package main

import (
	"fmt"
	"log"
	"os"
	"time"

	raven "github.com/getsentry/raven-go"
	tb "gopkg.in/tucnak/telebot.v2"
)

func recuperavariabile(variabile string) (result string, err error) {
	if result, ok := os.LookupEnv(variabile); ok && len(result) != 0 {
		return result, nil
	}
	err = fmt.Errorf("la variabile %s non esiste o è vuota", variabile)
	fmt.Fprintln(os.Stderr, err.Error())
	raven.CaptureError(err, nil)
	return "", err
}

func main() {

	//Recupera la variabile d'ambiente
	TELEGRAMTOKEN, err := recuperavariabile("TELEGRAMTOKEN")

	b, err := tb.NewBot(tb.Settings{
		Token:  TELEGRAMTOKEN,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/pd", func(m *tb.Message) {
		b.Send(m.Sender, "mannaggia al pd in eterno")
	})

	b.Handle("/saluto", func(m *tb.Message) {
		b.Send(m.Sender, "ciao a te straniero")
	})

	b.Handle("alberto", func(m *tb.Message) {
		b.Send(m.Sender, "è un genio")
	})

	b.Handle("/sarumann", func(m *tb.Message) {
		// photos only
		p := &tb.Photo{File: tb.FromDisk("./sarumann.jpg")}
		b.Send(m.Sender, p)
	})

	b.Start()
}
