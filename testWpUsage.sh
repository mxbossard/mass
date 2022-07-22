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

cat <<EOF > wp/compose.yml
services:
  db:
    image: mariadb:10.6.4-focal
    command: '--default-authentication-plugin=mysql_native_password'
    volumes:
      - db_data:/var/lib/mysql
    restart: always
    environment:
      - MYSQL_ROOT_PASSWORD=somewordpress
      - MYSQL_DATABASE=wordpress
      - MYSQL_USER=wordpress
      - MYSQL_PASSWORD=wordpress
    expose:
      - 3306
      - 33060
  wordpress:
    image: wordpress:latest
    ports:
      - 80:80
    restart: always
    environment:
      - WORDPRESS_DB_HOST=db
      - WORDPRESS_DB_USER=wordpress
      - WORDPRESS_DB_PASSWORD=wordpress
      - WORDPRESS_DB_NAME=wordpress
volumes:
  db_data:
EOF

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

#$massCmd build --no-cache p/wp
$massCmd build p/wp

$massCmd up p/wp
