package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/myl7/tgchan2tw/pkg/conf"
	"log"
)

func getDB() (*sql.DB, error) {
	return sql.Open("sqlite3", conf.DBPath)
}

const schema = `
CREATE TABLE IF NOT EXISTS msg (
  id INTEGER(8) PRIMARY KEY,
  guid TEXT UNIQUE NOT NULL
);
`

func init() {
	db, err := getDB()
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec(schema)
	if err != nil {
		log.Fatalln(err)
	}
}
