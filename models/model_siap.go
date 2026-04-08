package models

type Personal struct {
	PID      string `json:"pid"`
	Nama     string `json:"nama"`
	Password string `json:"password"`
	Bagian   string `json:"bagian"`
	Pass     string `json:"pswd"`
}

type Barang struct {
	ID         string  `json:"id"`
	Namabarang string  `json:"namabarang"`
	Nobarcode  string  `json:"nobarcode"`
	Urlgbr     *string `json:"url_gbr"`
	Image      string  `json:"image"`
}
