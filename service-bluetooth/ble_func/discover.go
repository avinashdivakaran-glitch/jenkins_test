package ble_func

import (
	"fmt"

	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/sirupsen/logrus"
)

// Makes the adapter visible to other devices
func Make_discoverable(a *adapter.Adapter1, timeout uint32) error {

	logrus.Info("Making adapter discoverable...")

	if err := a.SetProperty("DiscoverableTimeout", timeout); err != nil {
		return fmt.Errorf("failed to set discoverable timeout: %v", err)
	}

	if err := a.SetProperty("Discoverable", true); err != nil {
		return fmt.Errorf("failed to make adapter discoverable: %v", err)
	}

	logrus.Infof("Adapter is now discoverable for %d seconds 🔍\n", timeout)

	return nil
}

func Stop_discoverable(a *adapter.Adapter1) error {
	logrus.Info("Stopping discoverable mode...")

	if err := a.SetProperty("Discoverable", false); err != nil {
		return fmt.Errorf("failed to stop discoverable mode: %v", err)
	}

	logrus.Info("Adapter is no longer discoverable ❌")

	return nil
}
