---
id: upgrading
---

# Upgrading

Updating to a New Version

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
   pm2 start --name pocket-id-frontend --node-args="--env-file .env" build/index.js

   # Optional: Start Caddy (You can use any other reverse proxy)
   cd ..
   pm2 start caddy --name pocket-id-caddy -- run --config reverse-proxy/Caddyfile
   ```
