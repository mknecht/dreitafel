package dreitafel

import (
	"fmt"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
)

func CompileToSvgToStdout(input string) {
	lines := make(chan *string, 500) // lines are independent; statements don't span lines yet
	fmcdiagrams := make(chan *FmcBlockDiagram)
	dotdiagrams := make(chan DotGenerator)
	errors := make(chan error, 500) // errors are independent

	var waitGroup sync.WaitGroup

	go KeepParsing(lines, fmcdiagrams, errors, &waitGroup)
	go forwardFmcToDot(fmcdiagrams, dotdiagrams)
	go GenerateDot(dotdiagrams, errors, &waitGroup)
	go printUntilNil(errors, &waitGroup)

	// feed the compiler the input
	lines <- &input

	log.Debug("waiting for parser to finish.")
	waitGroup.Add(1)
	lines <- nil // end token
	waitGroup.Wait()

	log.Debug("waiting for diagram builder to finish.")
	waitGroup.Add(1)
	fmcdiagrams <- nil // end token
	waitGroup.Wait()

	log.Debug("waiting for error handler to finish.")
	waitGroup.Add(1)
	errors <- nil
	waitGroup.Wait()
	log.Debug("Compilation done.")
}

func forwardFmcToDot(in chan *FmcBlockDiagram, out chan DotGenerator) {
	for {
		fmcdiagram := <-in
		if fmcdiagram == nil {
			// cannot put fmcdiagram, since that pointer is typed
			// https://golang.org/doc/faq#nil_error
			out <- nil
			return
		}
		out <- fmcdiagram
	}
}

func printUntilNil(errors <-chan error, waitGroup *sync.WaitGroup) {
	for {
		err := <-errors
		if err == nil {
			break
		}
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	waitGroup.Done()
}
