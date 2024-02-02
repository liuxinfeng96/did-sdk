#!/bin/bash

contractName=$1
targetARCH=$2
crypto=""

if [ "$(uname)" == "Linux" ];then
  crypto="-tags crypto"
fi

if  [[ ! -n $contractName ]] ;then
    echo "contractName is empty. use as: ./build.sh contractName."
    exit 1
fi

if  [[ ! -n $targetARCH ]] ;then
    targetARCH=amd64
fi

echo "[CMD] ./build.sh $contractName $targetARCH"

GOOS=linux GOARCH=$targetARCH go build $crypto -ldflags="-s -w" -o $contractName

echo "[OK] Compiled project to contract bin $contractName."

7z a $contractName $contractName -sdel > /dev/null

echo "[OK] Compressed contract bin to $contractName.7z."

echo "[OK] Completed!"

echo -e "[NOTE] The default ARCH is amd64, it needs to be the same with the vm-engine host machine's ARCH.
You can execute \"go tool dist list -json\" to get all ARCHs from GOARCH."
