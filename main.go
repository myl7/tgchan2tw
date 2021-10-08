package tgchan2tw

import (
	"github.com/myl7/tgchan2tw/pkg/fetch"
)

func Run() error {
	return fetch.Poll()
}
