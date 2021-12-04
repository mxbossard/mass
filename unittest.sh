#! /bin/bash -e
scriptDir=$( dirname $( readlink -f $0 ) )

packages="src/mby.fr/mass/cmd src/mby.fr/mass/internal src/mby.fr/utils"

for pkg in $packages; do
	>&2 echo ""
	>&2 echo "Testing package(s) $pkg ..."
	cd $scriptDir/$pkg
	go test -cover "$@" ./...
done
