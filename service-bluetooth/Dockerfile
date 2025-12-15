

# Stage 1: Build Go binary
FROM golang:1.25.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOARCH=arm64 GOOS=linux go build -o ble_fn_mqtt .



# Stage 2: Final Image
FROM docker.io/library/debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/ble_fn_mqtt /usr/local/bin/
COPY certs /app/certs/

CMD [ "ble_fn_mqtt" ]



# 1. Build img
# sudo podman build --platform linux/arm64 -t service-bluetooth:v1.0.0 .
# 2. Build and run container
# sudo podman run -d   --name bluez_server   --net=host   --privileged   -v /var/run/dbus:/var/run/dbus   service-bluetooth
