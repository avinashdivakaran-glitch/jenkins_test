package main

import (
	"ble_fn_mqtt/ble_func"
	"ble_fn_mqtt/ble_models"
	"context"
	"fmt"

	"github.com/muka/go-bluetooth/bluez/profile/adapter"
)

var TNN_bleAdapter *adapter.Adapter1

func Init_tnnBLE() (err error) {
	TNN_bleAdapter, err = ble_func.Get_adapter()
	if err != nil {
		return fmt.Errorf("Error ble_func/Get_adapter: %v", err)
	}

	return nil
}

func Update_tnnBLE() (adapterName string, powerState bool, discoverable bool, err error) {
	adapterName, powerState, discoverable, err = ble_func.Prop_adapter(TNN_bleAdapter)
	if err != nil {
		return "", false, false, fmt.Errorf("Error ble_func/Prop_adapter: %v", err)
	}

	return adapterName, powerState, discoverable, nil
}

func Bound_tnnBLE_devie() (deviceList []ble_models.PairDeviceInfo, err error) {
	deviceList, err = ble_func.Bound_devices(TNN_bleAdapter)
	if err != nil {
		return nil, fmt.Errorf("Error ble_func/Bound_devices: %v", err)
	}

	return deviceList, nil
}

func Turn_tnnBLE_on(state bool) error {
	err := ble_func.PowerOn_adapter(TNN_bleAdapter, state)
	if err != nil {
		return fmt.Errorf("Error ble_func/PowerOn_adapter: %v", err)
	}

	return nil
}

func Discover_tnnBLE(state bool) error {
	if state {
		if err := ble_func.Make_discoverable(TNN_bleAdapter, 0); err != nil {
			return fmt.Errorf("Error ble_func/Make_discoverable %v", err)
		}
	} else {
		if err := ble_func.Stop_discoverable(TNN_bleAdapter); err != nil {
			return fmt.Errorf("Error ble_func/Stop_discoverable: %v", err)
		}
	}
	return nil
}

func Scanning_tnnBLE(ctx context.Context, scan_DeviceList *ble_models.ScanDeviceList, scan *ble_models.BluetoothStatus, scanStatusCh chan bool) error {
	err := ble_func.Scan_tnnBLE(ctx, TNN_bleAdapter, &scan_DeviceList.Data.Devices, &scan_DeviceList.Data.Error, &scan.ScanMsg, scanStatusCh)
	if err != nil {
		return fmt.Errorf("Error ble_func/Scan_tnnBLE: %v", err)
	}

	return nil
}

func Pair_tnnBLE_device(mac string) error {
	err := ble_func.Pair_device(TNN_bleAdapter, mac, 20)
	if err != nil {
		return fmt.Errorf("Error ble_func/Pair_device: %v", err)
	}

	return nil
}

func Remove_tnnBLE_device(mac string) error {
	// err := ble_func.Pair_device(TNN_bleAdapter, mac, 20)

	err := ble_func.Remove_device(TNN_bleAdapter, mac)
	if err != nil {
		return fmt.Errorf("Error ble_func/Remove_device: %v", err)
	}

	return nil
}

func Connect_tnnBLE_device(mac string) error {
	err := ble_func.Connect_device(TNN_bleAdapter, mac)
	if err != nil {
		return fmt.Errorf("Error ble_func/Connect_device: %v", err)
	}

	return nil
}

func Disconnect_tnnBLE_device(mac string) error {
	err := ble_func.Disconnect_device(TNN_bleAdapter, mac)
	if err != nil {
		return fmt.Errorf("Error ble_func/Disconnect_device: %v", err)
	}

	return nil
}
