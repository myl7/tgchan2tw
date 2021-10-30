package pub

import (
	"github.com/myl7/tgchan2tw/pkg/conf"
	twtext "github.com/myl7/twitter-text-parse-go"
	"unicode/utf16"
)

func splitTweetBody(body string) ([]string, error) {
	var bodies []string
	remain := body
	for {
		if remain == "" {
			break
		}

		res, err := twtext.Parse(remain)
		if err != nil {
			return nil, err
		}

		if res.IsValid {
			bodies = append(bodies, remain)
			break
		}

		b, r := splitTweetBodyOnce(remain, res.ValidTextRange.End+1)
		bodies = append(bodies, b)
		remain = r
	}
	return bodies, nil
}

// body is ensured to be not empty
func splitTweetBodyOnce(body string, end int) (string, string) {
	s := utf16.Encode([]rune(body))
	b := ""
	r := ""
	start := 0

	if conf.TwTextSplitBackDisableRate == "" {
		start = end * conf.TwTextSplitBackRate / 100
		sep, ok := findInSplitRange(s, end, start)
		if ok {
			b = string(utf16.Decode(s[:sep]))
			r = string(utf16.Decode(s[sep:]))
			return b, r
		}
	}

	if conf.TwTextSplitBackDisableLen == "" {
		start = end - conf.TwTextSplitBackLen
		if start < 0 {
			start = 0
		}
		sep, ok := findInSplitRange(s, end, start)
		if ok {
			b = string(utf16.Decode(s[:sep]))
			r = string(utf16.Decode(s[sep:]))
			return b, r
		}
	}

	sep := end
	b = string(utf16.Decode(s[:sep]))
	r = string(utf16.Decode(s[sep:]))
	return b, r
}

func findInSplitRange(s []uint16, end int, start int) (int, bool) {
	puncs := utf16.Encode([]rune(" ，。？！；：、"))
	for i := end - 1; i >= start; i-- {
		for j := range puncs {
			if s[i] == puncs[j] {
				return i + 1, true
			}
		}
	}
	return 0, false
}
