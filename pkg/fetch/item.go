package fetch

import (
	"github.com/mmcdole/gofeed"
	"github.com/myl7/tgchan2tw/pkg/db"
	"github.com/myl7/tgchan2tw/pkg/pub"
)

func handleItems(items []*gofeed.Item) error {
	var msgs []pub.TweetMsg
	var itemList [][]*gofeed.Item
	var isForwardList []bool
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

		body, err := filterText(item.Description)
		if err != nil {
			return err
		}

		replyTo := int64(0)
		if body.ReplyUrl != "" {
			replyId, err := Guid2Id(body.ReplyUrl)
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
			Body:      body.Text,
			ImageUrls: body.ImageUrls,
			ReplyTo:   replyTo,
		}
		msgs = append(msgs, msg)
		itemList = append(itemList, []*gofeed.Item{item})
		isForwardList = append(isForwardList, body.IsForward)
	}

	// Merge two messages when forwarding with comment in Telegram
	var mergeList []int
	for i := 0; i < len(msgs); i++ {
		// If two messages are sent simultaneously and the latter one is forwarded message
		if i < len(msgs)-1 && itemList[i][0].Published == itemList[i+1][0].Published && isForwardList[i+1] {
			mergeList = append(mergeList, i)
			i++
		}
	}
	for j := len(mergeList) - 1; j >= 0; j-- {
		i := mergeList[j]
		msgs[i].Body = msgs[i].Body + "\n" + msgs[i+1].Body
		msgs = append(msgs[:i+1], msgs[i+2:]...)
		itemList[i] = append(itemList[i], itemList[i+1]...)
		itemList = append(itemList[:i+1], itemList[i+2:]...)
	}

	for i := range msgs {
		createdMsgIds, err := pub.Tweet(msgs[i])
		if err != nil {
			return err
		}

		var itemIds []int
		for j := range itemList[i] {
			itemId, err := Guid2Id(itemList[i][j].GUID)
			if err != nil {
				return err
			}

			itemIds = append(itemIds, itemId)
		}

		err = db.SetMsgs(createdMsgIds, itemIds)
		if err != nil {
			return err
		}
	}

	return nil
}
