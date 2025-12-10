package main

import (
	"health_monitor/sensors"
	"time"
)

// Channels for sensor data
var (
	deviceDataChan  = make(chan sensors.DeviceData, 10)
	TsensorDataChan = make(chan sensors.SensorData, 10)
	HsensorDataChan = make(chan sensors.SensorData, 10)
	OsensorDataChan = make(chan sensors.SensorData, 10)
)

func main() {

	// Init mqtt client
	mqttclient := Init_mqttClient()

	time.Sleep(2 * time.Second)

	go sensorReader()

	for {
		select {
		case devData := <-deviceDataChan:
			Publish_mqttBroker(mqttclient, topics["device"], devData)

		case tData := <-TsensorDataChan:
			Publish_mqttBroker(mqttclient, topics["temperature"], tData)

		case hData := <-HsensorDataChan:
			Publish_mqttBroker(mqttclient, topics["heartrate"], hData)

		case oData := <-OsensorDataChan:
			Publish_mqttBroker(mqttclient, topics["bodyoxygen"], oData)

		}
	}

}
