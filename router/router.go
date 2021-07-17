package router

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/marcogregorius/url-shortener/handler"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			log.Println(r.RequestURI, time.Since(start))
		}()
		next.ServeHTTP(w, r)
	})
}

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/ping", handler.PingHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/shortlinks", handler.CreateShortlinkHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/shortlinks/{id}", handler.GetShortlinkHandler).Methods(http.MethodGet)
	r.HandleFunc("/{id}", handler.RedirectHandler).Methods(http.MethodGet)
	r.Use(loggingMiddleware)

	return r
}
