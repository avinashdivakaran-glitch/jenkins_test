

# Stage 1: Build Go binary
FROM golang:1.25.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOARCH=arm64 GOOS=linux go build -o health_monitor .



# Stage 2: Final Image
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/health_monitor /usr/local/bin/
COPY certs /app/certs/

CMD [ "health_monitor" ]



# 1. Build img
# sudo podman build --platform linux/arm64 -t service-blesensors:v1.0.0 .
# 2. Build and run container, map to podman local location
# sudo podman run -d --name sensor_server  --net=host  --privileged  -v /var/run/dbus:/var/run/dbus debian-sensor
