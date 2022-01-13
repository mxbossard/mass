#! /bin/bash -e
scriptDir=$( dirname $( readlink -f $0 ) )

workspaceDir="/tmp/massWorkspace"

cd $scriptDir/src/mby.fr/mass
go install

rm -rf -- "$workspaceDir"

# Init a workspace
mass init workspace $workspaceDir
cd $workspaceDir

# Init some projects
mass init project p1 p2 p3

# Init some images
mass init image p1/i11 p1/i12 p1/i13 p2/i21 p2/i22 p2/i23 p3/i31 p3/i32 p3/i33

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
EOF
	echo "FROM alpine" > $name/Dockerfile
	done
}

# Init env configs
configEnvs dev

# Init project configs
configProjects p1 p2 p3

# Init image configs
configImages p1/i11 p1/i12 p1/i13 p2/i21 p2/i22 p2/i23 p3/i31 p3/i32 p3/i33

tree -Ca $workspaceDir

# Display config for env
echo "##### Testing mass config ..."
mass config e/dev
mass config p/p1 i/p1/i11
mass config -e stage p/p1 i/p1/i11
mass config p,i p1 p1/i11 notExist || true

echo "##### Testing mass build ..."
mass build e/dev i/p1/i11 || true
mass build i/p3/i31
mass build p/p1 p/p2
