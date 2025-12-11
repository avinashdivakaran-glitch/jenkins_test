package wifi_fn

import (
	"fmt"

	"github.com/Wifx/gonetworkmanager"
	"github.com/sirupsen/logrus"
)

// Return the netowk manger
func Get_networkManger() (gonetworkmanager.NetworkManager, error) {
	nm, err := gonetworkmanager.NewNetworkManager()
	if err != nil {
		return nil, fmt.Errorf("error connecting to NetworkManager: %v", err)
	}

	return nm, nil
}

// Return Service Setting
func Get_serviceSettings() (gonetworkmanager.Settings, error) {
	settingsService, err := gonetworkmanager.NewSettings()
	if err != nil {
		return nil, fmt.Errorf("failed to get settings service: %v", err)
	}

	return settingsService, nil
}

// return the wireless device
func Get_wirelessDevice(nm gonetworkmanager.NetworkManager) (gonetworkmanager.DeviceWireless, string, error) {
	devices, err := nm.GetPropertyAllDevices()
	if err != nil {
		return nil, "", fmt.Errorf("error getting devices: %v", err)
	}

	for _, device := range devices {
		devType, err := device.GetPropertyDeviceType()
		if err != nil {
			continue
		}
		if devType == gonetworkmanager.NmDeviceTypeWifi {
			wifiDevice, err := gonetworkmanager.NewDeviceWireless(device.GetPath())
			if err != nil {
				continue
			}
			ifaceName, err := wifiDevice.GetPropertyInterface()
			if err != nil {
				return nil, "", fmt.Errorf("error getting device interfaace name, %v", err)
			}
			logrus.Infof("Found wireless device interface %v", ifaceName)
			return wifiDevice, ifaceName, nil
		}
	}

	return nil, "", fmt.Errorf("not any wireless device found")
}

// Poer on adapter
func PowerOn_adapter(nm gonetworkmanager.NetworkManager, on bool) error {
	wifi_state, err := nm.GetPropertyWirelessEnabled()
	if err != nil {
		return fmt.Errorf("error getting wifi status: %v", err)
	}

	if on {
		if !wifi_state {
			logrus.Info("Wifi adapter is OFF - Powering on ...")
			err = nm.SetPropertyWirelessEnabled(true)
			if err != nil {
				return fmt.Errorf("failed to enable WiFi: %v", err)
			}
			logrus.Infoln("Adapter powered ON.")
		} else {
			logrus.Infoln("Adapter already ON.")
		}
	} else {
		if wifi_state {
			logrus.Info("Wifi adapter is ON - Powering off ...")
			err = nm.SetPropertyWirelessEnabled(false)
			if err != nil {
				return fmt.Errorf("failed to disable WiFi: %v", err)
			}
			logrus.Infoln("Adapter powered OFF.")
		} else {
			logrus.Infoln("Adapter already OFF.")
		}
	}

	return nil

}
