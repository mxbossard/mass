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
mass init project p1
mass init project p2
mass init project p3

# Init some images
mass init image p1/i11
mass init image p1/i12
mass init image p1/i13

mass init image p2/i21
mass init image p2/i22
mass init image p2/i23

mass init image p3/i31
mass init image p3/i32
mass init image p3/i33

# Init env configs

cat <<EOF > envs/dev/config.yaml
labels:
  lkey1: ldev
tags:
  tkey1: tdev
environment:
  ekey1: edev
  ctx: env
EOF

# Init project configs

cat <<EOF > p1/config.yaml
labels:
  lkey2: lproject1
tags:
  tkey2: tproject1
environment:
  ekey2: eproject2
  ctx: project
EOF

cat <<EOF > p1/config.yaml
labels:
  lkey2: lproject2
tags:
  tkey2: tproject2
environment:
  ekey2: eproject2
  ctx: project
EOF

# Init image configs

cat <<EOF > p1/i11/config.yaml
labels:
  lkey3: limage11
tags:
  tkey3: timage11
environment:
  ekey3: eimage11
  ctx: image
EOF

cat <<EOF > p2/i21/config.yaml
labels:
  lkey3: limage21
tags:
  tkey3: timage21
environment:
  ekey3: eimage21
  ctx: image
EOF

cat <<EOF > p3/i31/config.yaml
labels:
  lkey3: limage31
tags:
  tkey3: timage31
environment:
  ekey3: eimage31
  ctx: image
EOF


# Display config for env
mass config dev
mass config e/stage


