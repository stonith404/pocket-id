---
id: caddy
---

### Caddy

With [caddy-security](https://github.com/greenpau/caddy-security) you can easily protect your services with Pocket ID.

#### 1. Create a new OIDC client in Pocket ID.

Create a new OIDC client in Pocket ID by navigating to `https://<your-domain>/settings/admin/oidc-clients`. Now enter `https://<domain-of-proxied-service>/auth/oauth2/generic/authorization-code-callback` as the callback URL. After adding the client, you will obtain the client ID and client secret, which you will need in the next step.

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

  # Configure caddy-security.
	order authenticate before respond
	security {
		oauth identity provider generic {
			realm generic
			driver generic
			client_id client-id-from-pocket-id # Replace with your own client ID
			client_secret client-secret-from-pocket-id # Replace with your own client secret
			scopes openid email profile
			base_auth_url http://localhost
			metadata_url http://localhost/.well-known/openid-configuration
		}

		authentication portal myportal {
			crypto default token lifetime 3600 # Seconds until you have to re-authenticate
			enable identity provider generic
			cookie insecure off # Set to "on" if you're not using HTTPS

			transform user {
				match realm generic
				action add role user
			}
		}

		authorization policy mypolicy {
			set auth url /auth/oauth2/generic
			allow roles user
			inject headers with claims
		}
	}
}

https://<domain-of-your-service> {
	@auth {
		path /auth/oauth2/generic
		path /auth/oauth2/generic/authorization-code-callback
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
