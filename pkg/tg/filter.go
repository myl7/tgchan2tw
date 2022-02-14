// Copyright 2021-2022 myl7
// SPDX-License-Identifier: Apache-2.0

package tg

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"github.com/myl7/tgchan2tw/pkg/cfg"
	"github.com/myl7/tgchan2tw/pkg/db"
	"github.com/myl7/tgchan2tw/pkg/mdl"
	"golang.org/x/net/html"
	"regexp"
	"strings"
)

func filterItems(items []*gofeed.Item) []*mdl.Msg {
	var msgs []*mdl.Msg
	for i := range items {
		item := items[i]

		exist := db.GetTgIn(item.GUID)
		if exist {
			continue
		}

		itemBody := FilterText(item.Description, item.Link)
		msg := mdl.Msg{
			ID:        item.GUID,
			Body:      itemBody.Text,
			ImageUrls: itemBody.ImageUrls,
			ReplyTo:   itemBody.ReplyUrl,
			FwdFrom:   itemBody.ForwardUrl,
			InIDs:     []string{item.GUID},
		}
		msgs = append(msgs, &msg)
	}

	msgs = mergeFwdMsgs(msgs, items)

	return msgs
}

// Merge two messages when forwarding with comment in Telegram
func mergeFwdMsgs(msgs []*mdl.Msg, items []*gofeed.Item) []*mdl.Msg {
	var mergeList []int
	for i := 0; i < len(msgs); i++ {
		// If two messages are sent simultaneously and the former one is forwarded message
		if i < len(msgs)-1 && items[i].Published == items[i+1].Published && msgs[i].FwdFrom != "" {
			mergeList = append(mergeList, i)
			i++
		}
	}
	for j := len(mergeList) - 1; j >= 0; j-- {
		i := mergeList[j]
		msgs[i].Body = msgs[i].Body + "\n" + msgs[i+1].Body
		msgs[i].InIDs = append(msgs[i].InIDs.([]string), msgs[i+1].InIDs.([]string)...)
		msgs = append(msgs[:i+1], msgs[i+2:]...)
	}

	return msgs
}

type ItemBody struct {
	Text       string
	ImageUrls  []string
	ReplyUrl   string
	ForwardUrl string
}

func FilterText(body string, selfUrl string) ItemBody {
	b := bytes.NewBufferString("<body>" + body + "</body>")
	h, err := html.Parse(b)
	if err != nil {
		panic(err)
	}

	d := goquery.NewDocumentFromNode(h)
	imageUrls := filterImageUrls(d)
	quoteUrl, isForward := filterHeadQuote(d)

	var blocks []string
	d.Find("body").Children().Each(func(_ int, s *goquery.Selection) {
		t := ""
		if goquery.NodeName(s) == "p" {
			t = filterBodyP(s)
		} else if goquery.NodeName(s) == "pre" {
			t = filterBodyPre(s)
		}
		if t != "" {
			blocks = append(blocks, t)
		}
	})

	replyUrl := ""
	forwardUrl := ""
	res := ""
	if isForward {
		if quoteUrl == "" {
			// If the original message can not be referred by an url, use your Telegram forward message url
			quoteUrl = selfUrl
		}
		res = quoteUrl
		forwardUrl = quoteUrl
	} else {
		res = strings.Join(blocks, "\n")

		// quote url may be of other channel, filter out the situation
		reg := regexp.MustCompile(`t[.]me/([^/]+)/`)
		m := reg.FindStringSubmatch(quoteUrl)
		if len(m) > 1 && m[1] == cfg.Cfg.TgChanName {
			replyUrl = quoteUrl
		}
	}

	return ItemBody{
		Text:       res,
		ImageUrls:  imageUrls,
		ReplyUrl:   replyUrl,
		ForwardUrl: forwardUrl,
	}
}

func filterImageUrls(d *goquery.Document) []string {
	var imageUrls []string
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

		s.Remove()
		imageUrls = append(imageUrls, src)
	})
	return imageUrls
}

func filterHeadQuote(d *goquery.Document) (string, bool) {
	r := d.Find("body").Children().First()
	quoteUrl, _ := r.Find("a").Attr("href")

	if goquery.NodeName(r) == "blockquote" {
		// Ensure it is reply message
		r.Remove()
		return quoteUrl, false
	}

	isForward := false
	if strings.HasPrefix(r.Text(), "Forwarded From") {
		// Ensure it is forward message
		// Notice: Forward message may also have empty quote url when it is forwarded from a private channel or chat.
		isForward = true
	}

	if !isForward {
		// It is a plain message, so reset quoteUrl
		quoteUrl = ""
	}

	return quoteUrl, isForward
}

func filterBodyP(s *goquery.Selection) string {
	t := ""
	s.Contents().Each(func(i int, p *goquery.Selection) {
		if goquery.NodeName(p) == "br" {
			t += "\n"
		} else {
			t += p.Text()
		}
	})
	ts := strings.Split(t, "\n")
	for i := range ts {
		ts[i] = strings.TrimSpace(ts[i])
	}
	t = strings.Join(ts, "\n")
	t = strings.TrimSpace(t)
	return t
}

func filterBodyPre(s *goquery.Selection) string {
	t := ""
	s.Contents().Each(func(_ int, s *goquery.Selection) {
		if goquery.NodeName(s) == "br" {
			t += "\n"
		} else {
			t += s.Text()
		}
	})
	t = strings.TrimRight(t, "\n")
	return t
}
