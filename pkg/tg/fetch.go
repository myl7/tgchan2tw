// Copyright 2021-2022 myl7
// SPDX-License-Identifier: Apache-2.0

package tg

import (
	"github.com/mmcdole/gofeed"
	"github.com/myl7/tgchan2tw/pkg/cfg"
	"github.com/myl7/tgchan2tw/pkg/mdl"
	"net/url"
	"path"
	"strconv"
)

func Fetch() []*mdl.Msg {
	items := reqRsshub()
	msgs := filterItems(items)
	return msgs
}

func reqRsshub() []*gofeed.Item {
	u, err := url.Parse(cfg.Cfg.RsshubUrl)
	if err != nil {
		panic(err)
	}

	u.Path = path.Join("/telegram/channel", cfg.Cfg.TgChanName)

	q := u.Query()
	q.Set("filter_time", strconv.Itoa(cfg.Cfg.PollRange))
	q.Set("filterout", cfg.Cfg.PostFilterOut)
	u.RawQuery = q.Encode()

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(u.String())
	if err != nil {
		panic(err)
	}

	return feed.Items
}
