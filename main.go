package main

import (
	"flag"

	"github.com/marcogregorius/url-shortener/app"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	var port string
	flag.StringVar(&port, "port", ":8080", "Listening port")
	flag.Parse()

	a := app.App{}

	a.Initialize()
	a.Run(port)

	/*
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
	*/
}
