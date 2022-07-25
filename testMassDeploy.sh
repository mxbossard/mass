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
$massCmd init project p1 p2 p3

# Init some images
$massCmd init image p1/i11 p1/i12 p1/i13 p2/i21 p2/i22 p2/i23 p3/i31 p3/i32 p3/i33

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
	done
}

configImages() {
	for name in "$@"; do
		projactName=$( echo $name | cut -d'/' -f1 )
		imageName=$( echo $name | cut -d'/' -f2 )
		cat <<EOF > $name/config.yaml
labels:
  lkey3: l$name
tags:
  tkey3: t$name
environment:
  ekey3: e$name
  ctx: image
buildargs:
  barg3: b$name
runargs:
  - hello world from runargs 3
EOF
	cat <<EOF > $name/Dockerfile
FROM alpine
ARG barg3 notDefinedBuildArg
ENV envBarg3 \${barg3:-missingBuildArg}
ENV ekey2 NotDefinedEnvVar
RUN echo "foo"
RUN echo "\$barg3"
RUN echo "\${ekey2}"
RUN echo -e "#!/bin/sh\n" >> entrypoint.sh \
    echo -e 'echo "buildArg: \$barg3"\n' >> entrypoint.sh \
    echo -e 'echo "envBarg : \$envBarg3"\n' >> entrypoint.sh \
    echo -e 'echo "envParam: \$ekey2"\n' >> entrypoint.sh \
    echo -e 'echo "runArgs : \$@"\n' >> entrypoint.sh

ENTRYPOINT ["/bin/sh", "entrypoint.sh", "runArgFromDockerfile"]
CMD ["NotDefinedRunArg"]
EOF
	done
}

# Init env configs
configEnvs dev

# Init project configs
configProjects p1 p2 p3

# Init image configs
configImages p1/i11 p1/i12 p1/i13 p2/i21 p2/i22 p2/i23 p3/i31 p3/i32 p3/i33

#tree -Ca $workspaceDir

# Display config for env
echo "##### Testing mass config ..."
$massCmd config i/p2/i21

echo "##### Testing mass build -vv (DEBUG) ..."
$massCmd build i/p2/i21 -vv

echo "##### Testing mass up -vvv (TRACE) ..."
$massCmd up i/p2/i21 -vvv || true

echo "##### Testing mass test ..."
$massCmd test i/p2/i21 || true
$massCmd test p/p2 || true

echo "##### Testing mass down -v (INFO) ..."
$massCmd down i/p2/i21 -v
