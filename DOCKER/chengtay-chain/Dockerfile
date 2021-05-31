FROM alpine:3.12

# Tendermint will be looking for the genesis file in /chengtay/config/genesis.json
# (unless you change `genesis_file` in config.toml). You can put your config.toml and
# private validator file into /chengtay/config.
#
# The /chengtay/data dir is used by tendermint to store state.
ENV TMHOME /chengtay
WORKDIR $TMHOME

# p2p, rpc and prometheus port
EXPOSE 26656 26657 26660

ENTRYPOINT ["/usr/bin/entrypoint.sh"]
CMD ["node"]
STOPSIGNAL SIGTERM

COPY ["chengtay-chain", "entrypoint.sh", "/usr/bin/"]

# OS environment setup
# jq and curl used for extracting `pub_key` from private validator while
# deploying tendermint with Kubernetes.

# China specific
#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

#RUN apk update && \
#    apk upgrade && \
#    apk --no-cache add curl jq bash

RUN chmod +x /usr/bin/chengtay-chain && \
    chmod +x /usr/bin/entrypoint.sh

# Expose the data directory as a volume since there's mutable state in there
VOLUME [ $TMHOME ]


