
## Setup go in Docker
mkdir -p ~/bin
ln -s $PWD/goid.sh ~/bin/go

## Your first program
see https://go.dev/doc/code
modulePath=<module-path> # Example: mby.fr/mass
mkdir -p src/$modulePath
cd src/$modulePath
go mod init $modulePath

## Using cobra
go get github.com/spf13/cobra/cobra
go install github.com/spf13/cobra/cobra
ln -s $PWD/bin/cobra ~/bin/cobra
export GOPATH=$PWD
cobra
