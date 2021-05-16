# !/bin/bash

# SET GOPATH AND GO BIN
currentFolder=$(pwd)
parent=$(dirname "$currentFolder")
gopath=$(dirname "$parent")

export GOPATH="$gopath" &&
export GOBIN=$GOPATH/bin &&
PATH=$PATH:$GOPATH:$GOBIN &&
export PATH &&

    echo "Running in MAC-OS..."


    # RUN FOR MACOS
    export FC_ENABLE=1 && export FC_PARTIALS=$currentFolder/config/partials && export FC_SETTINGS=$currentFolder/config/settings && export GODEBUG="x509ignoreCN=0" && $currentFolder/binaries/integration-hub run -c  $currentFolder/config/krakend.json

