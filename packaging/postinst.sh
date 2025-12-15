#!/bin/bash
set -e

APP_PATH="/opt/tnn_backend/bundles"
SERVICES_BLUETOOTH="service-bluetooth"
SERVICES_SENSORS="service-blesensors"
SERVICES_WIFI="service-wifi"
SERVICES_MQTT="service-mqtt"

echo "1. Loading OCI images into Podman storage..."

# Pull from the local installed OCI bundles into system storage
sudo podman pull "oci:${APP_PATH}/${SERVICES_BLUETOOTH}"
sudo podman pull "oci:${APP_PATH}/${SERVICES_SENSORS}"
sudo podman pull "oci:${APP_PATH}/${SERVICES_WIFI}"
sudo podman pull "oci:${APP_PATH}/${SERVICES_MQTT}"


echo "2. Reloading systemd daemon ..."

# scans /etc/containers/systemd/ and creates /run/systemd/generator/service-mqtt.service
systemctl daemon-reload

echo "3. Starting Services..."
systemctl enable --now "${SERVICES_MQTT}"
# systemctl enable --now service-bluetooth
# systemctl enable --now service-wifi
# systemctl enable --now service-blesensors

echo "Installation Complete!"
exit 0