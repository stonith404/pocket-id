
echo "Creating user and group..."

PUID=${PUID:-1000}
PGID=${PGID:-1000}

# Check if the group with PGID exists; if not, create it
if ! getent group appgroup > /dev/null 2>&1; then
    addgroup -g "$PGID" appgroup
fi

# Check if a user with PUID exists; if not, create it
if ! id -u appuser > /dev/null 2>&1; then
    if ! getent passwd "$PUID" > /dev/null 2>&1; then
        adduser -D -u "$PUID" -G appgroup -s /bin/sh appuser
    else
        # If a user with the PUID already exists, use that user
        existing_user=$(getent passwd "$PUID" | cut -d: -f1)
        echo "Using existing user: $existing_user"
    fi
fi

# Change ownership of the /app directory
chown -R "$PUID:$PGID" /app

# Switch to the non-root user
exec su-exec "$PUID:$PGID" "$@"