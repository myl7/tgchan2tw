// Copyright 2021-2022 myl7
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"fmt"
	"github.com/myl7/tgchan2tw/pkg/cfg"
	"github.com/myl7/tgchan2tw/pkg/db"
	"github.com/myl7/tgchan2tw/pkg/mdl"
	"github.com/myl7/tgchan2tw/pkg/tg"
	"github.com/myl7/tgchan2tw/pkg/tw"
	"log"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
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
			errMsg := fmt.Sprintln(r) + fmt.Sprint(string(debug.Stack()))
			errs := reportErr(errMsg)
			if len(errs) > 0 {
				for i := range errs {
					log.Println(errs[i])
				}
			}
		}
	}()

	msgs := tg.Fetch()
	for i := len(msgs) - 1; i >= 0; i-- {
		msg := msgs[i]
		handleMsg(msg)
	}
}

func handleMsg(msg *mdl.Msg) {
	tx, err := db.DB.Begin()
	if err != nil {
		panic(err)
	}

	ok := false
	defer func() {
		if !ok {
			err := tx.Rollback()
			if err != nil {
				panic(err)
			}
		}
	}()

	var replyTo int64
	if msg.ReplyTo != "" {
		msgIds := db.GetTgInToTwOut(tx, msg.ReplyTo)
		if len(msgIds) > 0 {
			replyTo = msgIds[len(msgIds)-1]
		}
	}

	images, tmpDir := tmpDl(msg.ImageUrls)

	log.Println("to forward msg:", "msg =", msg, "reply to tw ID =", replyTo, "len of images =", len(images))

	var twOutIDs []int64
	if cfg.Cfg.DryRun {
		twOutIDs = getFakeTwOutIDs(msg.ID)
	} else {
		twOutIDs = tw.Tweet(msg, images, replyTo)
	}

	tgInIDs := msg.InIDs.([]string)
	db.SetTgInAndTwOut(tx, twOutIDs, tgInIDs)

	log.Println(fmt.Sprintf("forwarded tg in msg %s to tw out msgs:", msg.ID), twOutIDs)

	err = os.RemoveAll(tmpDir)
	if err != nil {
		panic(err)
	}

	ok = true
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

func getFakeTwOutIDs(twID string) []int64 {
	r := regexp.MustCompile(`/\d+$`)
	s := r.Find([]byte(twID))
	i, err := strconv.Atoi(string(s)[1:])
	if err != nil {
		panic(err)
	}

	return []int64{int64(i)}
}
