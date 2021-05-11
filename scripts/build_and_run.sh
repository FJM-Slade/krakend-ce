# !/bin/bash

# SET GOPATH AND GO BIN
currentFolder=$(pwd)
parent=$(dirname "$currentFolder")
gopath=$(dirname "$parent")

export GOPATH="$gopath" &&
export GOBIN=$GOPATH/bin &&
PATH=$PATH:$GOPATH:$GOBIN &&
export PATH &&

if [[ "$OSTYPE" == "darwin19" ]] || [[ "$OSTYPE" == "darwin20" ]]; then
    echo "Building for MAC-OS..."
    rm $currentFolder/binaries/integration-hub
    # BUILD FOR MACOS
    env GOOS=darwin GOARCH=amd64 go build -o $currentFolder/binaries/integration-hub $currentFolder/cmd/krakend-ce/main.go &&

    echo "binary built! You can find it in:"
    echo "$currentFolder/binaries/integration-hub"

    # RUN FOR MACOS
    export FC_ENABLE=1 && export FC_PARTIALS=$currentFolder/config/partials && export GODEBUG="x509ignoreCN=0" && $currentFolder/binaries/integration-hub run -c  $currentFolder/config/krakend.json
else
    echo "Building for LINUX"
    rm $currentFolder/binaries/integration-hub
    # BUILD AND RUN FOR LINUX
    env GOOS=linux GOARCH=amd64 go build -o $currentFolder/binaries/integration-hub $currentFolder/cmd/krakend-ce/main.go &&

    echo "binary built! You can find it in:"
    echo "$currentFolder/binaries/integration-hub"

    # RUN FOR LINUX
    export FC_ENABLE=1 && export FC_PARTIALS=$currentFolder/config/partials && $currentFolder/binaries/integration-hub run -c  $currentFolder/config/krakend.json
fi
