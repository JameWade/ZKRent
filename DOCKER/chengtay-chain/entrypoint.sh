#!/bin/sh -e

if [ ! -d "$TMHOME" ]; then
	echo \$TMHOME Directory does not exist. "$TMHOME" >&2
	exit 1
fi

#if [ ! -f "$TMHOME/config/priv_validator_key.json" ] || [ ! -f "$TMHOME/data/priv_validator_state.json" ] || [ ! -f "$TMHOME/config/node_key.json" ]; then
#	/usr/bin/chengtay-chain init
#fi

/usr/bin/chengtay-chain init

sed -i \
	-e 's/^addr_book_strict\s*=.*/addr_book_strict = false/' \
	-e 's/^timeout_commit\s*=.*/timeout_commit = "500ms"/' \
	-e 's/^index_all_tags\s*=.*/index_all_tags = true/' \
	-e 's,^laddr = "tcp://127.0.0.1:26657",laddr = "tcp://0.0.0.0:26657",' \
	-e 's/^prometheus\s*=.*/prometheus = true/' \
	"$TMHOME/config/config.toml"

if [ -n "$MONIKER" ]; then
	sed -i \
		-e 's/^moniker\s*=.*/moniker = "'"$MONIKER"'"/' \
		"$TMHOME/config/config.toml"
fi

if [ -n "$BOOT_NODE_STR" ]; then
	/usr/bin/chengtay-chain "$@" --p2p.persistent_peers "$BOOT_NODE_STR"
else
	/usr/bin/chengtay-chain "$@"
fi
