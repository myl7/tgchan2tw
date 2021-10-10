package fetch

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"strings"
)

type ItemBody struct {
	Text      string
	ImageUrls []string
	ReplyUrl  string
	IsForward bool
}

func filterText(body string) (ItemBody, error) {
	b := bytes.NewBufferString("<body>" + body + "</body>")
	h, err := html.Parse(b)
	if err != nil {
		return ItemBody{}, err
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
		t := ""
		s.Contents().Each(func(i int, p *goquery.Selection) {
			if goquery.NodeName(p) == "br" {
				t += "\n"
			} else if t == "" {
				t += p.Text()
			} else {
				t += " " + p.Text()
			}
		})
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

	replyUrl := ""
	r := d.Find("body :first-child")
	if goquery.NodeName(r) == "blockquote" {
		href, ok := r.Find("a").Attr("href")
		if ok {
			replyUrl = href
		}
	}

	var res string
	if isForward {
		res = forwardLink
	} else {
		res = strings.Join(blocks, "\n")
	}

	return ItemBody{
		Text:      res,
		ImageUrls: imageUrls,
		ReplyUrl:  replyUrl,
		IsForward: isForward,
	}, nil
}
