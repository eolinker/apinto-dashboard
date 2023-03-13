#!/bin/sh
# ===========================================================================
# File: build.sh
# Description: usage: ./build.sh [outdir]
# ===========================================================================

# exit when any command fails
set -e

cd "$(dirname "$0")/../"
. ./scripts/init.sh

OUTPUT_DIR=$(mkdir_output "$1")
OUTPUT_BINARY=$OUTPUT_DIR/apserver

GO_VERSION=`go version | { read _ _ v _; echo ${v#go}; }`

if [ "$(version ${GO_VERSION})" -lt "$(version 1.18)" ];
then
   echo "${RED}Precheck failed.${NC} Require go version >= 1.19. Current version ${GO_VERSION}."; exit 1;
fi
#
#if ! command -v npm > /dev/null
#then
#   echo "${RED}Precheck failed.${NC} npm is not installed."; exit 1;
#fi
#
## Step 1 - Build the frontend release version into the backend/server/dist folder
## Step 2 - Build the monolithic app by building backend release version together with the backend/server/dist (leveraing embed introduced in Golang 1.19).
#echo "Start building apinto enterprise monolithic ${VERSION}..."
#
#echo ""
#echo "Step 1 - building apinto enterprise frontend..."
#
#if command -v pnpm > /dev/null
#then
#    pnpm --dir ./frontend i && pnpm --dir ./frontend build
#else
#    npm --prefix ./frontend run build
#fi
echo "当前脚本未编译前端文件、留意留意"


echo "Completed building apinto enterprise frontend."
if [[ "$2" != "" ]];then
  VERSION=$2
fi
echo "${VERSION}"
echo "Step 2 - building apinto enterprise backend..."

flags="-X 'github.com/eolinker/apinto-dashboard/app/apserver/version.Version=${VERSION}'
-X 'github.com/eolinker/apinto-dashboard/app/apserver/version.goversion=$(go version)'
-X 'github.com/eolinker/apinto-dashboard/app/apserver/version.gitcommit=$(git rev-parse HEAD)'
-X 'github.com/eolinker/apinto-dashboard/app/apserver/version.buildtime=$(date -u +"%Y-%m-%dT%H:%M:%SZ")'
-X 'github.com/eolinker/apinto-dashboard/app/apserver/version.builduser=$(id -u -n)'"


# -ldflags="-w -s" means omit DWARF symbol table and the symbol table and debug information
go build --tags "release,mysql" -ldflags "-w -s $flags" -o ${OUTPUT_BINARY} ./app/apserver



mkdir -p apinto-dashboard

cp ./apinto-build/apserver ./apinto-dashboard
cp ./scripts/docker_run.sh ./apinto-dashboard
cp ./scripts/resource/* ./apinto-dashboard

echo "Completed building apinto enterprise backend."

echo ""
echo "Step 3 - printing version..."

tar -czvf apinto-dashboard.tar.gz apinto-dashboard

cp apinto-dashboard.tar.gz ./scripts/

rm -rf apinto-dashboard

echo "docker build start"
docker build -t "docker.eolinker.com/docker/apinto-dashboard:${VERSION}" -f ./scripts/Dockerfile ./scripts/
echo "docker build end"
#


