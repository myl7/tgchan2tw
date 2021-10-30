package pub

import (
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

		s := utf16.Encode([]rune(remain))
		b := string(utf16.Decode(s[:res.ValidTextRange.End+1]))
		bodies = append(bodies, b)
		remain = string(utf16.Decode(s[res.ValidTextRange.End+1:]))
	}
	return bodies, nil
}
