package configs

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func LoadDB() {
	cfg := AppConfig.DB
	if AppConfig == nil {
		log.Fatal("App Config Masih nill ! Load Config Blm Di Panggil")
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.User,
		cfg.Pass,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Tambahkan pengecekan koneksi agar yakin sudah terhubung ke remote IP
	err = DB.Ping()
	if err != nil {
		log.Fatal("Koneksi DB Gagal: ", err)
	}
}
