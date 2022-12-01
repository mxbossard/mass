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

#ctDurationInSec=3600
rootDir=$scriptDir
rootDirInCt=/home/goid
workDir=$PWD
workDirInCt=$rootDirInCt/$( relpath $rootDir $workDir/. )
goUser="$( id -u )"
goGroup="$( id -g )"
goUserAndGroup="$goUser:$goGroup"

export GOBIN="$rootDirInCt/bin"
export GOCACHE="$rootDirInCt/.cache/go-build"
export GOENV="$rootDirInCt/.config/go/env"
export GOMODCACHE="/go/pkg/mod"
export GOPATH="$rootDirInCt"

hostGoDir="$HOME/go"
hostModCacheDir="$hostGoDir/pkg/mod"
mkdir -p "$hostModCacheDir"

#ctImage="golang:1.18-alpine"
#ctImage="golang:1.18"
ctImage="goid-dev:latest"
#ctImage="go-dev-image:latest"
ctName="goid_$( basename $scriptDir )_$( id -u )"

buildCmd="$scriptDir/goid/build.sh"

#runCmd="docker run --rm -d --name=$ctName --user=$( id -u ):$( id -g ) --workdir=$rootDirInCt -e GOBIN -e GOCACHE -e GOENV -e GOMODCACHE -e GOPATH --volume=$hostGoDir:/go:rw --volume=$scriptDir:$rootDirInCt:rw $ctImage"
#runCmd="docker run --rm -d --name=$ctName -e GOBIN -e GOCACHE -e GOENV -e GOMODCACHE -e GOPATH --volume=$hostGoDir:/go:rw --volume=$scriptDir:$rootDirInCt:rw --volume=/var/run/docker.sock:/var/run/docker.sock $ctImage sleep 3600"
runCmd="docker run --rm -d --name=$ctName -e GOBIN -e GOCACHE -e GOENV -e GOMODCACHE -e GOPATH --volume=$hostGoDir:/go:rw --volume=$scriptDir:$scriptDir:rw --volume=/var/run/docker.sock:/var/run/docker.sock $ctImage sleep 3600"
#>&2 echo "docker run: $runCmd"

#>&2 echo "Executing go in a container on workspace: $rootDir within container workdir: $workDirInCt ..."
#>&2 echo "GOMODCACHE: $hostModCacheDir ..."

ctId=$( docker ps -f name=$ctName -q )
if [ -z "$ctId" ]; then
    #ctId=$( $runCmd sleep $ctDurationInSec )
    $buildCmd
    ctId=$( $runCmd )
    #sleep 1
    docker exec "$ctId" addgroup --gid "$goGroup" gofer
    docker exec "$ctId" adduser --system --no-create-home --home $scriptDir --uid "$goUser" --gid "$goGroup" gofer
    #docker exec "$ctId" addgroup docker
    #docker exec "$ctId" addgroup gofer docker
    docker exec "$ctId" chmod a+rw /var/run/docker.sock
    echo "Wait Docker dameon started ..."
    while ! 2>&1 >/dev/null docker exec "--user=$goUserAndGroup" "$ctId" docker ps; do
        sleep 1
        #echo -n "."
    done
fi

#echo "workdir=$workDirInCt"
docker exec "--workdir=$( pwd )" "--user=$goUser" "$ctId" go "$@"

