#!/bin/bash
VERSION=`awk '/version :=/ { gsub("\"", "", $3); print $3 }' main.go`
fpm -s dir -t deb -n sesmail -v $VERSION --prefix /usr/local/bin sesmail
