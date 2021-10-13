package test

import "github.com/myl7/tgchan2tw/pkg/fetch"

type Items struct {
	Items []Item `yaml:"items"`
}

type Item struct {
	Title string   `yaml:"title"`
	Body  string   `yaml:"body"`
	Info  ItemInfo `yaml:"info"`
}

type ItemInfo struct {
	Text      string   `yaml:"text"`
	ImageUrls []string `yaml:"image_urls"`
	ReplyUrl  string   `yaml:"reply_url"`
	IsForward bool     `yaml:"is_forward"`
}

func (lhs ItemInfo) EqItemBody(rhs fetch.ItemBody) bool {
	if lhs.Text != rhs.Text {
		return false
	}
	if lhs.ReplyUrl != rhs.ReplyUrl {
		return false
	}
	if lhs.IsForward != rhs.IsForward {
		return false
	}
	if len(lhs.ImageUrls) != len(rhs.ImageUrls) {
		return false
	}
	for i := range lhs.ImageUrls {
		if lhs.ImageUrls[i] != rhs.ImageUrls[i] {
			return false
		}
	}
	return true
}
