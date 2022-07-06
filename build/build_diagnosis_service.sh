#!/bin/sh
cd ..

echo create folder for build package ...
if [ ! -e package ]; then
  mkdir package
fi

echo build the code ...
cd service/diagnosis_service
#添加参数CGO_ENABLED=0，关闭CGO,是为了是编译后的程序可以在alpine中运行
CGO_ENABLED=0 go build
cd ../..

if [ ! -e package/service ]; then
  mkdir package/service
fi

if [ -e package/service/diag_service ]; then
  rm -rf package/service/diag_service
fi

mv service/diagnosis_service/diagnosis ./package/service/diag_service

echo diag_service package build over.
