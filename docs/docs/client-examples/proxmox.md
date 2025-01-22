---
id: proxmox
---

# Proxmox

The following example variables are used, and should be replaced with your actual URLS.

- proxmox.example.com (The url of your proxmox instance.)
- id.example.com (The url of your Pocket ID instance.)

## Pocket ID Setup

1. In Pocket-ID create a new OIDC Client, name it i.e. `Proxmox`.
2. Set a logo for this OIDC Client if you would like too.
3. Set the callback URL to: `https://proxmox.example.com`.
4. Copy the `Client ID`, and the `Client Secret` for use in the next steps.

## Proxmox Setup

1. Open the Proxmox Console and navigate to: `Datacenter - Realms`
2. Add a new `Open ID Connect Server` Realm
3. Enter `https://id.example.com` for the `Issuer URL`
4. Enter a name for the realm of your choice ie. `PocketID`
5. Paste the `Client ID` from Pocket ID into the `Client ID` field in Proxmox.
6. Paste the `Client Secret` from Pocket ID into the `Client Key` field in Proxmox.
7. You can check the `Default` box if you want this to be the deafult realm proxmox uses when signing in.
8. Check the `Autocreate Users` checkbox. (This will automaitcally create users in Proxmox if they dont exsist.)
9. Select `username` for the `Username Claim` dropdown (This is personal preference and controls how the username is shown, ie: `username = username@PocketID` or `email = username@example@PocketID`).
10. Leave the rest as defaults and click `OK` to save the new realm.
