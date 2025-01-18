# <div align="center"><img  src="https://github.com/user-attachments/assets/4ceb2708-9f29-4694-b797-be833efce17d" width="100"/> </br>Pocket ID</div>

Pocket ID is a simple OIDC provider that allows users to authenticate with their passkeys to your services.

→ Try out the [Demo](https://pocket-id.eliasschneider.com)

<img src="https://github.com/user-attachments/assets/96ac549d-b897-404a-8811-f42b16ea58e2" width="1200"/>

The goal of Pocket ID is to be a simple and easy-to-use. There are other self-hosted OIDC providers like [Keycloak](https://www.keycloak.org/) or [ORY Hydra](https://www.ory.sh/hydra/) but they are often too complex for simple use cases.

Additionally, what makes Pocket ID special is that it only supports [passkey](https://www.passkeys.io/) authentication, which means you don’t need a password. Some people might not like this idea at first, but I believe passkeys are the future, and once you try them, you’ll love them. For example, you can now use a physical Yubikey to sign in to all your self-hosted services easily and securely.

## Table of Contents

- [ Pocket ID](#-pocket-id)
  - [Table of Contents](#table-of-contents)
  - [Setup](#setup)
    - [Before you start](#before-you-start)
    - [Installation with Docker (recommended)](#installation-with-docker-recommended)
    - [Unraid](#unraid)
    - [Stand-alone Installation](#stand-alone-installation)
    - [Nginx Reverse Proxy](#nginx-reverse-proxy)
  - [Proxy Services with Pocket ID](#proxy-services-with-pocket-id)
  - [Update](#update)
      - [Docker](#docker)
      - [Stand-alone](#stand-alone)
  - [Environment variables](#environment-variables)
  - [Account recovery](#account-recovery)
  - [Contribute](#contribute)

## Setup

> [!WARNING]
> Pocket ID is in its early stages and may contain bugs. There might be OIDC features that are not yet implemented. If you encounter any issues, please open an issue.

### Before you start

Pocket ID requires a [secure context](https://developer.mozilla.org/en-US/docs/Web/Security/Secure_Contexts), meaning it must be served over HTTPS. This is necessary because Pocket ID uses the [WebAuthn API](https://developer.mozilla.org/en-US/docs/Web/API/Web_Authentication_API).

### Installation with Docker (recommended)

1. Download the `docker-compose.yml` and `.env` file:

   ```bash
    curl -O https://raw.githubusercontent.com/stonith404/pocket-id/main/docker-compose.yml

    curl -o .env https://raw.githubusercontent.com/stonith404/pocket-id/main/.env.example
   ```

2. Edit the `.env` file so that it fits your needs. See the [environment variables](#environment-variables) section for more information.
3. Run `docker compose up -d`

You can now sign in with the admin account on `http://localhost/login/setup`.

### Unraid

Pocket ID is available as a template on the Community Apps store.

### Stand-alone Installation

Required tools:

- [Node.js](https://nodejs.org/en/download/) >= 20
- [Go](https://golang.org/doc/install) >= 1.23
- [Git](https://git-scm.com/downloads)
- [PM2](https://pm2.keymetrics.io/)
- [Caddy](https://caddyserver.com/docs/install) (optional)

1. Copy the `.env.example` file in the `frontend` and `backend` folder to `.env` and change it so that it fits your needs.

   ```bash
   cp frontend/.env.example frontend/.env
   cp backend/.env.example backend/.env
   ```

2. Run the following commands:

   ```bash
   git clone https://github.com/stonith404/pocket-id
   cd pocket-id

   # Checkout the latest version
   git fetch --tags && git checkout $(git describe --tags `git rev-list --tags --max-count=1`)

   # Start the backend
   cd backend/cmd
   go build -o ../pocket-id-backend
   cd ..
   pm2 start pocket-id-backend --name pocket-id-backend

   # Start the frontend
   cd ../frontend
   npm install
   npm run build
   pm2 start --name pocket-id-frontend --node-args="--env-file .env" build/index.js

   # Optional: Start Caddy (You can use any other reverse proxy)
   cd ..
   pm2 start --name pocket-id-caddy caddy -- run --config reverse-proxy/Caddyfile
   ```

You can now sign in with the admin account on `http://localhost/login/setup`.

### Nginx Reverse Proxy

To use Nginx as a reverse proxy for Pocket ID, update the configuration to increase the header buffer size. This adjustment is necessary because SvelteKit generates larger headers, which may exceed the default buffer limits.

```nginx
proxy_busy_buffers_size   512k;
proxy_buffers   4 512k;
proxy_buffer_size   256k;
```

## Proxy Services with Pocket ID

The goal of Pocket ID is to function exclusively as an OIDC provider. As such, we don't have a built-in proxy provider. However, you can use other tools that act as a middleware to protect your services and support OIDC as an authentication provider.

See the [guide](docs/proxy-services.md) for more information.

## Update

#### Docker

```bash
docker compose pull
docker compose up -d
```

#### Stand-alone

1. Stop the running services:
   ```bash
   pm2 delete pocket-id-backend pocket-id-frontend pocket-id-caddy
   ```
2. Run the following commands:

   ```bash
   cd pocket-id

   # Checkout the latest version
   git fetch --tags && git checkout $(git describe --tags `git rev-list --tags --max-count=1`)

   # Start the backend
   cd backend/cmd
   go build -o ../pocket-id-backend
   cd ..
   pm2 start pocket-id-backend --name pocket-id-backend

   # Start the frontend
   cd ../frontend
   npm install
   npm run build
   pm2 start build/index.js --name pocket-id-frontend

   # Optional: Start Caddy (You can use any other reverse proxy)
   cd ..
   pm2 start caddy --name pocket-id-caddy -- run --config reverse-proxy/Caddyfile
   ```

## Environment variables

| Variable                     | Default Value             | Recommended to change | Description                                                                                                                                                                                                                                                                                                                                                               |
| ---------------------------- | ------------------------- | --------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `PUBLIC_APP_URL`             | `http://localhost`        | yes                   | The URL where you will access the app.                                                                                                                                                                                                                                                                                                                                    |
| `TRUST_PROXY`                | `false`                   | yes                   | Whether the app is behind a reverse proxy.                                                                                                                                                                                                                                                                                                                                |
| `MAXMIND_LICENSE_KEY`        | `-`                       | yes                   | License Key for the GeoLite2 Database. The license key is required to retrieve the geographical location of IP addresses in the audit log. If the key is not provided, IP locations will be marked as "unknown." You can obtain a license key for free [here](https://www.maxmind.com/en/geolite2/signup).                                                                |
| `PUID` and `PGID`            | `1000`                    | yes                   | The user and group ID of the user who should run Pocket ID inside the Docker container and owns the files that are mounted with the volume. You can get the `PUID` and `GUID` of your user on your host machine by using the command `id`. For more information see [this article](https://docs.linuxserver.io/general/understanding-puid-and-pgid/#using-the-variables). |
| `DB_PROVIDER`                | `sqlite`                  | no                    | The database provider you want to use. Currently `sqlite` and `postgres` are supported.                                                                                                                                                                                                                                                                                   |
| `SQLITE_DB_PATH`             | `data/pocket-id.db`       | no                    | The path to the SQLite database. This gets ignored if you didn't set `DB_PROVIDER` to `sqlite`.                                                                                                                                                                                                                                                                           |
| `POSTGRES_CONNECTION_STRING` | `-`                       | no                    | The connection string to your Postgres database. This gets ignored if you didn't set `DB_PROVIDER` to `postgres`. A connection string can look like this: `postgresql://user:password@host:5432/pocket-id`.                                                                                                                                                               |
| `UPLOAD_PATH`                | `data/uploads`            | no                    | The path where the uploaded files are stored.                                                                                                                                                                                                                                                                                                                             |
| `INTERNAL_BACKEND_URL`       | `http://localhost:8080`   | no                    | The URL where the backend is accessible.                                                                                                                                                                                                                                                                                                                                  |
| `GEOLITE_DB_PATH`            | `data/GeoLite2-City.mmdb` | no                    | The path where the GeoLite2 database should be stored.                                                                                                                                                                                                                                                                                                                    |
| `CADDY_PORT`                 | `80`                      | no                    | The port on which Caddy should listen. Caddy is only active inside the Docker container. If you want to change the exposed port of the container then you sould change this variable.                                                                                                                                                                                     |
| `PORT`                       | `3000`                    | no                    | The port on which the frontend should listen.                                                                                                                                                                                                                                                                                                                             |
| `BACKEND_PORT`               | `8080`                    | no                    | The port on which the backend should listen.                                                                                                                                                                                                                                                                                                                              |


## Account recovery

There are two ways to create a one-time access link for a user:

1. **UI**: An admin can create a one-time access link for the user in the admin panel under the "Users" tab by clicking on the three dots next to the user's name and selecting "One-time link".
2. **Terminal**: You can create a one-time access link for a user by running the `scripts/create-one-time-access-token.sh` script. To execute this script with Docker you have to run the following command:
   ```bash
   docker compose exec pocket-id sh "sh scripts/create-one-time-access-token.sh <username or email>"
   ```

## Contribute

You're very welcome to contribute to Pocket ID! Please follow the [contribution guide](/CONTRIBUTING.md) to get started.
