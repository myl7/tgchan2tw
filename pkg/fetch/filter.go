package fetch

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"strings"
)

type ItemBody struct {
	Text       string
	ImageUrls  []string
	ReplyUrl   string
	ForwardUrl string
}

func FilterText(body string, selfUrl string) (ItemBody, error) {
	b := bytes.NewBufferString("<body>" + body + "</body>")
	h, err := html.Parse(b)
	if err != nil {
		return ItemBody{}, err
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
		replyUrl = quoteUrl
	}

	return ItemBody{
		Text:       res,
		ImageUrls:  imageUrls,
		ReplyUrl:   replyUrl,
		ForwardUrl: forwardUrl,
	}, nil
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
