#! /bin/bash
scriptDir=$( dirname $( readlink -f $0 ) )

export GOBIN="./bin"

ctDurationInSec=60
ctId=$( docker run --rm -d --workdir=/root -e GOBIN="/root/$GOBIN" --volume=$scriptDir:/root:rw golang:1.17 sleep $ctDurationInSec )

docker exec -it $ctId go "$@"
