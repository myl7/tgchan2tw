package db

import "database/sql"

func SetMsg(id int64, guid string) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	s := "INSERT INTO msg VALUES ($1, $2)"
	_, err = db.Exec(s, id, guid)
	if err != nil {
		return err
	}

	return nil
}

func CheckMsg(guid string) (int64, error) {
	db, err := getDB()
	if err != nil {
		return 0, err
	}

	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	s := "SELECT id FROM msg WHERE guid = $1"
	q, err := db.Query(s, guid)
	if err != nil {
		return 0, err
	}

	if !q.Next() {
		return 0, nil
	}

	var id int64
	err = q.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
