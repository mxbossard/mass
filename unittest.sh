#! /bin/bash -e
scriptDir=$( dirname $( readlink -f $0 ) )

#packages="src/mby.fr/utils src/mby.fr/mass/cmd src/mby.fr/mass/internal"
packages="src/mby.fr/utils src/mby.fr/mass"

success=true
for pkg in $packages; do
	>&2 echo ""
	>&2 echo "Testing package(s) $pkg ..."
	cd $scriptDir/$pkg
	$go mod tidy || true
	cmd="$go test -cover $@"
	# check if go files or pathes were supplied in args or launch all tests
	echo "$cmd" | egrep '.*\.go|./' > /dev/null || cmd="$cmd ./..."
	>&2 echo "Running [ $cmd ] in dir [ $pkg ] ..."
	$cmd || success=false
done

echo
if $success; then
	echo SUCCESS
else
	echo FAILURE
fi

$success
