package ble_mqtt

import (
	"ble_fn_mqtt/ble_models"
	"encoding/json"
	"sync"

	"github.com/sirupsen/logrus"
)

var mu sync.Mutex

func Handle_adapterSub(tnn_bleStatus *ble_models.BluetoothStatus, payload []byte) {
	mu.Lock()
	defer mu.Unlock()

	var data ble_models.AdapterMsg
	err := json.Unmarshal(payload, &data)
	if err != nil {
		logrus.Warningf("Error unmarshaling JSON: %s\n", err)
		return
	}

	tnn_bleStatus.AdapterMsg = data
	logrus.Infof("Received Handle_adapterSub : %+v\n", data)
}

func Handle_discoverSub(tnn_bleStatus *ble_models.BluetoothStatus, payload []byte) {
	mu.Lock()
	defer mu.Unlock()

	var data ble_models.DiscoverMsg
	err := json.Unmarshal(payload, &data)
	if err != nil {
		logrus.Warningf("Error unmarshaling JSON: %s\n", err)
		return
	}

	tnn_bleStatus.DiscoverMsg = data
	logrus.Infof("Received Handle_discoverSub : %+v\n", data)
}

func Handle_scanSub(tnn_bleStatus *ble_models.BluetoothStatus, payload []byte) {
	mu.Lock()
	defer mu.Unlock()

	var data ble_models.ScanMsg
	err := json.Unmarshal(payload, &data)
	if err != nil {
		logrus.Warningf("Error unmarshaling JSON: %s\n", err)
		return
	}

	tnn_bleStatus.ScanMsg = data
	logrus.Infof("Received Handle_scanSub : %+v\n", data)
}

func Handle_pairSub(tnn_bleStatus *ble_models.BluetoothStatus, payload []byte) {
	mu.Lock()
	defer mu.Unlock()

	var data ble_models.PairMsg
	err := json.Unmarshal(payload, &data)
	if err != nil {
		logrus.Warningf("Error unmarshaling JSON: %s\n", err)
		return
	}

	tnn_bleStatus.PairMsg = data
	logrus.Infof("Received Handle_pairSub : %+v\n", data)
}

func Handle_connectSub(tnn_bleStatus *ble_models.BluetoothStatus, payload []byte) {
	mu.Lock()
	defer mu.Unlock()

	var data ble_models.ConnectMsg
	err := json.Unmarshal(payload, &data)
	if err != nil {
		logrus.Warningf("Error unmarshaling JSON: %s\n", err)
		return
	}

	tnn_bleStatus.ConnectMsg = data
	logrus.Infof("Received Handle_connectSub : %+v\n", data)
}
