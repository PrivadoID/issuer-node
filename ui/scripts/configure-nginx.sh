#!/bin/sh
# Writes the nginx config of the Issuer Node UI. Used by deploy.sh and deploy-eks.sh at container start.
#
# The UI reaches the issuer API through nginx at /api/, and nginx adds the API credentials itself, so they never
# reach the browser. Required in both modes:
#   ISSUER_API_AUTH_USER, ISSUER_API_AUTH_PASSWORD  the issuer API credentials
#   ISSUER_UI_API_PROXY_URL                         where nginx sends API calls, e.g. http://api:3001
#                                                   (defaults to ISSUER_SERVER_URL)
#
# ISSUER_UI_INSECURE (default false):
#   false  The UI and /api/ require the UI login: ISSUER_UI_AUTH_USERNAME and ISSUER_UI_AUTH_PASSWORD are required.
#   true   No login at all. Anyone who can reach the UI can use the whole issuer API. PUBLIC DEMOS ONLY.
#
# Usage: configure-nginx.sh <directory with nginx.conf, nginx_insecure.conf and api-proxy.conf.template>
set -eu

DEPLOYMENT_DIR="${1:?usage: configure-nginx.sh <deployment dir>}"

fail() {
  echo "ERROR: $*" >&2
  exit 1
}

[ -n "${ISSUER_API_AUTH_USER:-}" ] || fail "ISSUER_API_AUTH_USER is required: the UI uses it to call the issuer API."
[ -n "${ISSUER_API_AUTH_PASSWORD:-}" ] || fail "ISSUER_API_AUTH_PASSWORD is required: the UI uses it to call the issuer API."

ISSUER_UI_API_PROXY_URL="${ISSUER_UI_API_PROXY_URL:-${ISSUER_SERVER_URL:-}}"
[ -n "$ISSUER_UI_API_PROXY_URL" ] || fail "Set ISSUER_UI_API_PROXY_URL (or ISSUER_SERVER_URL) to the issuer API address."
case "$ISSUER_UI_API_PROXY_URL" in
  http://*|https://*) ;;
  *) fail "ISSUER_UI_API_PROXY_URL must start with http:// or https://, got: $ISSUER_UI_API_PROXY_URL" ;;
esac
# The template appends the slash itself.
while [ "${ISSUER_UI_API_PROXY_URL%/}" != "$ISSUER_UI_API_PROXY_URL" ]; do
  ISSUER_UI_API_PROXY_URL="${ISSUER_UI_API_PROXY_URL%/}"
done
echo "UI: proxying /api/ to $ISSUER_UI_API_PROXY_URL"

ISSUER_UI_API_BASIC_AUTH="$(printf '%s:%s' "$ISSUER_API_AUTH_USER" "$ISSUER_API_AUTH_PASSWORD" | base64 | tr -d '\n')"
export ISSUER_UI_API_PROXY_URL ISSUER_UI_API_BASIC_AUTH

API_PROXY_CONF=/etc/nginx/issuer-ui-api-proxy.conf
# The rendered file holds the API credentials.
(umask 077 && envsubst '${ISSUER_UI_API_PROXY_URL} ${ISSUER_UI_API_BASIC_AUTH}' \
  < "$DEPLOYMENT_DIR/api-proxy.conf.template" > "$API_PROXY_CONF")

if [ "${ISSUER_UI_INSECURE:-false}" = "true" ]; then
  cat >&2 <<'EOF'
#####################################################################################
#  WARNING: ISSUER_UI_INSECURE=true - THE UI HAS NO LOGIN. FOR PUBLIC DEMOS ONLY.
#
#  Anyone who can reach this UI can use the whole issuer API through /api/:
#  create identities, issue and revoke credentials, manage keys and payment options.
#  Never use this mode for an issuer with real users or real data.
#  Set ISSUER_UI_INSECURE=false and ISSUER_UI_AUTH_USERNAME / ISSUER_UI_AUTH_PASSWORD
#  to require a login.
#####################################################################################
EOF
  cp "$DEPLOYMENT_DIR/nginx_insecure.conf" /etc/nginx/conf.d/default.conf
else
  [ -n "${ISSUER_UI_AUTH_USERNAME:-}" ] || fail "ISSUER_UI_AUTH_USERNAME is required for the UI login."
  [ -n "${ISSUER_UI_AUTH_PASSWORD:-}" ] || fail "ISSUER_UI_AUTH_PASSWORD is required for the UI login."
  # -i reads the password from stdin, so it does not show up in the process list.
  printf '%s\n' "$ISSUER_UI_AUTH_PASSWORD" | htpasswd -c -i /etc/nginx/.htpasswd "$ISSUER_UI_AUTH_USERNAME"
  cp "$DEPLOYMENT_DIR/nginx.conf" /etc/nginx/conf.d/default.conf
fi
