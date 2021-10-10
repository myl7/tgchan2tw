package fetch

import (
	"github.com/mmcdole/gofeed"
	"github.com/myl7/tgchan2tw/pkg/db"
	"github.com/myl7/tgchan2tw/pkg/pub"
)

func handleItems(items []*gofeed.Item) error {
	for i := len(items) - 1; i >= 0; i-- {
		item := items[i]

		itemId, err := Guid2Id(item.GUID)
		if err != nil {
			return err
		}

		msgIds, err := db.CheckItem(itemId)
		if err != nil {
			return err
		}

		if len(msgIds) != 0 {
			continue
		}

		body, imageUrls, replyGuid, err := filterText(item.Description)
		if err != nil {
			return err
		}

		replyTo := int64(0)
		if replyGuid != "" {
			replyId, err := Guid2Id(replyGuid)
			if err != nil {
				return err
			}

			msgIds, err = db.CheckItem(replyId)
			if err != nil {
				return err
			}

			if len(msgIds) > 0 {
				replyTo = msgIds[len(msgIds)-1]
			}
		}

		msg := pub.TweetMsg{
			Body:      body,
			ImageUrls: imageUrls,
			ReplyTo:   replyTo,
		}

		createdMsgIds, err := pub.Tweet(msg)
		if err != nil {
			return err
		}

		err = db.SetMsgs(createdMsgIds, []int{itemId})
		if err != nil {
			return err
		}
	}

	return nil
}
