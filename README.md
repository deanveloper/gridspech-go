# gridspech-go

`gridspech-go` is an implementation of the puzzle game [gridspech](https://krackocloud.itch.io/gridspech) written in Go.

This implementation was written to make a solver, which the package can be found at [`gridspech-go/solve`](solve). The solver was written to find unintended solutions to puzzles, so don't use the solver until you have figured out the puzzle on your own :)

## CLI Installation

Binaries for common environments can be found on the [releases page](https://github.com/deanveloper/gridspech-go/releases). After this, you can add it to your PATH (ie in `/usr/local/bin` on \*nix systems).

To install the CLI from source, first make sure you have the [Go compiler](https://golang.org/dl/) installed. Then, run `go install "github.com/deanveloper/gridspech-go/solve/cmd/gs-solve@latest"`. The binary will appear in $GOBIN, which by default is `~/go/bin`
