package main

import (
	"log"

	"github.com/zyablitsev/testapp.backend"
)

func main() {
	var err error

	// populate log with dummy data
	if err = testapp.GenerateLogFile(); err != nil {
		log.Fatal(err)
	}
}
