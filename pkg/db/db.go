package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/myl7/tgchan2tw/pkg/conf"
	"log"
)

func GetDB() (*sql.DB, error) {
	return sql.Open("sqlite3", conf.DBPath)
}

const schema = `
CREATE TABLE IF NOT EXISTS items (
  id INTEGER(4) PRIMARY KEY
);
CREATE TABLE IF NOT EXISTS msgs (
  id INTEGER(8) PRIMARY KEY
);
CREATE TABLE IF NOT EXISTS item2msg (
  item_id INTEGER(4) REFERENCES items (id),
  msg_id INTEGER(8) REFERENCES msgs (id)
);
`

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
