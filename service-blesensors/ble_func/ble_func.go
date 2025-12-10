package ble_func

import (
	"fmt"
	"path"

	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	"github.com/sirupsen/logrus"
)

// Adapter represents a Bluetooth adapter on the system.
type Adapter struct {
	Path dbus.ObjectPath
}

var cachedAdapter *adapter.Adapter1

// ListAdapters returns all the available adapter.
func listAdapters() ([]*Adapter, error) {
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
func get_adapter() (*adapter.Adapter1, error) {
	// List all available adapters
	adapters, err := listAdapters()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	fullPath := adapters[0].Path
	lastPart := path.Base(string(fullPath))
	a, err := api.GetAdapter(lastPart)
	if err != nil {
		return nil, fmt.Errorf("failed to get adapter: %v", err)
	}
	return a, nil
}

func InitAdapter() error {
	a, err := get_adapter()
	if err != nil {
		return fmt.Errorf("Fails to get adapter: %s", err)
	}
	cachedAdapter = a
	return nil
}

func GetService_Device(ServiceUUID string) (*device.Device1, error) {
	devices, err := cachedAdapter.GetDevices()
	if err != nil {
		return nil, fmt.Errorf("Failed to get devices: %s", err)
	}

	for _, dev := range devices {
		// Check if it matches target ServiceUUID
		uuids, err := dev.GetUUIDs()
		if err != nil {
			return nil, fmt.Errorf("failed to GetServiceData: %v", err)
		}
		for _, uuid := range uuids {
			if len(uuid) >= 23 && uuid[:23] == ServiceUUID {
				return dev, nil
			}
		}
	}
	return nil, nil
}
