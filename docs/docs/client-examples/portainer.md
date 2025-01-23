---
id: portainer
---

# Portainer

**This requires Portainers Business Edition**

The following example variables are used, and should be replaced with your actual URLS.

- portainer.example.com (The url of your Portainer instance.)
- id.example.com (The url of your Pocket ID instance.)

## Pocket ID Setup

1. In Pocket-ID create a new OIDC Client, name it i.e. `Portainer`.
2. Set a logo for this OIDC Client if you would like too.
3. Set the callback URL to: `https://portainer.example.com/`.
4. Copy the `Client ID`, `Client Secret`, `Authorization URL`, `Userinfo URL`, and `Token URL` for use in the next steps.

# Portainer Setup

- While initally setting up OAuth in Portainer, its recommended to keep the `Hide internal authentication prompt` set to `Off` incase you need a fallback login
- This guide does **NOT** cover how to setup group claims in Portainer.

1. Open the Portainer web interface and navigate to: `Settings > Authentication`
2. Select `Custom OAuth Provider`
3. Paste the `Client ID` from Pocket ID into the `Client ID` field in Portainer.
4. Paste the `Client Secret` from Pocket ID into the `Client Secret` field in Portainer.
5. Paste the `Authorization URL` from Pocket ID into the `Authorization URL` field in Portainer.
6. Paste the `Token URL` from Pocket ID into the `Access token URL` field in Portainer.
7. Paste the `Userinfo URL` from Pocket ID into the `Resource URL` field in Portainer.
8. Set the `Redirect URL` to `https://portainer.example.com`
9. Set the `Logout URL` to `https://portainer.example.com`
10. Set the `User identifier` field to `preferred_username`. (This will use the users username vs the email)
11. Set the `Scopes` field to: `email openid profile`
12. Set `Auth Style` to `Auto detect`
13. Save the settings and test the new OAuth Login.
