package api

import (
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

func addResponseHeaders(fn http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		fn.ServeHTTP(w, r)
	})
}

func addRequestLog(fn http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			startResponseTime time.Time = time.Now()
			stopResponseTime  float64
		)
		fn.ServeHTTP(w, r)
		stopResponseTime = float64(time.Since(startResponseTime).Seconds() * 1000.0)
		log.Printf("%f\t%s\t%s", stopResponseTime, r.Method, r.RequestURI)
	})
}

func HandleWrap(h http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var (
			ctx context.Context = context.WithValue(r.Context(), "params", p)
			h   http.Handler    = h
		)

		// apply middlewares
		h = addResponseHeaders(h)
		h = addRequestLog(h)

		h.ServeHTTP(w, r.WithContext(ctx))
	}
}
