package test

import (
	"github.com/myl7/tgchan2tw/pkg/pub"
	twtext "github.com/myl7/twitter-text-parse-go"
	"testing"
)

func TestSplitAtSpace(t *testing.T) {
	body := `昨天惊觉 mailgun 因为 student 已过，receiving 功能已经没法用了，赶紧搭了个 postfix+dovecot 来收邮件，现在 myl@myl.moe 应该能 work 了
顺带再领略了一遍 postfix 的干涩文档，dovecot 的倒是还够舒服了
另外 docker 下的 complete solution 还挺多的，就是带了一堆垃圾过滤和其他功能后 hardware requirement free memory 1G 起步，着实受不住，外加自己搭的发邮件时挺容易进垃圾箱的，和 mailgun aws 之类的特地做邮件服务的应该没法比，在发还是有免费额度的情况下还是靠他们对抗垃圾箱吧`
	block := `昨天惊觉 mailgun 因为 student 已过，receiving 功能已经没法用了，赶紧搭了个 postfix+dovecot 来收邮件，现在 myl@myl.moe 应该能 work 了
顺带再领略了一遍 postfix 的干涩文档，dovecot 的倒是还够舒服了
另外 docker 下的 complete solution 还挺多的，就是带了一堆垃圾过滤和其他功能后`
	remain := `hardware requirement free memory 1G 起步，着实受不住，外加自己搭的发邮件时挺容易进垃圾箱的，和 mailgun aws 之类的特地做邮件服务的应该没法比，在发还是有免费额度的情况下还是靠他们对抗垃圾箱吧`
	res, err := twtext.Parse(body)
	if err != nil {
		t.Error(err)
	}

	end := res.ValidTextRange.End + 1
	b, r := pub.SplitTweetBodyOnce(body, end)
	if block != b {
		t.Errorf("splitting at space failed: block required:\n%s\nvs got:\n%s\n", block, b)
	}
	if remain != r {
		t.Errorf("splitting at space failed: remain required:\n%s\nvs got:\n%s\n", remain, r)
	}
}
