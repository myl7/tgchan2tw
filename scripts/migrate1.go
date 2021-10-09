package main

import (
	"github.com/myl7/tgchan2tw/pkg/db"
	"github.com/myl7/tgchan2tw/pkg/fetch"
	"log"
)

func main() {
	d, err := db.GetDB()
	if err != nil {
		log.Fatalln(err)
	}

	s := "SELECT (id, guid) FROM msg"
	q, err := d.Query(s)
	if err != nil {
		log.Fatalln(err)
	}

	var msgId int64
	var itemGuid string

	var msgIds []int64
	var itemIds []int
	for q.Next() {
		err := q.Scan(&msgId, &itemGuid)
		if err != nil {
			log.Fatalln(err)
		}

		msgIds = append(msgIds, msgId)
		itemId, err := fetch.Guid2Id(itemGuid)
		if err != nil {
			log.Fatalln(err)
		}

		itemIds = append(itemIds, itemId)
	}

	_, err = d.Exec("DROP TABLE IF EXISTS msg")
	if err != nil {
		log.Fatalln(err)
	}

	for i := range msgIds {
		s := "INSERT OR IGNORE INTO msgs VALUES ($1)"
		_, err := d.Exec(s, msgIds[i])
		if err != nil {
			log.Fatalln(err)
		}
	}

	for i := range itemIds {
		s := "INSERT OR IGNORE INTO items VALUES ($1)"
		_, err := d.Exec(s, itemIds[i])
		if err != nil {
			log.Fatalln(err)
		}
	}

	for i := range msgIds {
		s := "INSERT OR IGNORE INTO item2msg VALUES ($1, $2)"
		_, err := d.Exec(s, itemIds[i], msgIds[i])
		if err != nil {
			log.Fatalln(err)
		}
	}
}
