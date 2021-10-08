package fetch

import (
	"bufio"
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

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

func downloadImages(imageUrls []string) ([]io.ReadCloser, string, error) {
	dir, err := ioutil.TempDir("/tmp", "tgchan2tw")
	if err != nil {
		return nil, "", err
	}

	var images []io.ReadCloser
	for i := range imageUrls {
		url := imageUrls[i]
		res, err := http.Get(url)
		if err != nil {
			return nil, "", err
		}

		f, err := ioutil.TempFile(dir, "image")
		if err != nil {
			return nil, "", err
		}

		_, err = bufio.NewReader(res.Body).WriteTo(f)
		if err != nil {
			return nil, "", err
		}

		_, err = f.Seek(0, 0)
		if err != nil {
			return nil, "", err
		}

		images = append(images, f)
	}

	return images, dir, nil
}
