package wifi_fn

import (
	"fmt"

	"github.com/Wifx/gonetworkmanager"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func WiFi_hotspotOn(nm gonetworkmanager.NetworkManager, w_dev gonetworkmanager.DeviceWireless, hotspotSSID string, hotspotPW string) error {
	// Wireless device check
	if w_dev == nil {
		return fmt.Errorf("no Wireless device found")
	}

	logrus.Infof("Initilizing Hotspot '%s' ...\n", hotspotSSID)

	// Define Connection Settings
	settings := make(map[string]map[string]interface{})

	// metadata of the profile
	settings["connection"] = map[string]interface{}{
		"id":          hotspotSSID,         // Connection name
		"type":        "802-11-wireless",   // Wireless type
		"uuid":        uuid.New().String(), // Unique ID for this profile
		"autoconnect": false,               // auto  reconnect option
	}

	// configures the physical radio layer
	settings["802-11-wireless"] = map[string]interface{}{
		"ssid": []byte(hotspotSSID), // SSID must be a byte array
		"mode": "ap",
		"band": "bg",
	}

	// handles
	settings["802-11-wireless-security"] = map[string]interface{}{
		"key-mgmt": "wpa-psk", // WPA2 - Personal
		"proto":    []string{"rsn"},
		"psk":      hotspotPW,
	}

	// ipv4 and ipv6 sections (DHCP)
	settings["ipv4"] = map[string]interface{}{
		"method": "shared",
	}
	settings["ipv6"] = map[string]interface{}{
		"method": "ignore",
	}

	settingsService, err := gonetworkmanager.NewSettings()
	if err != nil {
		return fmt.Errorf("failed to get settings service: %v", err)
	}

	connection, err := settingsService.AddConnection(settings)
	if err != nil {
		return fmt.Errorf("failed to create connection profile: %v", err)
	}
	logrus.Infof("Profile Created: %s", connection.GetPath())

	activeConnection, err := nm.ActivateConnection(connection, w_dev, nil)
	if err != nil {
		_ = connection.Delete()
		return fmt.Errorf("failed to activate hotspot: %v", err)
	}

	logrus.Infof("Hotspot started successfully!")
	logrus.Infof("Active Connection Object: %s\n", activeConnection.GetPath())

	return nil
}
