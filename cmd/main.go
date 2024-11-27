package main

import (
	"log"
	"os"

	"github.com/rahul0tripathi/pipetg/app/piped"
)

func main() {
	// Check if running in CLI mode
	if len(os.Args) > 1 {
		err := piped.RunCLI(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// Default server mode
	err := piped.Run()
	if err != nil {
		log.Fatal(err)
	}
}
