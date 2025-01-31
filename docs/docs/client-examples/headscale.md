---
id: headscale
---
# Headscale

## Create OIDC Client in Pocket ID
1. Create a new OIDC Client in Pocket ID (e.g., `Headscale`).
2. Set the callback URL:  `https://<HEADSCALE-DOMAIN>/oidc/callback`
3. Enable `PKCE`.
4. Copy the **Client ID** and **Client Secret**.

## Configure Headscale
> Refer to the example [`config.yaml`](https://github.com/juanfont/headscale/blob/main/config-example.yaml) for full OIDC configuration options.

Add the following to `config.yaml`:

```yaml
oidc:
  issuer: "https://<POCKET-ID-DOMAIN>"
  client_id: "<CLIENT-ID>"
  client_secret: "<CLIENT-SECRET>"
  pkce:
    enabled: true
    method: S256
```

### (Optional) Restrict Access to Certain Groups
To allow only specific groups, add:

```yaml
  scope: ["openid", "profile", "email", "groups"]
  allowed_groups:
    - <POCKET-ID-GROUP-NAME> #example: headscale 
```
