package client

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFields(t *testing.T) {

	Convey("Given a notification", t, func() {

		Convey("The cell must be at least 10 digits", func() {
			So(verificaCell("12345678"), ShouldNotBeNil)
			So(verificaCell("1234567890"), ShouldBeNil)
			So(verificaCell("+391234567890"), ShouldBeNil)
			So(verificaCell("+391234567890a"), ShouldNotBeNil)

		})

		Convey("And the message must be less than 50 characters", func() {
			So(verificaMsg("Problema con il web server, chiamare pincopallo"), ShouldBeNil)
			So(verificaMsg("Problema con il web server, chiamare pincopallo,Problema con il web server, chiamare pincopalloProblema con il web server, chiamare pincopallo"), ShouldNotBeNil)
		})

	})

}

/*

//TestSendPost verifica il numero di parametri passati
func TestSendPost(t *testing.T) {

	Parametri := []parametri{
		{endpoint: "http://127.0.0.1/create", ok: true, notifica: Notifica{Hostname: "host1", Service: "www", Piattaforma: "CDN"}},
		{endpoint: "http://127.0.0.1/create", ok: false, notifica: Notifica{Hostname: "host1", Service: "www", Piattaforma: ""}},
	}

	for _, elements := range Parametri {
		err := SendPost(elements.endpoint, elements.notifica.Hostname,
			elements.notifica.Service, elements.notifica.Piattaforma,
			elements.notifica.Reperibile,
			elements.notifica.Cellulare, elements.notifica.Messaggio)

			if err != nil && {
				ok.
			}

	}
}
*/
