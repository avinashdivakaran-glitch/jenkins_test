package main

import (
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
	tcp_broker = "tcp://localhost:1884"
	ssl_broker = "ssl://localhost:8884"
	clientID   = "DEVICE_01_sensor"
	username   = ""
	password   = ""
)

// Paths to TLS certificates
var (
	caCertFile     = "app/certs/ca.crt"
	clientCertFile = "app/certs/tnn-ble_sensor.crt"
	clientKeyFile  = "app/certs/tnn-ble_sensor.key"
)

var (
	topics = map[string]string{
		"service":     "tnn_server/body_sensor/service",
		"device":      "tnn_server/body_sensor/device",
		"temperature": "tnn_server/body_sensor/data/temperature",
		"heartrate":   "tnn_server/body_sensor/data/heartrate",
		"bodyoxygen":  "tnn_server/body_sensor/data/bodyoxygen",
	}
)

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

// Test TLS coneection
func testTLSConnection(brokerURL string, config *tls.Config) error {
	u, err := url.Parse(brokerURL)
	if err != nil {
		return fmt.Errorf("invalid broker URL for TLS test; %v", err)
	}

	conn, err := tls.Dial("tcp", u.Host, config)
	if err != nil {
		return fmt.Errorf("TLS dial test error: %v", err)
	}
	conn.Close()

	return nil
}

// init mqtt
func Init_mqttClient() (client mqtt.Client) {
	var broker string
	var tlsConfig *tls.Config

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

	// MQTT client options
	opts := mqtt.NewClientOptions()

	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	if cert {
		opts.SetTLSConfig(tlsConfig)
	}

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

// Publish sensor data
func Publish_mqttBroker(client mqtt.Client, topic string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Error marshalling sensor data: %v", err)
	}

	token := client.Publish(topic, 0, false, payload)
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("Error publishing to %s: %v", topic, token.Error())
	}
	logrus.Infof("🚀 Published to %s: %v\n", topic, data)

	return nil
}
