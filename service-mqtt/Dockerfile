

FROM debian:bookworm-slim

# 1. Update and install necessary packages
RUN apt-get update && \
    apt-get install -y mosquitto mosquitto-clients && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# 2. Copy Mosquitto config into container
COPY mosquitto.conf /etc/mosquitto/

# 3. Copy our new startup script
COPY start_mosquitto.sh /start_mosquitto.sh

# 4. Make script executable
RUN chmod +x /start_mosquitto.sh

# 5. Run the script directly
CMD ["/start_mosquitto.sh"]


# 1. Build img
# sudo podman build --platform linux/arm64 -t service_mqtt:v1.0.0 .
# 2. Build and run container
# sudo podman run -d --name mosquitto_con  --net=host  --privileged   debian-mqtt