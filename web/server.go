package web

import (
	"dreitafel"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"fmt"
	"net/http"
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

	if r.Form["format"][0] == "png" {
		w.Header().Add("Content-Type", "image/png")
		writeDotGeneratedImage("png", w, dotSrcLines, errors, &wg)
	} else if r.Form["format"][0] == "svg" {
		w.Header().Add("Content-Type", "image/svg+xml")
		writeDotGeneratedImage("svg", w, dotSrcLines, errors, &wg)
	} else {
		writeDotSrcLines(w, dotSrcLines, errors, &wg)
	}

	fmcSrcLines <- &r.Form["diagram"][0]
	close(fmcSrcLines)

	wg.Wait()
}

func writeDotSrcLines(w io.Writer, dotSrcLines chan *string, errors chan error, wg *sync.WaitGroup) {
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
}

func writeDotGeneratedImage(format string, w io.Writer, dotSrcLines chan *string, errors chan error, wg *sync.WaitGroup) {
	cmd := exec.Command("dot", "-T"+format)
	in, err := cmd.StdinPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	errp, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	wg.Add(1)
	go func() {
		for line := range dotSrcLines {
			fmt.Printf("%v\n", *line)
			fmt.Fprintf(in, "%v\n", *line)
		}

		in.Close()
		fmt.Println("Done with writing src to stdin of dot.")
		wg.Done()
	}()

	wg.Add(1)
	go func() {

		fmt.Printf("Starting to read.\n")
		b, err := ioutil.ReadAll(out)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(b)
		fmt.Printf("Done with writing %v PNG bytes as response.\n", len(b))
		wg.Done()
	}()

	wg.Add(1)
	go func() {

		fmt.Printf("Starting to read stderr.\n")
		b, err := ioutil.ReadAll(errp)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(b)
		fmt.Printf("Done with writing %v stderr as response.\n", len(b))
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for err := range errors {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}

		wg.Done()
	}()

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func ListenAndServe() {
	http.HandleFunc("/fmc/", compileFmcBlockDiagramFromQueryString)
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}
