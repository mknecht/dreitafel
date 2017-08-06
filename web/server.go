package web

import (
	"dreitafel"

	"fmt"
	"net/http"
	"os"
	"sync"
)

func index(w http.ResponseWriter, r *http.Request) {
	txt := `
<p>Welcome to <a href="https://github.com/mknecht/dreitafel">Dreitafel</a>!</p>

<p>Try this: <a href="/fmc/q?format=png&diagram=%5BActor%5D%20-%3E%20%28Storage%29</p>

`
	fmt.Fprintf(w, txt)
}

func compileFmcBlockDiagramFromQueryString(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// TODO input validation & error handling :]
	fmcSrcLines := make(chan *string) // lines are independent; statements don't span lines yet

	dotSrcLines := make(chan *string, 500)
	errors := make(chan error, 500) // errors are independent

	var wg sync.WaitGroup

	go dreitafel.CompileFmcBlockDiagramToDot(fmcSrcLines, dotSrcLines, errors)

	wg.Add(1)
	go func() {
		for line := range dotSrcLines {
			fmt.Printf("%v\n", *line)
			fmt.Fprintf(w, "%v\n", *line)
		}

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for err := range errors {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}

		wg.Done()
	}()

	fmcSrcLines <- &r.Form["diagram"][0]
	close(fmcSrcLines)

	wg.Wait()
}

func ListenAndServe() {
	http.HandleFunc("/fmc/q", compileFmcBlockDiagramFromQueryString)
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}
