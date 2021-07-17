package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/marcogregorius/url-shortener/models"
	"github.com/marcogregorius/url-shortener/router"
)

func main() {
	var port string
	flag.StringVar(&port, "port", ":8080", "Listening port")
	flag.Parse()

	db := models.InitDb()
	defer db.Close()
	r := router.Router()

	srv := &http.Server{
		Addr:    port,
		Handler: r,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	go func() {
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
