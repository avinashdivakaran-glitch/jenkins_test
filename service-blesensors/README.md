# 🔌 IoT Sensor Data Emulator

A lightweight simulator for generating random IoT sensor data (Temperature, Heart Rate, Oxygen) and publishing it via MQTT over TLS.  
Built with **Go (backend)** and **Flutter (frontend)**.


---

## 🚀 Features
- Generate random sensor data periodically.
- Publish via MQTT with TLS encryption.
- Frontend Flutter app subscribes and displays real-time data.
- Configurable broker settings.

---

## 🧰 Tech Stack
- **Backend:** Go (Golang)
- **Frontend:** Flutter
- **Protocol:** MQTT (with TLS)
- **Broker:** Mosquitto

---

## ⚙️ Installation

### Prerequisites
- Go 1.21+
- Flutter 3.x
- Mosquitto MQTT broker (with TLS setup)



## ⚙️ How It Works (Backend Script)


### 1. **Data Structure**

All sensor data is represented using a common structure:

```go
type SensorData struct {
    Available bool      `json:"available"`
    SensorID  string    `json:"sensor_id"`
    Timestamp time.Time `json:"timestamp"`
    Data      struct {
        Value float64 `json:"value"`
        Unit  string  `json:"unit"`
        State string  `json:"state"`
    } `json:"data"`
}
```

### 2. **Go Script Structure**
``` go
1. sensorReader() runs continuously in a goroutine, generating random sensor values and sending them to channels.
2. mqttPublisher() runs in another goroutine, listens to those channels, and publishes data to MQTT topics.
3. Both goroutines run concurrently using Go’s channels and select {}.
4. TLS ensures secure MQTT communication.
```

### 3. **Go const define**
```go
// MQTT broker settings
const (
	broker   = "ssl://192.168.4.107:8883"
	clientID = "DEVICE_01"
)

// Paths to TLS certificates
const (
	caCertFile     = "certs/ca.crt"
	clientCertFile = "certs/client.crt"
	clientKeyFile  = "certs/client.key"
)
```

### 4. **MQTT Topic and structure**
```
tnn/device_01/sensors/Temperature
{
    "available" :   true,
    "sensor_id" :   "T001",
    "timestamp" :   time.Now(),
    "data"      :   {
        "value"     :   28.5,
        "unit"      :   "C",
        "state"     "   "normal"
    }
}

tnn/device_01/sensors/HeartRate
{
    "available" :   true,
    "sensor_id" :   "H001",
    "timestamp" :   time.Now(),
    "data"      :   {
        "value"     :   78,
        "unit"      :   "bpm",
        "state"     "   "normal"
    }
}

tnn/device_01/sensors/BodyOxygen
{
    "available" :   true,
    "sensor_id" :   "O001",
    "timestamp" :   time.Now(),
    "data"      :   {
        "value"     :   98.2,
        "unit"      :   "%",
        "state"     "   "normal"
    }
}
```


## ⚙️ Certificate generation

- ca.crt → Certificate Authority (trusted by both server and client)
- broker.crt & broker.key → Broker certificate and private key
- client.crt & client.key → Client certificate and private key

### Step 1: Create a Certificate Authority (CA)
```bash
# 1. Generate private key for CA
openssl genrsa -out ca.key 2048
# 2. Create self-signed CA certificate
openssl req -x509 -new -nodes -key ca.key -sha256 -days 365 -out ca.crt -subj "/CN=TNN_CA"
```

### Step 2: Generate Broker Certificate
```bash
# 1. Generate broker private key
openssl genrsa -out broker.key 2048
# 2. Create a certificate signing request (CSR) for the broker
openssl req -new -key broker.key -out broker.csr -config broker.cnf
# 3. Sign broker CSR with CA to generate broker certificate
openssl x509 -req -in broker.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out broker.crt -days 365 -sha256 -extfile broker.cnf -extensions v3_req
```

### Step 3: Generate Client Certificate
```bash
# 1. Generate client private key
openssl genrsa -out client.key 2048
# 2. Create client CSR
openssl req -new -key client.key -out client.csr -subj "/CN=DEVICE_01"
# 3. Sign client CSR with CA
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 365 -sha256
# (extra) adding SAN to client
openssl req -new -key client.key -out client.csr -subj "/CN=DEVICE_01" -reqexts SAN -config <(cat /etc/ssl/openssl.cnf <(printf "\n[SAN]\nsubjectAltName=DNS:DEVICE_01,IP:192.168.4.107"))
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 365 -sha256 -extfile <(printf "subjectAltName=DNS:DEVICE_01,IP:192.168.4.107")
```

### Step 4: Verify Certificates
```bash
# Verify broker certificate
openssl verify -CAfile ca.crt broker.crt
# Verify client certificate
openssl verify -CAfile ca.crt client.crt
openssl x509 -in broker.crt -text -noout | grep -A1 "Subject Alternative Name"
```

### Step 5: Make sure all certificates are accessible by Mosquitto
```bash
# Make sure the key, cert, and CA are owned by Mosquitto
sudo chown mosquitto:mosquitto broker.key broker.crt ca.crt
# Ensure proper file permissions (readable only by owner)
sudo chmod 600 broker.key
sudo chmod 644 broker.crt ca.crt
```


## mosquitto configuration

```conf
# MQTT over TCP (MQTTS)
listener 8883
cafile certs/ca.crt
certfile certs/broker.crt
keyfile certs/broker.key
require_certificate true
allow_anonymous true
tls_version tlsv1.2

# Verbose logging
log_dest stdout
log_type error
log_type warning
log_type notice
log_type information
log_type debug
connection_messages true
```

## runing instructions
```bash
    sudo mosquitto -c mosquitto.conf
    go run health_monitor.go
```

## Reference
- 1. https://www.emqx.com/en/blog/how-to-use-mqtt-in-golang
- 2. https://medium.com/@somanathtv/mqttnet-broker-and-client-with-tls-openssl-in-c-d2b328416992
- 3. https://openest.io/non-classe-en/mqtts-how-to-use-mqtt-with-tls/