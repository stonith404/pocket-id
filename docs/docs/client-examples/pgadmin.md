---
id: pgadmin
---

# pgAdmin

The following example variables are used, and should be replaced with your actual URLS.

- pgadmin.example.com (The url of your pgAdmin instance.)
- id.example.com (The url of your Pocket ID instance.)

## Pocket ID Setup

1. In Pocket-ID create a new OIDC Client, name it i.e. `pgAdmin`.
2. Set a logo for this OIDC Client if you would like too.
3. Set the callback URL to: `https://pgadmin.example.com/oauth2/authorize`.
4. Copy the `Client ID`, `Client Secret`, `Authorization URL`, `Userinfo URL`, `Token URL`, and `OIDC Discovery URL` for use in the next steps.

# pgAdmin Setup

1. Add the following to the `config_local.py` file for pgAdmin:

**Make sure to replace https://id.example.com with your actual Pocket ID URL**

```python
AUTHENTICATION_SOURCES = ['oauth2', 'internal'] # This keeps internal authentication enabled as well as oauth2
OAUTH2_AUTO_CREATE_USER = True
OAUTH2_CONFIG = [{
        'OAUTH2_NAME' : 'pocketid',
        'OAUTH2_DISPLAY_NAME' : 'Pocket ID',
        'OAUTH2_CLIENT_ID' : '<client id from the earlier step>',
        'OAUTH2_CLIENT_SECRET' : '<client secret from the earlier step>',
        'OAUTH2_TOKEN_URL' : 'https://id.example.com/api/oidc/token',
        'OAUTH2_AUTHORIZATION_URL' : 'https://id.example/authorize',
        'OAUTH2_API_BASE_URL' : 'https://id.example.com',
        'OAUTH2_USERINFO_ENDPOINT' : 'https://id.example.com/api/oidc/userinfo',
        'OAUTH2_SERVER_METADATA_URL' : 'https://id.example.com/.well-known/openid-configuration',
        'OAUTH2_SCOPE' : 'openid email profile',
        'OAUTH2_ICON' : 'fa-openid',
        'OAUTH2_BUTTON_COLOR' : '#fd4b2d' # Can select any color you would like here.
}]
```
