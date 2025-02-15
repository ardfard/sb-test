package main

import (
	"log"

	"github.com/ardfard/sb-test/cmd/server/command"
)

func main() {
	if err := command.Execute(); err != nil {
		log.Fatal(err)
	}
}
