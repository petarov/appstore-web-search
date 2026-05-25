#!/bin/bash

if [ -z "$1" ]; then
  echo "error: ❌ specify target server hostname"
  exit 1
fi
if [ -z "$2" ] || [ ! -f "$2" ]; then
  echo "error: ❌ invalid or missing path to .ssh/id_<name> file"
  exit 1
fi
if [ -z "$3" ]; then
  echo "error: ❌ specify server user"
  exit 1
fi

# === CONFIG ===
KEY_PATH="$2"
RSYNC="/opt/homebrew/bin/rsync"

REMOTE_USER="$3"
REMOTE_HOST="$1"
REMOTE_DIR="/srv/asws"

if ssh -o BatchMode=yes -o ConnectTimeout=3 "$REMOTE_HOST" 'exit' 2>/dev/null; then
  echo "✅ Host $1 is reachable"
else
  echo "error: ❌ host $1 is NOT reachable"
  exit 1
fi

# === PROJECT DIRS ===
PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"

# === RSYNC OPTIONS ===
RSYNC_OPTS="-az --chown=www-data:www-data --progress --delete"
SSH_OPTS="ssh -i $KEY_PATH -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null"

echo "$RSYNC_OPTS -e $SSH_OPTS"

# === DEPLOY ===

echo "Syncing /dist..."
$RSYNC $RSYNC_OPTS -e "$SSH_OPTS" "$PROJECT_ROOT/dist/" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/"

echo "✅ Deployment complete"