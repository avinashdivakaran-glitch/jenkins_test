package main

import (
	"main/wifi_models"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var TNN_wifiStatus = wifi_models.WiFiStatus{
	AdapterMsg: wifi_models.AdapterMsg{
		HeaderID: "adapter",
	},
	HotspotMsg: wifi_models.HotspotMsg{
		HeaderID: "hotspot",
	},
	ScanMsg: wifi_models.ScanMsg{
		HeaderID: "scan",
	},
	ConnectMsg: wifi_models.ConnectMsg{
		HeaderID: "connect",
	},
	DisconnectMsg: wifi_models.DisconnectMsg{
		HeaderID: "disconnect",
	},
}
var TNN_wifiProperties = wifi_models.PropertiesMsg{HeaderID: "wifi properties"}
var TNN_wifiScanDevice = wifi_models.ScanDeviceList{HeaderID: "scan device"}

func init() {

	logrus.SetFormatter(&logrus.TextFormatter{})

	logrus.SetOutput(os.Stdout)

	level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)
}

func main() {
	logrus.Info("Starting WiFi MQTT serice...")

	mqttClient := Init_mqttClient()
	if mqttClient == nil {
		logrus.Fatal("Failed to initialize MQTT client. Exiting.")
	}

	resultsChan := make(chan []wifi_models.DeviceInfo)
	var stopChan chan struct{}
	errChan := make(chan error)

	TNN_WiFiconfig(mqttClient)
	TNN_wifiStatusFnuc(mqttClient)

	Subscribe_mqttBroker(mqttClient, &TNN_wifiStatus)

	Turn_tnnWiFi_on(false)

	for {
		TNN_wifiStatusFnuc(mqttClient)

		TNN_adapterFunc(mqttClient)

		TNN_DeleteFunc(mqttClient)

		TNN_DisconnectFunc(mqttClient)

		TNN_ConnectFunc(mqttClient)

		TNN_ScanFunc(mqttClient, resultsChan, stopChan, errChan)

		TNN_HotspotFunc(mqttClient)

		time.Sleep(1 * time.Second)
	}

}
