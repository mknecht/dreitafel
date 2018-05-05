package web

import (
	"bufio"
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

var responseFormats = map[string]string{
	"dot": "text/dot",
	"png": "image/png",
	"svg": "image/svg+xml",
}

var validFormats []string

func index(w http.ResponseWriter, r *http.Request) {
	txt := `
<p>Welcome to <a href="https://github.com/mknecht/dreitafel">Dreitafel</a> viewer!</p>

<p>See this example:</p>

<pre>[Actor] -&gt; (Storage)</pre>

<p>A diagram can be generated dynamically using this URL: /fmc/?format=svg&diagram=[Actor]->(Storage)</p>

<p>And here's the diagram:</p>

<p><img src="/fmc/?format=svg&diagram=[Actor]->(Storage)"/></p>

`
	fmt.Fprintf(w, txt)
}

func compileFmcBlockDiagramFromQueryString(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	format := "png"
	if len(r.Form["format"]) > 0 {
		format = r.Form["format"][len(r.Form["format"])-1]
	}

	contentType, ok := responseFormats[format]
	if !ok {
		http.Error(w, fmt.Sprintf("Supported values for 'format' param are: %v", validFormats), http.StatusBadRequest)
		return
	}

	if len(r.Form["diagram"]) == 0 {
		http.Error(w, "Need 'diagram' query param with Dreitafel source code", http.StatusBadRequest)
		return
	}
	diagram := r.Form["diagram"][len(r.Form["diagram"])-1]

	fmcSrcLines := make(chan *string) // lines are independent; statements don't span lines yet

	dotSrcLines := make(chan *string, 500)
	errors := make(chan error, 500) // errors are independent

	var wg sync.WaitGroup

	go dreitafel.CompileFmcBlockDiagramToDot(fmcSrcLines, dotSrcLines, errors)

	w.Header().Add("Content-Type", contentType)
	if format == "dot" {
		writeDotSrcLines(w, dotSrcLines, errors, &wg)
	} else {
		writeDotGeneratedImage(format, w, dotSrcLines, errors, &wg)
	}

	fmcSrcLines <- &diagram
	close(fmcSrcLines)

	wg.Wait()
}

func compileFmcBlockDiagramFromGithubUrl(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	format := "png"
	if len(r.Form["format"]) > 0 {
		format = r.Form["format"][len(r.Form["format"])-1]
	}

	contentType, ok := responseFormats[format]
	if !ok {
		http.Error(w, fmt.Sprintf("Supported values for 'format' param are: %v", validFormats), http.StatusBadRequest)
		return
	}

	if len(r.Form["repository"]) == 0 {
		http.Error(w, "Need 'repository' query param with full name of the Github repository, for example 'mknecht/dreitafel'", http.StatusBadRequest)
		return
	}
	repository := r.Form["repository"][len(r.Form["repository"])-1]

	if len(r.Form["path"]) == 0 {
		http.Error(w, "Need 'path' query param with path to file in the Github repository, for example 'examples/car.fmc'", http.StatusBadRequest)
		return
	}
	path := r.Form["path"][len(r.Form["path"])-1]

	response, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/%v/master/%v", repository, path))
	if err != nil {
		return
	}
	defer response.Body.Close()

	fmcSrcLines := make(chan *string) // lines are independent; statements don't span lines yet
	dotSrcLines := make(chan *string, 500)
	errors := make(chan error, 500) // errors are independent

	var wg sync.WaitGroup

	go dreitafel.CompileFmcBlockDiagramToDot(fmcSrcLines, dotSrcLines, errors)

	w.Header().Add("Content-Type", contentType)
	if format == "dot" {
		writeDotSrcLines(w, dotSrcLines, errors, &wg)
	} else {
		writeDotGeneratedImage(format, w, dotSrcLines, errors, &wg)
	}

	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		line := scanner.Text()
		fmcSrcLines <- &line
	}
	if scanner.Err() != nil {
		errors <- scanner.Err()
	}
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
	validFormats = make([]string, len(responseFormats))
	i := 0
	for key, _ := range responseFormats {
		validFormats[i] = key
		i += 1
	}

	http.HandleFunc("/fmc/", compileFmcBlockDiagramFromQueryString)
	http.HandleFunc("/gh/", compileFmcBlockDiagramFromGithubUrl)
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}
