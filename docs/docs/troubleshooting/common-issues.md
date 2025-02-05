---
id: common-issues
---

# Common Issues

## Unable to Add a Passkey

Ensure that the `PUBLIC_APP_URL` is set correctly to the public URL of the Pocket ID instance.

Example:

```ini
PUBLIC_APP_URL=https://id.example.com
```

## Unable to Access the Admin UI After Setup

To set up the initial passkey for the admin user, navigate to:

```
https://id.example.com/login/setup
```

## Invalid Callback URL

One of the most common issues with OIDC clients is a misconfigured `Callback URL`.

If the `redirect_uri` URL parameter starts with `http` but `https` is expected, the client is the issue. If you can’t resolve the issue on the client side, you can add a secondary callback URL using both `http` and `https` versions.
