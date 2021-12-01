#! /bin/bash -e
scriptDir=$( dirname $( readlink -f $0 ) )

modules="src/mby.fr/mass/cmd src/mby.fr/mass/internal"

for mod in $modules; do
	>&2 echo "Testing module $mod ..."
	cd $scriptDir/$mod
	go test -cover "$@" ./...
done
