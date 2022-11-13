#!/usr/bin/env bash

. $(dirname $0)/common.sh
cd ${BasePath}/

#创建docker镜像
function dockerBuild(){
    APP=$1
    VERSION=$2
    cd "${BasePath}/out/${APP}-${VERSION}"
    pwd
    docker build -t eolinker/${APP}:${VERSION} .
    cd "${BasePath}"
}

VERSION=$(genVersion $1)
folder="${BasePath}/out/apinto-dashboard-${VERSION}"
if [[ ! -d "$folder" ]]
then
  mkdir "$folder"
  ${CMD}/build.sh $1
  if [[ "$?" != "0" ]]
  then
    exit 1
  fi
fi

dockerBuild apinto-dashboard $VERSION

cd ${ORGPATH}