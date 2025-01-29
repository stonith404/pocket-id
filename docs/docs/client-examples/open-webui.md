---
id: open-webui
---

# Open WebUI

1. In Pocket-ID, create a new OIDC Client, name it i.e. `Open WebUI`.
2. Set the callback URL to: `https://openwebui.domain/oauth/oidc/callback`
3. Add the following to your docker `.env` file for Open WebUI:

```ini
  ENABLE_OAUTH_SIGNUP=true
  OAUTH_CLIENT_ID=<client id from pocket ID>
  OAUTH_CLIENT_SECRET=<client secret from pocket ID>
  OAUTH_PROVIDER_NAME=Pocket ID
  OPENID_PROVIDER_URL=https://<your pocket id url>/.well-known/openid-configuration
```
