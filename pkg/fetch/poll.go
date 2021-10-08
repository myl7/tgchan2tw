package fetch

import (
	"github.com/mmcdole/gofeed"
	"github.com/myl7/tgchan2tw/pkg/conf"
	"github.com/myl7/tgchan2tw/pkg/db"
	"github.com/myl7/tgchan2tw/pkg/pub"
	"log"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

func Poll() error {
	for true {
		err := pollRound()
		if err != nil {
			log.Println(err)
		}

		time.Sleep(time.Duration(conf.PollInterval) * time.Second)
	}

	return nil
}

func pollRound() error {
	items, err := reqRsshub()
	if err != nil {
		return err
	}

	for i := len(items) - 1; i >= 0; i-- {
		item := items[i]
		err := handleItem(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func reqRsshub() ([]*gofeed.Item, error) {
	u, err := url.Parse(conf.RsshubHost)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join("/telegram/channel", conf.TgChanName)

	q := u.Query()
	q.Set("filter_time", strconv.Itoa(conf.PollInterval))
	q.Set("filterout", strings.Replace(conf.PostFilterOut, "#", "%23", -1))
	u.RawQuery = q.Encode()

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(u.String())
	if err != nil {
		return nil, err
	}

	return feed.Items, nil
}

func handleItem(item *gofeed.Item) error {
	id, err := db.CheckMsg(item.GUID)
	if err != nil {
		return err
	}

	if id != 0 {
		return nil
	}

	body, imageUrls, replyGuid, err := filterText(item.Description)
	if err != nil {
		return err
	}

	replyTo := int64(0)
	if replyGuid != "" {
		var err error
		replyTo, err = db.CheckMsg(replyGuid)
		if err != nil {
			return err
		}
	}

	id, err = pub.Tweet(body, imageUrls, replyTo)
	if err != nil {
		return err
	}

	err = db.SetMsg(id, item.GUID)
	if err != nil {
		return err
	}

	return nil
}
