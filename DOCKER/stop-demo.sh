#!/bin/sh -e

[ -f "./run/lock" ] || (echo The demo is not running. >&2 && exit 3)
NUM=$(cat ./run/validator_num)
NETWORK_NAME=$(cat ./run/network_name)

for i in $(seq 0 $((NUM - 1))); do
	CONTAINER_NAME=chengtay_"$NETWORK_NAME"_"$i"
	(docker stop "$CONTAINER_NAME" && docker rm "$CONTAINER_NAME") &
done

wait
docker network rm "$NETWORK_NAME"
rm ./run/lock
