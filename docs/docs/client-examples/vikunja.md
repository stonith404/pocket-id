---
id: vikunja
---

# Vikunja

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
