---
id: semaphore-ui
---

# Semaphore UI

1. In Pocket-ID create a new OIDC Client, name it i.e. `Semaphore UI`.
2. Set the callback URL to: `https://<your-semaphore-ui-url>/api/auth/oidc/pocketid/redirect/`.
3. Add the following to your `config.json` file for Semaphore UI:

```json
"oidc_providers": {
    "pocketid": {
        "display_name": "Sign in with PocketID",
        "provider_url": "https://<your-pocket-id-url>",
        "client_id": "<client-id-from-pocket-id>",
        "client_secret": "<client-secret-from-pocket-id>",
        "redirect_url": "https://<your-semaphore-ui-url>/api/auth/oidc/pocketid/redirect/",
        "scopes": [
            "openid",
            "profile",
            "email"
        ],
        "username_claim": "email",
        "name_claim": "given_name"
    }
}
```
