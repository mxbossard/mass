
## Setup go in Docker
ln -s ~/.local/bin/go $PWD/goid.sh

## Your first program
see https://go.dev/doc/code

## Using cobra
go install github.com/spf13/cobra/cobra
ln -s ~/.local/bin/ $PWD/bin/cobra
export GOPATH=$PWD
cobra
