package main

import (
	"fmt"
	"health_monitor/ble_func"
	"health_monitor/sensors"
	"strconv"
	"time"

	"github.com/muka/go-bluetooth/bluez/profile/device"
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
	"github.com/sirupsen/logrus"
)

var (
	// ServiceUUID = "12345678-1234-5678-1234-111111111111"
	ServiceUUID = "12345678-1234-5678-1234"
	TempUUID    = "11111111-1111-1111-1111-111111111111"
	HRUUID      = "22222222-2222-2222-2222-222222222222"
	OxyUUID     = "33333333-3333-3333-3333-333333333333"
)

// BLE sensor data subscribe
func sensorReader() {

	// 1. Initialize Adapter ONCE.
	if err := ble_func.InitAdapter(); err != nil {
		logrus.Fatalln("Could not init adapter:", err)
	}

	logrus.Infoln("Waiting for sensor device...")

	// loop at every 5 seconds if nodive connect
	for {
		serviceDev, err := ble_func.GetService_Device(ServiceUUID)
		if err != nil {
			logrus.Error("Error GetService_Device", err)
			continue
		}
		if serviceDev == nil {
			logrus.Info("No Device service. Wait 5s ...")
			time.Sleep(5 * time.Second)
			continue
		}
		logrus.Infof("Found sensor device: %s \n", serviceDev.Properties.Address)

		// Continue with device
		for {
			if !serviceDev.Properties.Connected {
				logrus.Warn("Device is not connected. Retry after 5s ...")
				time.Sleep(5 * time.Second)
				break
			} else {
				logrus.Info("Already connected, starting service.")
			}

			// Subscribe to DBus property signals for this specific device
			propChan, err := serviceDev.WatchProperties()
			if err != nil {
				logrus.Error("Failed to watch device properties:", err)
				break
			}
			defer serviceDev.UnwatchProperties(propChan)

			// A ticker for BLE data subscribe
			subTicker := time.NewTicker(1 * time.Second)
			defer subTicker.Stop()

			ConnectionTry := 0
			for {
				select {
				// Event 1: DBus Signal received
				case change := <-propChan:
					if change.Interface == "org.bluez.Device1" && change.Name == "Connected" {
						isConnected := change.Value.(bool)
						if isConnected {
							ConnectionTry = 0
							logrus.Info("✅ Device connected (Signal received)")
							// Start the data reading ticker
							subTicker.Reset(1 * time.Second)
						} else {
							logrus.Warn("❌ Device disconnected (Signal received)")
							pair, _ := serviceDev.GetPaired()
							logrus.Infoln(" PAIR ", pair)
							if !pair {
								logrus.Infoln("device removed")
							} else {
								// call device connect func
								go serviceDev.Connect()
							}
							break
						}
					}

				// Event 2: The 1-second ticker (Only fires when BLE connected)
				case <-subTicker.C:
					if serviceDev.Properties.Connected {
						err = update_sensorData(serviceDev)
						if err != nil {
							logrus.Errorf(" Error update_sensorData : %v", err)
						}
					} else {
						ConnectionTry = ConnectionTry + 1
						logrus.Infoln("Try to connect device ....")
					}
				}

				// if connect time out 5 times, break the loop
				if ConnectionTry >= 5 {
					subTicker.Stop()
					serviceDev.Disconnect()
					break
				}
			}

		}

	}

}

func update_sensorData(device *device.Device1) error {
	// Discover the services
	uuids, err := device.GetUUIDs()
	if err != nil {
		return fmt.Errorf("failed to GetServiceData: %v", err)
	}
	var deviceData sensors.DeviceData
	for _, uuid := range uuids {
		if len(uuid) >= 23 && uuid[:23] == ServiceUUID {
			deviceData.Available = true
			deviceData.DeviceID = uuid
			deviceData.Timestamp = time.Now()
		}
	}
	deviceDataChan <- deviceData
	logrus.Infof(" 📟 Device Id : %v \n", deviceData.DeviceID)

	charList, err := device.GetCharacteristicsList()
	if err != nil {
		return fmt.Errorf("Failed to discover services: %v", err)
	}

	for _, chars := range charList {
		charProp, _ := gatt.NewGattCharacteristic1(chars)
		var value sensors.SensorData
		switch charProp.Properties.UUID {
		case TempUUID:
			value, err = reStruct_data(charProp, "℃", "")
			if err != nil {
				return fmt.Errorf(" TempUUID : %v", err)
			}
			value.SensorID = "Temperature"
			TsensorDataChan <- value
			// logrus.Infof("🌡️ Temperature : %v \n", value.Data.Value)

		case HRUUID:
			value, err = reStruct_data(charProp, "bpm", "")
			if err != nil {
				return fmt.Errorf(" HRUUID : %v", err)
			}
			value.SensorID = "Heart Rate"
			HsensorDataChan <- value
			// logrus.Infof("❤️ Heart Rate  : %v \n", value.Data.Value)

		case OxyUUID:
			value, err = reStruct_data(charProp, "%", "")
			if err != nil {
				return fmt.Errorf(" OxyUUID : %v", err)
			}
			value.SensorID = "SpO₂"
			OsensorDataChan <- value
			// logrus.Infof("🫁 SpO₂     : %v \n", value.Data.Value)
		}

	}

	return nil
}

func reStruct_data(charProp *gatt.GattCharacteristic1, unit string, state string) (sensors.SensorData, error) {
	var sensorData sensors.SensorData
	sensorData.Available = false
	options := map[string]interface{}{
		"offset": uint16(0),
	}
	valueByte, err := charProp.ReadValue(options)
	if err != nil {
		return sensorData, fmt.Errorf(" charProp.ReadValue, %v\n", err)
	}

	sensorData.Timestamp = time.Now()
	sensorData.Data.Value, err = strconv.ParseFloat(string(valueByte), 64)
	if err != nil {
		return sensorData, fmt.Errorf(" strconv.ParseFloat, %v\n", err)
	}
	sensorData.Data.Unit = unit
	sensorData.Data.State = state
	sensorData.Available = true

	return sensorData, nil
}
