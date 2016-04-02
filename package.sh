#!/bin/bash
VERSION=`awk 'BEGIN { FS="\"" } /version/ { print $2 }'`
fpm -s dir -t deb -n sesmail -v $VERSION --prefix /usr/local/bin sesmail
