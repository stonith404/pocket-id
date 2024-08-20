# Proxy Services through Pocket ID

The goal of Pocket ID is to stay simple. Because of that we don't have a built-in proxy provider. However, you can use [OAuth2 Proxy](https://oauth2-proxy.github.io/) to add authentication to your services that don't support OIDC. This guide will show you how to set up OAuth2 Proxy with Pocket ID.

## Docker Setup

#### 1. Add OAuth2 proxy to the service that should be proxied.

To configure OAuth2 Proxy with Pocket ID, you have to add the following service to the service that should be proxied. E.g., [Uptime Kuma](https://github.com/louislam/uptime-kuma) should be proxied, you can add the following service to the `docker-compose.yml` of Uptime Kuma:

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

Create a new OIDC client in Pocket ID by navigating to `https://<your-domain>/settings/admin/oidc-clients`. After adding the client, you will obtain the client ID and client secret.

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
scope = "openid email profile"

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

## Standalone Installation

Setting up OAuth2 Proxy with Pocket ID without Docker is similar to the Docker setup. As the setup depends on your environment, you have to adjust the steps accordingly but is should be similar to the Docker setup.

You can visit the official [OAuth2 Proxy documentation](https://oauth2-proxy.github.io/oauth2-proxy/installation) for more information.
