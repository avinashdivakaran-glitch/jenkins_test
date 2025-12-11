package main

import (
	"fmt"
	"main/wifi_models"
	"reflect"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

func TNN_WiFiconfig(mqttClient mqtt.Client) {
	dev, err := Init_tnnWifi()
	if err != nil {
		TNN_wifiStatus.AdapterMsg.Data.IsAvailable = false
		TNN_wifiStatus.AdapterMsg.Data.IsPowered = false
		TNN_wifiStatus.AdapterMsg.Data.Adapter = ""
		TNN_wifiStatus.AdapterMsg.Data.Error = err.Error()
	} else {
		TNN_wifiStatus.AdapterMsg.Data.IsAvailable = true
		TNN_wifiStatus.AdapterMsg.Data.Adapter = dev
		TNN_wifiStatus.AdapterMsg.Data.Error = ""
	}

	TNN_wifiStatus.AdapterMsg.Timestamp = time.Now()
	err = Publish_mqttBroker(mqttClient, topics["adapter"], TNN_wifiStatus.AdapterMsg)
	if err != nil {
		logrus.Warning(err)
	}
}

func TNN_wifiStatusFnuc(mqttClient mqtt.Client) {
	status, err := WiFi_status()
	if err != nil {
		TNN_wifiStatus.PropertiesMsg.Data.Error = err.Error()
		logrus.Errorf("Error TNN_wifiStatusFnuc : %v \n", err)
	} else {
		TNN_wifiStatus.PropertiesMsg.Data.Error = ""
	}

	if !reflect.DeepEqual(status.Data, TNN_wifiStatus.PropertiesMsg.Data) {
		TNN_wifiStatus.PropertiesMsg.Timestamp = time.Now()
		TNN_wifiStatus.PropertiesMsg.Data = status.Data
		err = Publish_mqttBroker(mqttClient, topics["property"], TNN_wifiStatus.PropertiesMsg)
		if err != nil {
			logrus.Warning(err)
		}
	}

}

func TNN_adapterFunc(mqttClient mqtt.Client) {
	if TNN_wifiStatus.AdapterMsg.SetPower != TNN_wifiStatus.AdapterMsg.Data.IsPowered {
		if err := Turn_tnnWiFi_on(TNN_wifiStatus.AdapterMsg.SetPower); err != nil {
			TNN_wifiStatus.AdapterMsg.Data.Error = err.Error()
			TNN_wifiStatus.AdapterMsg.SetPower = TNN_wifiStatus.AdapterMsg.Data.IsPowered
			logrus.Errorf("Error Turn_tnnWiFi_on : %v \n", err)
		} else {
			TNN_wifiStatus.AdapterMsg.Data.IsPowered = TNN_wifiStatus.AdapterMsg.SetPower
			TNN_wifiStatus.AdapterMsg.Data.Error = ""
		}

		TNN_wifiStatus.AdapterMsg.Timestamp = time.Now()
		err := Publish_mqttBroker(mqttClient, topics["adapter"], TNN_wifiStatus.AdapterMsg)
		if err != nil {
			logrus.Warning(err)
		}
	}
}

func TNN_HotspotFunc(mqttClient mqtt.Client) {
	if TNN_wifiStatus.HotspotMsg.SetHotspot != TNN_wifiStatus.HotspotMsg.Data.IsHotspot {
		if TNN_wifiStatus.HotspotMsg.SetHotspot {
			if err := Hotspot_tnnOn("jetson_tnn", "12345678"); err != nil {
				TNN_wifiStatus.HotspotMsg.Data.Error = err.Error()
				TNN_wifiStatus.HotspotMsg.SetHotspot = TNN_wifiStatus.HotspotMsg.Data.IsHotspot
				logrus.Errorf("Error Hotspot_tnnOn : %v \n", err)
			} else {
				TNN_wifiStatus.HotspotMsg.Data.IsHotspot = true
				TNN_wifiStatus.HotspotMsg.Data.Error = ""
			}
		} else {
			if err := Disconnect_tnnWiFi(); err != nil {
				TNN_wifiStatus.HotspotMsg.Data.Error = err.Error()
				TNN_wifiStatus.HotspotMsg.SetHotspot = TNN_wifiStatus.HotspotMsg.Data.IsHotspot
				logrus.Errorf("Error Disconnect_tnnWiFi : %v \n", err)
			} else {
				TNN_wifiStatus.HotspotMsg.Data.IsHotspot = false
				TNN_wifiStatus.HotspotMsg.Data.Error = ""
			}
		}

		time.Sleep(2 * time.Second)

		TNN_wifiStatus.HotspotMsg.Timestamp = time.Now()
		err := Publish_mqttBroker(mqttClient, topics["hotspot"], TNN_wifiStatus.HotspotMsg)
		if err != nil {
			logrus.Warning(err)
		}
	}
}

func TNN_ScanFunc(mqttClient mqtt.Client, resultsChan chan []wifi_models.DeviceInfo, stopChan chan struct{}, errChan chan error) {
	var currentList []wifi_models.DeviceInfo

	if TNN_wifiStatus.ScanMsg.SetScan {
		if TNN_wifiStatus.ScanMsg.Data.IsScan {
			currentList = <-resultsChan
			fmt.Println(currentList)
			TNN_wifiScanDevice.Data.Devices = currentList
			TNN_wifiScanDevice.Data.Error = ""
		} else {
			stopChan = make(chan struct{})
			go Scan_tnnWiFi_dev(resultsChan, stopChan, errChan)
			TNN_wifiStatus.ScanMsg.Data.IsScan = true
		}

		// err := <-errChan
		// TNN_wifiScanDevice.Data.Error = err.Error()
		TNN_wifiScanDevice.Timestamp = time.Now()
		err := Publish_mqttBroker(mqttClient, topics["scan_dev"], TNN_wifiScanDevice)
		if err != nil {
			logrus.Warning(err)
		}

		TNN_wifiStatus.ScanMsg.Timestamp = time.Now()
		err = Publish_mqttBroker(mqttClient, topics["scan"], TNN_wifiStatus.ScanMsg)
		if err != nil {
			logrus.Warning(err)
		}

	} else {
		if TNN_wifiStatus.ScanMsg.Data.IsScan {
			if stopChan != nil {
				close(stopChan)
				TNN_wifiStatus.ScanMsg.Data.IsScan = false

				TNN_wifiStatus.ScanMsg.Timestamp = time.Now()
				err := Publish_mqttBroker(mqttClient, topics["scan"], TNN_wifiStatus.ScanMsg)
				if err != nil {
					logrus.Warning(err)
				}

				TNN_wifiScanDevice.Timestamp = time.Now()
				TNN_wifiScanDevice.Data.Devices = nil
				err = Publish_mqttBroker(mqttClient, topics["scan_dev"], TNN_wifiScanDevice)
				if err != nil {
					logrus.Warning(err)
				}
			}
		} else {
			// nothing
		}

	}

}

func TNN_ConnectFunc(mqttClient mqtt.Client) {
	if TNN_wifiStatus.ConnectMsg.SetConnect {
		if len(TNN_wifiStatus.ConnectMsg.Data.SSID) == 0 || len(TNN_wifiStatus.ConnectMsg.Data.Password) == 0 {
			TNN_wifiStatus.ConnectMsg.Data.Error = "No SSID or Password"
		} else {
			err := Connect_tnnWiFi_dev(TNN_wifiStatus.ConnectMsg.Data.SSID, TNN_wifiStatus.ConnectMsg.Data.Password, TNN_wifiStatus.ConnectMsg.Data.Recon)
			if err != nil {
				TNN_wifiStatus.ConnectMsg.Data.Error = err.Error()
				logrus.Errorf("Error TNN_ConnectFunc : %v", err)
			} else {
				TNN_wifiStatus.ConnectMsg.Data.IsConnected = true
				TNN_wifiStatus.ConnectMsg.Data.Error = ""
			}
		}

		TNN_wifiStatus.ConnectMsg.SetConnect = false
		TNN_wifiStatus.ConnectMsg.Timestamp = time.Now()
		err := Publish_mqttBroker(mqttClient, topics["connect"], TNN_wifiStatus.ConnectMsg)
		if err != nil {
			logrus.Warning(err)
		}
	}
}

func TNN_DisconnectFunc(mqttClient mqtt.Client) {
	if TNN_wifiStatus.DisconnectMsg.SetDisconnect {
		err := Disconnect_tnnWiFi()
		if err != nil {
			TNN_wifiStatus.DisconnectMsg.Data.Error = err.Error()
			logrus.Errorf("Error Disconnect_tnnWiFi : %v", err)
		} else {
			TNN_wifiStatus.DisconnectMsg.Data.IsDisconnected = true
			TNN_wifiStatus.DisconnectMsg.Data.Error = ""
		}

		TNN_wifiStatus.DisconnectMsg.SetDisconnect = false
		TNN_wifiStatus.DisconnectMsg.Timestamp = time.Now()
		err = Publish_mqttBroker(mqttClient, topics["disconnect"], TNN_wifiStatus.DisconnectMsg)
		if err != nil {
			logrus.Warning(err)
		}
	}
}

func TNN_DeleteFunc(mqttClient mqtt.Client) {
	if TNN_wifiStatus.DisconnectMsg.SetRemove {
		if len(TNN_wifiStatus.DisconnectMsg.Data.SSID) == 0 {
			TNN_wifiStatus.DisconnectMsg.Data.Error = "No SSID"
		} else {
			err := Delete_tnnWiFi_con(TNN_wifiStatus.DisconnectMsg.Data.SSID)
			if err != nil {
				TNN_wifiStatus.DisconnectMsg.Data.Error = err.Error()
				logrus.Errorf("Error Delete_tnnWiFi_con : %v", err)
			} else {
				TNN_wifiStatus.DisconnectMsg.Data.IsDisconnected = true
				TNN_wifiStatus.DisconnectMsg.Data.Error = ""
			}

			err = Disconnect_tnnWiFi()
			if err != nil {
				TNN_wifiStatus.DisconnectMsg.Data.Error = err.Error()
				logrus.Errorf("Error Disconnect_tnnWiFi : %v", err)
			} else {
				TNN_wifiStatus.DisconnectMsg.Data.IsDisconnected = true
				TNN_wifiStatus.DisconnectMsg.Data.Error = ""
			}
		}

		TNN_wifiStatus.DisconnectMsg.SetRemove = false
		TNN_wifiStatus.DisconnectMsg.Timestamp = time.Now()
		err := Publish_mqttBroker(mqttClient, topics["disconnect"], TNN_wifiStatus.DisconnectMsg)
		if err != nil {
			logrus.Warning(err)
		}
	}
}
