package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	var (
		genlogflag bool

		config = getConfig()
		data   = getModel()

		handler = &httpHandler{data: data}

		err error
	)

	flag.BoolVar(&genlogflag, "gendata", false, "generate data logfile")
	flag.Parse()

	if genlogflag {
		// populate log with dummy data
		if err = generateLogFile(); err != nil {
			log.Fatal(err)
		}
		return
	}

	go func() {
		data.readLog()

		c := time.Tick(1 * time.Second)
		for _ = range c {
			data.readLog()
		}
	}()

	// start http server
	log.Println("Now listening…")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), handler))
}
