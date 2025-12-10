package ble_models

import (
	"time"
)

type AdapterMsg struct {
	HeaderID  string    `json:"HeaderID"`
	Timestamp time.Time `json:"timestamp"`
	SetPower  bool      `json:"SetPower"`
	Data      struct {
		IsAvailable bool   `json:"IsAvailable"`
		Adapter     string `json:"Adapter"`
		IsPowered   bool   `json:"IsPowered"`
		Error       string `json:"Error,omitempty"`
	} `json:"data"`
}

type DiscoverMsg struct {
	HeaderID    string    `json:"HeaderID"`
	Timestamp   time.Time `json:"timestamp"`
	SetDiscover bool      `json:"SetDiscover"`
	Data        struct {
		IsDiscoverable bool   `json:"IsDiscoverable"`
		Error          string `json:"Error,omitempty"`
	} `json:"data"`
}

type ScanMsg struct {
	HeaderID  string    `json:"HeaderID"`
	Timestamp time.Time `json:"timestamp"`
	SetScan   bool      `json:"SetScan"`
	Data      struct {
		IsScan bool   `json:"IsScan"`
		Error  string `json:"Error,omitempty"`
	} `json:"data"`
}

type PairMsg struct {
	HeaderID  string    `json:"HeaderID"`
	Timestamp time.Time `json:"timestamp"`
	SetPair   bool      `json:"SetPair"`
	SetRemove bool      `json:"SetRemove"`
	Data      struct {
		IsPaired  bool   `json:"IsPaired"`
		IsRemoved bool   `json:"IsRemoved"`
		DeviceMAC string `json:"DeviceMAC"`
		Error     string `json:"Error,omitempty"`
	} `json:"data"`
}

type ConnectMsg struct {
	HeaderID      string    `json:"HeaderID"`
	Timestamp     time.Time `json:"timestamp"`
	SetConnect    bool      `json:"SetConnect"`
	SetDisconnect bool      `json:"SetDisconnect"`
	Data          struct {
		IsConnected    bool   `json:"IsConnected"`
		IsDisconnected bool   `json:"IsDisconnected"`
		DeviceMAC      string `json:"DeviceMAC"`
		Error          string `json:"Error,omitempty"`
	} `json:"data"`
}

type BluetoothStatus struct {
	AdapterMsg  AdapterMsg  `json:"AdapterMsg"`
	DiscoverMsg DiscoverMsg `json:"DiscoverMsg"`
	ScanMsg     ScanMsg     `json:"ScanMsg"`
	PairMsg     PairMsg     `json:"PairMsg"`
	ConnectMsg  ConnectMsg  `json:"ConnectMsg"`
}

type PairDeviceInfo struct {
	Name        string   `json:"Name"`
	UUIDs       []string `json:"UUID"`
	MAC         string   `json:"MAC"`
	Signal      int16    `json:"Signal"`
	Class       uint32   `json:"Class"`
	Type        string   `json:"Type"`
	IsAvailable bool     `json:"IsAvailable"`
}

// Information about a Bluetooth device
type DeviceInfo struct {
	Name   string   `json:"Name"`
	UUIDs  []string `json:"UUID"`
	MAC    string   `json:"MAC"`
	Signal int16    `json:"Signal"`
	Class  uint32   `json:"Class"`
	Type   string   `json:"Type"`
}

type ScanDeviceList struct {
	HeaderID  string    `json:"HeaderID"`
	Timestamp time.Time `json:"timestamp"`
	Data      struct {
		Devices []DeviceInfo `json:"Devices"`
		Error   string       `json:"Error,omitempty"`
	} `json:"data"`
}

type PairedDeviceList struct {
	HeaderID  string    `json:"HeaderID"`
	Timestamp time.Time `json:"timestamp"`
	Data      struct {
		Devices []PairDeviceInfo `json:"Devices"`
		Error   string           `json:"Error,omitempty"`
	} `json:"data"`
}
