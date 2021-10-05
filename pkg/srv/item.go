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
	d.Find("img").Each(func(_ int, s *goquery.Selection) {
		ok := true
		s.Parents().Each(func(_ int, s *goquery.Selection) {
			if goquery.NodeName(s) == "blockquote" {
				ok = false
			}
		})
		if !ok {
			return
		}

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
		t, err := s.Html()
		if err != nil {
			return
		}

		t = strings.Replace(t, "<br/>", "\n", -1)
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

	replyGuid := ""
	r := d.Find("body :first-child")
	if goquery.NodeName(r) == "blockquote" {
		href, ok := r.Find("a").Attr("href")
		if ok {
			replyGuid = href
		}
	}

	var res string
	if isForward {
		res = forwardLink
	} else {
		res = strings.Join(blocks, "\n")
	}

	return res, imageUrls, replyGuid, nil
}
