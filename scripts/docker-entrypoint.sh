
echo "Starting frontend..."
node frontend/build &

echo "Starting backend..."
cd backend && ./pocket-id-backend &

echo "Starting Caddy..."
caddy start --config /etc/caddy/Caddyfile &

wait