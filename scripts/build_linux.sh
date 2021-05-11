# !/bin/bash

# SET GOPATH AND GO BIN
currentFolder=$(pwd)
parent=$(dirname "$currentFolder")
gopath=$(dirname "$parent")

export GOPATH="$gopath" &&
export GOBIN=$GOPATH/bin &&
PATH=$PATH:$GOPATH:$GOBIN &&
export PATH &&

echo "Building for LINUX"
# BUILD AND RUN FOR LINUX
rm $currentFolder/binaries/integration-hub
env GOOS=linux GOARCH=amd64 GODEBUG=x509ignorecn=0 go build -o $currentFolder/binaries/integration-hub $currentFolder/cmd/krakend-ce/main.go &&

echo "binary built! You can find it in:"
echo "$currentFolder/binaries/integration-hub"
