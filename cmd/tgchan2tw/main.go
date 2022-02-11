// Copyright 2021-2022 myl7
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/myl7/tgchan2tw/pkg"
	"github.com/myl7/tgchan2tw/pkg/cfg"
	"github.com/myl7/tgchan2tw/pkg/db"
	"log"
)

func main() {
	err := cfg.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	err = db.LoadDB()
	if err != nil {
		log.Fatalln(err)
	}

	err = db.InitDB()
	if err != nil {
		log.Fatalln(err)
	}

	pkg.Poll()
}
