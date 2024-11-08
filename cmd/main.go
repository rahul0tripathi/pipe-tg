package main

import (
	"log"

	"github.com/rahul0tripathi/pipetg/app/piped"
)

func main() {
	err := piped.Run()
	if err != nil {
		log.Fatal(err)
	}
}
