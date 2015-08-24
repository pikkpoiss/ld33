#!/usr/bin/env bash

GITROOT=`git rev-parse --show-toplevel`

cd $GITROOT
iconutil -c icns assets/icons.iconset
cd -
