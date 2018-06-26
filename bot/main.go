package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	raven "github.com/getsentry/raven-go"
	tb "gopkg.in/tucnak/telebot.v2"
)

func emme3(cache string) (image string, err error) {

	ora := time.Now().Unix()
	ore24 := time.Now().Add(-24 * time.Hour).Unix()

	orat := strconv.FormatInt(ora, 10)
	ore24t := strconv.FormatInt(ore24, 10)
	fmt.Println(orat, ore24)

	URL := "http://localhost/cdn/pnp4nagios/index.php/image?host=" + cache + "&srv=FILE_DESCRIPTORS&theme=multisite&baseurl=..%2Fcheck_mk%2F&view=2&source=1&start=" + ore24t + "&end=" + orat

	fmt.Println(URL)
	//Prepara il client http
	client := &http.Client{}

	//Imposta la request
	req, err := http.NewRequest("GET", URL, nil)

	//Aggiunge username e password impostate sul server web di arrivo
	req.SetBasicAuth("omdadmin", "omd")

	//Avvia il client e riceve la response
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	//open a file for writing
	image = "/tmp/" + cache + ".jpg"
	file, err := os.Create(image)
	if err != nil {
		log.Fatal(err)
	}
	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, res.Body)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	fmt.Println("Success!")
	return
}

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
	if err != nil {
		log.Fatal(err)
		return
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  TELEGRAMTOKEN,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

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

	b.Handle("fd", func(m *tb.Message) {
		// photos only
		image, err := emme3("se-rm3-14")
		if err != nil {
			log.Println(err.Error())
		}

		p := &tb.Photo{File: tb.FromDisk(image)}
		b.Send(m.Sender, p)
	})

	b.Start()
}
