// Copyright 2021-2022 myl7
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"github.com/myl7/tgchan2tw/pkg/cfg"
	"github.com/myl7/tgchan2tw/pkg/db"
	"github.com/myl7/tgchan2tw/pkg/tg"
	"github.com/myl7/tgchan2tw/pkg/tw"
	"log"
	"os"
	"time"
)

func Poll() {
	for true {
		pollRound()
		time.Sleep(time.Duration(cfg.Cfg.PollInterval) * time.Second)
	}
}

func pollRound() {
	defer func() {
		r := recover()
		if r != nil {
			log.Println(r)
		}
	}()

	msgs := tg.Fetch()
	for i := len(msgs) - 1; i >= 0; i-- {
		msg := msgs[i]

		msgIds := db.CheckTgIn(msg.ReplyTo)
		if len(msgIds) <= 0 {
			continue
		}
		replyTo := msgIds[len(msgIds)-1]

		images, tmpDir := tmpDl(msg.ImageUrls)
		twOutIDs := tw.Tweet(msg, images, replyTo)

		tgInIDs := msg.InIDs.([]string)
		db.SetTwOut(twOutIDs, tgInIDs)

		err := os.RemoveAll(tmpDir)
		if err != nil {
			panic(err)
		}
	}
}
