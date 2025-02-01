---
id: fly-io
---

# Deploying Pocket ID on Fly.io

This guide walks you through deploying Pocket ID on Fly.io using the Hobby Plan. You'll set up an account, install the CLI, configure your project, and deploy the app.


---

🚀 Getting Started

1️⃣ Sign Up for [Fly.io](https://fly.io/)

Create an account on Fly.io and choose the Hobby Plan. This plan includes:
✅ 3 instances
✅ 256MB RAM per instance
✅ Shared CPU
✅ Up to 3GB of storage (minimum 1GB per volume)


---

2️⃣ Install Fly.io CLI

Run the following script on your local machine to install the Fly CLI:

```bash
curl -L https://na01.safelinks.protection.outlook.com/?url=https%3A%2F%2Ffly.io%2Finstall.sh&data=05%7C02%7C%7C8d0fc5ee609b4c3b2a6908dd40a9d08b%7C84df9e7fe9f640afb435aaaaaaaaaaaa%7C1%7C0%7C638737820269863137%7CUnknown%7CTWFpbGZsb3d8eyJFbXB0eU1hcGkiOnRydWUsIlYiOiIwLjAuMDAwMCIsIlAiOiJXaW4zMiIsIkFOIjoiTWFpbCIsIldUIjoyfQ%3D%3D%7C0%7C%7C%7C&sdata=2vneeC0lDAhVceP5ZWvPFQlMVIPdkzYxbyRN8gj%2Fqus%3D&reserved=0 | sh
```

After installation, update your shell environment:

```bash
export PATH="$HOME/.fly/bin:$PATH"
source ~/.bashrc  # or source ~/.zshrc
```

Verify the installation:

```bash
fly version
```

---

3️⃣ Authenticate

Login to your Fly.io account:

```bash
fly auth login
```

---

4️⃣ Configure Pocket ID

Create a fly.toml file in your project directory with the following configuration:

**You may need to modify some variables depending on region and other settings**

```toml
app = 'pocket-id'
primary_region = 'yyz'

[build]
  image = 'stonith404/pocket-id'

[env]
  PGID = '0'
  PUBLIC_APP_URL = 'https://pocketid.example.com'
  PUID = '0'
  TRUST_PROXY = 'false'

[[mounts]]
  source = 'data'
  destination = '/app/backend/data'

[http_service]
  internal_port = 80
  force_https = true
  auto_stop_machines = 'stop'
  auto_start_machines = true
  min_machines_running = 0
  processes = ['app']

[[services]]
  protocol = ''
  internal_port = 80

  [[services.ports]]
    port = 80
    handlers = ['http']

  [[services.ports]]
    port = 443
    handlers = ['tls', 'http']

[[vm]]
  size = 'shared-cpu-1x'

[experimental]
  entrypoint = ["sh", "/app/scripts/docker/create-user.sh"]
  cmd = ["sh", "/app/scripts/docker/entrypoint.sh"]
```

---

5️⃣ Launch Pocket ID

Run the following command to initialize the app:

```bash
fly launch
```

When prompted, confirm using the existing fly.toml file.


---

6️⃣ Adjust Configuration via the Web Dashboard

After running fly launch, you'll see a link to the Fly.io configuration page in the output.

1. Open that link in a browser


2. Adjust the settings as needed (e.g., confirm 256MB RAM and Shared CPU)


3. Save changes



To apply the updates, redeploy:

```bash
fly deploy
```

---

7️⃣ Set Up a Custom Domain (Optional)

If you want a custom domain:

1. On the Fly.io dashboard, configure your domain


2. Add a CNAME record in your DNS provider to verify ownership



Once verified, you can access your app via:

✅ Fly.io default URL: pocket-id.fly.dev

✅ Your custom domain: https://pocketid.example.com


---

🎉 Done!

Your Pocket ID instance is now running on Fly.io! 🚀