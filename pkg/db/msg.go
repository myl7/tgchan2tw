// Copyright 2021-2022 myl7
// SPDX-License-Identifier: Apache-2.0

package db

import "database/sql"

func CheckItem(id int) ([]int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	defer func(tx *sql.Tx) {
		_ = tx.Commit()
	}(tx)

	s := "SELECT id FROM items WHERE id = $1"
	q, err := tx.Query(s, id)
	if err != nil {
		return nil, err
	}

	if !q.Next() {
		s = "INSERT INTO items VALUES ($1)"
		_, err = tx.Exec(s, id)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	s = "SELECT msg_id FROM item2msg WHERE item_id = $1 ORDER BY msg_id"
	q, err = tx.Query(s, id)
	if err != nil {
		return nil, err
	}

	var ids []int64
	var i int64
	for q.Next() {
		err := q.Scan(&i)
		if err != nil {
			return nil, err
		}

		ids = append(ids, i)
	}
	return ids, nil
}

func SetMsgs(msgIds []int64, itemIds []int) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func(tx *sql.Tx) {
		_ = tx.Commit()
	}(tx)

	for i := range msgIds {
		s := "INSERT INTO msgs VALUES ($1)"
		_, err = tx.Exec(s, msgIds[i])
		if err != nil {
			return err
		}

		s = "INSERT INTO item2msg VALUES ($1, $2)"
		for j := range itemIds {
			_, err = tx.Exec(s, itemIds[j], msgIds[i])
			if err != nil {
				return err
			}
		}
	}

	return nil
}
