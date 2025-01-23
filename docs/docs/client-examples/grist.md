---
id: grist
---

# Grist

## Pocket ID Setup
1. In Pocket-ID create a new OIDC Client, name it i.e. `Grist`
2. Set the callback url to: `https://<Grist Host>/oauth2/callback`
3. In Grist (Docker/Docker Compose/etc), set these environment variables:

```ini
GRIST_OIDC_IDP_ISSUER="https://<Pocket ID Host>/.well-known/openid-configuration"
GRIST_OIDC_IDP_CLIENT_ID="<Client ID from the OIDC Client created in Pocket ID>"
GRIST_OIDC_IDP_CLIENT_SECRET="<Client Secret from the OIDC Client created in Pocket ID>"
GRIST_OIDC_SP_HOST="https://<Grist Host>"
GRIST_OIDC_IDP_SCOPES="openid email profile"  # Default
GRIST_OIDC_IDP_SKIP_END_SESSION_ENDPOINT=true  # Default=false, needs to be true for Pocket Id b/c end_session_endpoint is not implemented
GRIST_OIDC_IDP_END_SESSION_ENDPOINT="https://<Pocket ID Host>/api/webauthn/logout" # Only set this if GRIST_OIDC_IDP_SKIP_END_SESSION_ENDPOINT=false and you need to define a custom endpoint
```
4. Also ensure that the `GRIST_DEFAULT_EMAIL` env variable is set to the same email address as your user profile within Pocket ID
5. Start/Restart Grist