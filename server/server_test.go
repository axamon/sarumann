package server

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type contatto struct {
	num    string
	Valido bool
}

var Numeri []contatto

var okcell bool

func TestVerificaCell(t *testing.T) {
	Convey("Given a cellphone", t, func() {

		Convey("The cell number must be at least 10 digits", func() {
			_, err := verificaCell("12345678")
			So(err, ShouldNotBeNil)
			_, err = verificaCell("1234567890")
			So(err, ShouldBeNil)
			cell, err := verificaCell("+391234567890")
			So(cell, ShouldEqual, "1234567890")
			So(err, ShouldBeNil)
			_, err = verificaCell("+391234567890a")
			So(err, ShouldNotBeNil)

		})
	})

}
