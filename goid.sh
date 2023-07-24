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
goUid="$( id -u )"
goGid="$( id -g )"

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
goVersion="1.20"
ctImage="goid-dev:${goVersion}-latest"
#ctImage="go-dev-image:latest"
ctName="goid_$( basename $scriptDir )_go${goVersion}_$goUid"

buildCmd="$scriptDir/goid/build.sh $goVersion"

dockerSocketVolArg=""
if [ -e "/var/run/docker.sock" ]; then
	# Mount host docker socket in container
	dockerSocketVolArg="--volume=/var/run/docker.sock:/var/run/docker.sock"
fi

#runCmd="docker run --rm -d --name=$ctName --user=$( id -u ):$( id -g ) --workdir=$rootDirInCt -e GOBIN -e GOCACHE -e GOENV -e GOMODCACHE -e GOPATH --volume=$hostGoDir:/go:rw --volume=$scriptDir:$rootDirInCt:rw $ctImage"
#runCmd="docker run --privileged --rm -d --name=$ctName -e GOBIN -e GOCACHE -e GOENV -e GOMODCACHE -e GOPATH --volume=$hostGoDir:/go:rw --volume=$scriptDir:$rootDirInCt:rw --volume=/var/run/docker.sock:/var/run/docker.sock $ctImage"
runCmd="docker run --rm -d --name=$ctName --entrypoint=/bin/sh -e GOBIN -e GOCACHE -e GOENV -e GOMODCACHE -e GOPATH --volume=$hostGoDir:/go:rw --volume=$scriptDir:$rootDirInCt:rw $dockerSocketVolArg $ctImage -c 'sleep inf'"
#echo $runCmd

#>&2 echo "Executing go in a container on workspace: $rootDir within container workdir: $workDirInCt ..."
#>&2 echo "GOMODCACHE: $hostModCacheDir ..."

ctId=$( docker ps -f name=$ctName -q )
if [ -z "$ctId" ]; then
    #ctId=$( $runCmd sleep $ctDurationInSec )
    $buildCmd
    ctId=$( eval "$runCmd" )
    >&2 echo "goid container running with id: $ctId"
    #sleep 1
    # Create go user and add it to docker group
    dockerGid="$( getent group docker | awk -F: '{printf $3}' )"
    docker exec "$ctId" addgroup --gid $dockerGid docker
    docker exec "$ctId" adduser --quiet --system --no-create-home --shell /bin/sh --home $workDirInCt --uid $goUid go
    docker exec "$ctId" usermod -a -G docker go
    echo "Wait Docker dameon started ..."
    #while ! 2>&1 >/dev/null docker exec "--user=$goUser" "$ctId" docker ps; do
    while ! 2>&1 >/dev/null docker ps; do
        sleep 1
        #echo -n "."
    done
fi

#echo "workdir=$workDirInCt"
docker exec "--workdir=$workDirInCt" "--user=$goUid" "$ctId" go "$@"
