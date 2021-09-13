package srv

import (
	"github.com/mmcdole/gofeed"
	"github.com/myl7/tgchan2tw/pkg/conf"
	"log"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

func Start() error {
	for true {
		err := poll()
		if err != nil {
			log.Println(err)
		}

		time.Sleep(time.Duration(conf.PollInterval) * time.Second)
	}

	return nil
}

func poll() error {
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
