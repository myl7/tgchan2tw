package fetch

import (
	"regexp"
	"strconv"
)

func Guid2Id(guid string) (int, error) {
	r := regexp.MustCompile(`/\d+$`)
	s := r.Find([]byte(guid))
	i, err := strconv.Atoi(string(s)[1:])
	if err != nil {
		return 0, err
	}

	return i, nil
}
