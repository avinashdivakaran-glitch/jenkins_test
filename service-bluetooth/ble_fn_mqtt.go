package main

import (
	"ble_fn_mqtt/ble_models"
	"context"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var TNN_bleStatus = ble_models.BluetoothStatus{
	AdapterMsg: ble_models.AdapterMsg{
		HeaderID: "adapter",
	},
	DiscoverMsg: ble_models.DiscoverMsg{
		HeaderID: "func_discovery",
	},
	ScanMsg: ble_models.ScanMsg{
		HeaderID: "func_scan",
	},
	PairMsg: ble_models.PairMsg{
		HeaderID: "func_pair",
	},
	ConnectMsg: ble_models.ConnectMsg{
		HeaderID: "func_connect",
	},
}

var TNN_bleBondDevice = ble_models.PairedDeviceList{
	HeaderID: "paired_devices",
}
var TNN_bleScanDevice = ble_models.ScanDeviceList{
	HeaderID: "scaned_devices",
}

var scan_option bool

func init() {
	// Set logrus to output JSON.
	// This is the standard for services.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Send logs to Stdout
	// systemd will capture this.
	logrus.SetOutput(os.Stdout)

	// Set the logging level.
	// Read from an env var, but default to 'info'.
	level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)
}

func main() {
	logrus.Info("Starting BLE MQTT Service...")

	mqttClient := Init_mqttClient()
	if mqttClient == nil {
		// Init_mqttClient should log its own error
		logrus.Fatal("Failed to initialize MQTT client. Exiting.")
	}

	time.Sleep(2 * time.Second)

	TNN_BLEconfig(mqttClient)
	TNN_BLEStateFunc(mqttClient)

	Subscribe_mqttBroker(mqttClient, &TNN_bleStatus)

	// Create a cancellable context to control scanning
	scanCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	scanStatusCh := make(chan bool)
	scan_option = true
	go Scanning_tnnBLE(scanCtx, &TNN_bleScanDevice, &TNN_bleStatus, scanStatusCh)

	realTime_ticker := time.NewTicker(1 * time.Second)
	defer realTime_ticker.Stop()

	for range realTime_ticker.C {

		TNN_BLEStateFunc(mqttClient)

		TNN_adapterFunc(mqttClient)

		TNN_discoverFunc(mqttClient)

		TNN_ScanFunc(mqttClient, scanStatusCh)

		TNN_PairFunc(mqttClient)

		TNN_RemoveFunc(mqttClient)

		TNN_ConnectFunc(mqttClient)

		TNN_DisconnectFunc(mqttClient)

	}

}
