// Copyright 2021-2022 myl7
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"database/sql"
	_ "embed"
	_ "github.com/mattn/go-sqlite3"
	"github.com/myl7/tgchan2tw/pkg/cfg"
)

func GetDB() (*sql.DB, error) {
	return sql.Open("sqlite3", cfg.Cfg.DBPath)
}

var DB *sql.DB

func LoadDB() error {
	if DB != nil {
		return nil
	}

	var err error
	DB, err = GetDB()
	if err != nil {
		return err
	}
	return nil
}

//go:embed schema.sql
var schema string

func InitDB() error {
	_, err := DB.Exec(schema)
	if err != nil {
		return err
	}
	return nil
}
