package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type App struct {
	DB     *sql.DB
	Router *mux.Router
}

func (a *App) Initialize() {
	dbHost := os.Getenv("DB_HOST")
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal("DB_PORT must be integer")
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connString := fmt.Sprintf("port=%d host=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbPort, dbHost, dbUser, dbPassword, dbName)

	a.DB, err = sql.Open(
		"postgres",
		connString,
	)
	if err != nil {
		log.Fatal(err)
	}

	a.InitializeRouter()
}

func (a *App) Run(port string) {
	srv := &http.Server{
		Addr:    port,
		Handler: a.Router,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	go func() {
		//if err := http.ListenAndServe(port, r); err != nil {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	log.Printf("Listening on port %s", port)

	// make channel that will block until we receive SIGINT or SIGTERM
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c
	log.Println("Shutting down")
}

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

func (a *App) InitializeRouter() {
	r := mux.NewRouter()
	r.HandleFunc("/ping", a.PingHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/shortlinks", a.CreateShortlinkHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/shortlinks/{id}", a.GetShortlinkHandler).Methods(http.MethodGet)
	r.HandleFunc("/{id}", a.RedirectHandler).Methods(http.MethodGet)
	r.Use(LoggingMiddleware)

	a.Router = r
}
