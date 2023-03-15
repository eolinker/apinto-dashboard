#!/bin/sh
# ===========================================================================
# File: build.sh
# Description: usage: ./build.sh [outdir]
# ===========================================================================

# exit when any command fails
set -e

cd "$(dirname "$0")/../"
. ./scripts/init.sh

OUTPUT_BINARY=apsescrver

GO_VERSION=`go version | { read _ _ v _; echo ${v#go}; }`

if [ "$(version ${GO_VERSION})" -lt "$(version 1.18)" ];
then
   echo "${RED}Precheck failed.${NC} Require go version >= 1.18. Current version ${GO_VERSION}."; exit 1;
fi

if ! command -v npm > /dev/null
then
   echo "${RED}Precheck failed.${NC} npm is not installed."; exit 1;
fi

# Step 1 - Build the frontend release version into the backend/server/dist folder
# Step 2 - Build the monolithic app by building backend release version together with the backend/server/dist (leveraing embed introduced in Golang 1.19).
echo "Start building apinto enterprise monolithic ${VERSION}..."

echo ""
echo "Step 1 - building apinto enterprise frontend..."

if command -v pnpm > /dev/null
then
    pnpm --dir ./frontend i && pnpm --dir ./frontend build
elif command -v yarn > /dev/null
then
    cd frontend && yarn install --registry http://172.18.65.55:4873/ --legacy-peer-deps && yarn build
    cd ../
else
    npm --prefix ./frontend run build

fi

echo "Completed building apinto enterprise frontend."

echo "${VERSION}"
echo "Step 2 - building apinto enterprise backend..."

GOVERSION=$(go version) goreleaser release --skip-publish --skip-validate --rm-dist

echo "Completed building apinto enterprise backend."

echo ""
echo "Step 3 - printing version..."


if [ -f "./dist/${OUTPUT_BINARY}_${GOOS}_${GOARCH}/${OUTPUT_BINARY}" ]; then
  ./dist/${OUTPUT_BINARY}_${GOOS}_${GOARCH}/${OUTPUT_BINARY} version
elif [ -f "./dist/${OUTPUT_BINARY}_${GOOS}_${GOARCH}_v1/${OUTPUT_BINARY}" ]; then
  "./dist/${OUTPUT_BINARY}_${GOOS}_${GOARCH}_v1/${OUTPUT_BINARY}" version
else
  exit 1
fi

echo ""
echo "${GREEN}Completed building apinto enterprise monolithic ${VERSION}."

#read -n 1 -p "Do you want to upload packages to Qiniu Cloud? (y/n)" char
#printf "\n"
#if [ "${char}" == 'y' ]; then
#  packages=$(ls ./dist/*.tar.gz)
#  for package in ${packages[@]}
#  do
#    name=${package##*/}
#    qshell rput --overwrite  eolinker-main private/apinto/apinto-dashboard/"${VERSION}"/"${name}" "${name}" && echo "upload OK"
#  done
#fi
