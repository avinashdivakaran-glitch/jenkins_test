package ble_func

import (
	"ble_fn_mqtt/ble_models"
	"fmt"
	"sort"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/agent"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	"github.com/sirupsen/logrus"
)

// RegisterAgent creates and registers a pairing agent
func registerAgent(agentPath string) error {
	// Connect to system bus
	conn, err := dbus.SystemBus()
	if err != nil {
		return fmt.Errorf("failed to connect to system bus: %v", err)
	}
	_ = conn

	myAgent := NewMyAgent(agentPath)

	if err := agent.ExposeAgent(conn, myAgent, "KeyboardDisplay", true); err != nil {
		return fmt.Errorf("ExposeAgent failed: %v", err)
	}
	logrus.Infof("Agent exposed at %s and set as default agent\n", agentPath)
	logrus.Infoln("Bluetooth agent registered")

	return nil
}

// Pair with given MAC address
func Pair_device(adapter *adapter.Adapter1, mac string, timeout time.Duration) error {
	// Create device path from MAC
	devPath := dbus.ObjectPath(string(adapter.Path()) + "/dev_" + formatMAC(mac))
	// Ensure agent is registered
	registerAgent("/com/example/agent" + formatMAC(mac))
	// ; err != nil {
	// 	return err
	// }

	// Get the device object
	dev, err := device.NewDevice1(devPath)
	if err != nil {
		return fmt.Errorf("failed to get device object: %v", err)
	}

	// Check if already paired
	props, err := dev.GetProperties()
	if err != nil {
		return fmt.Errorf("failed to get properties: %v", err)
	}
	if props.Paired {
		logrus.Infof("Device %s is already paired\n", mac)
		return nil
	}

	logrus.Infof("Pairing with device %s ...\n", mac)
	err = dev.Pair()
	if err != nil {
		return fmt.Errorf("pairing failed: %v", err)
	}

	// Wait until device is paired
	start := time.Now()
	for {
		props, _ = dev.GetProperties()
		if props.Paired {
			logrus.Infof("Device %s paired successfully\n", mac)
			break
		}
		if time.Since(start) > timeout {
			return fmt.Errorf("pairing timeout after %s", timeout)
		}
		time.Sleep(500 * time.Millisecond)
	}

	return nil

}

func unregisterAgent(agentPath dbus.ObjectPath) error {
	manager, err := agent.NewAgentManager1()
	if err != nil {
		return fmt.Errorf("failed to get agent manager: %v", err)
	}

	if err := manager.UnregisterAgent(agentPath); err != nil {
		return fmt.Errorf("failed to unregister agent: %v", err)
	}

	logrus.Infoln("Bluetooth agent unregistered")
	return nil
}

// Remove the given MAC from pair list
func Remove_device(adapter *adapter.Adapter1, mac string) error {
	// Create device path from MAC
	devPath := dbus.ObjectPath(string(adapter.Path()) + "/dev_" + formatMAC(mac))
	agentPath := dbus.ObjectPath("/com/example/agent" + formatMAC(mac))

	// registerAgent("/com/example/agent" + formatMAC(mac))

	// Ensure agent is unregistered
	unregisterAgent(agentPath)
	// ; err != nil {
	// 	return err
	// }

	// Remove the device from BlueZ
	if err := adapter.RemoveDevice(devPath); err != nil {
		return fmt.Errorf("failed to forget device: %v", err)
	}

	logrus.Infof("Device %s forgotten successfully\n", mac)
	return nil
}

func Bound_devices(adapter *adapter.Adapter1) ([]ble_models.PairDeviceInfo, error) {
	// Get all devices known to the adapter
	devices, err := adapter.GetDevices()
	if err != nil {
		return nil, fmt.Errorf("failed to get devices: %v", err)
	}

	logrus.Info("🔗 Bonded (Paired) Bluetooth Devices:")
	logrus.Info("------------------------------------")

	var pairedDevices []ble_models.PairDeviceInfo

	for _, d := range devices {
		props, err := d.GetProperties()
		if err != nil {
			continue
		}

		// Filter only devices that are Paired
		if props.Paired {
			pairedDevice := ble_models.PairDeviceInfo{
				Name:        props.Name,
				UUIDs:       props.UUIDs,
				MAC:         props.Address,
				Signal:      props.RSSI,
				Class:       props.Class,
				Type:        "",
				IsAvailable: props.Connected,
			}
			pairedDevices = append(pairedDevices, pairedDevice)

			// Output the paired device details
			logrus.Infof("Name: %s, MAC: %s, Connected: %v\n", props.Name, props.Address, props.Connected)
		}
	}

	// Sort the pairedDevices slice by Name (alphabetically)
	sort.Slice(pairedDevices, func(i, j int) bool {
		return pairedDevices[i].Name < pairedDevices[j].Name
	})

	logrus.Info("------------------------------------")

	if len(pairedDevices) == 0 {
		return nil, nil
	}

	return pairedDevices, nil
}
