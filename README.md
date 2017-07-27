# Dreitafel
Discussing software architecture, made simple

## Trying it

Convert a diagram to dot:

```
./try.sh "[Engine] (Gasoline) (Oil)"
```

Convert dot to an image:

```
dot -Tpng testground/simple.dot > testground/simple.png
```

## The (Planned) Architecture of Dreitafel

Dreitafel will consist of three main components:

```
[ Playground webapp ]      -o- [          ]
                               [ Compiler ]
[                        ] -o- [          ]
[ Markdown-viewer webapp ]
[                        ] -o- [ GitHub ]
```

* Compiler: Reads the text source for a diagram and produces the same block diagram nicely rendered as PNG or SVG.
* Playground webapp: Gist/JS-Fiddle clone to play with and link diagrams
* Markdown-viewer webapp: Reads diagram sources, or GH-flavored Markdown from the web
  and renders them as PNG/SVG or HTML, respectively. This will allow you to read a README.md from
  GitHub with beautiful graphs!


### The compiler

The compiler reads Dreitafel source code
and spits out dot-generated graphs of the same FMC Block diagrams.
It consists of the following components:

```dreitafel:fmcblock
[ Reader ] -o- [ Lexer ] -o- [ Parser ] -o- [ DotGenerator ]

[ ErrorHandler ]
```

Each of these components is running as its own goroutine;
communication between them happens with channels.

First, the [Lexer](lexer.go) reads the source code of the diagram.
Its job is to recognize individual diagram elements, such as Actors and Storages.

Any syntax element it recognizes, i.e. `AST Element` from Abstract Syntax Tree, is put into a queue (a Go channel) to be read by the Parser.

```dreitafel:fmcblock
( Source ) -> [ Lexer ] -> ( AstElement ) -> [ Parser ]
                        -> ( Unrecognizable text ) -> [ ErrorHandler ]
```

Anything the Lexer cannot recognize results in an error.
The rest of the line will be skipped and forwarded to the ErrorHandler.

Then, the **Parser** reads the diagram elements recognized by the Lexer.
It's job is to assemble a valid diagram from the individual parts.

```dreitafel:fmcblock
( AST element ) -> [ Parser ] -> ( Diagram )      -> [ Dot Generator ]
                              -> ( Invalid AST Elements ) -> [ ErrorHandler ]
```

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

## Next technical steps

To remember where I left off:

* Use logging system, and log to stderr, so that the diagram (stdout) can be piped
* Read from stdin, the goal being pipelines like: `cat simple.fmc | dreitafel | dot -Tsvg`
* Write some simple system tests

### Roadmap (Ideas)

* Minimal deployment
  * [X] Diagram elements
    * [X] Actor
    * [X] Storage
    * [X] Actor reads to Storage
    * [X] Actor writes to Storage.
  * [ ] Compiler for FMC block diagram text-syntax to graphviz' dot.
  * [ ] Compiler GH-flavored markdown with FMC block diagrams to HTML.
  * [ ] Webapp live-rendering this.
  * [ ] Deploy this service on a host.
  * [ ] Create badge. :)
  * [ ] Create FMC syntax guide, as documentation, and eat-your-own-dogfood.
  * [ ] Build & deploy gist/jsfiddle/play equivalent.
* Add diagram elements and statements.
  * [ ] support for arbitrary number of actors & storages
  * [ ] reading access
  * [ ] modifying access
  * [ ] unidirectional channel
  * [ ] bidirectional channel
  * [ ] Human Actor
  * [ ] Support for IDs to centralize common attributes when an element re-occurs
  * [ ] multi-line elements
  * [ ] 2d connections (vertical, diagonal and freeflowing)
  * [ ] structure variance
  * [ ] U-formed actor
  * [ ] Nest elements inside an actor or storage to group them.
* Proofs of concept:
  * [ ] Model a physical machine
  * [ ] Document the engageSPARK software architecture.
  * [ ] Document an existing open-source project. (maybe gunicorn or sanic?)
* Possible features
  * Drawing styles: whiteboard, formal, chalk, business-flashy
  * Links: Within diagram, within page, external  - org syntax?
  * Comments: Using GH issues?
  * Changes: Using GH PRs?
  * Composability: include diagrams
  * Layout hints
  * Zooming: Step “into” an element to view its details.
  * Printable version
