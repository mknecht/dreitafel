package dreitafel

import "sync"

func CompileFmcBlockDiagramToDot(fmcSrcLines chan *string, dotSrcLines chan *string, errors chan error) {
	fmcdiagrams := make(chan *FmcBlockDiagram)
	dotdiagrams := make(chan DotGenerator)

	var wg sync.WaitGroup

	go KeepParsing(fmcSrcLines, fmcdiagrams, errors)
	go forwardFmcToDot(fmcdiagrams, dotdiagrams)
	wg.Add(1)
	go GenerateDot(dotdiagrams, dotSrcLines, errors, &wg)

	wg.Wait()
	close(errors)

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
