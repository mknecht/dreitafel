package main

import (
	log "github.com/Sirupsen/logrus"

	"dreitafel/web"
)

func main() {
	log.SetLevel(log.DebugLevel)

	web.ListenAndServe()
}
