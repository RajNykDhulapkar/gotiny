#!/bin/sh
set -e

echo "Starting GoTiny Service..."

check_required_vars() {
  if [ -z "$GOTINY_PORT" ]; then
    echo "ERROR: GOTINY_PORT not set"
    exit 1
  fi
  if [ -z "$REDIS_URL" ]; then
    echo "ERROR: REDIS_URL not set"
    exit 1
  fi
  if [ -z "$RANGE_ALLOCATOR_ADDRESS" ]; then
    echo "ERROR: RANGE_ALLOCATOR_ADDRESS not set"
    exit 1
  fi
  if [ -z "$SERVICE_ID" ]; then
    echo "ERROR: SERVICE_ID not set"
    exit 1
  fi
  if [ -z "$MONGODB_URI" ]; then
    echo "ERROR: MONGODB_URI not set"
    exit 1
  fi
  if [ -z "$MONGODB_DATABASE" ]; then
    echo "ERROR: MONGODB_DATABASE not set"
    exit 1
  fi
}

check_required_vars

echo "Configuration:"
echo "- HTTP Port: $GOTINY_PORT"
echo "- Redis URL: $REDIS_URL"
echo "- Range Allocator: $RANGE_ALLOCATOR_ADDRESS"
echo "- MongoDB URI: $MONGODB_URI"
echo "- Service ID: $SERVICE_ID"

exec /app/gotiny
