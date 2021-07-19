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
}
