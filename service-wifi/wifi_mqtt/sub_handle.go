package wifi_mqtt

import (
	"encoding/json"
	"main/wifi_models"
	"sync"

	"github.com/sirupsen/logrus"
)

var mu sync.Mutex

func Handle_adapterSub(tnn_wifiStatus *wifi_models.WiFiStatus, payload []byte) {
	mu.Lock()
	defer mu.Unlock()

	var data wifi_models.AdapterMsg
	err := json.Unmarshal(payload, &data)
	if err != nil {
		logrus.Warningf("Error unmarshaling JSON: %s\n", err)
		return
	}

	tnn_wifiStatus.AdapterMsg.SetPower = data.SetPower
	logrus.Infof("Received Handle_adapterSub : %+v\n", data)
}

// func Handle_propertiesSub(tnn_wifiStatus chan wifi_models.WiFiStatus, payload []byte) {
// 	var data wifi_models.PropertiesMsg
// 	err := json.Unmarshal(payload, &data)
// 	if err != nil {
// 		logrus.Warningf("Error unmarshaling JSON: %s\n", err)
// 		return
// 	}

// 	tnn_wifiStatus <- wifi_models.WiFiStatus{
// 		PropertiesMsg: data,
// 	}
// 	logrus.Infof("Received Handle_propertiesSub : %+v\n", data)
// }

func Handle_scanSub(tnn_wifiStatus *wifi_models.WiFiStatus, payload []byte) {
	mu.Lock()
	defer mu.Unlock()

	var data wifi_models.ScanMsg
	err := json.Unmarshal(payload, &data)
	if err != nil {
		logrus.Warningf("Error unmarshaling JSON: %s\n", err)
		return
	}

	tnn_wifiStatus.ScanMsg.SetScan = data.SetScan
	logrus.Infof("Received Handle_ScanSub : %+v\n", data)
}

func Handle_connectSub(tnn_wifiStatus *wifi_models.WiFiStatus, payload []byte) {
	mu.Lock()
	defer mu.Unlock()

	var data wifi_models.ConnectMsg
	err := json.Unmarshal(payload, &data)
	if err != nil {
		logrus.Warningf("Error unmarshaling JSON: %s\n", err)
		return
	}

	tnn_wifiStatus.ConnectMsg.SetConnect = data.SetConnect
	tnn_wifiStatus.ConnectMsg.Data.SSID = data.Data.SSID
	tnn_wifiStatus.ConnectMsg.Data.Password = data.Data.Password
	tnn_wifiStatus.ConnectMsg.Data.Recon = data.Data.Recon
	logrus.Infof("Received Handle_ConnectSub : %+v\n", data)
}

func Handle_disconnectSub(tnn_wifiStatus *wifi_models.WiFiStatus, payload []byte) {
	mu.Lock()
	defer mu.Unlock()

	var data wifi_models.DisconnectMsg
	err := json.Unmarshal(payload, &data)
	if err != nil {
		logrus.Warningf("Error unmarshaling JSON: %s\n", err)
		return
	}

	tnn_wifiStatus.DisconnectMsg.SetDisconnect = data.SetDisconnect
	tnn_wifiStatus.DisconnectMsg.SetRemove = data.SetRemove
	tnn_wifiStatus.DisconnectMsg.Data.SSID = data.Data.SSID
	logrus.Infof("Received Handle_adapterSub : %+v\n", data)
}

func Handle_hotspotSub(tnn_wifiStatus *wifi_models.WiFiStatus, payload []byte) {
	mu.Lock()
	defer mu.Unlock()

	var data wifi_models.HotspotMsg
	err := json.Unmarshal(payload, &data)
	if err != nil {
		logrus.Warningf("Error unmarshaling JSON: %s\n", err)
		return
	}

	tnn_wifiStatus.HotspotMsg.SetHotspot = data.SetHotspot
	logrus.Infof("Received Handle_adapterSub : %+v\n", data)
}
