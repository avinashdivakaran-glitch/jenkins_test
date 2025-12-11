package wifi_fn

import (
	"fmt"
	"main/wifi_models"
	"time"

	"github.com/Wifx/gonetworkmanager"
	"github.com/sirupsen/logrus"
)

func ScanRoutine(w_dev gonetworkmanager.DeviceWireless, results chan<- []wifi_models.DeviceInfo, stopSignal <-chan struct{}) error {
	// Ticker to control hardware scan frequency
	scanTicker := time.NewTicker(2 * time.Second)
	defer scanTicker.Stop()

	// Run one scan immediately upon start
	if err := performScan(w_dev, results, stopSignal); err != nil {
		logrus.Error("Initial scan failed:", err)
	}

	for {
		select {
		case <-stopSignal:
			logrus.Infof("Scan Stoped ....")
			return nil

		case <-scanTicker.C:
			// Trigger periodic scan
			if err := performScan(w_dev, results, stopSignal); err != nil {
				logrus.Warnf("Periodic scan error: %v", err)
			}
		}
	}
}

func performScan(w_dev gonetworkmanager.DeviceWireless, results chan<- []wifi_models.DeviceInfo, stopSignal <-chan struct{}) error {
	// Scan request
	err := w_dev.RequestScan()
	if err != nil {
		return fmt.Errorf("error requesting scan: %v", err)
	}

	// // use a select here to allow interrupting the sleep immediately if stop is triggered
	select {
	case <-stopSignal:
		return nil
	case <-time.After(2 * time.Second):
		// Continue

	}

	// Fetch Results
	aps, err := w_dev.GetPropertyAccessPoints()
	if err != nil {
		return fmt.Errorf("error Get Property AccessPoints: %v", err)
	}

	var list []wifi_models.DeviceInfo
	for _, ap := range aps {
		ssid, _ := ap.GetPropertySSID()
		if ssid == "" {
			continue
		} // Skip hidden networks
		bssid, _ := ap.GetPropertyHWAddress()
		strength, _ := ap.GetPropertyStrength()

		list = append(list, wifi_models.DeviceInfo{
			SSID:     ssid,
			BSSID:    bssid,
			Strength: strength,
		})
	}

	// Send data to Main Loop
	select {
	case <-stopSignal:
		return nil
	case results <- list:
	}

	return nil
}

// inter, _ := dev.GetPropertyInterface()
// fmt.Println(inter)

// resultsChan := make(chan []wifi_models.DeviceInfo)
// scanStart := true
// var stopChan chan struct{}

// var currentList []wifi_models.DeviceInfo

// displayTicker := time.NewTicker(1 * time.Second)
// defer displayTicker.Stop()

// userButtonTicker := time.NewTicker(15 * time.Second)
// defer userButtonTicker.Stop()

// for {
// 	select {
// 	case newData := <-resultsChan:
// 		currentList = newData

// 	case <-displayTicker.C:
// 		fmt.Println(currentList)

// 	case <-userButtonTicker.C:
// 		scanStart = !scanStart

// 		if scanStart {
// 			logrus.Info("Start wifi scanning ..................")
// 			stopChan = make(chan struct{})
// 			go wifi_fn.ScanRoutine(dev, resultsChan, stopChan)
// 		} else {
// 			logrus.Info("Stop wifi scanning .. ................")
// 			if stopChan != nil {
// 				close(stopChan)
// 			}
// 			currentList = nil
// 		}
// 	}
// }

// ----------------------------------------------------------
// err := Init_tnnWifi()
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	err = Turn_tnnWiFi_on(true)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	resultsChan := make(chan []wifi_models.DeviceInfo)
// 	var stopChan chan struct{}
// 	errChan := make(chan error)

// 	var currentList []wifi_models.DeviceInfo

// 	displayTicker := time.NewTicker(1 * time.Second)
// 	defer displayTicker.Stop()

// 	userButtonTicker := time.NewTicker(10 * time.Second)
// 	defer userButtonTicker.Stop()

// 	stopChan = make(chan struct{})
// 	go Scan_tnnWiFi_dev(resultsChan, stopChan, errChan)

// 	for {
// 		select {
// 		case newData := <-resultsChan:
// 			currentList = newData

// 		case err := <-errChan:
// 			if err != nil {
// 				fmt.Printf("Scanner crashed: %v\n", err)
// 			}
// 			return

// 		case <-displayTicker.C:
// 			fmt.Println(currentList)

// 		case <-userButtonTicker.C:
// 			if stopChan != nil {
// 				close(stopChan)
// 			}

// 		}
// 	}
