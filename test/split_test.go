package test

import (
	"github.com/myl7/tgchan2tw/pkg/tw"
	twtext "github.com/myl7/twitter-text-parse-go"
	"gopkg.in/yaml.v2"
	"os"
	"testing"
)

type split struct {
	Texts []splitText `yaml:"texts"`
}

type splitText struct {
	Title  string `yaml:"title"`
	Body   string `yaml:"body"`
	Block  string `yaml:"block"`
	Remain string `yaml:"remain"`
}

func TestSplit(t *testing.T) {
	f, err := os.ReadFile("data/split.yaml")
	if err != nil {
		t.Error(err)
	}

	var s split
	err = yaml.UnmarshalStrict(f, &s)
	if err != nil {
		t.Error(err)
	}

	for i := range s.Texts {
		text := s.Texts[i]
		body := text.Body
		block := text.Block
		remain := text.Remain
		res, err := twtext.Parse(body)
		if err != nil {
			t.Error(err)
		}

		end := res.ValidTextRange.End + 1
		b, r := tw.SplitTweetBodyOnce(body, end)
		if block != b {
			t.Errorf("splitting %s failed: block required:\n%s\nvs got:\n%s\n", text.Title, block, b)
		}
		if remain != r {
			t.Errorf("splitting %s failed: remain required:\n%s\nvs got:\n%s\n", text.Title, remain, r)
		}
	}
}
