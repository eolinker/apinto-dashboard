#!/usr/bin/env bash

echo $0
. $(dirname $0)/common.sh

#echo ${BasePath}
#echo ${CMD}
#echo ${Hour}

VERSION=$(genVersion $1)
OUTPATH="${BasePath}/out/apinto-dashboard-${VERSION}"
buildApp apinto-dashboard $VERSION

cp -a ${BasePath}/builds/account.yml  ${OUTPATH}/
cp -a ${BasePath}/builds/config.yml  ${OUTPATH}/
cp -a ${BasePath}/builds/Dockerfile  ${OUTPATH}/
cp -a ${BasePath}/builds/static  ${OUTPATH}/
cp -a ${BasePath}/builds/tpl  ${OUTPATH}/


cd ${ORGPATH}
