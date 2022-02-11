// Copyright 2021-2022 myl7
// SPDX-License-Identifier: Apache-2.0

package db

import "database/sql"

func CheckTgIn(id string) []int64 {
	tx, err := DB.Begin()
	if err != nil {
		panic(err)
	}

	defer func(tx *sql.Tx) {
		_ = tx.Commit()
	}(tx)

	s := "SELECT id FROM tg_in WHERE id = $1"
	q, err := tx.Query(s, id)
	if err != nil {
		panic(err)
	}

	if !q.Next() {
		s = "INSERT INTO tg_in VALUES ($1)"
		_, err = tx.Exec(s, id)
		if err != nil {
			panic(err)
		}

		return nil
	}

	s = "SELECT tw_out_id FROM tg_in_to_tw_out WHERE tg_in_id = $1 ORDER BY tw_out_id"
	q, err = tx.Query(s, id)
	if err != nil {
		panic(err)
	}

	var ids []int64
	var i int64
	for q.Next() {
		err := q.Scan(&i)
		if err != nil {
			panic(err)
		}

		ids = append(ids, i)
	}
	return ids
}

func SetTwOut(twOutIDs []int64, tgInIDs []string) {
	tx, err := DB.Begin()
	if err != nil {
		panic(err)
	}

	defer func(tx *sql.Tx) {
		_ = tx.Commit()
	}(tx)

	for i := range twOutIDs {
		s := "INSERT INTO tg_in VALUES ($1)"
		_, err = tx.Exec(s, twOutIDs[i])
		if err != nil {
			panic(err)
		}

		s = "INSERT INTO tg_in_to_tw_out VALUES ($1, $2)"
		for j := range tgInIDs {
			_, err = tx.Exec(s, tgInIDs[j], twOutIDs[i])
			if err != nil {
				panic(err)
			}
		}
	}
}
