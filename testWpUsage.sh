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
  component: frontend
tags:
  component: frontend
buildArgs:
EOF

cat <<EOF > wp/db/config.yaml
labels:
  component: db
tags:
  component: db
buildArgs:
EOF

cat <<EOF > wp/config.yaml
labels:
  app: wordpress
tags:
  app: wordpress
environment:
  MYSQL_ROOT_PASSWORD: mypassword
  WORDPRESS_DB_USER: wordpress
  WORDPRESS_DB_PASSWORD: wordpress
  WORDPRESS_DB_NAME: wordpress
EOF

cat <<EOF > wp/compose.yaml
services:
  db:
    #image: mariadb:10.6.4-focal
    image: wp/db:0.0.1
    command: '--default-authentication-plugin=mysql_native_password'
    volumes:
      - db_data:/var/lib/mysql
    restart: always
    environment:
      - MYSQL_ROOT_PASSWORD
      - MYSQL_DATABASE=\${WORDPRESS_DB_NAME}
      - MYSQL_USER=\${WORDPRESS_DB_USER}
      - MYSQL_PASSWORD=\${WORDPRESS_DB_PASSWORD}
    expose:
      - 3306
      - 33060
  wordpress:
    #image: wordpress:latest
    image: wp/wordpress:0.0.1
    ports:
      - 80:80
    restart: always
    environment:
      - WORDPRESS_DB_HOST=db
      - WORDPRESS_DB_USER
      - WORDPRESS_DB_PASSWORD
      - WORDPRESS_DB_NAME
volumes:
  db_data:
EOF

tree -Ca $workspaceDir

# Display configs
$massCmd config i/wp/wordpress i/wp/db

# Add DB initialization
initDbDir="wp/db/src/docker-entrypoint-initdb.d"
mkdir -p "$initDbDir"
cp ~/Documents/ede_backup_2020-01/ecrindes_phtest.sql $initDbDir

#$massCmd build --no-cache p/wp
$massCmd build p/wp

$massCmd up p/wp
