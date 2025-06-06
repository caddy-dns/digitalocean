# DigitalOcean module for Caddy

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy). It can be used to manage DNS records with DigitalOcean accounts.

## Caddy module name

```
dns.providers.digitalocean
```

## Config examples

To use this module for the ACME DNS challenge, [configure the ACME issuer in your Caddy JSON](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuer/acme/) like so:

```json
{
	"module": "acme",
	"challenges": {
		"dns": {
			"provider": {
				"name": "digitalocean",
				"auth_token": "{env.YOUR_DIGITALOCEAN_API_TOKEN}"
			}
		}
	}
}
```

or with the Caddyfile:

```
your.domain.com {
	respond "Hello World"	# replace with whatever config you need...
	tls {
		dns digitalocean {env.YOUR_DIGITALOCEAN_API_TOKEN}
	}
}
```

You can replace `{env.YOUR_DIGITALOCEAN_API_TOKEN}` with the actual auth token if you prefer to put it directly in your config instead of an environment variable.

## Authenticating

See [the associated README in the libdns package](https://github.com/libdns/digitalocean) for important information about credentials.

**NOTE**: If migrating from Caddy v1, you will need to change from using a DigitalOcean API Key to a scoped API Token. Please see link above for more information.

## Building

Don't forget to run `go test -v`.

### Build a Caddy Docker image

> NOTE: this will build using the GitHub version, not a local clone

```bash
docker build -t caddy-with-digitalocean-dns .
```

### Building a Caddy binary using xcaddy

install [xcaddy](https://github.com/caddyserver/xcaddy): `go install github.com/caddyserver/xcaddy/cmd/xcaddy@latest`

Then build caddy:

> NOTE: this will build using the GitHub version, not a local clone


```
xcaddy build \
    --with github.com/caddy-dns/digitalocean@master
```