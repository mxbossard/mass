#! /bin/bash -e
scriptDir=$( dirname $( readlink -f $0 ) )

#packages="src/mby.fr/utils src/mby.fr/mass/cmd src/mby.fr/mass/internal"
packages="src/mby.fr/utils src/mby.fr/mass"

success=true
for pkg in $packages; do
	>&2 echo ""
	>&2 echo "Testing package(s) $pkg ..."
	cd $scriptDir/$pkg
	go mod tidy || true
	go test -cover "$@" ./... || success=false
done

$success