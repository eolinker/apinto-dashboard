#!/usr/bin/env bash

echo $0
. $(dirname $0)/common.sh

#echo ${BasePath}
#echo ${CMD}
#echo ${Hour}

VERSION=$(genVersion $1)
OUTPATH="${BasePath}/out/apinto-dashboard-${VERSION}"
buildApp apinto-dashboard $VERSION

cp -a ${BasePath}/builds/resources/*  ${OUTPATH}/

cd ${ORGPATH}
