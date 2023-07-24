#! /bin/bash -e
scriptDir=$( dirname $( readlink -f $0 ) )
cd $scriptDir

goVersion="$1"

docker build --build-arg "GO_VERSION=$goVersion" -t goid-dev:$goVersion-latest .
