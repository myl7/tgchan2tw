package fetch

import (
	"github.com/mmcdole/gofeed"
	"github.com/myl7/tgchan2tw/pkg/conf"
	"log"
	"net/url"
	"path"
	"strconv"
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

	err = handleItems(items)
	if err != nil {
		return err
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
	q.Set("filterout", conf.PostFilterOut)
	u.RawQuery = q.Encode()

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(u.String())
	if err != nil {
		return nil, err
	}

	return feed.Items, nil
}
