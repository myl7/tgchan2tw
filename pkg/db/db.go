package db

import (
	"database/sql"
	_ "embed"
	_ "github.com/mattn/go-sqlite3"
	"github.com/myl7/tgchan2tw/pkg/cfg"
	"log"
)

func GetDB() (*sql.DB, error) {
	return sql.Open("sqlite3", cfg.DBPath)
}

//go:embed schema.sql
var schema string

func init() {
	db, err := GetDB()
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec(schema)
	if err != nil {
		log.Fatalln(err)
	}
}
