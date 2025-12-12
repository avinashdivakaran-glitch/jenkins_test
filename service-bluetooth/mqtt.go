package main

import (
	"ble_fn_mqtt/ble_models"
	"ble_fn_mqtt/ble_mqtt"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

// MQTT broker settings
var (
	cert       = true
	broker     = ""
	tcp_broker = "tcp://localhost:1884"
	ssl_broker = "ssl://localhost:8884"
	clientID   = "TNN_bluetoothService"
	username   = ""
	password   = ""
)

// Paths to TLS certificates
var (
	caCertFile     = "/app/certs/ca.crt"
	clientCertFile = "/app/certs/ble_app.crt"
	clientKeyFile  = "/app/certs/ble_app.key"
)

var (
	topics = map[string]string{
		"service":  "tnn_server/bluetooth/service",
		"adapter":  "tnn_server/bluetooth/adapter",
		"discover": "tnn_server/bluetooth/func/discover",
		"scan":     "tnn_server/bluetooth/func/scan",
		"pair":     "tnn_server/bluetooth/func/pair",
		"connect":  "tnn_server/bluetooth/func/connect",
		"scan_dev": "tnn_server/bluetooth/device/scanned",
		"save_dev": "tnn_server/bluetooth/device/paired",
	}
)

func init() {
	if b := os.Getenv("TNN_MQTT_BROKER"); b != "" {
		broker = b
	}
	if id := os.Getenv("TNN_MQTT_CLIENT_ID"); id != "" {
		clientID = id
	}
	if u := os.Getenv("TNN_MQTT_USERNAME"); u != "" {
		username = u
	}
	if p := os.Getenv("TNN_MQTT_PASSWORD"); p != "" {
		password = p
	}
}

// Generate TLS config for mutual TLS
func newTLSConfig(caFile, certFile, keyFile string) (*tls.Config, error) {
	// Load CA certificate
	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to read CA file: %v", err)
	}
	caPool := x509.NewCertPool()
	if ok := caPool.AppendCertsFromPEM(caCert); !ok {
		return nil, fmt.Errorf("Failed to append CA certificate")
	}

	// Load client certificate and key
	clientCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to load client certificate/key: %v", err)
	}

	return &tls.Config{
		Certificates:       []tls.Certificate{clientCert},
		RootCAs:            caPool,
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: false,
	}, nil
}

// testTLSConnection performs a "pre-flight check" to test the TLS connection.
func testTLSConnection(brokerURL string, config *tls.Config) error {
	// just the host and port, e.g., "192.168.1.34:8883"
	u, err := url.Parse(brokerURL)
	if err != nil {
		return fmt.Errorf("invalid broker URL for TLS test: %v", err)
	}

	// u.Host already contains "hostname:port"
	conn, err := tls.Dial("tcp", u.Host, config)
	if err != nil {
		return fmt.Errorf("TLS dial test error: %v", err)
	}
	conn.Close()

	return nil
}

func Init_mqttClient() (client mqtt.Client) {
	if !cert {
		broker = tcp_broker
	} else {
		broker = ssl_broker

		// 1. Create the TLS config
		tlsConfig, err := newTLSConfig(caCertFile, clientCertFile, clientKeyFile)
		if err != nil {
			logrus.Fatalf("Failed to create TLS config: %v", err)
		}

		// 2. set ServerName
		brokerURL, err := url.Parse(broker)
		if err != nil {
			logrus.Fatalf("Failed to parse broker URL '%s': %v", broker, err)
		}
		hostname, _, err := net.SplitHostPort(brokerURL.Host)
		if err != nil {
			hostname = brokerURL.Host
		}
		tlsConfig.ServerName = hostname
		logrus.Infof("Setting TLS ServerName to: %s", hostname)

		// 3. TLS config test
		logrus.Info("Running TLS connection pre-flight check...")
		if err := testTLSConnection(broker, tlsConfig); err != nil {
			// If this fails, we stop.
			logrus.Fatalf("TLS connection pre-flight check FAILED: %v", err)
		}
		logrus.Info("TLS pre-flight check succeeded.")
	}

	// 4. MQTT client options
	opts := mqtt.NewClientOptions()

	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	// opts.SetTLSConfig(tlsConfig)

	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetAutoReconnect(true)

	if username != "" {
		opts.SetUsername(username)
		opts.SetPassword(password)
	}

	// Set the last will and testament
	opts.SetWill(topics["service"], "offline", 0, true)

	// Optional: connection lost handler
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		logrus.Warnf("⚠️ MQTT connection lost: %v", err)
	}

	// Optional: connection handler
	opts.OnConnect = func(client mqtt.Client) {
		logrus.Info("✅ Connected to MQTT broker")

		token := client.Publish(topics["service"], 0, true, []byte("online"))
		token.Wait()
		logrus.Info("Published 'online' status")
	}

	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logrus.Warnf("Error connecting to MQTT broker: %v", token.Error())
	}

	return client

}

func Publish_mqttBroker(client mqtt.Client, topic string, data interface{}) error {
	// Serialize the data to JSON
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling data: %v", err)
	}

	// Publish the data to the given topic
	token := client.Publish(topic, 0, false, payload)
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("error publishing message: %v", token.Error())
	}
	logrus.Infof("Publish data on %s : %v\n", topic, data)

	return nil
}

func Subscribe_mqttBroker(mqttClient mqtt.Client, tnn_bleStatus *ble_models.BluetoothStatus) {
	for _, topic := range topics {
		mqttClient.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
			handleMessage(msg, tnn_bleStatus)
		})
	}
}

func handleMessage(msg mqtt.Message, tnn_bleStatus *ble_models.BluetoothStatus) {
	topic := msg.Topic()

	switch topic {
	case topics["adapter"]:
		ble_mqtt.Handle_adapterSub(tnn_bleStatus, msg.Payload())
	case topics["discover"]:
		ble_mqtt.Handle_discoverSub(tnn_bleStatus, msg.Payload())
	case topics["scan"]:
		ble_mqtt.Handle_scanSub(tnn_bleStatus, msg.Payload())
	case topics["pair"]:
		ble_mqtt.Handle_pairSub(tnn_bleStatus, msg.Payload())
	case topics["connect"]:
		ble_mqtt.Handle_connectSub(tnn_bleStatus, msg.Payload())
	}
}
