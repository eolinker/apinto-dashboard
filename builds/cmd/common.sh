#!/usr/bin/env bash

ORGPATH=$(pwd) #原始目录

cd  $(dirname $0) # 当前位置跳到脚本位置
CMD=$(pwd) # 脚本所在位置
cd ../..
BasePath=$(pwd) ## 项目根目录


# 生成版本号
function genVersion(){

    if [[ "$1" = "" ]]
    then
       v=$(git rev-parse --short HEAD)
       time=$(date "+%Y%m%d%H")
       echo "$time-$v"
       exit 0
    fi
    echo "$1"
}

# 构建app
function buildApp(){
    APP=$1
    VERSION=$2
    OUTPATH="${BasePath}/out/${APP}-${VERSION}"
    rm -rf ${OUTPATH}
    mkdir -p ${OUTPATH}
    buildCMD="go build  -o ${OUTPATH}/$APP ${BasePath}/app/$APP"
    echo "build $APP:${buildCMD}"
    ${buildCMD}

    if [[ "$?" != "0" ]]
    then
        rm -rf $OUTPATH
        exit 1
    fi
    echo "$VERSION" > ${OUTPATH}/version

}
#打包app
function packageApp(){
    APP=$1
    VERSION=$2
    cd "${BasePath}/out/${APP}-${VERSION}"
    tar -zcf "${BasePath}/out/${APP}-${VERSION}.linux.x64.tar.gz" --xform 's#^#apinto-dashboard/#'  *
    cd "${BasePath}"
}

#创建docker镜像
function dockerBuild(){
    APP=$1
    VERSION=$2
    cd "${BasePath}/out/${APP}-${VERSION}"
    docker build -t eolinker/${APP}:${VERSION} .
    cd "${BasePath}"
}