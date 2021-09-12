package main

import (
	"github.com/myl7/tgchan2tw/pkg/srv"
	"log"
)

func main() {
	err := srv.Start()
	if err != nil {
		log.Fatalln(err)
	}
}
