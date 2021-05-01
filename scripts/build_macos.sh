# !/bin/bash

# SET GOPATH AND GO BIN
currentFolder=$(pwd)
parent=$(dirname "$currentFolder")
gopath=$(dirname "$parent")

echo "$currentFolder"
echo "$parent"
echo "$gopath"



export GOPATH="$gopath" &&
export GOBIN=$GOPATH/bin &&
PATH=$PATH:$GOPATH:$GOBIN &&
export PATH &&

echo "Building for MAC-OS..."
# BUILD FOR MACOS
rm $currentFolder/binaries/integration-hub
env GOOS=darwin GOARCH=amd64 go build -o $currentFolder/binaries/integration-hub $currentFolder/cmd/krakend-ce/main.go &&

echo "binary built! You can find it in:"
echo "$currentFolder/binaries/integration-hub"
