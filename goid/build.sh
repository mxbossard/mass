#! /bin/bash -e
scriptDir=$( dirname $( readlink -f $0 ) )

cd $scriptDir

docker build -t goid-dev:latest .
