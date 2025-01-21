---
id: cloudflare-zero-trust
---
# Cloudflare Zero Trust

**Note: Cloudflare will need to be able to reach your Pocket ID instance and vice versa for this to work correctly**

## Pocket ID Setup

1. In Pocket-ID create a new OIDC Client, name it i.e. `Cloudflare Zero Trust`.
2. Set a logo for this OIDC Client if you would like too.
3. Set the callback URL to: `https://<your-team-name>.cloudflareaccess.com/cdn-cgi/access/callback`.
4. Copy the Client ID, Client Secret, Authorization URL, Token URL, and Certificate URL for the next steps.

## Cloudflare Zero Trust Setup

1. Login to Cloudflare Zero Trust [Dashboard](https://one.dash.cloudflare.com/).
2. Navigate to Settings > Authentication > Login Methods.
3. Click `Add New` under login methods.
4. Create a name for the new login method.
5. Paste in the `Client ID` from Pocket ID into the `App ID` field.
6. Paste the `Client Secret` from Pocket ID into the `Client Secret` field.
7. Paste the `Authorization URL` from Pocket ID into the `Auth URL` field.
8. Paste the `Token URL` from Pocket ID into the `Token URL` field.
9. Paste the `Certificate URL` from Pocket ID into the `Certificate URL` field.
10. Save the new login method and test to make sure it works with cloudflare.
