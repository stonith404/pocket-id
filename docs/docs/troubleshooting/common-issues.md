---
id: common-issues
---

# Common Issues

# Pocket ID doesn’t load after setup 

Make sure the `PUBLIC_APP_URL` is set correctly to the public host users would access it from.

Example:  
```ini
PUBLIC_APP_URL=https://id.example.com  
```

## Unable to Access the admin UI after setup  

To setup the initial passkey for the admin user, navigate to `https://id.example.com/login/setup`.  

## Invalid callback URL  

One of the most common issues with OIDC Clients in missconfigured `Callback URLs`

If the `redirect_uri` URL param starts with `http` but `https` is expected, the client is the issue. If you can’t solve the issue on the client side you can just add a secondary callback URL using `http` as well as the `https` URL.
