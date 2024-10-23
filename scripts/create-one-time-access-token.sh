# Default database path
DB_PATH="./backend/data/pocket-id.db"

# Parse command-line arguments for the -d flag (database path)
while getopts ":d:" opt; do
    case $opt in
    d)
        DB_PATH="$OPTARG"
        ;;
    \?)
        echo "Invalid option -$OPTARG" >&2
        exit 1
        ;;
    esac
done

shift $((OPTIND - 1))

# Ensure username or email is provided as a parameter
if [ -z "$1" ]; then
    echo "Usage: $0 [-d <database_path>] <username or email>"
    echo "  -d   Specify the database path (optional, defaults to ./backend/data/pocket-id.db)"
    exit 1
fi

USER_IDENTIFIER="$1"

# Check and try to install the required commands
check_and_install() {
    local cmd=$1
    local pkg=$2

    if ! command -v "$cmd" &>/dev/null; then
        if command -v apk &>/dev/null; then
            echo "$cmd not found. Installing..."
            apk add "$pkg" --no-cache
        else
            echo "$cmd is not installed, please install it manually."
            exit 1
        fi
    fi
}

check_and_install sqlite3 sqlite
check_and_install uuidgen uuidgen

# Generate a 16-character alphanumeric secret token
SECRET_TOKEN=$(LC_ALL=C tr -dc 'A-Za-z0-9' </dev/urandom | head -c 16)

# Get the current Unix timestamp for creation and expiration (1 hour from now)
CREATED_AT=$(date +%s)
EXPIRES_AT=$((CREATED_AT + 3600))

# Retrieve user_id from the users table based on username or email
USER_ID=$(sqlite3 "$DB_PATH" "SELECT id FROM users WHERE username='$USER_IDENTIFIER' OR email='$USER_IDENTIFIER';")

# Check if user exists
if [ -z "$USER_ID" ]; then
    echo "User not found for username/email: $USER_IDENTIFIER"
    exit 1
fi

# Insert the one-time token into the one_time_access_tokens table
sqlite3 "$DB_PATH" <<EOF
INSERT INTO one_time_access_tokens (id, created_at, token, expires_at, user_id)
VALUES ('$(uuidgen)', '$CREATED_AT', '$SECRET_TOKEN', '$EXPIRES_AT', '$USER_ID');
EOF

if [ $? -eq 0 ]; then
    echo "A one-time access token valid for 1 hour has been created for \"$USER_IDENTIFIER\"."
    echo "Use the following URL to sign in once: ${PUBLIC_APP_URL:=https://<your-pocket-id-domain>}/login/$SECRET_TOKEN"
else
    echo "Error creating access token."
    exit 1
fi
