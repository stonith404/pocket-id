---
id: netbox
---

# Netbox

**This guide does not currently show how to map groups in netbox from OIDC claims**

The following example variables are used, and should be replaced with your actual URLS.

- netbox.example.com (The url of your netbox instance.)
- id.example.com (The url of your Pocket ID instance.)

## Pocket ID Setup

1. In Pocket-ID create a new OIDC Client, name it i.e. `Netbox`.
2. Set a logo for this OIDC Client if you would like too.
3. Set the callback URL to: `https://netbox.example.com/oauth/complete/oidc/`.
4. Copy the `Client ID`, and the `Client Secret` for use in the next steps.

## Netbox Setup

This guide assumes you are using the git based install of netbox.

1. On your netbox server navigate to `/opt/netbox/netbox/netbox`
2. Add the following to your `configuration.py` file:

```python
# Remote authentication support
REMOTE_AUTH_ENABLED = True
REMOTE_AUTH_BACKEND = 'social_core.backends.open_id_connect.OpenIdConnectAuth'
REMOTE_AUTH_HEADER = 'HTTP_REMOTE_USER'
REMOTE_AUTH_USER_FIRST_NAME = 'HTTP_REMOTE_USER_FIRST_NAME'
REMOTE_AUTH_USER_LAST_NAME = 'HTTP_REMOTE_USER_LAST_NAME'
REMOTE_AUTH_USER_EMAIL = 'HTTP_REMOTE_USER_EMAIL'
REMOTE_AUTH_AUTO_CREATE_USER = True
REMOTE_AUTH_DEFAULT_GROUPS = []
REMOTE_AUTH_DEFAULT_PERMISSIONS = {}

SOCIAL_AUTH_OIDC_ENDPOINT = 'https://id.example.com'
SOCIAL_AUTH_OIDC_KEY = '<client id from the first part of this guide>'
SOCIAL_AUTH_OIDC_SECRET = '<client id from the first part of this guide>'
LOGOUT_REDIRECT_URL = 'https://netbox.example.com'
```

3. Save the file and restart netbox: `sudo systemctl start netbox netbox-rq`