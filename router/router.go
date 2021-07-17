package router

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gorilla/mux"
	"github.com/marcogregorius/url-shortener/handler"
	log "github.com/sirupsen/logrus"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	return
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Error(err, string(debug.Stack()))
			}
		}()
		start := time.Now()
		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r)
		log.Println(wrapped.status, r.Method, r.RequestURI, time.Since(start))
	})
}

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/ping", handler.PingHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/shortlinks", handler.CreateShortlinkHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/shortlinks/{id}", handler.GetShortlinkHandler).Methods(http.MethodGet)
	r.HandleFunc("/{id}", handler.RedirectHandler).Methods(http.MethodGet)
	r.Use(LoggingMiddleware)

	return r
}
