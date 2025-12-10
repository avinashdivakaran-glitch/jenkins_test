package sensors

import "time"

type DeviceData struct {
	Available bool      `json:"available"`
	DeviceID  string    `json:"device_id"`
	Timestamp time.Time `json:"timestamp"`
}

type SensorData struct {
	Available bool      `json:"available"`
	SensorID  string    `json:"sensor_id"`
	Timestamp time.Time `json:"timestamp"`
	Data      struct {
		Value float64 `json:"value"`
		Unit  string  `json:"unit"`
		State string  `json:"state"`
	} `json:"data"`
}
