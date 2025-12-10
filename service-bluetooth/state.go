package main

import (
	"reflect"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

func TNN_BLEconfig(mqttClient mqtt.Client) {
	if err := Init_tnnBLE(); err != nil {
		TNN_bleStatus.AdapterMsg.Data.IsAvailable = false
		TNN_bleStatus.AdapterMsg.Data.Error = err.Error()
		logrus.Errorf("Error Init_tnnBLE : %v \n", err)
	} else {
		TNN_bleStatus.AdapterMsg.Data.IsAvailable = true
		TNN_bleStatus.AdapterMsg.Data.Error = "no error"
	}

	TNN_bleStatus.AdapterMsg.Timestamp = time.Now()
	err := Publish_mqttBroker(mqttClient, topics["adapter"], TNN_bleStatus.AdapterMsg)
	if err != nil {
		logrus.Warning(err)
	}

}

func TNN_BLEStateFunc(mqttClient mqtt.Client) {
	if adapterName, powerState, discoverable, err := Update_tnnBLE(); err != nil {
		TNN_bleStatus.AdapterMsg.Data.Error = err.Error()
		logrus.Errorf("Error Update_tnnBLE : %v \n", err)
	} else {
		TNN_bleStatus.AdapterMsg.Data.Adapter = adapterName
		TNN_bleStatus.AdapterMsg.Data.IsPowered = powerState
		TNN_bleStatus.DiscoverMsg.Data.IsDiscoverable = discoverable
		TNN_bleStatus.AdapterMsg.Data.Error = "no error"
	}

	if !TNN_bleStatus.AdapterMsg.Data.IsPowered {
		return
	}

	if TNN_bleStatus.AdapterMsg.SetPower != TNN_bleStatus.AdapterMsg.Data.IsPowered {
		TNN_bleStatus.AdapterMsg.Timestamp = time.Now()
		err := Publish_mqttBroker(mqttClient, topics["adapter"], TNN_bleStatus.AdapterMsg)
		if err != nil {
			logrus.Warning(err)
		}
	}

	if TNN_bleStatus.DiscoverMsg.SetDiscover != TNN_bleStatus.DiscoverMsg.Data.IsDiscoverable {
		TNN_bleStatus.DiscoverMsg.Timestamp = time.Now()
		err := Publish_mqttBroker(mqttClient, topics["discover"], TNN_bleStatus.DiscoverMsg)
		if err != nil {
			logrus.Warning(err)
		}
	}

	if deviceList, err := Bound_tnnBLE_devie(); err != nil {
		TNN_bleBondDevice.Data.Devices = nil
		TNN_bleBondDevice.Data.Error = err.Error()
		logrus.Errorf("Error Bound_tnnBLE_devie : %v \n", err)
	} else {
		if !reflect.DeepEqual(TNN_bleBondDevice.Data.Devices, deviceList) {
			// if len(TNN_bleBondDevice.Data.Devices) != len(deviceList) {
			TNN_bleBondDevice.Data.Devices = deviceList
			TNN_bleBondDevice.Data.Error = "no error"
			TNN_bleBondDevice.Timestamp = time.Now()
			err := Publish_mqttBroker(mqttClient, topics["save_dev"], TNN_bleBondDevice)
			if err != nil {
				logrus.Warning(err)
			}
		}

	}

}

func TNN_adapterFunc(mqttClient mqtt.Client) {
	if TNN_bleStatus.AdapterMsg.SetPower != TNN_bleStatus.AdapterMsg.Data.IsPowered {
		if err := Turn_tnnBLE_on(TNN_bleStatus.AdapterMsg.SetPower); err != nil {
			TNN_bleStatus.AdapterMsg.Data.Error = err.Error()
			logrus.Errorf("Error Turn_tnnBLE_on : %v \n", err)
		} else {
			TNN_bleStatus.AdapterMsg.Data.IsPowered = TNN_bleStatus.AdapterMsg.SetPower
			TNN_bleStatus.AdapterMsg.Data.Error = "no error"
		}

		TNN_bleStatus.AdapterMsg.Timestamp = time.Now()
		err := Publish_mqttBroker(mqttClient, topics["adapter"], TNN_bleStatus.AdapterMsg)
		if err != nil {
			logrus.Warning(err)
		}
	}
}

func TNN_discoverFunc(mqttClient mqtt.Client) {
	if !TNN_bleStatus.AdapterMsg.Data.IsPowered {
		return
	}

	if TNN_bleStatus.DiscoverMsg.SetDiscover != TNN_bleStatus.DiscoverMsg.Data.IsDiscoverable {
		if err := Discover_tnnBLE(TNN_bleStatus.DiscoverMsg.SetDiscover); err != nil {
			TNN_bleStatus.DiscoverMsg.Data.Error = err.Error()
			logrus.Errorf("Error Discover_tnnBLE : %v \n", err)
		} else {
			TNN_bleStatus.DiscoverMsg.Data.IsDiscoverable = TNN_bleStatus.DiscoverMsg.SetDiscover
			TNN_bleStatus.DiscoverMsg.Data.Error = "no error"
		}

		TNN_bleStatus.DiscoverMsg.Timestamp = time.Now()
		err := Publish_mqttBroker(mqttClient, topics["discover"], TNN_bleStatus.DiscoverMsg)
		if err != nil {
			logrus.Warning(err)
		}
	}
}

func TNN_ScanFunc(mqttClient mqtt.Client, scanStatusCh chan bool) {
	if !TNN_bleStatus.AdapterMsg.Data.IsPowered {
		return
	}

	if TNN_bleStatus.ScanMsg.SetScan {
		scanStatusCh <- true

		TNN_bleStatus.ScanMsg.Data.Error = "no error"
		TNN_bleScanDevice.Data.Error = "no error"
		TNN_bleScanDevice.Timestamp = time.Now()
		err := Publish_mqttBroker(mqttClient, topics["scan_dev"], TNN_bleScanDevice)
		if err != nil {
			logrus.Warning(err)
		}

		if scan_option {
			for !TNN_bleStatus.ScanMsg.Data.IsScan {
				scanStatusCh <- false
				scanStatusCh <- true
			}

			TNN_bleStatus.ScanMsg.Timestamp = time.Now()
			err = Publish_mqttBroker(mqttClient, topics["scan"], TNN_bleStatus.ScanMsg)
			if err != nil {
				logrus.Warning(err)
			}
			scan_option = false
		}

	} else {
		scanStatusCh <- false
		scan_option = true
		if TNN_bleStatus.ScanMsg.Data.IsScan {
			TNN_bleStatus.ScanMsg.Data.IsScan = false
			TNN_bleStatus.ScanMsg.Timestamp = time.Now()
			err := Publish_mqttBroker(mqttClient, topics["scan"], TNN_bleStatus.ScanMsg)
			if err != nil {
				logrus.Warning(err)
			}
		}

	}

}

func TNN_PairFunc(mqttClient mqtt.Client) {
	if !TNN_bleStatus.AdapterMsg.Data.IsPowered {
		return
	}

	if TNN_bleStatus.PairMsg.SetPair {
		if len(TNN_bleStatus.PairMsg.Data.DeviceMAC) == 0 {
			TNN_bleStatus.PairMsg.Data.Error = "MAC not available"
		} else {
			// if err := Remove_tnnBLE_device(TNN_bleStatus.PairMsg.Data.DeviceMAC); err != nil {
			// 	TNN_bleStatus.PairMsg.Data.Error = err.Error()
			// } else {
			// 	TNN_bleStatus.PairMsg.Data.IsRemoved = true
			// }
			if err := Pair_tnnBLE_device(TNN_bleStatus.PairMsg.Data.DeviceMAC); err != nil {
				TNN_bleStatus.PairMsg.Data.Error = err.Error()
				logrus.Errorf("Error Pair_tnnBLE_device : %v \n", err)
			} else {
				TNN_bleStatus.PairMsg.Data.IsPaired = true
				TNN_bleStatus.PairMsg.Data.Error = "no error"
			}
		}

		TNN_bleStatus.PairMsg.SetPair = false
		TNN_bleStatus.PairMsg.Timestamp = time.Now()
		err := Publish_mqttBroker(mqttClient, topics["pair"], TNN_bleStatus.PairMsg)
		if err != nil {
			logrus.Warning(err)
		}
	}
}

func TNN_RemoveFunc(mqttClient mqtt.Client) {
	if !TNN_bleStatus.AdapterMsg.Data.IsPowered {
		return
	}

	if TNN_bleStatus.PairMsg.SetRemove {
		if len(TNN_bleStatus.PairMsg.Data.DeviceMAC) == 0 {
			TNN_bleStatus.PairMsg.Data.Error = "MAC not available"
		} else {
			// if err := Pair_tnnBLE_device(TNN_bleStatus.PairMsg.Data.DeviceMAC); err != nil {
			// 	TNN_bleStatus.PairMsg.Data.Error = err.Error()
			// } else {
			// 	TNN_bleStatus.PairMsg.Data.IsPaired = true
			// }

			if err := Remove_tnnBLE_device(TNN_bleStatus.PairMsg.Data.DeviceMAC); err != nil {
				TNN_bleStatus.PairMsg.Data.Error = err.Error()
				logrus.Errorf("Error Remove_tnnBLE_device : %v \n", err)
			} else {
				TNN_bleStatus.PairMsg.Data.IsRemoved = true
				TNN_bleStatus.PairMsg.Data.Error = "no error"
			}
		}

		TNN_bleStatus.PairMsg.SetRemove = false
		TNN_bleStatus.PairMsg.Timestamp = time.Now()
		err := Publish_mqttBroker(mqttClient, topics["pair"], TNN_bleStatus.PairMsg)
		if err != nil {
			logrus.Warning(err)
		}
	}
}

func TNN_ConnectFunc(mqttClient mqtt.Client) {
	if !TNN_bleStatus.AdapterMsg.Data.IsPowered {
		return
	}

	if TNN_bleStatus.ConnectMsg.SetConnect {
		if len(TNN_bleStatus.ConnectMsg.Data.DeviceMAC) == 0 {
			TNN_bleStatus.ConnectMsg.Data.Error = "MAC not available"
		} else {
			if err := Connect_tnnBLE_device(TNN_bleStatus.ConnectMsg.Data.DeviceMAC); err != nil {
				TNN_bleStatus.ConnectMsg.Data.Error = err.Error()
				logrus.Errorf("Error Connect_tnnBLE_device : %v \n", err)
			} else {
				TNN_bleStatus.ConnectMsg.Data.IsConnected = true
				TNN_bleStatus.ConnectMsg.Data.Error = "no error"
			}
		}

		TNN_bleStatus.ConnectMsg.SetConnect = false
		TNN_bleStatus.ConnectMsg.Timestamp = time.Now()
		err := Publish_mqttBroker(mqttClient, topics["connect"], TNN_bleStatus.ConnectMsg)
		if err != nil {
			logrus.Warning(err)
		}
	}
}

func TNN_DisconnectFunc(mqttClient mqtt.Client) {
	if !TNN_bleStatus.AdapterMsg.Data.IsPowered {
		return
	}

	if TNN_bleStatus.ConnectMsg.SetDisconnect {
		if len(TNN_bleStatus.ConnectMsg.Data.DeviceMAC) == 0 {
			TNN_bleStatus.ConnectMsg.Data.Error = "MAC not available"
		} else {
			if err := Disconnect_tnnBLE_device(TNN_bleStatus.ConnectMsg.Data.DeviceMAC); err != nil {
				TNN_bleStatus.ConnectMsg.Data.Error = err.Error()
				logrus.Errorf("Error Disconnect_tnnBLE_device : %v \n", err)
			} else {
				TNN_bleStatus.ConnectMsg.Data.IsConnected = true
				TNN_bleStatus.ConnectMsg.Data.Error = "no error"
			}
		}

		TNN_bleStatus.ConnectMsg.SetDisconnect = false
		TNN_bleStatus.ConnectMsg.Timestamp = time.Now()
		err := Publish_mqttBroker(mqttClient, topics["connect"], TNN_bleStatus.ConnectMsg)
		if err != nil {
			logrus.Warning(err)
		}
	}
}
