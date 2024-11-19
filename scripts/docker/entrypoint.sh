
echo "Starting frontend..."
node frontend/build &

echo "Starting backend..."
cd backend && ./pocket-id-backend &

echo "Starting Caddy..."

# Check if TRUST_PROXY is set to true and use the appropriate Caddyfile
if [ "$TRUST_PROXY" = "true" ]; then
  caddy start --config /etc/caddy/Caddyfile.trust-proxy &
else
  caddy start --config /etc/caddy/Caddyfile &
fi

wait