#! /bin/bash
set -e
scriptDir=$( dirname $( readlink -f $0 ) )

workspaceDir="/tmp/massWorkspace"

echo "##### Building mass binary ..."
export GOBIN="$scriptDir/bin"
cd $scriptDir/src/mby.fr/mass
go install

rm -rf -- "$workspaceDir"

massCmd="$GOBIN/mass $@"
ls -lh "$GOBIN/mass"

############## TEST FRAMEWORK ############
STOP_ON_FAILURE=false
LOGS_ONLY_ON_FAILURE=true

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

report() {
	if [ -f "$reportFile" ]; then
		>&2 echo -e "${RED_COLOR}FAILURE${RESET_COLOR}"
		if [ "true" = "$STOP_ON_FAILURE" ]; then
			>&2 echo "Stop after first failure."
		else
			>&2 echo "Encountered $( cat "$reportFile" | wc -l ) failure(s) :"
			>&2 cat "$reportFile"
			>&2 echo -e "${RED_COLOR}FAILURE${RESET_COLOR}"
		fi
		exit 1
	else
		>&2 echo -e "${GREEN_COLOR}SUCCESS${RESET_COLOR}"
		exit 0
	fi
}

shouldSucceed() {
	title="$1"
	test -n "$title" || die "$0 must be passed a title !"

	>&2 echo -e "${CYAN_COLOR}## Test [$title] should succeed ...${RESET_COLOR}"
	if [ "true" = "$LOGS_ONLY_ON_FAILURE" ]; then
		logFile="$( mktemp "$scriptDir/.tmp/XXXXXX.log" )"
		cat > "$logFile"
	else
		cat
	fi
	
	#rc="${PIPESTATUS[0]}"
	rc="$?"
	if [ 0 -lt "$rc" ] ;then
		if [ "true" = "$LOGS_ONLY_ON_FAILURE" ]; then
			>&2 cat "$logFile"
			>&2 echo
		fi
		echo -e "${RED_COLOR} > Test [$title] should succeed !${RESET_COLOR}" | >&2 tee -a "$reportFile"

		test "false" = "$STOP_ON_FAILURE" || report
	fi

	test -f "$logFile" && rm -- "$logFile" || true
}

shouldFail() {
	title="$1"
	expectedRc="$2"
	test -n "$title" || die "$0 must be passed a title !"
	test -n "$expectedRc" || die "$0 must be passed an expected RC !"
	test 0 -lt "$expectedRc" || die "$0 expected RC must be > 0 !"

	>&2 echo -e "${CYAN_COLOR}## Test [$title] should fail ...${RESET_COLOR}"
	if [ "true" = "$LOGS_ONLY_ON_FAILURE" ]; then
		logFile="$( mktemp "$scriptDir/.tmp/XXXXXX.log" )"
		cat > "$logFile"
	else
		cat
	fi

	#rc="${PIPESTATUS[0]}"
	rc="$?"
	if [ 0 -eq "$rc" ] ;then
		if [ "true" = "$LOGS_ONLY_ON_FAILURE" ]; then
			>&2 cat "$logFile"
			>&2 echo
		fi
		echo -e "${RED_COLOR} > Test [$title] should fail !${RESET_COLOR}" | >&2 tee -a "$reportFile"
		if [ "$expectedRc" -ne "$rc" ] ;then
			echo -e "${RED_COLOR} > Test [$title] expect RC=$expectedRc but received RC=$rc !${RESET_COLOR}" | >&2 tee -a "$reportFile"

			test "false" = "$STOP_ON_FAILURE" || report
		fi
		test "false" = "$STOP_ON_FAILURE" || report
	fi


	test -f "$logFile" && rm -- "$logFile" || true
}

################ TESTS #################

# Init a workspace
echo "##### Testing mass init ..."
$massCmd init workspace $workspaceDir | shouldSucceed "init workspace"

cd /tmp
$massCmd init project p1 | shouldFail "init resource outside project" 1

cd $workspaceDir

# Init some projects
$massCmd init project p1 p2 p3 p444444444444 | shouldSucceed "init Resources"

# Init some images
$massCmd init image p1/i11 p1/i12 p1/i13 p2/i21 p2/i22 p2/i23 p3/i31 p3/i32 p3/i33 p444444444444/i444444444444 | shouldSucceed "init images"

configEnvs() {
	for name in "$@"; do	
		cat <<EOF > envs/$name/config.yaml
labels:
  lkey1: l$name
tags:
  tkey1: t$name
environment:
  ekey1: e$name
  ctx: env
EOF
	done
}

configProjects() {
	for name in "$@"; do	
		cat <<EOF > $name/config.yaml
labels:
  lkey2: l$name
tags:
  tkey2: t$name
environment:
  ekey2: e$name
  ctx: project
EOF

	cat <<EOF > $name/compose.yaml
services:
EOF
	done
}

configImages() {
	for name in "$@"; do
		projectName=$( echo $name | cut -d'/' -f1 )
		imageName=$( echo $name | cut -d'/' -f2 )
		imageDir="$projectName/img-$imageName"
		cat <<EOF > "$imageDir/config.yaml"
labels:
  lkey3: l$name
tags:
  tkey3: t$name
environment:
  ekey3: e$name
  ctx: image
entrypoint: ["/bin/sh", "-c"]
commandArgs:
  - echo display args configured in config.yaml file
EOF

	cat <<EOF > "$imageDir/Dockerfile"
FROM alpine
RUN echo foo
#RUN echo bar
#RUN echo baz
ENTRYPOINT ["/bin/sh"]
CMD ["-c", "echo hello world !"]
EOF

	cat <<EOF >> "$projectName/compose.yaml"
  $imageName:
    image: $name:0.1.0-dev
    entrypoint: ["/bin/sh", "-c"]
    command:
      - sleep inf
EOF
	done
}

echo "##### Populating test workspace ..."
# Init env configs
configEnvs dev

# Init project configs
configProjects p1 p2 p3 p444444444444

# Init image configs
configImages p1/i11 p1/i12 p1/i13 p2/i21 p2/i22 p2/i23 p3/i31 p3/i32 p3/i33 p444444444444/i444444444444

#tree -Ca $workspaceDir

# Display config for env
echo "##### Testing mass config ..."
$massCmd config e/dev | shouldSucceed "config env"
$massCmd config p/p1 i/p1/i11 | shouldSucceed "config resources 1"
$massCmd config -e stage p/p1 i/p1/i11 | shouldSucceed "config resources 2"
$massCmd config p,i p1 p1/i11 notExist | shouldFail "config resources 3" 4

echo "##### Testing mass build ..."
$massCmd build e/dev i/p1/i11 2>&1 | shouldFail "build not buildable" 2
$massCmd build i/p3/i31 2>&1 | shouldSucceed "build image"
$massCmd build p/p3 2>&1 | shouldSucceed "build project"
$massCmd build --no-cache p/p1 p/p2 2>&1 | shouldSucceed "build without cache"
$massCmd build p444444444444/i444444444444 2>&1 | shouldSucceed "build image 2"

#echo "##### p2/compose.yaml :"
#cat p2/compose.yaml

echo "##### Testing mass down ..."
$massCmd down i/p3/i31 p/p2 | shouldSucceed "down resources"

echo "##### Testing mass up ..."
$massCmd up i/p3/i31 | shouldSucceed "up image"

$massCmd up p/p2 | shouldSucceed "up project"

echo $errorCount
report
