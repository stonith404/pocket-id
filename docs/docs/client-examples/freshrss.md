---
id: freshrss
---

# FreshRSS

The following example variables are used, and should be replaced with your actual URLs.

- `freshrss.example.com` (The URL of your Proxmox instance.)
- `id.example.com` (The URL of your Pocket ID instance.)

## Pocket ID Setup

1. In Pocket ID create a new OIDC Client, name it, for example, `FreshRSS`.
2. Set a logo for this OIDC Client if you would like to.
3. Set the callback URL to: `https://freshrss.example.com/i/oidc/`.
4. Copy the `Client ID`, `Client Secret`, and `OIDC Discovery URL` for use in the next steps.

## FreshRSS Setup

See [FreshRSS’ OpenID Connect documentation](16_OpenID-Connect.md) for general OIDC settings.

This is an example docker-compose file for FreshRSS with OIDC enabled.

```yaml
services:
  freshrss:
    image: freshrss/freshrss:1.25.0
    container_name: freshrss
    ports:
      - 8080:80
    volumes:
      - /freshrss_data:/var/www/FreshRSS/data
      - /freshrss_extensions:/var/www/FreshRSS/extensions
    environment:
      CRON_MIN: 1,31
      TZ: Etc/UTC
      OIDC_ENABLED: 1
      OIDC_CLIENT_ID: <POCKET_ID_CLIENT_ID>
      OIDC_CLIENT_SECRET: <POCKET_ID_SECRET>
      OIDC_PROVIDER_METADATA_URL: https://id.example.com/.well-known/openid-configuration
      OIDC_SCOPES: openid email profile
      OIDC_X_FORWARDED_HEADERS: X-Forwarded-Proto X-Forwarded-Host
      OIDC_REMOTE_USER_CLAIM: preferred_username
    restart: unless-stopped
    networks:
      - freshrss
networks:
  freshrss:
    name: freshrss
```

:::important
The Username used in Pocket ID must match the Username used in FreshRSS **exactly**. This also applies to case sensitivity. As of version `0.24` of Pocket ID all Usernames are required to be entirely lowercase. FreshRSS allows for uppercase. If a Pocket ID Username is `amanda` and your FreshRSS Username is `Amanda`, you will get a 403 error in FreshRSS and be unable to login. As of version `1.25` of FreshRSS, you are unable to change your username in the GUI. To change your FreshRSS username to lowercase or to match your Pocket ID username, you must nagivate to your FreshRSS volume location. Go to `data/users/` and change the folder for your user to the matching username in Pocket ID, then restart the FreshRSS container to apply the changes.
:::

## Complete OIDC Setup

If you are setting up a new instance of FreshRSS, simply start the container with the OIDC variables and navigate to your FreshRSS URL.

If you are adding OIDC to an existing FreshRSS instance, recreate the container with the docker-compose file with the OIDC variables in it and navigate to your FreshRSS URL. Go to `Settings > Authentication` and change the Authentication method to **HTTP** and hit Submit. Logout to test your OIDC connection.

If you have an error with Pocket ID or are unable to login to your FreshRSS account, you can revert to password login by editing your `config.php` file for FreshRSS. Find the value for `auth_type` and change from `http_auth` to `form`. Restart the FreshRSS container to revert to password login.
