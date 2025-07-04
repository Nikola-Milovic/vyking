#!/usr/bin/env bash

set -e

until docker inspect vyking-mysql --format='{{.State.Health.Status}}' | grep -q "healthy"; do
  echo "Waiting for MySQL to be ready..."
  sleep 1
done

echo "MySQL is ready"