#!/bin/sh
# ===========================================================================
# File: build.sh
# Description: usage: ./build.sh [outdir]
# ===========================================================================

# exit when any command fails
set -e

cd "$(dirname "$0")/../"
. ./scripts/init.sh

docker tag docker.eolinker.com/docker/apinto/apserver:${VERSION} docker.eolinker.com/docker/apinto/apserver:latest
docker push docker.eolinker.com/docker/apinto/apserver:${VERSION}
docker push docker.eolinker.com/docker/apinto/apserver:latest
