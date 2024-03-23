#! /bin/bash 
set -e
scriptDir=$( dirname $( readlink -f $0 ) )

#packages="src/mby.fr/utils src/mby.fr/mass/cmd src/mby.fr/mass/internal"
#packages="src/mby.fr/k8s2docker src/mby.fr/utils src/mby.fr/mass src/mby.fr/scribble"
packages="src/mby.fr/utils src/mby.fr/mass"

success=true
for pkg in $packages; do
	>&2 echo ""
	>&2 echo "Testing package(s) $pkg ..."
	cd $scriptDir/$pkg
	go mod tidy || true
	cmd="go test -cover"

	# check if go files or pathes were supplied in args or launch all tests
	options=""
	pathes=""
	for p in "$@"; do
		if echo "$p" | egrep '.*\.go|./' > /dev/null; then
			#test -e "$p" && pathes="$pathes $p" || true
			pathes="$pathes $p"
		else
			options="$options $p"
		fi
	done

	if [ -n "$pathes" ]; then
		for p in $pathes; do
			test -e "$p" || continue
			testCmd="$cmd $options $p"
			>&2 echo "Running [ $testCmd ] in dir [ $pkg ] ..."
			$testCmd || success=false
		done
	else
		testCmd="$cmd $options ./..."
		>&2 echo "Running [ $testCmd ] in dir [ $pkg ] ..."
		$testCmd || success=false
	fi

done

echo
if $success; then
	echo SUCCESS
else
	echo FAILURE
fi

$success
