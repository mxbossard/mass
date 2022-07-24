#! /bin/bash -e
scriptDir=$( dirname $( readlink -f $0 ) )

cd $scriptDir/src/mby.fr/mass
go mod tidy
go install

