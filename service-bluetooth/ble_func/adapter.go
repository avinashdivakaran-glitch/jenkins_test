package ble_func

import (
	"fmt"
	"path"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/sirupsen/logrus"
)

// Adapter represents a Bluetooth adapter on the system.
type Adapter struct {
	Path dbus.ObjectPath
}

// Convert MAC
func formatMAC(mac string) string {
	// Example: 68:5F:4A:5C:10:31 -> 68_5F_4A_5C_10_31
	return strings.ReplaceAll(mac, ":", "_")
}

// ListAdapters returns all the available adapter.
func ListAdapters() ([]*Adapter, error) {
	objManager, err := bluez.GetObjectManager()
	if err != nil {
		return nil, err
	}

	managedObjects, err := objManager.GetManagedObjects()
	if err != nil {
		return nil, err
	}

	var adapters []*Adapter
	// scanning that map to find objects that implement the org.bluez.Adapter1 interface
	for path, ifaces := range managedObjects {
		if _, ok := ifaces["org.bluez.Adapter1"]; ok {
			adapters = append(adapters, &Adapter{
				Path: path,
			})
			logrus.Info("Found adapter 📡 :", path)
		}
	}

	if len(adapters) == 0 {
		return nil, fmt.Errorf("no Bluetooth adapters found")
	}

	return adapters, nil
}

// Returns the default Bluetooth adapter (usually hci0)
func Get_adapter() (*adapter.Adapter1, error) {

	// List all available adapters
	adapters, err := ListAdapters()
	if err != nil {
		logrus.Fatal(err)
	}

	fullPath := adapters[0].Path
	lastPart := path.Base(string(fullPath))
	a, err := api.GetAdapter(lastPart)
	if err != nil {
		return nil, fmt.Errorf("failed to get adapter: %v", err)
	}
	return a, nil
}

// Return adapter properties
func Prop_adapter(a *adapter.Adapter1) (name string, power bool, discoverable bool, err error) {
	props, err := a.GetProperties()
	if err != nil {
		return "", false, false, fmt.Errorf("failed to get adapter properties: %v", err)
	}

	name = string(a.Path())
	if props.Powered {
		power = true
	} else {
		power = false
	}

	if props.Discoverable {
		discoverable = true
	} else {
		discoverable = false
	}

	return name, power, discoverable, nil
}

// Enable or Disable the adpater
func PowerOn_adapter(a *adapter.Adapter1, on bool) error {
	props, err := a.GetProperties()
	if err != nil {
		return fmt.Errorf("failed to get adapter properties: %v", err)
	}

	if on {
		if !props.Powered {
			logrus.Infoln("Bluetooth adapter is OFF — powering on...")
			// sudo rfkill unblock bluetooth
			if err := a.SetProperty("Powered", true); err != nil {
				return fmt.Errorf("failed to power on adapter: %v (sudo rfkill unblock bluetooth)", err)
			}
			logrus.Infoln("Adapter powered ON")
		} else {
			logrus.Infoln("Adapter already ON")
		}
	} else {
		if props.Powered {
			logrus.Infoln("Bluetooth adapter is ON — powering off...")
			if err := a.SetProperty("Powered", false); err != nil {
				return fmt.Errorf("failed to power off adapter: %v", err)
			}
			logrus.Infoln("Adapter powered OFF")
		} else {
			logrus.Infoln("Adapter already OFF")
		}
	}

	return nil
}
