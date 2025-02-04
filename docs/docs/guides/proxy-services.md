---
id: proxy-services
---

# Proxy Services

The goal of Pocket ID is to function exclusively as an OIDC provider. As such, we don't have a built-in proxy provider. However, you can use other tools that act as a middleware to protect your services and support OIDC as an authentication provider.

There are two ways to do this:

- Implement OIDC into your reverse proxy
- Use [OAuth2 Proxy](https://oauth2-proxy.github.io/oauth2-proxy/)

## Reverse Proxy

Almost every reverse proxy somehow supports protecting your services with OIDC. Currently only Caddy is documented but you can search on Google for your reverse proxy and OIDC.

We would really appreciate if you contribute to this documentation by adding your reverse proxy and how to configure it with Pocket ID.

### Caddy

With [caddy-security](https://github.com/greenpau/caddy-security) you can easily protect your services with Pocket ID.

#### 1. Create a new OIDC client in Pocket ID.

Create a new OIDC client in Pocket ID by navigating to `https://<your-domain>/settings/admin/oidc-clients`. Now enter `https://<domain-of-proxied-service>/caddy-security/oauth2/generic/authorization-code-callback` as the callback URL. After adding the client, you will obtain the client ID and client secret, which you will need in the next step.

#### 2. Install caddy-security

Run the following command to install caddy-security:

```bash
caddy add-package github.com/greenpau/caddy-security
```

#### 3. Create your Caddyfile

```bash
{
  	# Port to listen on
	http_port 443

  	# Configure caddy-security.
	order authenticate before respond
	security {
		oauth identity provider generic {
			delay_start 3
			realm generic
			driver generic
			client_id client-id-from-pocket-id # Replace with your own client ID
			client_secret client-secret-from-pocket-id # Replace with your own client secret
			scopes openid email profile
			base_auth_url http://localhost
			metadata_url http://localhost/.well-known/openid-configuration
		}

		authentication portal myportal {
			crypto default token lifetime 3600 # Seconds until you have to re-authenticate
			enable identity provider generic
			cookie insecure off # Set to "on" if you're not using HTTPS

			transform user {
				match realm generic
				action add role user
			}
		}

		authorization policy mypolicy {
			set auth url /caddy-security/oauth2/generic
			allow roles user
			inject headers with claims
		}
	}
}

https://<domain-of-your-service> {
	@auth {
		path /caddy-security/*
    }

	route @auth {
		authenticate with myportal
	}

	route /* {
		authorize with mypolicy
		reverse_proxy http://<service-to-be-proxied>:<port> # Replace with your own service
	}
}
```

For additional configuration options, refer to the official [caddy-security documentation](https://docs.authcrunch.com/docs/intro).

#### 4. Start Caddy

```bash
caddy run --config Caddyfile
```

#### 5. Access the service

Your service should now be protected by Pocket ID.

## OAuth2 Proxy

### Docker Installation

#### 1. Add OAuth2 proxy to the service that should be proxied.

To configure OAuth2 Proxy with Pocket ID, you have to add the following service to the service that should be proxied. E.g., if [Uptime Kuma](https://github.com/louislam/uptime-kuma) should be proxied, you can add the following service to the `docker-compose.yml` of Uptime Kuma:

```yaml
# Example with Uptime Kuma
# uptime-kuma:
#  image: louislam/uptime-kuma
oauth2-proxy:
  image: quay.io/oauth2-proxy/oauth2-proxy:v7.6.0
  command: --config /oauth2-proxy.cfg
  volumes:
    - "./oauth2-proxy.cfg:/oauth2-proxy.cfg"
  ports:
    - 4180:4180
```

#### 2. Create a new OIDC client in Pocket ID.

Create a new OIDC client in Pocket ID by navigating to `https://<your-domain>/settings/admin/oidc-clients`. Now enter `https://<domain-of-proxied-service>/oauth2/callback` as the callback URL. After adding the client, you will obtain the client ID and client secret, which you will need in the next step.

#### 3. Create a configuration file for OAuth2 Proxy.

Create a configuration file named `oauth2-proxy.cfg` in the same directory as your `docker-compose.yml` file of the service that should be proxied (e.g. Uptime Kuma). This file will contain the necessary configurations for OAuth2 Proxy to work with Pocket ID.

Here is the recommend `oauth2-proxy.cfg` configuration:

```cfg
# Replace with your own credentials
client_id="client-id-from-pocket-id"
client_secret="client-secret-from-pocket-id"
oidc_issuer_url="https://<your-pocket-id-domain>"

# Replace with a secure random string
cookie_secret="random-string"

# Upstream servers (e.g http://uptime-kuma:3001)
upstreams="http://<service-to-be-proxied>:<port>"

# Additional Configuration
provider="oidc"
scope = "openid email profile groups"

# If you are using a reverse proxy in front of OAuth2 Proxy
reverse_proxy = true

# Email domains allowed for authentication
email_domains = ["*"]

# If you are using HTTPS
cookie_secure="true"

# Listen on all interfaces
http_address="0.0.0.0:4180"
```

For additional configuration options, refer to the official [OAuth2 Proxy documentation](https://oauth2-proxy.github.io/oauth2-proxy/configuration/overview).

#### 4. Start the services.

After creating the configuration file, you can start the services using Docker Compose:

```bash
docker compose up -d
```

#### 5. Access the service.

You can now access the service through OAuth2 Proxy by visiting `http://localhost:4180`.

### Standalone Installation

Setting up OAuth2 Proxy with Pocket ID without Docker is similar to the Docker setup. As the setup depends on your environment, you have to adjust the steps accordingly but is should be similar to the Docker setup.

You can visit the official [OAuth2 Proxy documentation](https://oauth2-proxy.github.io/oauth2-proxy/installation) for more information.
