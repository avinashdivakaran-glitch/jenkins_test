package sensors

import (
	"math/rand"
	"time"
)

// GenerateRandomTemperature creates random Temperature sensor reading
func Generate_temp() SensorData {
	value := 36.0 + rand.Float64()*3 // 36–39°C

	var s SensorData
	s.Available = true
	s.SensorID = "T001"
	s.Timestamp = time.Now()
	s.Data.Value = value
	s.Data.Unit = "C"
	s.Data.State = getTemp_state(value)

	return s
}

// GenerateRandomHeartRate creates random HeartRate sensor reading
func Generate_heartRate() SensorData {
	value := float64(60 + rand.Intn(50)) // 60–110 bpm

	var s SensorData
	s.Available = true
	s.SensorID = "H001"
	s.Timestamp = time.Now()
	s.Data.Value = value
	s.Data.Unit = "bpm"
	s.Data.State = getHeart_state(int(value))

	return s
}

// GenerateRandomBodyOxygen creates random BodyOxygen sensor reading
func Generate_bodyOxygen() SensorData {
	value := 90.0 + rand.Float64()*10.0 // 90–100%

	var s SensorData
	s.Available = true
	s.SensorID = "O001"
	s.Timestamp = time.Now()
	s.Data.Value = value
	s.Data.Unit = "%"
	s.Data.State = getOxy_state(value)

	return s
}

// -------------------
// State functions
// -------------------

func getTemp_state(v float64) string {
	switch {
	case v < 36.5:
		return "Warning"
	case v >= 36.5 && v <= 37.5:
		return "Stable"
	default:
		return "Critical"
	}
}

func getHeart_state(v int) string {
	switch {
	case v < 60 || v > 100:
		return "Warning"
	case v >= 60 && v <= 90:
		return "Stable"
	default:
		return "Critical"
	}
}

func getOxy_state(v float64) string {
	switch {
	case v >= 95.0:
		return "Stable"
	case v >= 90.0 && v < 95.0:
		return "Warning"
	default:
		return "Critical"
	}
}

func init() {
	rand.Seed(time.Now().UnixNano()) // ensure randomness
}
