#! /bin/bash -e
scriptDir=$( dirname $( readlink -f $0 ) )

workspaceDir="/tmp/massWorkspace"

cd $scriptDir/src/mby.fr/mass
go install

rm -rf -- "$workspaceDir"

massCmd="mass"

# Init a workspace
$massCmd init workspace $workspaceDir
cd $workspaceDir

# Init some projects
$massCmd init project p1 p2 p3 p444444444444

# Init some images
mass init image p1/i11 p1/i12 p1/i13 p2/i21 p2/i22 p2/i23 p3/i31 p3/i32 p3/i33 p444444444444/i444444444444

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
		cat <<EOF > $name/config.yaml
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

	cat <<EOF > $name/Dockerfile
FROM alpine
RUN echo foo
#RUN echo bar
#RUN echo baz
ENTRYPOINT ["/bin/sh"]
CMD ["-c", "echo hello world !"]
EOF

	cat <<EOF >> $projectName/compose.yaml
  $imageName:
    image: $name:0.1.0-dev
    entrypoint: ["/bin/sh", "-c"]
    command:
      - sleep inf
EOF
	done
}

# Init env configs
configEnvs dev

# Init project configs
configProjects p1 p2 p3 p444444444444

# Init image configs
configImages p1/i11 p1/i12 p1/i13 p2/i21 p2/i22 p2/i23 p3/i31 p3/i32 p3/i33 p444444444444/i444444444444

tree -Ca $workspaceDir

# Display config for env
echo "##### Testing mass config ..."
$massCmd config e/dev
$massCmd config p/p1 i/p1/i11
$massCmd config -e stage p/p1 i/p1/i11
$massCmd config p,i p1 p1/i11 notExist || true

echo "##### Testing mass build ..."
$massCmd build e/dev i/p1/i11 || true
$massCmd build i/p3/i31
$massCmd build p/p3
$massCmd build --no-cache p/p1 p/p2
$massCmd build p444444444444/i444444444444

echo "##### p2/compose.yaml :"
cat p2/compose.yaml

echo "##### Testing mass down ..."
mass down i/p3/i31 p/p2

echo "##### Testing mass up ..."
mass up i/p3/i31

mass up p/p2

echo
echo SUCCESS
