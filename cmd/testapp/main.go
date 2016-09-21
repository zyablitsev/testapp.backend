package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/zyablitsev/testapp.backend/api"
	"github.com/zyablitsev/testapp.backend/model"
	"github.com/zyablitsev/testapp.backend/settings"
)

func main() {
	var (
		config = settings.GetInstance()
		data   = model.GetDataInstance()
	)

	go func() {
		data.ReadLog()

		c := time.Tick(1 * time.Second)
		for _ = range c {
			data.ReadLog()
		}
	}()

	router := httprouter.New()
	router.NotFound = api.NotFound()
	router.Handle("GET", "/:user_id1/:user_id2", api.HandleWrap(api.GetIsDupes))

	// start http server
	log.Println("Now listening…")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), router))
}
