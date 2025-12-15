#!/bin/bash
set -e

APP_PATH="/opt/tnn-backend/bundles"

echo "1. Loading OCI images into Podman storage..."

# Pull from the local installed OCI bundles into system storage
podman pull oci:${APP_PATH}/${SERVICES_BLUETOOTH}
podman pull oci:${APP_PATH}/${SERVICES_SENSORS}
podman pull oci:${APP_PATH}/${SERVICES_WIFI}
podman pull oci:${APP_PATH}/${SERVICES_MQTT}


echo "2. Triggering Quadlet Generation..."

# scans /etc/containers/systemd/ and creates /run/systemd/generator/service-mqtt.service
systemctl daemon-reload

echo "3. Starting Services..."
systemctl enable --now service-mqtt
# systemctl enable --now service-bluetooth
# systemctl enable --now service-wifi
# systemctl enable --now service-blesensors

echo "Installation Complete!"
exit 0