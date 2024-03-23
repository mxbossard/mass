#! /bin/bash -e
scriptDir=$( dirname $( readlink -f $0 ) )

cd $scriptDir
go mod tidy
go install ./embedded/embedded.go
go install ./client/client.go

