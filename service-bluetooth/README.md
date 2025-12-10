# Bluetooth Control Terminal

## Overview

A Bluetooth control terminal for managing Bluetooth adapters and devices using BlueZ and Go Bluetooth libraries. This tool allows users to interact with Bluetooth adapters, make them discoverable, scan for devices, and perform device management tasks like connecting, disconnecting, pairing, and removing devices.

The program provides an interactive CLI interface to control Bluetooth operations on a system running BlueZ. It uses Go's concurrency features and BlueZ's D-Bus API to manage devices efficiently.

## Features

- Power on/off Bluetooth adapter
- Toggle discoverability of the Bluetooth adapter
- Scan for available Bluetooth devices (both BLE and Classic)
- Connect, disconnect, pair, and remove Bluetooth devices
- Display current Bluetooth adapter properties

## Installation

### Install Go
Ensure Go (1.16 or newer) is installed on your machine. You can install Go from the official Go website

### Install dependencies
The project depends on the go-bluetooth package for interacting with Bluetooth. Install it with:
```bash
go get github.com/muka/go-bluetooth
```

## Usage
Run the program with the following command:
```bash
go run .
```

## Main Commands
1. Power ON Adapter: Powers on the Bluetooth adapter.
2. Power OFF Adapter: Powers off the Bluetooth adapter.
3. Make Discoverable: Makes the Bluetooth adapter discoverable for 30 seconds.
4. Stop Discoverable: Disables discoverable mode.
5. Scan Devices: Scans for nearby Bluetooth devices. You'll be prompted to input a scan duration in seconds.
6. Connect Device: Connect to a device by entering its MAC address.
7. Disconnect Device: Disconnect a device by entering its MAC address.
8. Pair Device: Pair with a device by entering its MAC address.
9. Remove Device: Remove a device from the paired list by entering its MAC address.
10. Exit: Exits the program.


## Functions

### 1. Get_adapter()
Returns the default Bluetooth adapter (usually hci0).

### 2. Prop_adapter()
Retrieves and returns properties of the Bluetooth adapter, including its name, power state, and discoverability.

### 3. PowerOn_adapter()
Powers the Bluetooth adapter on or off.

### 4. Make_discoverable()
Enables or disables discoverability for the Bluetooth adapter.

### 5. Scan_allDevices()
Scans for all available Bluetooth devices (both BLE and Classic) for a specified duration.

### 6. Connect_device()
Connects to a device by its MAC address.

### 7. Disconnect_device()
Disconnects a device by its MAC address.

### 8. Pair_device()
Pairs with a device by its MAC address.

### 9. Remove_device()
Removes a device from the list of paired devices by its MAC address.