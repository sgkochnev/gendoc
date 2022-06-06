package handler

import (
	"log"
	"net/http"
	"os"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)

		loginfo := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
		loginfo.Printf("%s %s %s", req.Method, req.RequestURI, time.Since(start))
	})
}
