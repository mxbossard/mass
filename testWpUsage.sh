#! /bin/bash -e
scriptDir=$( dirname $( readlink -f $0 ) )

workspaceDir="/tmp/myMassWordpressWorkspace"

cd $scriptDir/src/mby.fr/mass
go install

rm -rf -- "$workspaceDir"

massCmd="mass"

# Init a workspace
$massCmd init workspace $workspaceDir
cd $workspaceDir

# Init git repo
git init .

# Init some projects
$massCmd init project wp

# Init some images
mass init image wp/wordpress wp/db

cat <<EOF > wp/wordpress/Dockerfile
FROM wordpress:6.0-fpm-alpine
RUN echo foo
RUN echo bar
RUN echo baz
EOF

cat <<EOF > wp/db/Dockerfile
FROM mariadb:10.7-focal
RUN echo pif
RUN echo paf
RUN echo pouf
EOF

cat <<EOF > wp/wordpress/config.yaml
labels:
  app: wordpress
tags:
  app: wordpress
environment:
  WORDPRESS_DB_HOST: db
  WORDPRESS_DB_USER: wordpress
  WORDPRESS_DB_PASSWORD: wordpress
  WORDPRESS_DB_NAME: wordpress
EOF

cat <<EOF > wp/db/config.yaml
labels:
  app: mariadb
tags:
  app: mariadb
environment:
  MYSQL_ROOT_PASSWORD: mypassword
  MYSQL_DATABASE: wordpress
  MYSQL_USER: wordpress
  MYSQL_PASSWORD: wordpress
EOF

tree -Ca $workspaceDir

# Display configs
$massCmd config i/wp/wordpress i/wp/db

$massCmd build --no-cache p/wp

$massCmd up p/wp
