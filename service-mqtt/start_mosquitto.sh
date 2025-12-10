#!/bin/bash

# 1. Setup Directories & Permissions
mkdir -p /var/log/mosquitto /var/lib/mosquitto
chown -R mosquitto:mosquitto /var/log/mosquitto /var/lib/mosquitto

echo "Starting Mosquitto Broker..."

# 2. Start Mosquitto
# 'exec' replaces the shell process with the mosquitto process.
# This ensures Mosquitto becomes PID 1 and receives shutdown signals correctly.
exec stdbuf -oL /usr/sbin/mosquitto -c /etc/mosquitto/mosquitto.conf