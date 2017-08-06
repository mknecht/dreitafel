package main

import (
	"bufio"
	"dreitafel"
	"fmt"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	fmcSrcLines := make(chan *string, 500) // lines are independent; statements don't span lines yet
	dotSrcLines := make(chan *string, 500)
	errors := make(chan error, 500) // errors are independent

	var wg sync.WaitGroup

	go readStdinIntoChannel(fmcSrcLines)
	go dreitafel.CompileFmcBlockDiagramToDot(fmcSrcLines, dotSrcLines, errors)
	wg.Add(1)
	go printToStdout(dotSrcLines, &wg)
	wg.Add(1)
	go printToStderr(errors, &wg)

	wg.Wait()
}

func readStdinIntoChannel(lines chan *string) {
	defer close(lines)
	reader := bufio.NewReader(os.Stdin)

	for line, err := reader.ReadString('\n'); err == nil; line, err = reader.ReadString('\n') {
		newline := line // need new memory location
		log.Debugf("Read line: '%v' @%v", line, &line)
		lines <- &newline
	}
	log.Debugf("Input reader done.")
}

func printToStdout(lines <-chan *string, wg *sync.WaitGroup) {
	for line := range lines {
		fmt.Fprintf(os.Stdout, "%v\n", *line)
	}

	wg.Done()
}

func printToStderr(errors <-chan error, wg *sync.WaitGroup) {
	for err := range errors {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	wg.Done()
}
