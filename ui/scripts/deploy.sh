#!/bin/sh
set -eu

# Configure nginx first: it fails fast on missing credentials, before the slow build.
/bin/sh /app/scripts/configure-nginx.sh /app/deployment

ENV_FILENAME="/app/.env"

# Mapping of issuer-level env variables to app-level ones. Everything written here is compiled into the public JS
# bundle, so never add secrets (the API credentials stay in nginx, see configure-nginx.sh).
{
  echo "VITE_ISSUER_LOGO=${ISSUER_ISSUER_LOGO:-}"
  echo "VITE_ISSUER_NAME=${ISSUER_ISSUER_NAME:-}"
  echo "VITE_BUILD_TAG=${ISSUER_UI_BUILD_TAG:-}"
  echo "VITE_WARNING_MESSAGE=${ISSUER_UI_WARNING_MESSAGE:-}"
  echo "VITE_IPFS_GATEWAY_URL=${ISSUER_UI_IPFS_GATEWAY_URL:-}"
  echo "VITE_SCHEMA_EXPLORER_AND_BUILDER_URL=${ISSUER_UI_SCHEMA_EXPLORER_AND_BUILDER_URL:-}"
  echo "VITE_DISPLAY_METHOD_BUILDER_URL=${ISSUER_UI_DISPLAY_METHOD_BUILDER_URL:-}"
  echo "VITE_BASE_URL=${ISSUER_UI_BASE_URL:-}"
  echo "VITE_PAYMENT_PAGES_ENABLED=${ISSUER_UI_PAYMENT_PAGES_ENABLED:-}"
} > "$ENV_FILENAME"

# Build app
cd /app && npm run build

# Copy app dist
cp -r /app/dist/. /usr/share/nginx/html

# Delete build files
rm -rf /app/dist

# Run nginx
exec nginx -g 'daemon off;'
