package dreitafel

import (
	"fmt"
	"sync"
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

	fmt.Println("waiting for parser to finish.")
	waitGroup.Add(1)
	lines <- nil // end token
	waitGroup.Wait()

	fmt.Println("waiting for diagram builder to finish.")
	waitGroup.Add(1)
	fmcdiagrams <- nil // end token
	waitGroup.Wait()

	fmt.Println("waiting for error handler to finish.")
	waitGroup.Add(1)
	errors <- nil
	waitGroup.Wait()
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
		fmt.Println(err)
	}

	waitGroup.Done()
}
