#!/bin/bash
set -e

SERVICES_BLUETOOTH="service-bluetooth"
SERVICES_SENSORS="service-blesensors"
SERVICES_WIFI="service-wifi"
SERVICES_MQTT="service-mqtt"


echo "1. Loading OCI images into Podman storage..."

# Pull from the local installed OCI bundles into system storage
podman pull "oci:${SERVICES_MQTT}"

podman pull "oci:${SERVICES_BLUETOOTH}"
podman pull "oci:${SERVICES_SENSORS}"

# podman pull "oci:${SERVICES_WIFI}"




echo "2. Reloading systemd daemon ..."

# scans /etc/containers/systemd/ and creates /run/systemd/generator/service-mqtt.service
systemctl daemon-reload



echo "3. Starting Services..."

systemctl enable --now "${SERVICES_MQTT}"

systemctl enable --now "${SERVICES_BLUETOOTH}"
systemctl enable --now "${SERVICES_SENSORS}"

# systemctl enable --now "${SERVICES_WIFI}"



echo "Installation Complete!"
exit 0