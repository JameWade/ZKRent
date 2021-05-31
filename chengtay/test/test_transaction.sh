#! /bin/bash
set -ex

#- kvstore over socket, curl
#- counter over socket, curl
#- counter over grpc, curl
#- counter over grpc, grpc

# TODO: install everything

export PATH="$GOBIN:$PATH"
export TMHOME=$HOME/.chengtaychain

# start tendermint first
function chengtay_start(){
    rm -rf $TMHOME
    echo "new gen_genesis"
    echo `pwd`
    echo `ls ./build`
    ./build/chengtay-chain gen_genesis 1
    echo "copy gen_genesis and init"
    mv $TMHOME/config/priv_validator_key.json.0.json $TMHOME/config/priv_validator_key.json
    mv $TMHOME/data/priv_validator_state.json.0.json $TMHOME/data/priv_validator_state.json
    ./build/chengtay-chain init
    echo "chengtay start"
    ./build/chengtay-chain node > tendermint.log &
    pid_tendermint=$!

    echo "running test"
    go test ./chengtay/test/chengtay_test.go
    sleep 5
    kill -9 $pid_cli $pid_tendermint
}

chengtay_start
