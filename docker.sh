#!/bin/bash

docker run -d --rm \
  --name canonical \
	--mount type=bind,src=$PWD/docker-scripts/,dst=/docker-entrypoint-initdb.d/ \
	-e  MYSQL_USER=sql \
  -e  MYSQL_PASSWORD=password \
	-e  MYSQL_DATABASE=bookmanager \
	-p  3306:3306 \
  mysql/mysql-server:8.0.23-1.1.19
