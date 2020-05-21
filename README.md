# Gorun (for Go 1.11+ with Go Modules enabled)

The standard `go run` method to run commands comes with one big problem: it compiles the code into a temporary directory which is meant to be transient and deleted. `gorun` will just run `go build && ./[name of the executable]` in the root directory of your `main.go` source code.

The tool also builds in some useful things like measuring `go build` times (makes it easy to sell `Go` to your company if you consistently show the compile times compared with the other technologies you use) as well as a built in split stream to both standard output and `output.txt` file.

## Installation

```
go get -u github.com/gilgameshskytrooper/gorun
```
