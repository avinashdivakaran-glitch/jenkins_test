package wifi_fn

import (
	"fmt"
	"time"

	"github.com/Wifx/gonetworkmanager"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Connect to the wifi network
func Connect_device(nm gonetworkmanager.NetworkManager, w_dev gonetworkmanager.DeviceWireless, ssid string, psword string, recon bool) error {
	// Wireless device check
	if w_dev == nil {
		return fmt.Errorf("no Wireless device found")
	}

	logrus.Infof("Connecting process to '%s'\n", ssid)

	// Define Connection Settings
	settings := make(map[string]map[string]interface{})

	// metadata of the profile
	settings["connection"] = map[string]interface{}{
		"id":          ssid,                // Connection name
		"type":        "802-11-wireless",   // Wireless type
		"uuid":        uuid.New().String(), // Unique ID for this profile
		"autoconnect": true,                // auto  reconnect option
	}

	// configures the physical radio layer
	settings["802-11-wireless"] = map[string]interface{}{
		"ssid": []byte(ssid), // SSID must be a byte array
		"mode": "infrastructure",
	}

	// handles authentication (Skip for Open networks)
	if psword != "" {
		settings["802-11-wireless-security"] = map[string]interface{}{
			"key-mgmt": "wpa-psk", // WPA2 - Personal
			"psk":      psword,
		}
	}

	// ipv4 and ipv6 sections (DHCP)
	settings["ipv4"] = map[string]interface{}{"method": "auto"}
	settings["ipv6"] = map[string]interface{}{"method": "auto"}

	// Add and Activate the Connection
	activeConnection, err := nm.AddAndActivateConnection(settings, w_dev)
	if err != nil {
		return fmt.Errorf("failed to initiate connection: %v", err)
	}

	logrus.Infof("Connection initiated!")
	logrus.Infof("Active Connection Object: %s\n", activeConnection.GetPath())

	// Wait for Connection Result
	logrus.Infof("Waiting for connection...")
	for i := 0; i < 30; i++ {
		state, err := activeConnection.GetPropertyState()
		if err != nil {
			return fmt.Errorf("\nError reading state: %v", err)
		}

		// NmActiveConnectionStateActivated = 2
		if state == 2 {
			logrus.Info("\nSUCCESS: Connected successfully!")
			return nil
		}

		fmt.Print(".")
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("timeout: Connection took too long or failed")
}

// Disconnect the connnect device
func Disconnect_con(w_dev gonetworkmanager.DeviceWireless) error {
	// Wireless device check
	if w_dev == nil {
		return fmt.Errorf("no Wireless device found")
	}

	// Check if currently active before disconnecting
	state, err := w_dev.GetPropertyState()
	if err != nil {
		return fmt.Errorf("error getting property %s", err)
	}

	// If state is Activated (100) or Activating, we disconnect
	if state > gonetworkmanager.NmDeviceStateDisconnected {
		logrus.Info("Disconnecting device...")
		err = w_dev.Disconnect()
		if err != nil {
			return fmt.Errorf("failed to disconnect: %v", err)
		}
		logrus.Info("Disconnected successfully.")
	} else {
		logrus.Info("Device is already disconnected.")
	}

	return nil
}

// Delete saved connection
func Delete_Con(con_Name string) error {
	set_service, err := gonetworkmanager.NewSettings()
	if err != nil {
		return fmt.Errorf("failed to get settings service: %v", err)
	}
	// List all saved connections
	connections, err := set_service.ListConnections()
	if err != nil {
		return fmt.Errorf("failed to list connections: %v", err)
	}

	found := false

	for _, conn := range connections {
		settings, err := conn.GetSettings()
		if err != nil {
			continue
		}
		if conSection, ok := settings["connection"]; ok {
			id, _ := conSection["id"].(string)
			uuid, _ := conSection["uuid"].(string)
			connType, _ := conSection["type"].(string)

			if id == con_Name {
				logrus.Infof("Found Match:\n - Name: %s\n - UUID: %s\n - Type: %s\n", id, uuid, connType)
			}

			err = conn.Delete()
			if err != nil {
				return fmt.Errorf("failed to delete connection: %v", err)
			}

			logrus.Infof("Connection profile deleted successfully.")
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("connection profile '%s' not found", con_Name)
	}

	return nil

}
