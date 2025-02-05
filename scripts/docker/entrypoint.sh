echo "Starting frontend..."
node frontend/build &

echo "Starting backend..."
cd backend && ./pocket-id-backend &

if [ "$CADDY_DISABLED" != "true" ]; then
  echo "Starting Caddy..."

  # Check if TRUST_PROXY is set to true and use the appropriate Caddyfile
  if [ "$TRUST_PROXY" = "true" ]; then
    caddy start --adapter caddyfile --config /etc/caddy/Caddyfile.trust-proxy &
  else
    caddy start --adapter caddyfile --config /etc/caddy/Caddyfile &
  fi
else
  echo "Caddy is disabled. Skipping..."
fi

wait
