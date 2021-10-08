package main

import (
	"github.com/myl7/tgchan2tw"
	"log"
)

func main() {
	err := tgchan2tw.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
