
# use "--build-arg" to build with a specific version
# for example "2.9.1-" if you are building for a pre-v1.0.0 libdns version
ARG CADDY_VERSION=

FROM caddy:${CADDY_VERSION}builder AS builder
RUN caddy-builder \
    github.com/caddy-dns/digitalocean

FROM caddy:${CADDY_VERSION}alpine
COPY --from=builder /usr/bin/caddy /usr/bin/caddy