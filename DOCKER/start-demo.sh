#!/bin/sh -e

[ ! -f "./run/lock" ] || (echo The demo is already running. >&2 && exit 3)

touch ./run/lock

test -x ./build-docker.sh || chmod +x ./build-docker.sh
./build-docker.sh

DOCKER_IMAGE_NAME="chengtay_chain_node"

# get the number from input

printf "The number of validators? "
read -r num

case $num in
'' | *[!0-9]*) echo An integer is required. >&2 && exit 2 ;;
*) ;;
esac

echo "$num" | tee ./run/validator_num

# generating the genesis file and private keys

[ ! -d "./data/" ] || (cd "./data/" && rm -rf ./*)
mkdir -p ./data/common
TMHOME=./data/common ../build/chengtay-chain gen_genesis "$num"

for i in $(seq 0 $((num - 1))); do
	mkdir -p ./data/chain-node-"$i"/config/
	mkdir -p ./data/chain-node-"$i"/data/
	mkdir -p ./data/chain-node-"$i"/run/
	cp ./data/common/config/config.toml ./data/chain-node-"$i"/config/
	cp ./data/common/config/genesis.json ./data/chain-node-"$i"/config/
	mv ./data/common/config/priv_validator_key.json."$i".json ./data/chain-node-"$i"/config/priv_validator_key.json
	mv ./data/common/data/priv_validator_state.json."$i".json ./data/chain-node-"$i"/data/priv_validator_state.json
done

# create a new docker network

NETWORK_NAME=$(dd if=/dev/urandom bs=8 count=1 | xxd -p | head)
echo "$NETWORK_NAME" | tee ./run/network_name

docker network create "$NETWORK_NAME"

# start the node one by one

BOOT_NODE_STR=""
for i in $(seq 0 $((num - 1))); do
	CONTAINER_NAME=chengtay_"$NETWORK_NAME"_"$i"

	TMHOME=./data/chain-node-"$i" ../build/chengtay-chain init
	NODE_ID=$(TMHOME=./data/chain-node-"$i" ../build/chengtay-chain show_node_id)
	NODE_IP="$CONTAINER_NAME"
	NODE_PORT=26656
	NODE_STR="$NODE_ID@$NODE_IP:$NODE_PORT"

	docker run \
		--detach \
		--name "$CONTAINER_NAME" \
		-v "$PWD"/data/chain-node-"$i":/chengtay \
		--network "$NETWORK_NAME" \
		--network-alias "$CONTAINER_NAME" \
		-e BOOT_NODE_STR="$BOOT_NODE_STR" \
		-e MONIKER="$CONTAINER_NAME" \
		"$DOCKER_IMAGE_NAME"

	if [ -z "$BOOT_NODE_STR" ]; then
		BOOT_NODE_STR="$NODE_STR"
	else
		BOOT_NODE_STR="$BOOT_NODE_STR,$NODE_STR"
	fi

done
