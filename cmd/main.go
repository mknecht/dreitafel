package main

import (
	"bufio"
	"dreitafel"
	"os"

	log "github.com/Sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	lines := make(chan *string, 500) // lines are independent; statements don't span lines yet

	go func() {
		defer close(lines)
		reader := bufio.NewReader(os.Stdin)

		for line, err := reader.ReadString('\n'); err == nil; line, err = reader.ReadString('\n') {
			newline := line // need new memory location
			log.Debugf("Read line: '%v' @%v", line, &line)
			lines <- &newline
		}
		log.Debugf("Input reader done.")
	}()

	dreitafel.CompileToSvgToStdout(lines)
}
