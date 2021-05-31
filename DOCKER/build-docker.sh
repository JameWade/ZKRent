#!/bin/sh -e

DOCKER_IMAGE_NAME="chengtay_chain_node"

# checking the docker file

[ -f "./chengtay-chain/Dockerfile" ] || (echo Can not find the docker file at \"./chengtay-chain/Dockerfile\" >&2 && exit 1)

# checking the executable file

[ -f "../build/chengtay-chain" ] ||
	(echo Missing executable file \"chengtay-chain\". Try to build one. >&2 && (cd .. && make build-upx)) ||
	(echo \"make build-upx\" failed. Fallback to \"make build\" >&2 && (cd .. && make build)) ||
	(echo Build failed. Please place the executable file \"chengtay-chain\" at \"../build/chengtay-chain\". >&2 && exit 1)

test -x ../build/chengtay-chain || chmod +x ../build/chengtay-chain

# build docker image

cp ../build/chengtay-chain ./chengtay-chain/
test -x ./chengtay-chain/chengtay-chain || chmod +x ./chengtay-chain/chengtay-chain

(cd ./chengtay-chain && docker build --tag "$DOCKER_IMAGE_NAME" .)
