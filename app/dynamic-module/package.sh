#!/usr/bin/env bash

ORG_PATH=`pwd`
for file in `ls .`
do
  if [[ "$file" == "package.sh" ]];then
    continue
  fi
  cd $ORG_PATH/$file
  rm -f *.zip
  zip -r -o $file.zip *
done

cd $ORG_PATH