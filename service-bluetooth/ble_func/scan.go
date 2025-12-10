package ble_func

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	"github.com/sirupsen/logrus"

	"ble_fn_mqtt/ble_models"
)

// Scane available all (BLE + Classic) devices
func Scan_tnnBLE(ctx context.Context, a *adapter.Adapter1, scan_DeviceList *[]ble_models.DeviceInfo, scan_Error *string, scan *ble_models.ScanMsg, scanStatusCh chan bool) error {
	fmt.Println("🔍 Scan routine started")
	scanning := false
	var discoveryCancel context.CancelFunc
	var discoveryChan <-chan *adapter.DeviceDiscovered
	var done chan struct{}
	var err error
	deviceMap := make(map[string]ble_models.DeviceInfo)
	var deviceMu sync.Mutex
	var lastScanEndTime time.Time

	updateTicker := time.NewTicker(2 * time.Second)
	defer updateTicker.Stop()

	for {
		select {
		case scanActive, ok := <-scanStatusCh:
			if !ok {
				if discoveryCancel != nil {
					discoveryCancel()
				}
				*scan_Error = "Scan status channel closed, stopping scan routine."
				return fmt.Errorf("scan status channel closed, stopping scan routine")
			}

			if scanActive && !scanning {
				logrus.Infoln("Starting discovery...")
				scanning = true
				deviceMu.Lock()
				scan.Data.IsScan = true
				done = make(chan struct{})

				// Reset global device map
				if time.Since(lastScanEndTime) > 20*time.Second {
					for k := range deviceMap {
						delete(deviceMap, k)
					}
					*scan_DeviceList = nil
				}
				// deviceMap = make(map[string]ble_models.DeviceInfo)
				*scan_DeviceList = nil
				deviceMu.Unlock()

				// Start discovery
				discoveryChan, discoveryCancel, err = api.Discover(a, nil)
				if err != nil {
					logrus.Errorln("Failed to start discovery:", err)
					deviceMu.Lock()
					scan.Data.IsScan = false
					deviceMu.Unlock()
					scanning = false
					continue
				}

				// Start detected_devices goroutine
				go Detected_devices(discoveryChan, deviceMap, &deviceMu, done)

			} else if !scanActive && scanning {
				logrus.Infoln("🛑 Stopping discovery...")
				deviceMu.Lock()
				scan.Data.IsScan = false
				deviceMu.Unlock()
				scanning = false
				*scan_DeviceList = nil
				if discoveryCancel != nil {
					discoveryCancel()
					discoveryCancel = nil
				}
				if done != nil {
					<-done
					done = nil
				}

				lastScanEndTime = time.Now()

				break
			}
			if scanning {
				deviceMu.Lock()
				scan.Data.IsScan = true
				deviceMu.Unlock()
			}

		case <-updateTicker.C:
			if scanning {
				// Update the slice from map every 4s
				deviceMu.Lock()
				*scan_DeviceList = make([]ble_models.DeviceInfo, 0, len(deviceMap))
				for _, d := range deviceMap {
					if len(d.Name) > 0 {
						*scan_DeviceList = append(*scan_DeviceList, d)
					}

					// Sort the pairedDevices slice by Name (alphabetically)
					sort.Slice(*scan_DeviceList, func(i, j int) bool {
						return (*scan_DeviceList)[i].Name < (*scan_DeviceList)[j].Name
					})
				}
				deviceMu.Unlock()
			}

		case <-ctx.Done():
			logrus.Infoln("Scan stopped due to context cancellation")
			if discoveryCancel != nil {
				discoveryCancel()
			}
			if done != nil {
				<-done
			}
			*scan_Error = "Scan stopped due to context cancellation"
			return fmt.Errorf("Scan stopped due to context cancellation")
		}
	}

}

// Goroutine to process discovered devices
func Detected_devices(discoveryChan <-chan *adapter.DeviceDiscovered, devicesMap map[string]ble_models.DeviceInfo, mu *sync.Mutex, done chan struct{}) {
	// Continuously reads discovered devices until the channel is closed
	for ev := range discoveryChan {
		dev, err := device.NewDevice1(ev.Path)
		if err != nil {
			logrus.Errorf("Failed to get device %s: %s", ev.Path, err)
			continue
		}
		props := dev.Properties

		devType := "Unknown"
		if len(props.UUIDs) > 0 {
			devType = "BLE"
		} else if props.Class != 0 {
			devType = "Classic"
		}

		device := ble_models.DeviceInfo{
			Name:   props.Name,
			UUIDs:  props.UUIDs,
			MAC:    props.Address,
			Signal: props.RSSI,
			Class:  props.Class,
			Type:   devType,
		}

		// Store in map to avoid duplicates
		mu.Lock()
		// if _, exists := devicesMap[device.MAC]; !exists {
		devicesMap[device.MAC] = device
		logrus.Infof("🚀 [%s] Name: %s, MAC: %s, Signal: %d, Class: %d\n", device.Type, device.Name, device.MAC, device.Signal, device.Class)

		// }
		mu.Unlock()
	}
	// Sends a signal to the main routine
	done <- struct{}{}
}
