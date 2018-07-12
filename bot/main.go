package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	raven "github.com/getsentry/raven-go"
	tb "gopkg.in/tucnak/telebot.v2"
)

func emme3(cache, command string) (image string, err error) {

	ora := time.Now().Unix()
	ore24 := time.Now().Add(-24 * time.Hour).Unix()

	orat := strconv.FormatInt(ora, 10)
	ore24t := strconv.FormatInt(ore24, 10)
	fmt.Println(orat, ore24)
	var URL string
	switch {
	case command == "fd":

		URL = "http://localhost/cdn/pnp4nagios/index.php/image?host=" + cache + "&srv=FILE_DESCRIPTORS&theme=multisite&baseurl=..%2Fcheck_mk%2F&view=2&source=0&start=" + ore24t + "&end=" + orat

	case command == "we":
		URL = "http://localhost/cdn/pnp4nagios/index.php/image?host=" + cache + "&srv=WebEngine&theme=multisite&baseurl=..%2Fcheck_mk%2F&view=1&source=0&start=" + ore24t + "&end=" + orat

	}
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
		log.Println(err)
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
	/* 	o := tb.InlineButton{Text: "se-rm3-14", URL: "http://www.google.it", Data: "se-rm3-14", InlineQuery: "ooo",
	Action: func(emme3(string))} */

	//Recupera la variabile d'ambiente
	TELEGRAMTOKEN, err := recuperavariabile("TELEGRAMTOKEN")
	if err != nil {
		log.Fatal(err)
		return
	}

	TELEGRAMTOKEN = "608145657:AAEUUw27zd41mOiPBQJzgr1QKzYwataFQrM"

	b, err := tb.NewBot(tb.Settings{
		Token:  TELEGRAMTOKEN,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	b.Handle("/pd", func(m *tb.Message) {

		b.Send(m.Chat, "mannaggia al pd in eterno")
	})

	b.Handle("/saluto", func(m *tb.Message) {
		b.Send(m.Sender, "ciao a te straniero")
	})

	b.Handle("/alberto", func(m *tb.Message) {
		b.Send(m.Chat, "Alberto è geniale")
	})

	b.Handle("/version", func(m *tb.Message) {
		b.Send(m.Chat, "Sarumann_bot v2.4.1 beta")
	})

	b.Handle("/sarumann", func(m *tb.Message) {
		// photos only
		p := &tb.Photo{File: tb.FromDisk("./sarumann.jpg")}
		b.Send(m.Chat, p)
	})

	b.Handle("Aggiornamenti", func(m *tb.Message) {
		ticker := time.NewTicker(5 * time.Second)
		go func() {
			for range ticker.C {
				b.Send(m.Sender, "ciao")
			}
		}()
		time.Sleep(20 * time.Second)
		ticker.Stop()
	})
	/* b.Handle("/fd", func(m *tb.Message) {
		// photos only
		image, err := emme3("se-rm3-14")
		if err != nil {
			log.Println(err.Error())
		}

		p := &tb.Photo{File: tb.FromDisk(image)}
		b.Send(m.Chat, p)
	}) */

	// This button will be displayed in user's
	// reply keyboard.
	si := tb.ReplyButton{Text: "Si"}

	no := tb.ReplyButton{Text: "No"}

	replyKeys := [][]tb.ReplyButton{
		[]tb.ReplyButton{si},
		[]tb.ReplyButton{no},
		// ...
	}

	b.Handle(&si, func(m *tb.Message) {
		// on reply button pressed
		emme3(m.Text, "fd")
	})

	// And this one — just under the message itself.
	// Pressing it will cause the client to send
	// the bot a callback.
	//
	// Make sure Unique stays unique as it has to be
	// for callback routing to work.
	inlineBtn := tb.InlineButton{
		Unique: "sad_moon",
		Text:   "🌚 Button #2",
	}
	inlineKeys := [][]tb.InlineButton{
		[]tb.InlineButton{inlineBtn},
		// ...
	}

	b.Handle(&inlineBtn, func(c *tb.Callback) {
		// on inline button pressed (callback!)

		// always respond!
		//b.Respond(c, &tb.CallbackResponse{...})
	})

	// Command: /start <PAYLOAD>
	b.Handle("/start", func(m *tb.Message) {
		if !m.Private() {
			return
		}
		b.Send(m.Chat, "vuoi verificare i file descriptors della cache?", &tb.ReplyMarkup{
			InlineKeyboard:      inlineKeys,
			ReplyKeyboard:       replyKeys,
			OneTimeKeyboard:     true,
			ResizeReplyKeyboard: true,
		})

	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		//b.Send(m.Sender, m.Text)
		//cerca se nella stringa di testo è presente fd + cache
		fd, _ := regexp.Compile(`^[fF]d\sse-[a-z]+[0-9]-[0-9]?`)
		if fd.MatchString(m.Text) == true {
			cache := strings.Split(m.Text, " ")[1]

			image, err := emme3(cache, "fd")
			if err != nil {
				log.Println(err.Error())
			}

			p := &tb.Photo{File: tb.FromDisk(image)}

			msg := fmt.Sprintf("Ecco i file descriptors delle ultime 24 ore per: %s", cache)
			b.Reply(m, msg)
			b.Send(m.Chat, p)
		}

		we, _ := regexp.Compile(`^[wW]e\sse-[a-z]+[0-9]-[0-9]?`)
		if we.MatchString(m.Text) == true {
			cache := strings.Split(m.Text, " ")[1]

			image, err := emme3(cache, "we")
			if err != nil {
				log.Println(err.Error())
			}

			p := &tb.Photo{File: tb.FromDisk(image)}

			msg := fmt.Sprintf("Ecco l'andamento webengine delle ultime 24 ore per: %s", cache)
			b.Reply(m, msg)
			b.Send(m.Chat, p)
		}

	})

	b.Handle(tb.OnQuery, func(q *tb.Query) {
		urls := []string{
			"http://factpile.wikia.com/wiki/File:3-hobbit_saruman.jpg",
		}

		results := make(tb.Results, len(urls)) // []tb.Result
		for i, url := range urls {
			result := &tb.PhotoResult{
				URL:     url,
				Title:   "Test",
				Caption: "Prova",

				// required for photos
				ThumbURL: url,
			}

			results[i] = result
			results[i].SetResultID(strconv.Itoa(i)) // It's needed to set a unique string ID for each result
		}

		err := b.Answer(q, &tb.QueryResponse{
			Results:    results,
			IsPersonal: false,
			NextOffset: "bla",
			CacheTime:  60, // a minute
		})

		if err != nil {
			fmt.Println(err)
		}
	})

	b.Start()
}
