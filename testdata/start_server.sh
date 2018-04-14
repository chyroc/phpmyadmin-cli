#!/usr/bin/env bash

set -e

docker run --name mysql_db_server -e MYSQL_ROOT_PASSWORD=pass -d mysql

docker run --name myadmin -d --link mysql_db_server:db -p 8000:80 phpmyadmin/phpmyadmin