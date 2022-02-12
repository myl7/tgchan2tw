package test

import (
	"github.com/myl7/tgchan2tw/pkg/tg"
)

type Items struct {
	Items []Item `yaml:"items"`
}

type Item struct {
	Title string   `yaml:"title"`
	Body  string   `yaml:"body"`
	Info  ItemInfo `yaml:"info"`
}

type ItemInfo struct {
	Text       string   `yaml:"text"`
	ImageUrls  []string `yaml:"image_urls"`
	ReplyUrl   string   `yaml:"reply_url"`
	ForwardUrl string   `yaml:"forward_url"`
}

func (lhs ItemInfo) EqItemBody(rhs tg.ItemBody) bool {
	if lhs.Text != rhs.Text {
		return false
	}
	if lhs.ReplyUrl != rhs.ReplyUrl {
		return false
	}
	if lhs.ForwardUrl != rhs.ForwardUrl {
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
