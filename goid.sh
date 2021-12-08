#! /bin/bash -e
scriptDir=$( dirname $( readlink -f $0 ) )

# Prints out the relative path between to absolute paths. Trivial.
#
# Parameters:
# $1 = first path
# $2 = second path
#
# Output: the relative path between 1st and 2nd paths
# see: https://unix.stackexchange.com/questions/85060/getting-relative-links-between-two-paths
relpath() {
    local pos="${1%%/}" ref="${2%%/}" down=''

    while :; do
        test "$pos" = '/' && break
        case "$ref" in $pos/*) break;; esac
        down="../$down"
        pos=${pos%/*}
    done

    echo "$down${ref##$pos/}"
}

ctDurationInSec=7200
rootDir=$scriptDir
rootDirInCt=/tmp/goid
workDir=$PWD
workDirInCt=$rootDirInCt/$( relpath $rootDir $workDir/. )

export GOBIN="$rootDirInCt/bin"
export GOCACHE="$rootDirInCt/.cache/go-build"
export GOENV="$rootDirInCt/.config/go/env"
export GOMODCACHE="/go/pkg/mod"
export GOPATH="$rootDirInCt"

hostGoDir="$HOME/go"
hostModCacheDir="$hostGoDir/pkg/mod"
mkdir -p "$hostModCacheDir"

ctName="goid_$( basename $scriptDir )"

runCmd="docker run --rm -d --name=$ctName --user=$( id -u ):$( id -g ) --workdir=$rootDirInCt -e GOBIN -e GOCACHE -e GOENV -e GOMODCACHE -e GOPATH --volume=$hostGoDir:/go:rw --volume=$scriptDir:$rootDirInCt:rw golang:1.17"
#echo $runCmd

#>&2 echo "Executing go in a container on workspace: $rootDir within container workdir: $workDirInCt ..."
#>&2 echo "GOMODCACHE: $hostModCacheDir ..."

ctId=$( docker ps -f name=$ctName -q )
if [ -z "$ctId" ]; then
	ctId=$( $runCmd sleep $ctDurationInSec )
fi

#echo $ctId

# Debug
#docker exec -it $ctId "$@"
#exit 0

#docker exec $ctId go "$@"
docker exec --workdir=$workDirInCt $ctId go "$@"
