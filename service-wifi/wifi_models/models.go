package wifi_models

import "time"

type PropertiesMsg struct {
	HeaderID  string    `json:"HeaderID"`
	Timestamp time.Time `json:"timestamp"`
	Data      struct {
		WiFiType  string `json:"WiFiType"`
		WiFiMode  string `json:"WiFiMode"`
		Security  string `json:"Security"`
		SSID      string `json:"SSID"`
		Password  string `json:"Password"`
		Signal    uint8  `json:"Signal"`
		Frequency uint32 `json:"Frequency"`
		MAC       string `json:"MAC"`
		IPv4      string `json:"IPv4"`
		IPv6      string `json:"IPv6"`
		Error     string `json:"Error,omitempty"`
	} `json:"data"`
}

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

type ScanMsg struct {
	HeaderID  string    `json:"HeaderID"`
	Timestamp time.Time `json:"timestamp"`
	SetScan   bool      `json:"SetScan"`
	Data      struct {
		IsScan bool   `json:"IsScan"`
		Error  string `json:"Error,omitempty"`
	} `json:"data"`
}

type ConnectMsg struct {
	HeaderID   string    `json:"HeaderID"`
	Timestamp  time.Time `json:"timestamp"`
	SetConnect bool      `json:"SetConnect"`
	Data       struct {
		IsConnected bool   `json:"IsConnected"`
		SSID        string `json:"SSiD"`
		Password    string `json:"Password"`
		Recon       bool   `json:"Recon"`
		Error       string `json:"Error,omitempty"`
	} `json:"data"`
}

type DisconnectMsg struct {
	HeaderID      string    `json:"HeaderID"`
	Timestamp     time.Time `json:"timestamp"`
	SetDisconnect bool      `json:"SetDisconnect"`
	SetRemove     bool      `json:"SetRemove"`
	Data          struct {
		IsDisconnected bool   `json:"IsDisconnected"`
		SSID           string `json:"SSiD"`
		Error          string `json:"Error,omitempty"`
	} `json:"data"`
}

type HotspotMsg struct {
	HeaderID   string    `json:"HeaderID"`
	Timestamp  time.Time `json:"timestamp"`
	SetHotspot bool      `json:"SetHotspot"`
	Data       struct {
		IsHotspot bool   `json:"IsHotspot"`
		Error     string `json:"Error,omitempty"`
	} `json:"data"`
}

type WiFiStatus struct {
	AdapterMsg    AdapterMsg    `json:"AdapterMsg"`
	PropertiesMsg PropertiesMsg `json:"PropertiesMsg"`
	HotspotMsg    HotspotMsg    `json:"HotspotMsg"`
	ScanMsg       ScanMsg       `json:"ScanMsg"`
	ConnectMsg    ConnectMsg    `json:"ConnectMsg"`
	DisconnectMsg DisconnectMsg `json:"DisconnectMsg"`
}

type DeviceInfo struct {
	SSID     string `json:"SSID"`
	BSSID    string `json:"BSSID"`
	Strength uint8  `json:"Strength"`
}

type ScanDeviceList struct {
	HeaderID  string    `json:"HeaderID"`
	Timestamp time.Time `json:"timestamp"`
	Data      struct {
		Devices []DeviceInfo `json:"Devices"`
		Error   string       `json:"Error,omitempty"`
	} `json:"data"`
}
