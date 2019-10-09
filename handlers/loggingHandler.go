package handlers

import (
	"io"
	"net/http"
	"time"
)

// LoggingHandler is a handler/middleware to log requests. Inspired by the func from gorilla/handlers
func LoggingHandler(out io.Writer, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		out.Write([]byte(r.URL.String() + ", " + time.Now().String() + "\n"))
		h.ServeHTTP(w, r) // this runs handler h
	})
}
