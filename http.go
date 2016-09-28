package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type httpHandler struct {
	data *data
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		path       []string = strings.Split(r.URL.Path, "/")
		uid1, uid2 int
		isDupes    bool

		result map[string]interface{} = map[string]interface{}{
			"dupes": &isDupes,
		}
		response []byte

		startResponseTime time.Time = time.Now()
		stopResponseTime  float64

		err error
	)

	// add common header
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// check request method and pattern
	if r.Method != "GET" ||
		(len(path) < 3 && len(path) > 4) ||
		(len(path) == 4 && len(path[3]) > 0) {
		h.NotFound(w)
		return
	}

	// parse uids
	if uid1, err = strconv.Atoi(path[1]); err != nil {
		h.NotFound(w)
		return
	}

	if uid2, err = strconv.Atoi(path[2]); err != nil {
		h.NotFound(w)
		return
	}

	// check dupes
	if uid1 != uid2 {
		isDupes = h.data.isDupes(uid1, uid2)
	}

	// write response
	response, _ = json.Marshal(result)
	w.Write(response)

	// write response time log
	stopResponseTime = float64(time.Since(startResponseTime).Seconds() * 1000.0)
	log.Printf("%f\t%s\t%s", stopResponseTime, r.Method, r.RequestURI)
}

func (h *httpHandler) NotFound(w http.ResponseWriter) {
	r, _ := json.Marshal(map[string]interface{}{
		"message": http.StatusText(http.StatusNotFound),
		"status":  http.StatusNotFound,
	})
	w.WriteHeader(http.StatusNotFound)
	w.Write(r)
}
