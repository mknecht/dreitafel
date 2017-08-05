# Dreitafel
Dreitafel helps you document and discuss the architecture and design of your software.

This is how you could model the basics of a webserver serving static files:

```
[Browser] --o-- [Webserver] <- (static files)
```

The language used is a text version of the Block Diagrams of the [Fundamental Modeling Concepts (FMC)](http://fmc-modeling.org/),

FMC is not restricted to software. 
It really is about *systems*.
Have a look at this example of how we could model an oven:

```
(Wood) -> [Oven] -> (Heat)
```

What FMC *is* made for is communication.

It's meant for you and me to talk about how an oven works.
Or a compiler.
Or a database system.
Or a car.

Wait, what do we need Dreitafel for?
To make things look beautiful:

![You put wood and get heat: The basics of an oven :)](examples/oven.png)

Dreitafel will comprise of…

- [x] a **compiler** from a text DSL to graphviz dot.
- [ ] a **viewer** of Markdown documents. These may include  Dreitafel source code, which is then replaced by images
- [ ] a **paste-bin and playground** for diagrams

## Trying it

You need Docker installed for the simple version to work.

Check out the repository and convert a text diagram to PNG:

```
./try.sh "(Wood) -> [Oven] -> (Heat)" cozy.png
```

Or, manually:

```
make dreitafel

echo '(Wood) -> [Oven] -> (Heat)' | ./dreitafel | dot -Tpng > cozy.png
```

## The Current Architecture of Dreitafel

In its first stage, Dreitafel is a compiler from FMC source code to graphviz' dot. 
The latter is then used to generate the actual image.

The following diagram illustrates this (generated with Dreitafel and graphviz of course):

![Integration of dreitafel with graphviz](examples/dreitafel.png)

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
( AST element ) -> [ Parser ] -> ( Diagram )              -> [ Dot Generator ]
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

* Write some simple system tests

### Roadmap

* Minimal deployment
  * [X] Diagram elements
    * [X] Actor
    * [X] Storage
    * [X] Actor reads to Storage
    * [X] Actor writes to Storage.
  * [X] Compiler for FMC block diagram text-syntax to graphviz' dot.
  * [ ] Compiler GH-flavored markdown with FMC block diagrams to HTML.
  * [ ] Webapp live-rendering this README.
  * [ ] Deploy this service on a host.
* Publish
  * Add diagram elements and statements.
    * [ ] modifying access
    * [ ] unidirectional channel
    * [ ] bidirectional channel
  * [ ] Create FMC syntax guide, as documentation, and eat-your-own-dogfood.
  * [ ] Logo! (Of course)
* 1.0
  * [ ] CLI <3 — i/o with files/stdin/stdout, proper config
  * [ ] Styling of the diagram, making it look more hand-drawn
  * [ ] Beautiful titles, with spaces, all kinds of characters.
  * [ ] Syntax: Comments
  * [ ] Create badge. :)
  * [ ] Build & deploy gist/jsfiddle/play equivalent.
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
  * Simple flow diagrams, too.

## Commands to remember:


```
# Convert dot to a PNG:
dot -Tpng simple.dot > simple.png
# Convert dot to a SVG:
dot -Tsvg simple.dot > simple.svg
```

