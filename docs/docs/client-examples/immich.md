---
id: immich
---
# Immich

## Create OIDC Client in Pocket ID
1. Create a new OIDC Client in Pocket ID (e.g., `immich`).
2. Set the callback URLs:  
    ```
    https://<IMMICH-DOMAIN>/auth/login
    https://<IMMICH-DOMAIN>/user-settings
    app.immich:///oauth-callback
    ```
4. Copy the **Client ID**, **Client Secret**, and **OIDC Discovery URL**.

## Configure Immich
1. Open Immich and navigate to:
   **`Administration > Settings > Authentication Settings > OAuth`**
2. Enable **Login with OAuth**.
3. Fill in the required fields:
   - **Issuer URL**: Paste the `Authorization URL` from Pocket ID.
   - **Client ID**: Paste the `Client ID` from Pocket ID.
   - **Client Secret**: Paste the `Client Secret` from Pocket ID.
4. *(Optional)* Change `Button Text` to `Login with Pocket ID`.
5. Save the settings.
6. Test the OAuth login to ensure it works.