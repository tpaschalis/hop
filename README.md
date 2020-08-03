# hop

![gopher logo](gopher.png)

Hop is an interactive test runner for Go. It was inspired by Jest's *--watch* mode and [this](https://twitter.com/felixge/status/1286359708799062016) tweet. It's called Hop, because well, you just hop in and start running your tests! Coffee and music are optional, but recommended!


## Are you just wrapping some commands?

Yes, I'm candy wrapping the following commands in a for-loop. Plus colors! You could probably do this yourself using some bash + aliases, but where's the fun in that?!
```
go test -v ./...
go test -v ./... -count=1
go test -v ./... -list .
go test -v ./... -list <pattern>
go test -v ./... -run <pattern>
```



## ROADMAP
At this point, the output is *not* streaming, but is printed all at once, at the end of a run.

This is problematic, for long running tests or larger suites, but I'm working on fixing this soon!

