---
id: installation
---

# Installation

# Before you start

Pocket ID requires a [secure context](https://developer.mozilla.org/en-US/docs/Web/Security/Secure_Contexts), meaning it must be served over HTTPS. This is necessary because Pocket ID uses the [WebAuthn API](https://developer.mozilla.org/en-US/docs/Web/API/Web_Authentication_API).

### Installation with Docker (recommended)

1. Download the `docker-compose.yml` and `.env` file:

   ```bash
    curl -O https://raw.githubusercontent.com/stonith404/pocket-id/main/docker-compose.yml

    curl -o .env https://raw.githubusercontent.com/stonith404/pocket-id/main/.env.example
   ```

2. Edit the `.env` file so that it fits your needs. See the [environment variables](/configuration/environment-variables) section for more information.
3. Run `docker compose up -d`

You can now sign in with the admin account on `http://localhost/login/setup`.

### Install on Proxmox using Helper Scripts

Run the following script as root in your proxmox shell. 

See [Here](https://community-scripts.github.io/ProxmoxVE/scripts?id=pocketid) for more information.

**Configuration Paths**
- /opt/pocket-id/backend/.env
- /opt/pocket-id/frontend/.env

```bash
bash -c "$(wget -qLO - https://github.com/community-scripts/ProxmoxVE/raw/main/ct/pocketid.sh)"
```

### Unraid

Pocket ID is available as a template on the Community Apps store.

### Stand-alone Installation

Required tools:

- [Node.js](https://nodejs.org/en/download/) >= 22
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
