package main

import (
	"dreitafel"
	"os"

	log "github.com/Sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	if len(os.Args) < 2 {
		log.Errorf("Usage: pattern to parse\n")
		os.Exit(2)
	}

	log.Debugf("Parsing: '%v'", os.Args[1])
	dreitafel.CompileToSvgToStdout(os.Args[1])
}
