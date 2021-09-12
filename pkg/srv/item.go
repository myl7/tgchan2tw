package srv

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"github.com/myl7/tgchan2tw/pkg/db"
	"golang.org/x/net/html"
	"strings"
)

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

	id, err = tweet(body, imageUrls, replyTo)
	if err != nil {
		return err
	}

	err = db.SetMsg(id, item.GUID)
	if err != nil {
		return err
	}

	return nil
}

func filterText(body string) (string, []string, string, error) {
	b := bytes.NewBufferString("<body>" + body + "</body>")
	h, err := html.Parse(b)
	if err != nil {
		return "", nil, "", err
	}

	var imageUrls []string
	d := goquery.NewDocumentFromNode(h)
	d.Find("body > :not(blockquote) img").Each(func(_ int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if !ok {
			return
		}

		imageUrls = append(imageUrls, src)
	})

	isForward := false
	forwardLink := ""
	var blocks []string
	d.Find("body > p").Each(func(_ int, s *goquery.Selection) {
		t := s.Text()
		t = strings.Replace(t, "<br>", "\n", -1)
		ts := strings.Split(t, "\n")
		for i := range ts {
			ts[i] = strings.TrimSpace(ts[i])
		}
		t = strings.Join(ts, "\n")
		t = strings.TrimSpace(t)

		if strings.HasPrefix(t, "Forwarded From") {
			href, ok := s.Find("a").Attr("href")
			if ok {
				isForward = true
				forwardLink = href
			}
		}

		blocks = append(blocks, t)
	})

	var res string
	if isForward {
		res = "Forwarded from " + forwardLink
	} else {
		res = strings.Join(blocks, "\n")
	}

	return res, imageUrls, "", nil
}
