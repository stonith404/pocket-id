# Sample configrations for well-known self-hosted services
## [Vikunja](https://vikunja.io/)
1. In Pocket-ID create a new OIDC Client, name it i.e. `Vikunja`
2. Set the callback url to: `https://<your-vikunja-subdomain>.<your-domain>/auth/openid/pocketid`
3. In `Vikunja` ensure to map a config file to your container, see [here](https://vikunja.io/docs/config-options/#using-a-config-file-with-docker-compose)
4. Add or set the following content to the `config.yml` file:
 ```yml
 auth:
  openid:
    enabled: true
    redirecturl: https://<your-vikunja-subdomain>.<your-domain>/auth/openid/pocketid
    providers:
      - name: Pocket-Id
        authurl: https://<your-pocket-id-subdomain>.<your-domain>
        clientid: <client id from the created OIDC client>
        clientsecret: <client secret from the created OIDC client>

  ```
## [Hoarder](https://docs.hoarder.app/)
1. In Pocket-ID create a new OIDC Client, name it i.e. `Hoarder`
2. Set the callback url to: `https://<your-hoarder-subdomain>.<your-domain>/api/auth/callback/custom`
3. Open your  `.env` file from your Hoarder compose and add these lines: 
```ini
  OAUTH_WELLKNOWN_URL = https://<your-pocket-id-subdomain>.<your-domain>/.well-known/openid-configuration
  OAUTH_CLIENT_SECRET = <client secret from the created OIDC client>
  OAUTH_CLIENT_ID = <client id from the created OIDC client>
  OAUTH_PROVIDER_NAME = Pocket-Id
  NEXTAUTH_URL = https:///<your-hoarder-subdomain>.<your-domain>

```
4. Optional: If you like to disable password authentication and link your existing hoarder account with your pocket-id identity
```ini
  DISABLE_PASSWORD_AUTH	= true
  OAUTH_ALLOW_DANGEROUS_EMAIL_ACCOUNT_LINKING = true
```
