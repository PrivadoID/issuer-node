#!/bin/sh
# Runtime entrypoint for EKS — the app was already built at Docker build time.
# Only nginx configuration is handled here.
set -eu

/bin/sh /app/scripts/configure-nginx.sh /app/deployment

exec nginx -g 'daemon off;'
