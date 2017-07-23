package dreitafel

import (
	"fmt"
	"sync"
)

func CompileToSvgToStdout(input string) {
	lines := make(chan *string, 500) // lines are independent; statements don't span lines yet
	diagrams := make(chan *FmcBlockDiagram)
	errors := make(chan error, 500) // errors are independent

	var waitGroup sync.WaitGroup

	go KeepParsing(lines, diagrams, errors, &waitGroup)
	go DotGenerator{}.Generate(diagrams, errors, &waitGroup)
	go printUntilNil(errors, &waitGroup)

	// feed the compiler the input
	lines <- &input

	fmt.Println("waiting for parser to finish.")
	waitGroup.Add(1)
	lines <- nil // end token
	waitGroup.Wait()

	fmt.Println("waiting for diagram builder to finish.")
	waitGroup.Add(1)
	diagrams <- nil // end token
	waitGroup.Wait()

	fmt.Println("waiting for error handler to finish.")
	waitGroup.Add(1)
	errors <- nil
	waitGroup.Wait()
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
