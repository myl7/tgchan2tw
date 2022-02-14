// Copyright 2021-2022 myl7
// SPDX-License-Identifier: Apache-2.0

package db

import "database/sql"

func GetTgIn(id string) bool {
	s := "SELECT id FROM tg_in WHERE id = $1"
	q, err := DB.Query(s, id)
	if err != nil {
		panic(err)
	}

	if q.Next() {
		return true
	}
	return false
}

func GetTgInToTwOut(tx *sql.Tx, id string) []int64 {
	s := "SELECT tw_out_id FROM tg_in_to_tw_out WHERE tg_in_id = $1 ORDER BY tw_out_id"
	q, err := tx.Query(s, id)
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

func SetTgInAndTwOut(tx *sql.Tx, twOutIDs []int64, tgInIDs []string) {
	for i := range tgInIDs {
		s := "INSERT INTO tg_in (id) VALUES ($1)"
		_, err := tx.Exec(s, tgInIDs[i])
		if err != nil {
			panic(err)
		}
	}

	for i := range twOutIDs {
		s := "INSERT INTO tw_out (id) VALUES ($1)"
		_, err := tx.Exec(s, twOutIDs[i])
		if err != nil {
			panic(err)
		}

		s = "INSERT INTO tg_in_to_tw_out (tg_in_id, tw_out_id) VALUES ($1, $2)"
		for j := range tgInIDs {
			_, err = tx.Exec(s, tgInIDs[j], twOutIDs[i])
			if err != nil {
				panic(err)
			}
		}
	}
}
