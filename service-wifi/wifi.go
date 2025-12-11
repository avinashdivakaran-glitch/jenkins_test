package main

import (
	"fmt"
	"main/wifi_fn"
	"main/wifi_models"

	"github.com/Wifx/gonetworkmanager"
)

var TNN_wirelessDevice gonetworkmanager.DeviceWireless
var TNN_networkManager gonetworkmanager.NetworkManager

func Init_tnnWifi() (dev string, err error) {
	TNN_networkManager, err = wifi_fn.Get_networkManger()
	if err != nil {
		return "", fmt.Errorf("⚠️  Error wifi_func/Get_networkManger: %v", err)
	}

	TNN_wirelessDevice, dev, err = wifi_fn.Get_wirelessDevice(TNN_networkManager)
	if err != nil {
		return "", fmt.Errorf("⚠️  Error wifi_func/Get_wirelessDevice: %v", err)
	}

	return dev, nil
}

func WiFi_status() (wifi_models.PropertiesMsg, error) {
	status, err := wifi_fn.ConnectionProp(TNN_wirelessDevice)
	if err != nil {
		return wifi_models.PropertiesMsg{}, fmt.Errorf("⚠️  Error wifi_func/ConnectionProp: %v", err)
	}

	return status, nil
}

func Turn_tnnWiFi_on(state bool) error {
	err := wifi_fn.PowerOn_adapter(TNN_networkManager, state)
	if err != nil {
		return fmt.Errorf("⚠️  Error wifi_func/PowerOn_adapter: %v", err)
	}

	return nil
}

func Hotspot_tnnOn(ssid string, password string) error {
	err := wifi_fn.WiFi_hotspotOn(TNN_networkManager, TNN_wirelessDevice, ssid, password)
	if err != nil {
		return fmt.Errorf("⚠️  Error wifi_func/WiFi_hotspotOn: %v", err)
	}

	return nil
}

func Scan_tnnWiFi_dev(resultsChan chan []wifi_models.DeviceInfo, stopChan chan struct{}, errChan chan error) {

	// resultsChan := make(chan []wifi_models.DeviceInfo)
	// var stopChan chan struct{}
	// errChan := make(chan error)

	go func() {
		errChan <- wifi_fn.ScanRoutine(TNN_wirelessDevice, resultsChan, stopChan)
	}()
}

func Connect_tnnWiFi_dev(ssid string, password string, recon bool) error {
	err := wifi_fn.Connect_device(TNN_networkManager, TNN_wirelessDevice, ssid, password, recon)
	if err != nil {
		return fmt.Errorf("⚠️  Error wifi_func/Connect_device: %v", err)
	}

	return nil
}

func Disconnect_tnnWiFi() error {
	err := wifi_fn.Disconnect_con(TNN_wirelessDevice)
	if err != nil {
		return fmt.Errorf("⚠️  Error wifi_func/Disconnect_device: %v", err)
	}

	return nil
}

func Delete_tnnWiFi_con(ssid string) error {
	err := wifi_fn.Delete_Con(ssid)
	if err != nil {
		return fmt.Errorf("⚠️  Error wifi_func/Delete_Con: %v", err)
	}

	return nil
}
