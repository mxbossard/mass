#! /bin/bash
set -e -o pipefail
scriptDir=$( dirname $( readlink -f $0 ) )

workspaceDir="/tmp/cmdtWorkspace"

>&2 echo "##### Building cmdtest binary ..."
export GOBIN="$scriptDir/bin"
cd "$scriptDir"
go install
cd - > /dev/null

rm -rf -- "$workspaceDir"

cmd="TO_REPLACE"
cmdt="$GOBIN/cmdtest"
ls -lh "$cmdt"

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

#cd "$GOBIN"
#test -e cmdt || ln -s cmduest cmdt
#export PATH="$PATH:."
#cmd="cmdt"

### NOTES
# - Il est facile de sortir un test de la bonne test suite, et ce test ne sera jamais report !
# => Should @report report all opened tests suites by default ?

>&2 echo
>&2 echo "## Test cmdt basic assertions should passed"
$cmdt @init=should_succeed @stopOnFailure=false

$cmdt @test=should_succeed/ true
$cmdt @test=should_succeed/ true @success
$cmdt @test=should_succeed/ false @fail
$cmdt @test=should_succeed/ true @exit=0
$cmdt @test=should_succeed/ false @exit=1

$cmdt @test=should_succeed/ echo foo bar @stdout~foo @stderr=
$cmdt @test=should_succeed/ echo foo bar @stdout~bar @stderr=
$cmdt @test=should_succeed/ echo foo bar @stdout~foo @stdout~bar @stderr=
$cmdt @test=should_succeed/ echo foo bar @stdout="foo bar\n" @stderr=

$cmdt @test=should_succeed/ sh -c ">&2 echo foo bar" @stderr~foo @stdout=
$cmdt @test=should_succeed/ sh -c ">&2 echo foo bar" @stderr~bar @stdout=
$cmdt @test=should_succeed/ sh -c ">&2 echo foo bar" @stderr~foo @stderr~bar @stdout=
$cmdt @test=should_succeed/ sh -c ">&2 echo foo bar" @stderr="foo bar\n" @stdout=

>&2 echo
>&2 echo "## Test cmdt basic assertions should failed"
$cmdt @init=should_fail @stopOnFailure=false

$cmdt @test=should_fail/ false
$cmdt @test=should_fail/ true @fail
$cmdt @test=should_fail/ false @success
$cmdt @test=should_fail/ true @exit=1
$cmdt @test=should_fail/ false @exit=0

$cmdt @test=should_fail/ echo foo bar @stdout=
$cmdt @test=should_fail/ echo foo bar @stdout=foo
$cmdt @test=should_fail/ echo foo bar @stdout=foo bar
$cmdt @test=should_fail/ echo foo bar @stdout~baz
$cmdt @test=should_fail/ echo foo bar @stdout~foo @stdout~baz
$cmdt @test=should_fail/ echo foo bar @stderr~foo

$cmdt @test=should_fail/ sh -c ">&2 echo foo bar" @stderr=
$cmdt @test=should_fail/ sh -c ">&2 echo foo bar" @stderr=foo
$cmdt @test=should_fail/ sh -c ">&2 echo foo bar" @stderr=foo bar
$cmdt @test=should_fail/ sh -c ">&2 echo foo bar" @stderr~baz
$cmdt @test=should_fail/ sh -c ">&2 echo foo bar" @stderr~foo @stderr~baz
$cmdt @test=should_fail/ sh -c ">&2 echo foo bar" @stdout~foo

$cmdt @report=should_succeed
! $cmdt @report=should_fail 2>&1 | grep "0 success" || die "should_fail test suite should have no success"


>&2 echo
>&2 echo "## Test @report without test"
expectedNothingToReportStderr="you must perform some test prior to report"
$cmdt @test=meta1/ @fail @stderr~"$expectedNothingToReportStderr" -- $cmdt @report=foo
$cmdt @test=meta1/ @fail @stderr~"$expectedNothingToReportStderr" -- $cmdt @report=foo

>&2 echo
>&2 echo "## Met1a test context not shared without token"
$cmdt @test=meta1/ @stderr~"PASSED" @stderr~"#01" -- $cmdt true
$cmdt @test=meta1/ @stderr~"PASSED" @stderr~"#01" -- $cmdt true
$cmdt @test=meta1/ @fail -- $cmdt @report

>&2 echo
>&2 echo "## Test printed token"
tk0=$( $cmdt @init @printToken )
>&2 echo "token: $tk0"
$cmdt @test=meta2/ @stderr~"PASSED" @stderr~"#01" -- $cmdt true @token=$tk0
$cmdt @test=meta2/ @stderr~"PASSED" @stderr~"#02" -- $cmdt true @token=$tk0
$cmdt @test=meta2/ @fail -- $cmdt @report
$cmdt @test=meta2/ @stderr~"Successfuly ran" -- $cmdt @report @token=$tk0
$cmdt @report 2>&1 | grep -v "Failures"

>&2 echo
>&2 echo "## Test exported token"
eval $( $cmdt @init @exportToken )
>&2 echo "token: $__CMDT_TOKEN"
$cmdt @test=meta3/ @stderr~"PASSED" @stderr~"#01" -- $cmdt true
$cmdt @test=meta3/ @stderr~"PASSED" @stderr~"#02" -- $cmdt true
$cmdt @test=meta3/ @stderr~"Successfuly ran" -- $cmdt @report=main
$cmdt @test=meta3/ @fail @stderr~"$expectedNothingToReportStderr" -- $cmdt @report @token=$tk0

$cmdt @test=meta4/ @stderr~"PASSED" @stderr~"#01" -- $cmdt @test=sub4/ true
$cmdt @test=meta4/ @stderr~"PASSED" @stderr~"#02" -- $cmdt @test=sub4/ true
$cmdt @test=meta4/ @stderr~"Successfuly ran" -- $cmdt @report=sub4
$cmdt @test=meta4/ @fail @stderr~"$expectedNothingToReportStderr" -- $cmdt @report=sub4 @token=$tk0
$cmdt @report 2>&1 | grep -v "Failures"
export -n __CMDT_TOKEN

>&2 echo
>&2 echo "## Test usage"
$cmdt @test=meta/ @fail @stderr~"usage:" -- $cmdt

>&2 echo
>&2 echo "## Test cmdt basic assertions"
$cmdt @test=meta/ @stderr~"PASSED" -- $cmdt true
$cmdt @test=meta/ @stderr~"FAILED" -- $cmdt false
$cmdt @test=meta/ @stderr~"PASSED" -- $cmdt true @success
$cmdt @test=meta/ @stderr~"FAILED" -- $cmdt false @success
$cmdt @test=meta/ @stderr~"FAILED" -- $cmdt true @fail
$cmdt @test=meta/ @stderr~"PASSED" -- $cmdt false @fail
$cmdt @test=meta/ @stderr~"PASSED" -- $cmdt true @exit=0
$cmdt @test=meta/ @stderr~"FAILED" -- $cmdt false @exit=0
$cmdt @test=meta/ @stderr~"FAILED" -- $cmdt true @exit=1
$cmdt @test=meta/ @stderr~"PASSED" -- $cmdt false @exit=1

$cmdt @test=meta/ @stderr~"PASSED" -- $cmdt echo foo bar @stdout~foo @stderr=
$cmdt @test=meta/ @stderr~"FAILED" -- $cmdt echo foo bar @stdout~baz @stderr=
$cmdt @test=meta/ @stderr~"PASSED" -- $cmdt echo foo bar @stdout~foo @stdout~bar @stderr=
$cmdt @test=meta/ @stderr~"FAILED" -- $cmdt echo foo bar @stdout~baz @stdout~bar @stderr=
$cmdt @test=meta/ @stderr~"FAILED" -- $cmdt echo foo bar @stdout~foo @stdout~baz @stderr=

$cmdt @test=meta/ @stderr~"PASSED" -- $cmdt sh -c ">&2 echo foo bar" @stderr~foo @stdout=
$cmdt @test=meta/ @stderr~"FAILED" -- $cmdt sh -c ">&2 echo foo bar" @stderr~baz @stdout=
$cmdt @test=meta/ @stderr~"PASSED" -- $cmdt sh -c ">&2 echo foo bar" @stderr~foo @stderr~bar @stdout=
$cmdt @test=meta/ @stderr~"FAILED" -- $cmdt sh -c ">&2 echo foo bar" @stderr~baz @stderr~bar @stdout=
$cmdt @test=meta/ @stderr~"FAILED" -- $cmdt sh -c ">&2 echo foo bar" @stderr~foo @stderr~baz @stdout=


## Forge a token for remaining tests
eval $( $cmdt @test=meta/ @keepStdout -- $cmdt @init=t1 @exportToken)

>&2 echo
>&2 echo "## Test assertions outputs"
# Init the context used in t1 test suite
$cmdt @test=meta/ @stderr~"#01..." @stderr~"PASSED" -- $cmdt true @test=t1/
$cmdt @test=meta/ @stderr~"#02..." @stderr~"PASSED" -- $cmdt true @test=t1/
$cmdt @test=meta/ @stderr~"#03..." @stderr~"FAILED" -- $cmdt false @test=t1/
$cmdt @test=meta/ @stderr~"#04..." @stderr~"PASSED" -- $cmdt false @fail @test=t1/
$cmdt @test=meta/ @fail @stderr~"Failures in t1 test suite (3 success, 1 failures, 4 tests in" -- $cmdt @report=t1


>&2 echo
>&2 echo "## Test namings"
$cmdt @test=naming/ @stderr~"Test [main]/name1 #01..." @stderr~"PASSED" -- $cmdt true @test=name1
$cmdt @test=naming/ @stderr~"Test [main]/name2 #02..." @stderr~"PASSED" -- $cmdt true @test=name2
$cmdt @test=naming/ @stderr~"Test [main]/" @stderr~"true" @stderr~"#03..." @stderr~"PASSED" -- $cmdt true
$cmdt @test=naming/ @stderr~"Test [main]/" @stderr~"true" @stderr~"#04..." @stderr~"PASSED" -- $cmdt true
$cmdt @test=naming/ @stderr~"Test [suite1]/name1 #01..." @stderr~"PASSED" -- $cmdt true @test=suite1/name1
$cmdt @test=naming/ @stderr~"Test [suite1]/name2 #02..." @stderr~"PASSED" -- $cmdt true @test=suite1/name2
$cmdt @test=naming/ @stderr~"Test [suite2]/" @stderr~"#01..." @stderr~"PASSED" -- $cmdt true @test=suite2/
$cmdt @test=naming/ @stderr~"Test [suite2]/" @stderr~"#02..." @stderr~"PASSED" -- $cmdt true @test=suite2/
$cmdt @test=naming/ @stderr~"Successfuly ran suite1 test suite" -- $cmdt  @report=suite1
$cmdt @test=naming/ @stderr~"Successfuly ran suite2 test suite" -- $cmdt @report=suite2
$cmdt @test=naming/ @stderr~"Successfuly ran main test suite" -- $cmdt @report=main


>&2 echo
>&2 echo "## Test rules missusage"
$cmdt @test=failing_rule_missusage/ @fail -- $cmdt @init true
$cmdt @test=failing_rule_missusage/ @fail -- $cmdt @init @test
$cmdt @test=failing_rule_missusage/ @fail -- $cmdt @init @report
$cmdt @test=failing_rule_missusage/ @fail -- $cmdt @test @report true
$cmdt @test=failing_rule_missusage/ @fail -- $cmdt @init @test true
$cmdt @test=failing_rule_missusage/ @fail -- $cmdt @init @test true
$cmdt @test=failing_rule_missusage/ @fail -- $cmdt @init @fail
$cmdt @test=failing_rule_missusage/ @fail -- $cmdt @init @success
$cmdt @test=failing_rule_missusage/ @fail -- $cmdt @init @exit=0
$cmdt @test=failing_rule_missusage/ @fail -- $cmdt true @fail @success
$cmdt @test=failing_rule_missusage/ @fail -- $cmdt true @fail @exit=0
$cmdt @test=failing_rule_missusage/ @fail -- $cmdt true @success @exit=0
$cmdt @test=failing_rule_missusage/ @fail -- $cmdt @report @fail
$cmdt @test=failing_rule_missusage/ @fail -- $cmdt @report @success
$cmdt @test=failing_rule_missusage/ @fail -- $cmdt @report @exit=0
$cmdt @test=failing_rule_missusage/ @fail @stderr~donotexist -- $cmdt true @donotexist
$cmdt @test=failing_rule_missusage/ @fail @stderr~donotexist -- $cmdt true @donotexist


>&2 echo
>&2 echo "## Test assertion"
$cmdt @test=assertion/ @stderr~PASSED -- $cmdt true
$cmdt @test=assertion/ @stderr~PASSED -- $cmdt true @success
$cmdt @test=assertion/ @stderr~FAILED -- $cmdt true @fail

$cmdt @test=assertion/ @stderr~FAILED -- $cmdt false
$cmdt @test=assertion/ @stderr~FAILED -- $cmdt false @success
$cmdt @test=assertion/ @stderr~PASSED -- $cmdt false @fail

$cmdt @test=assertion/ @stderr~PASSED -- $cmdt true @exit=0
$cmdt @test=assertion/ @stderr~FAILED -- $cmdt true @exit=1

$cmdt @test=assertion/ @stderr~FAILED -- $cmdt false @exit=0
$cmdt @test=assertion/ @stderr~PASSED -- $cmdt false @exit=1

$cmdt @test=assertion/ @stderr~PASSED -- $cmdt echo foo bar @stdout="foo bar\n"
$cmdt @test=assertion/ @stderr~PASSED -- $cmdt echo foo bar @stdout~foo
$cmdt @test=assertion/ @stderr~PASSED -- $cmdt echo foo bar @stdout~bar
$cmdt @test=assertion/ @stderr~PASSED -- $cmdt echo foo bar @stderr=
$cmdt @test=assertion/ @stderr~FAILED -- $cmdt echo foo bar @stdout="foo bar"
$cmdt @test=assertion/ @stderr~FAILED -- $cmdt echo foo bar @stdout=
$cmdt @test=assertion/ @fail -- $cmdt echo foo bar @stdout~
$cmdt @test=assertion/ @stderr~FAILED -- $cmdt echo foo bar @stdout~baz
$cmdt @test=assertion/ @fail -- $cmdt echo foo bar @stderr~

$cmdt @test=assertion/ @stderr~PASSED -- $cmdt sh -c ">&2 echo foo bar" @stderr="foo bar\n"
$cmdt @test=assertion/ @stderr~PASSED -- $cmdt sh -c ">&2 echo foo bar" @stderr~foo
$cmdt @test=assertion/ @stderr~PASSED -- $cmdt sh -c ">&2 echo foo bar" @stderr~bar
$cmdt @test=assertion/ @stderr~PASSED -- $cmdt sh -c ">&2 echo foo bar" @stdout=
$cmdt @test=assertion/ @stderr~FAILED -- $cmdt sh -c ">&2 echo foo bar" @stderr="foo bar"
$cmdt @test=assertion/ @stderr~FAILED -- $cmdt sh -c ">&2 echo foo bar" @stderr=
$cmdt @test=assertion/ @fail -- $cmdt sh -c ">&2 echo foo bar" @stderr~
$cmdt @test=assertion/ @stderr~FAILED -- $cmdt sh -c ">&2 echo foo bar" @stderr~baz
$cmdt @test=assertion/ @fail -- $cmdt sh -c ">&2 echo foo bar" @stdout~

# TODO: @exists, @cmd

$cmdt @test=assertion/ @fail -- $cmdt @report=main 


>&2 echo
>&2 echo "## Test config"
$cmdt @test=test_config/ @stderr~Ignore -- $cmdt true @ignore
$cmdt @test=test_config/ @stderr~Ignore @stderr!~FAILED @stderr!~PASSED -- $cmdt true @ignore

>&2 echo
>&2 echo "## Test suite config"

$cmdt @report= ; >&2 echo SUCCESS ; exit 0

>&2 echo "## Test sleep"
$cmd sleep 0.3

>&2 echo "## Test stdin"
echo foo | $cmd "@test=stdin" cat @stdout="foo\n"

>&2 echo "## Test exported var"
export foo=bar
$cmd "@test=exported var" sh -c "export" @stdout~"foo='bar'\n"

#>&2 echo "## Test alias"
#alias foo=echo
#$cmd "@test=alias" sh -c "foo bar" @stdout~"bar\n"

>&2 echo "## Test @timeout"
! $cmd sleep 1 @timeout=50ms || false


>&2 echo "## Test @cmd"
$cmd false @cmd="true" @fail
! $cmd false @cmd="false" @fail || false
! $cmd true @cmd="false" || false
! $cmd true @cmd= || false
! $cmd true @cmd=notExist || false

>&2 echo "## Test @exists"
! $cmd echo "test non existance" @exists=/tmp/donotexistsmbd123 || false
$cmd touch /tmp/existsmbd123 @exists=/tmp/existsmbd123
! $cmd true @exists= || false

>&2 echo "## Change rule prefix"
$cmd @prefix=% %fail false
! $cmd @prefix=% %success false || false

>&2 echo "## Rules parsing stopper --"
$cmd @stdout="foo @success\n" -- echo foo @success

>&2 echo "## cmdt Meta Test"
# -- rules parsing stopper version
$cmd @keepOutputs @test=meta1/ @success -- $cmd @test=sub/ false @fail
$cmd @keepOutputs @test=meta1/ @fail -- $cmd @test=sub/ false @success
$cmd @keepOutputs @test=meta1/ @fail -- $cmd @test=sub/ false

# @prefix= version
$cmd @prefix=% %keepOutputs %test=meta2/ %success $cmd @test=sub/ false @fail
$cmd @prefix=% %keepOutputs %test=meta2/ %fail $cmd @test=sub/ false @success
$cmd @prefix=% %keepOutputs %test=meta2/ %fail $cmd @test=sub/ false

>&2 echo "## Test @report"
$cmd @report

# Should not report a second time
! $cmd @report || false

>&2 echo "## Interlaced tests"
eval $( $cmd @init="another" @keepStdout )
eval $( $cmd @init="another one" @keepOutputs )

$cmd echo not ignored @test="another/"
$cmd echo ignored @ignore @test="another/"

$cmd echo another test @test="another one/"
$cmd echo interlaced test @test="another/"

>&2 echo should have 2 success and 1 ignored
$cmd @report="another"

>&2 echo should have 1 success
$cmd @report="another one"

>&2 echo "## Mutually exclusive rules @init"
! $cmd @init @fail || false
! $cmd @init @test || false

>&2 echo "## Mutually exclusive rules @test"
! $cmd true @fail @success || false
! $cmd true @fail @success || false
! $cmd true @fail @exit=0 || false
! $cmd true @fail @suiteTimeout=1m || false

>&2 echo "## Mutually exclusive rules @report"
! $cmd @report @ignore || false
! $cmd @report @test || false
! $cmd @report @fail || false


>&2 echo SUCCESS
