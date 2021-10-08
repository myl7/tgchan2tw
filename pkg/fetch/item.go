package fetch

import (
	"github.com/mmcdole/gofeed"
	"github.com/myl7/tgchan2tw/pkg/db"
	"github.com/myl7/tgchan2tw/pkg/pub"
)

func handleItems(items []*gofeed.Item) error {
	for i := len(items) - 1; i >= 0; i-- {
		item := items[i]

		id, err := db.CheckMsg(item.GUID)
		if err != nil {
			return err
		}

		if id != 0 {
			continue
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

		msg := pub.TweetMsg{
			Body:      body,
			ImageUrls: imageUrls,
			ReplyTo:   replyTo,
		}

		id, err = pub.Tweet(msg)
		if err != nil {
			return err
		}

		err = db.SetMsg(id, item.GUID)
		if err != nil {
			return err
		}
	}

	return nil
}
