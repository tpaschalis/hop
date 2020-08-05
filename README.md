# hop

![gopher logo](gopher.png)

Hop is an interactive test runner for Go. It was inspired by Jest's [*--watch*](https://jestjs.io/docs/en/cli#--watch) mode and [this](https://twitter.com/felixge/status/1286359708799062016) tweet. It's called Hop, because well, you just hop in and start running your tests! Coffee and music are optional, but recommended!

## Installation and Usage

You can get `hop` by running.
```
git clone https://github.com/tpaschalis/hop.git
cd hop
go install .
```

It comes with two modes; the default plain old boring B&W, as well as fabulous 🌈🎨🎆🍭colourized output.  
Just provide the `--colour` flag!

## Here's Hop in action!

*Monochrome mode*
<a href="https://asciinema.org/a/351777" target="_blank"><img src="https://asciinema.org/a/351777.svg" /></a>

## Are you just wrapping some commands?

Yes, I'm candy wrapping the following commands in a for-loop. Plus colours! You could probably do this yourself using some bash + aliases, but where's the fun in that?!
```
go test -v ./...
go test -v ./... -count=1
go test -v ./... -list .
go test -v ./... -list <pattern>
go test -v ./... -run <pattern>
go test -v ./... -bench=. -run ^Benchmark
```

## How does it work?
Hop is a simple REPL, which runs the aforementioned `go test` commands. As we mentioned, it features two modes, monochrome (by default) which displays the stdout/stderr, and colourful (enabled by using `--colour`).

Hop uses two custom `io.Writer` types, *monochromeWriter* and *colourizedWriter*. After that, it's just using [`exec.Command`](https://golang.org/pkg/os/exec/#Command), and `io.MultiWriter` to stream the output to our custom writer.

The inspiration, Jest, works a little bit different. It watches for files for changes and reruns tests related to changed files only. I tried using [fsnotify](https://github.com/fsnotify/fsnotify) for the same effect, but ultimately I decided it isn't worth it. 

Changes to a specific file might have effects on a different package, so something like fsnotify would not catch all the cases. The good news is that the `go test` tool takes care of all this, by caching tests that wouldn't change, and rerunning any tests that are affected by changes, which is what we'd ultimately want! 

## These colors suck!
First off, I'm using *colour*, with an 'ou'. Mainly because 'monochrome' and 'colourized' have the same number of letters.

Secondly, that's very probable! You could change them by tinkering with the `colors.go` file (see the irony?).

I'm probably going to add some easy way to configure them from a text file. I'm a below average designer, so feel free to propose changes!