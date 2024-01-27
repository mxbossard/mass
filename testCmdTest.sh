#! /bin/bash
set -e
scriptDir=$( dirname $( readlink -f $0 ) )

workspaceDir="/tmp/cmdtWorkspace"

echo "##### Building cmdtest binary ..."
export GOBIN="$scriptDir/bin"
cd $scriptDir/src/mby.fr/cmdtest
go install

rm -rf -- "$workspaceDir"

cmd="$GOBIN/cmdtest"
ls -lh "$cmd"

mkdir -p "$scriptDir/.tmp"
reportFile="$( mktemp "$scriptDir/.tmp/XXXXXX.log" )"
rm -- "$scriptDir/.tmp/"*.log || true

RED_COLOR="\e[41m\e[30m"
GREEN_COLOR="\e[42m\e[37m"
CYAN_COLOR="\e[46m\e[30m"
RESET_COLOR="\e[0m"

die() {
	>&2 echo "$1"
	exit 1
}

cd "$workspace"

>&2 echo "## Test whithout @init"
! $cmd "test without init" true @success || false

>&2 echo "## Test @report without init"
! $cmd @report || false

>&2 echo "## Test suite @init"
eval $( $cmd @init )

>&2 echo "## Test suite @success"
$cmd "test succes ok" true @success
! $cmd "test success ko" false @success || false

>&2 echo "## Test suite @fail"
$cmd "test fail ok" false @fail
! $cmd "test fail ko" true @fail || false

>&2 echo "## Test suite @rc="
$cmd "test rc=0 ok" true @rc=0
! $cmd "test rc=0 ko" false @rc=0 || false

$cmd "test rc=1 ok" false @rc=1
! $cmd "test rc=1 ko" true @rc=1 || false

>&2 echo "## Test suite @stdout="
$cmd "test empty stdout ok" true @stdout=
! $cmd "test empty stdout ko" echo "bar" @stdout= || false

$cmd "test stdout= ok" echo foo "bar baz" @stdout="foobar baz"
! $cmd "test stdout= ko" echo foo "bar baz" @stdout="foobarbaz" || false

>&2 echo "## Test suite @stdout~"
$cmd "test stdout~ ok" echo foo "bar baz" @stdout~"bar"
! $cmd "test stdout~ ko" echo foo "bar baz" @stdout~"pif" || false

>&2 echo "## Test suite @stderr="
$cmd "test empty stderr" true @stderr=
! $cmd "test empty stderr" notExist @stderr= || false

$cmd "test stderr= ok" notExist @stderr="notExist"
! $cmd "test stderr= ko" notExist @stderr="notExist" || false

>&2 echo "## Test suite @stderr~"
$cmd "test stderr~ ok" notExist @stderr~"notExist"
! $cmd "test stderr~ ko" notExist @stderr~"foo" || false


>&2 echo "## Test @report"
$cmd @report

# Should not report a second time
! $cmd @report || false

