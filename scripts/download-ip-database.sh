#!/bin/bash

# Check if the license key environment variable is set
if [ -z "$MAXMIND_LICENSE_KEY" ]; then
  echo "Error: MAXMIND_LICENSE_KEY environment variable is not set."
  echo "Please set it using 'export MAXMIND_LICENSE_KEY=your_license_key' and try again."
  exit 1
fi
echo $MAXMIND_LICENSE_KEY
# GeoLite2 City Database URL
URL="https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=${MAXMIND_LICENSE_KEY}&suffix=tar.gz"

# Download directory
DOWNLOAD_DIR="./geolite2_db"
TARGET_PATH=./backend/GeoLite2-City.mmdb
mkdir -p $DOWNLOAD_DIR

# Download the database
echo "Downloading GeoLite2 City database..."
curl -L -o "$DOWNLOAD_DIR/GeoLite2-City.tar.gz" "$URL"

# Extract the downloaded file
echo "Extracting GeoLite2 City database..."
tar -xzf "$DOWNLOAD_DIR/GeoLite2-City.tar.gz" -C $DOWNLOAD_DIR --strip-components=1

mv "$DOWNLOAD_DIR/GeoLite2-City.mmdb" $TARGET_PATH

# Clean up
rm -rf "$DOWNLOAD_DIR"

echo "GeoLite2 City database downloaded and extracted to $TARGET_PATH"