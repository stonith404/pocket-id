---
id: headscale
---
# Headscale

## Step 1: Create OIDC Client in Pocket ID
1. Create a new OIDC Client in Pocket ID (e.g., `Headscale`).
2. Set the callback URL:  
   ```
   https://<your-headscale-domain>/oidc/callback
   ```
3. Enable `PKCE`.
4. Copy the **Client ID** and **Client Secret**.

## Step 2: Configure Headscale
Add the following to `config.yaml`:

```yaml
oidc:
  issuer: "https://<your-pocket-id-domain>"
  client_id: "<client id from the created OIDC client>"
  client_secret: "<client secret from the created OIDC client>"
  pkce:
    enabled: true
    method: S256
```

### (Optional) Restrict Access to Certain Groups
To allow only specific groups, add:

```yaml
  scope: ["openid", "profile", "email", "groups"]
  allowed_groups:
    - <pocket-id-group-name>
```
## Additional Resources
You can find an example `config.yaml` for Headscale on their [Github](https://github.com/juanfont/headscale/blob/main/config-example.yaml).
