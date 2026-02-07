package main

type Mutasi struct {
	Tanggal string
	Uraian  string
	Teller  string
	Debit   string
	Kredit  string
	Saldo   string
}

type Summary struct {
	SaldoAwal   float64
	TotalDebit  float64
	TotalKredit float64
	SaldoAkhir  float64
}
