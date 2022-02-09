#!/bin/bash

echo "git pull"
git pull

echo "soda migrate"
soda migrate

echo "go build"
go build -o sensor cmd/web/*.go && ./sensor -dbname=XXX -dbuser=XXX -dbpassword=XXX

echo "supervisor stop"
sudo supervisorctl stop sensor
echo "supervisor start"
sudo supervisorctl start sensor

echo "caddy stop"
sudo service caddy stop
sleep 2s
echo "caddy start"
sudo service caddy start