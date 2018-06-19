package server

import "testing"

type contatto struct {
	num    string
	Valido bool
}

var Numeri []contatto

var okcell bool

func TestVerificaCell(t *testing.T) {

	Numeri := []contatto{
		{num: "334233123", Valido: false},
		{num: "3342331230", Valido: true},
		{num: "+393342331230", Valido: true},
		{num: "+39 335 2331230", Valido: false},
	}

	for _, cell := range Numeri {

		cellulare, err := verificaCell(cell.num)
		switch {
		case len(cellulare) == 10:
			okcell = true
		case err != nil:
			t.Skip()
		default:
			//Se restituisce un numero diverso da 10 è sbagliato
			okcell = false
		}

		if okcell != cell.Valido {
			t.Fatal()

		}

	}
}
