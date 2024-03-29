#!/bin/sh
# ===========================================================================
# File: build.sh
# Description: usage: ./prefix.sh
# ===========================================================================
cd "$(dirname "$0")"
cd ../../
BASEPATH="$(pwd)"


# =========================================================================
# 更新 go-plugin
# =========================================================================
echo "更新 go-plugin"
cd "${BASEPATH}/"
if [ ! -d "./go-plugin" ]; then
   git clone git@gitlab.eolink.com:apinto/go-plugin.git
fi
cd "./go-plugin" && git pull


# =========================================================================
# 更新 eosc
# =========================================================================
echo "更新 eosc"
cd "${BASEPATH}/"
if [ ! -d "./eosc" ]; then
   git clone git@gitlab.eolink.com:goku/eosc.git
fi
cd "./eosc" && git pull

# =========================================================================
# 更新 business
# =========================================================================

echo "更新 business"
cd "${BASEPATH}/"
if [ ! -d "./apinto-business" ]; then

   git clone git@gitlab.eolink.com:apinto/business.git
   mv business apinto-business
fi
cd "./apinto-business" && git pull