# Dreitafel
**Heads up, this is still very much in development (2018-May-05)**, the compiler works for very simple diagrams, [try.dreitafel.org](http://try.dreitafel.org/) and [view.dreitafel.org](http://view.dreitafel.org/) are there, but that's it. See you later for something actually useful. :)

Dreitafel helps you document and discuss the architecture and design of your software.

At its core, Dreitafel is a modeling language plus some tools, such as a compiler.


For example, this is how you could model `grep`:

```
(text lines) -> [grep] -> (matches)
```

Dreitafel turns this text-based diagram into an image:

![grep searches text lines for matches](http://view.dreitafel.org/gh/?format=svg&repository=mknecht%2Fdreitafel&path=docs%2Freadme%2Fgrep.fmc)

(See the [source code](docs/readme/grep.fmc), it's [live rendered](http://try.dreitafel.org/).)

## Try it

Head over to [try.dreitafel.org](http://try.dreitafel.org/) to try it yourself!


### Trying locally

You need Docker installed for the simple version to work.
The Docker image is this one: [muratk/dreitafel](https://hub.docker.com/r/muratk/dreitafel/).

Check out the repository and convert a text diagram to PNG:

```
./try.sh "(Wood) -> [Oven] -> (Heat)" cozy.png
```

Or, with docker:

```
echo "(Wood) -> [Oven] -> (Heat)" | docker run -i --rm muratk/dreitafel sh -c "/usr/bin/dreitafel 2>/dev/null | dot -Tsvg" >cozy.svg
```

(For some reason, for the PNG version the characters get lost when piping them between Dreitafel and dot.)

Or, manually:

```
make dreitafel

echo '(Wood) -> [Oven] -> (Heat)' | ./dreitafel | dot -Tpng > cozy.png
```

The webcompiler you can run like so:

```
make dreitafel-web && ./dreitafel-web
```

and then visit [this url](http://localhost:8080/fmc/?format=svg&diagram=(staticFiles)%20-%3E%20[NGINX]%20----%3E%20(logFiles)) to see the SVG for `(staticFiles) -> [NGINX] ----> (logFiles)`

## Why Dreitafel?

To my mind, one of the major contributions to documentation in open-source projects
was GitHub's READMEs. Many projects are documented purely in this single file,
relying on the browser search for nagivation. And it works.

It works, because those Markdown files are *simple* and *easy to change*.

Modeling software using diagrams should be simple and easy to change, too.
To do that, Dreitafel relies an existing modeling language,
namely the Block Diagrams of the [Fundamental Modeling Concepts (FMC)](http://fmc-modeling.org/).

FMC is not restricted to software, though.
It really is about *modeling systems*.
Have a look at this example of how we could model an oven:

 ```
(Wood) -> [Oven] -> (Heat)
 ```

What FMC *is* made for is communication.

It's meant for you and me to talk about how an oven works.
Or a compiler.
Or a database system.
Or a cloud deployment.

Wait, what do we need Dreitafel for again?
To make things look beautiful:

![You put wood inside and get heat out: The basics of an oven :)](http://view.dreitafel.org/gh/?format=svg&repository=mknecht%2Fdreitafel&path=docs%2Freadme%2Foven.fmc)

The technology accessible to most people is text.
Dreitafel defines a text-based version of FMC Block Diagrams.

Additionally, Dreitafel will comprise of the following tools…

- [x] a **compiler** from a text DSL to graphviz dot.
- [x] a **web-version of the compiler** so you can send text and get back an image, see http://view.dreitafel.org
- [x] a **viewer of GitHub hosted Dreitafel source code**: You put your diagram source in a textfile on GitHub,
      and in your README link an image to the viewer. Whenever you change the textfile, the image will automatically
      by updated.
- [x] a **playground** for diagrams at http://try.dreitafel.org
- [ ] a **paste-bin** for diagrams
- [ ] a **viewer of Markdown documents**. These may include Dreitafel source code, which is then replaced by images

Eventually, it would be great to get Jekyll and Sphinx support. The viewers would then be somewhat superfluous.

### Related projects

* [PlantUML](http://plantuml.com/) renders text-based UML diagrams.
  There are [tons of plugins](http://plantuml.com/running) so if UML
  is what you want, your search may be at an end.
* http://www.asciidraw.com/#Draw
* [esimov/diagram](https://github.com/esimov/diagram) is a “CLI app to convert ascii arts into hand drawn diagrams”. Love the idea, it's one of the later features
  I'd like to see for FMC block diagrams, because they're *designed*
  to be hand-drawn.
* [Graphviz](http://www.graphviz.org/) is open source graph visualization software.
  It's used in Dreitafel to render FMC Block Diagrams.
* [Tinkerpop](https://tinkerpop.apache.org/) is a graph computing framework for both graph databases (OLTP) and graph analytic systems (OLAP).

See also the [list of Graphviz-related projects](http://www.graphviz.org/content/resources).

## The Current Architecture of Dreitafel

In its first stage, Dreitafel is a compiler from FMC source code to graphviz' dot.
The latter is then used to generate the actual image.

The following diagram illustrates this (generated with Dreitafel and graphviz of course):

![Integration of Dreitafel with graphviz](http://view.dreitafel.org/gh/?format=svg&repository=mknecht%2Fdreitafel&path=docs%2Freadme%2Fdreitafel-graphviz-flow.fmc)

## The (Planned) Architecture of Dreitafel

Dreitafel will consist of three main components:

![Main components of Dreitafel](http://view.dreitafel.org/gh/?format=svg&repository=mknecht%2Fdreitafel&path=docs%2Freadme%2Fmain-components.fmc)

* Compiler: Reads the text source for a diagram and produces the same block diagram nicely rendered as PNG or SVG.
* Playground webapp: Gist/JS-Fiddle clone to play with and link diagrams
* Markdown-viewer webapp: Reads diagram sources, or GH-flavored Markdown from the web
  and renders them as PNG/SVG or HTML, respectively. This will allow you to read a README.md from
  GitHub with beautiful graphs!


### The compiler

The compiler reads Dreitafel source code
and spits out dot-generated graphs of the same FMC Block diagrams.
It consists of the following components:

![Compiler Overview](http://view.dreitafel.org/gh/?format=svg&repository=mknecht%2Fdreitafel&path=docs%2Freadme%2Fcompiler-overview.fmc)

Each of these components is running as its own goroutine;
communication between them happens with channels.

First, the [Lexer](lexer.go) reads the source code of the diagram.
Its job is to recognize individual diagram elements, such as Actors and Storages.

Any syntax element it recognizes, i.e. `AST Element` from Abstract Syntax Tree, is put into a queue (a Go channel) to be read by the Parser.

![Dataflow Lexer to Parser](http://view.dreitafel.org/gh/?format=svg&repository=mknecht%2Fdreitafel&path=docs%2Freadme%2Fcompiler-dataflow-lexer.fmc)

Anything the Lexer cannot recognize results in an error.
The rest of the line will be skipped and forwarded to the ErrorHandler.

Then, the **Parser** reads the diagram elements recognized by the Lexer.
It's job is to assemble a valid diagram from the individual parts.

![Dataflow Parser to Dot Generator](http://view.dreitafel.org/gh/?format=svg&repository=mknecht%2Fdreitafel&path=docs%2Freadme%2Fcompiler-dataflow-parser.fmc)

The diagram is then forwarded to the Dot Generator.
The AST Elemeents that did not make sense,
i.e. that could not be used to make a valid diagram,
are ignored and forwarded to the error handler.

Example of errors:

* An actor reads from an actor (FMC block diagrams are bipartite graph)
* A channel is dangling, i.e. not connected to anything.

The **Dot Generator** takes the valid FMC block diagram model,
and produces a dot graph representing the FMC diagram.


## The road ahead

### Next technical steps

To remember where I left off:

* debug view with fmc source, dot source, output and compile errors at `/debug`
* errors for try.dreitafel.org
* Make sure “block” diagram is part of the url, for example view.dreitafel.org/fmc-blocks/, so that other FMC parts can be added
* subgraph for each line as simple layout hints

### How to build the markdown viewer

* blackfriday for parsing

### Roadmap

* Minimal deployment
  * [X] Diagram elements
    * [X] Actor
    * [X] Storage
    * [X] Actor reads to Storage
    * [X] Actor writes to Storage.
  * [X] Compiler for Dreitafel text-syntax to graphviz' dot.
  * [x] web-version of the compiler
  * [x] viewer for GH-hosted Dreitafel source files
  * [x] Deploy view.dreitafel.org/fmc-blocks/
  * [x] Deploy try.dreitafel.org
* Publish
  * Add diagram elements and statements.
    * [ ] modifying access
    * [ ] unidirectional channel
    * [x] bidirectional channel
  * [ ] Compiler GH-flavored markdown with FMC block diagrams to HTML.
  * [x] Deploy www.dreitafel.org
  * [ ] Deploy view.dreitafel.org/md/
  * [ ] Webapp live-rendering this README.
  * [ ] Create FMC syntax guide, as documentation, and eat-your-own-dogfood.
  * [ ] CLI documentation
  * [x] Docker image with binaries
  * [ ] Documentation for Docker image
  * [ ] Logo! (Of course)
  * [ ] Request addition to the [Graphviz list](http://www.graphviz.org/content/resources)
* 1.0
  * [ ] CLI <3 — i/o with files/stdin/stdout, proper config
  * [ ] Styling of the diagram, making it look more hand-drawn
  * [x] Beautiful titles, with spaces, all kinds of characters.
  * [ ] Syntax: Comments
  * [ ] Create badge. :)
  * [x] Build & deploy gist/jsfiddle/play equivalent.
  * [ ] Human Actor
  * [ ] Support for IDs to centralize common attributes when an element re-occurs
  * [ ] multi-line elements
  * [ ] structure variance
* The real deal
  * [ ] 2d connections (vertical, diagonal and freeflowing)
  * [ ] U-formed actor
  * [ ] Nest elements inside an actor or storage to group them.
* Proofs of concept:
  * [ ] Model a physical machine
  * [ ] Document the engageSPARK software architecture.
  * [ ] Document an existing open-source project. (maybe gunicorn or sanic?)
* Possible features
  * Drawing styles: whiteboard, formal, chalk, business-flashy
  * Links: Within diagram, within page, external  - org syntax?
  * Emacs mode <3
  * Comments: Using GH issues?
  * Changes: Using GH PRs?
  * Composability: include diagrams
  * Layout hints
  * Zooming: Step “into” an element to view its details.
  * Printable version
  * More than block diagrams, maybe simple flow diagrams, too.

## Milestones

* 2018-May-06 Added rendering of GitHub-hosted FMC source code!
* 2017-Sep-09 Added channels!

## Commands to remember:


```
# Convert dot to a PNG:
dot -Tpng simple.dot > simple.png
# Convert dot to a SVG:
dot -Tsvg simple.dot > simple.svg
```
