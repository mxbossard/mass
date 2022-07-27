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
FROM wordpress:6.0-php8.0-apache
RUN echo foo
RUN echo bar
RUN echo baz
EOF

cat <<EOF > wp/db/Dockerfile
FROM mariadb:10.7-focal
RUN echo pif
RUN echo paf
RUN echo pouf
ADD src/docker-entrypoint-initdb.d /docker-entrypoint-initdb.d
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
    volumes:
      - wp_site:/var/www/html:rw
    ports:
      - 8000:80
    restart: always
    environment:
      - WORDPRESS_DB_HOST=db
      - WORDPRESS_DB_USER
      - WORDPRESS_DB_PASSWORD
      - WORDPRESS_DB_NAME
      #- WORDPRESS_CONFIG_EXTRA=define( 'WP_CONTENT_URL', 'http://localhost/wp-content' ); define( 'WP_HOME', 'http://localhost' ); define( 'WP_SITEURL', 'http://localhost' );
      - WORDPRESS_CONFIG_EXTRA=define( 'WP_HOME', 'http://localhost:8000' ); define( 'WP_SITEURL', 'http://localhost:8000' );

  wait-db:
    image: mariadb:10.6.4-focal
    command: /bin/sh -c 'while ! mysql -h db -u \${WORDPRESS_DB_USER} --password=\${WORDPRESS_DB_PASSWORD}; do sleep 1; echo .; done'
    depends_on:
      - db

  cli:
    image: wordpress:cli-2.6
    command: |
      bash -c "
      while ! mysql -h db -u \${WORDPRESS_DB_USER} --password=\${WORDPRESS_DB_PASSWORD}; do sleep 1; echo .; done ; sleep 2 ;
      wp db repair ;
      wp theme delete --all --force ;
      wp db optimize ;
      wp core update-db ;
      wp theme install shuttle-clean ;
      wp theme activate shuttle-clean ;
      wp plugin install wp-optimize wordpress-seo disable-comments login-lockdown imsanity;
      "
    #wp search-replace old-site-url.co.uk new-site-url.co.uk ;
    working_dir: /var/www/html
    user: "33:33"
    volumes:
      - wp_site:/var/www/html:rw
    environment:
      - HOME=/tmp
      - WORDPRESS_DB_HOST=db
      - WORDPRESS_DB_USER
      - WORDPRESS_DB_PASSWORD
    depends_on: 
      - db
    

volumes:
  db_data:
  wp_site:
EOF

tree -Ca $workspaceDir

# Display configs
$massCmd config i/wp/wordpress i/wp/db

# Add DB initialization
initDbDir="wp/db/src/docker-entrypoint-initdb.d"
confDbDir="wp/db/src/conf.d"
mkdir -p "$initDbDir" "$confDbDir"
cp ~/Documents/ede_backup_2020-01/ecrindes_phtest.sql $initDbDir/10-init-ede-db.sql
echo "update wp_users set user_pass=md5('password') where user_login = 'nathalie';" > $initDbDir/20-change-admin-password.sql

#$massCmd build --no-cache p/wp
#$massCmd build p/wp

$massCmd up p/wp || true

echo "Will execute mass down in 10 seconds ..."
for k in $( seq 10 ); do
	echo -n .
	sleep 1
done
echo

$massCmd down p/wp

echo "Will execute mass down --volumes in 10 seconds ..."
for k in $( seq 10 ); do
	echo -n .
	sleep 1
done
echo

$massCmd down --volumes p/wp

echo
echo SUCCESS
