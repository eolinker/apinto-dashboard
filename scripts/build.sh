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

BUILD_MODE=$2

if [[ "$3" != "" ]];then
  VERSION=$3
fi

GO_VERSION=`go version | { read _ _ v _; echo ${v#go}; }`

if [ "$(version ${GO_VERSION})" -lt "$(version 1.18)" ];
then
   echo "${RED}Precheck failed.${NC} Require go version >= 1.19. Current version ${GO_VERSION}."; exit 1;
fi


# Step 1 - Build the frontend release version into the backend/server/dist folder
# Step 2 - Build the monolithic app by building backend release version together with the backend/server/dist (leveraing embed introduced in Golang 1.19).
echo "Start building apinto dashboard monolithic ${VERSION}..."

echo ""
echo "Step 1 - building apinto dashboard frontend..."

if [[ "$BUILD_MODE" == "all" || ! -d "controller/dist" ]];then
  echo "begin frontend building..."
  if command -v pnpm > /dev/null
  then
      pnpm --dir ./frontend i && pnpm --dir ./frontend build
  elif command -v yarn > /dev/null
  then
      echo "cd frontend && yarn install --registry https://registry.npmmirror.com --legacy-peer-deps "
      cd frontend && yarn install --registry https://registry.npmmirror.com --legacy-peer-deps
      echo "yarn build"
      yarn build
      cd ../
  else
      npm --prefix ./frontend run build
  fi
else
  echo "skip frontend building..."
fi

echo "Completed building apinto dashboard frontend."

echo "${VERSION}"
echo "Step 2 - building apinto dashboard backend..."

flags="-X 'github.com/eolinker/apinto-dashboard/app/apserver/version.Version=${VERSION}'
-X 'github.com/eolinker/apinto-dashboard/app/apserver/version.goversion=$(go version)'
-X 'github.com/eolinker/apinto-dashboard/app/apserver/version.gitcommit=$(git rev-parse HEAD)'
-X 'github.com/eolinker/apinto-dashboard/app/apserver/version.buildtime=$(date -u +"%Y-%m-%dT%H:%M:%SZ")'
-X 'github.com/eolinker/apinto-dashboard/app/apserver/version.builduser=$(id -u -n)'"


# -ldflags="-w -s" means omit DWARF symbol table and the symbol table and debug information
CGO_ENABLED=0 go build --tags "release,mysql" -ldflags "-w -s $flags" -o ${OUTPUT_BINARY} ./app/apserver

mkdir -p apserver_${VERSION}
#cp ./scripts/resource/config.yml.tpl ${OUTPUT_DIR}/config.yml

cp ./scripts/resource/config.yml.tpl ./apserver_${VERSION}
cp ${OUTPUT_BINARY} ./apserver_${VERSION}
cp ./scripts/resource/install.sh ./apserver_${VERSION}
cp ./scripts/resource/run.sh ./apserver_${VERSION}

echo "Completed building apinto dashboard backend."

echo ""
echo "Step 3 - printing version..."

tar -czvf apserver_${VERSION}_linux_amd64.tar.gz apserver_${VERSION}

rm -rf apserver_${VERSION}

cp apserver_${VERSION}_linux_amd64.tar.gz ${OUTPUT_DIR}

rm -rf apserver_${VERSION}_linux_amd64.tar.gz

echo "apserver_${VERSION}_linux_amd64.tar.gz 完成"

echo ""
echo "${GREEN}Completed building apinto dashboard monolithic ${VERSION} at ${OUTPUT_BINARY}.${NC}"

