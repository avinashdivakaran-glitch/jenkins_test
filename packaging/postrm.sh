#!/bin/bash
set -e

SERVICES_BLUETOOTH="service-bluetooth"
SERVICES_SENSORS="service-blesensors"
SERVICES_WIFI="service-wifi"
SERVICES_MQTT="service-mqtt"


echo "1. Stop and disable systemd services"

systemctl stop "${SERVICES_MQTT}"
systemctl stop "${SERVICES_BLUETOOTH}"
systemctl stop "${SERVICES_SENSORS}"
# systemctl stop "${SERVICES_WIFI}"

systemctl disable "${SERVICES_MQTT}"
systemctl disable "${SERVICES_BLUETOOTH}"
systemctl disable "${SERVICES_SENSORS}"
# systemctl disable "${SERVICES_WIFI}"


echo "2. Removing Podman containers"

podman rm -f "${SERVICES_MQTT}" || true
podman rm -f "${SERVICES_BLUETOOTH}" || true
podman rm -f "${SERVICES_SENSORS}" || true
# podman rm -f "${SERVICES_WIFI}" || true


echo "3. Removing Podman images"

podman rmi -f "${SERVICES_MQTT}" || true
podman rmi -f "${SERVICES_BLUETOOTH}" || true
podman rmi -f "${SERVICES_SENSORS}" || true
# podman rmi -f "${SERVICES_WIFI}" || true


echo "4. Removing systemd service file"

rm -f /etc/systemd/system/"${SERVICES_MQTT}".service
rm -f /etc/systemd/system/"${SERVICES_BLUETOOTH}".service
rm -f /etc/systemd/system/"${SERVICES_SENSORS}".service
# rm -f /etc/systemd/system/"${SERVICES_WIFI}".service


systemctl daemon-reload


echo "5. Cleaning OCI bundles"

rm -rf /opt/tnn_backend


echo "Uninstallation Complete!"
exit 0