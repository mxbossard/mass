#! /bin/bash -e
scriptDir=$( dirname $( readlink -f $0 ) )

ctDurationInSec=60
workDir=/tmp/goid

export GOBIN="$workDir/bin"
export GOCACHE="$workDir/.cache/go-build"
export GOENV="$workDir/.config/go/env"
export GOMODCACHE="/go/pkg/mod"

localModCacheDir="$HOME/go/pkg/mod"
mkdir -p "$localModCacheDir"

runCmd="docker run --rm -d --user=$( id -u ):$( id -g ) --workdir=$workDir -e GOBIN -e GOCACHE -e GOENV -e GOMODCACHE --volume=$localModCacheDir:$GOMODCACHE:rw --volume=$scriptDir:$workDir:rw golang:1.17"
#echo $runCmd

ctId=$( $runCmd sleep $ctDurationInSec )
#echo $ctId

# Debug
#docker exec -it $ctId "$@"
#exit 0

docker exec -it $ctId go "$@"
