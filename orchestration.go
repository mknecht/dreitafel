package dreitafel

import (
	"fmt"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
)

func CompileToSvgToStdout(lines chan *string) {
	fmcdiagrams := make(chan *FmcBlockDiagram)
	dotdiagrams := make(chan DotGenerator)
	errors := make(chan error, 500) // errors are independent

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)
	go KeepParsing(lines, fmcdiagrams, errors, &waitGroup)
	go forwardFmcToDot(fmcdiagrams, dotdiagrams)

	waitGroup.Add(1)
	go GenerateDot(dotdiagrams, errors, &waitGroup)

	go printUntilNil(errors, &waitGroup)

	log.Debug("waiting for parser and diagram builder.")
	waitGroup.Wait()
	log.Debug("Compilation done.")

	log.Debug("error handler to finish.")
	waitGroup.Add(1)
	close(errors)
	waitGroup.Wait()
}

func forwardFmcToDot(in chan *FmcBlockDiagram, out chan DotGenerator) {
	for fmcdiagram := range in {
		if fmcdiagram == nil {
			// cannot put fmcdiagram, since that pointer is typed
			// https://golang.org/doc/faq#nil_error
			close(out)
			return
		}
		out <- fmcdiagram
	}
	close(out)
}

func printUntilNil(errors <-chan error, waitGroup *sync.WaitGroup) {
	for err := range errors {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	waitGroup.Done()
}
