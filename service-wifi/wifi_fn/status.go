package wifi_fn

import (
	"fmt"
	"main/wifi_models"

	"github.com/Wifx/gonetworkmanager"
	"github.com/sirupsen/logrus"
)

func ConnectionProp(w_dev gonetworkmanager.DeviceWireless) (wifi_models.PropertiesMsg, error) {
	var properWiFi wifi_models.PropertiesMsg

	activeConn, err := w_dev.GetPropertyActiveConnection()
	if err != nil {
		return properWiFi, fmt.Errorf("error Get Property Active Connection %v", err)
	}
	if activeConn == nil || activeConn.GetPath() == "/" {
		return wifi_models.PropertiesMsg{}, nil
	}

	activeAccessPoint, err := w_dev.GetPropertyActiveAccessPoint()
	if err != nil {
		return properWiFi, fmt.Errorf("error Get Property Active Access Point %v", err)
	}

	connection, err := activeConn.GetPropertyConnection()
	if err != nil {
		return properWiFi, fmt.Errorf("error Get Property Connection %v", err)
	}

	seviceSettings, err := connection.GetSettings()
	if err != nil {
		return properWiFi, fmt.Errorf("error connection Get Settings %v", err)
	}

	if wifiSettings, ok := seviceSettings["connection"]; ok {
		properWiFi.Data.WiFiType, _ = wifiSettings["type"].(string)
	}

	properWiFi.Data.WiFiMode = "Client (Infrastructure)"
	if wifiSettings, ok := seviceSettings["802-11-wireless"]; ok {
		if m, ok := wifiSettings["mode"].(string); ok {
			if m == "ap" {
				properWiFi.Data.WiFiMode = "Hotspot (Access Point)"
			}
		}
	}

	properWiFi.Data.Security = "Open (None)"
	if _, ok := seviceSettings["802-11-wireless-security"]; ok {
		if secMap, ok := seviceSettings["802-11-wireless-security"]; ok {
			if mgmt, ok := secMap["key-mgmt"].(string); ok {
				properWiFi.Data.Security = mgmt
			}
			if psk, ok := secMap["psk"].(string); ok {
				properWiFi.Data.Password = psk
			}
		}
	}

	if wifiSettings, ok := seviceSettings["802-11-wireless"]; ok {
		if ssidBytes, ok := wifiSettings["ssid"].([]byte); ok {
			properWiFi.Data.SSID = string(ssidBytes)
		}
	}

	mac, err := activeAccessPoint.GetPropertyHWAddress()
	if err != nil {
		return properWiFi, fmt.Errorf("error connection Get HwAddress %v", err)
	}
	properWiFi.Data.MAC = mac

	if freq, ok := activeAccessPoint.GetPropertyFrequency(); ok != nil {
		properWiFi.Data.Frequency = freq
	}

	if stren, ok := activeAccessPoint.GetPropertyStrength(); ok != nil {
		properWiFi.Data.Signal = stren
	}

	ip4Config, err := activeConn.GetPropertyIP4Config()
	if err != nil || ip4Config == nil {
		logrus.Warn("Error getting IPv4 configuration or config is nil")
	} else {
		ip4addrs, err := ip4Config.GetPropertyAddressData()
		if err != nil || len(ip4addrs) == 0 {
			logrus.Warn("Error getting IPv4 address data or no addresses found")
		} else {
			properWiFi.Data.IPv4 = ip4addrs[0].Address
		}
	}

	ip6Config, err := activeConn.GetPropertyIP6Config()
	if err != nil || ip6Config == nil {
		logrus.Warn("Error getting IPv6 configuration or config is nil")
	} else {
		ip6addrs, err := ip6Config.GetPropertyAddressData()
		if err != nil || len(ip6addrs) == 0 {
			logrus.Warn("Error getting IPv6 address data or no addresses found")
		} else {
			properWiFi.Data.IPv6 = ip6addrs[0].Address
		}
	}

	// logrus.Infof("WiFi status : %v", properWiFi.Data)

	return properWiFi, nil

}
