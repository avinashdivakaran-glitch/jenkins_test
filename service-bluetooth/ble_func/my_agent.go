package ble_func

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/godbus/dbus/v5"
)

// The signatures use godbus/dbus and go-bluetooth style return types (string, *dbus.Error) etc.
type MyAgent struct {
	path dbus.ObjectPath
}

// constructor
func NewMyAgent(path string) *MyAgent {
	return &MyAgent{path: dbus.ObjectPath(path)}
}

// Path returns the object path (used by your ExposeAgent helper)
func (a *MyAgent) Path() dbus.ObjectPath {
	return a.path
}

// Interface returns the DBus interface name
func (a *MyAgent) Interface() string {
	return "org.bluez.Agent1"
}

// --- Agent methods (BlueZ will call these during pairing) ---

// RequestPinCode: called for legacy PIN pairing (we return a 6-digit PIN if asked)
func (a *MyAgent) RequestPinCode(device dbus.ObjectPath) (string, *dbus.Error) {
	// Prompt operator for the PIN (typed in the terminal)
	fmt.Printf("RequestPinCode for %s - enter PIN shown on peripheral: ", device)
	pin := readLineTrim()
	return pin, nil
}

// DisplayPinCode: called to display a PIN (we just print it)
func (a *MyAgent) DisplayPinCode(device dbus.ObjectPath, pincode string) *dbus.Error {
	log.Printf("[Agent] DisplayPinCode for %s : %s\n", device, pincode)
	return nil
}

// RequestPasskey: called when BlueZ wants the passkey (numeric entry)
func (a *MyAgent) RequestPasskey(device dbus.ObjectPath) (uint32, *dbus.Error) {
	fmt.Printf("RequestPasskey for %s - enter 6-digit passkey shown on peripheral: ", device)
	var pk uint32
	_, err := fmt.Sscanf(readLineTrim(), "%d", &pk)
	if err != nil {
		// return a DBus error if parse failed
		return 0, dbus.MakeFailedError(err)
	}
	return pk, nil
}

// DisplayPasskey: peripheral asks you to display the passkey (we log it)
func (a *MyAgent) DisplayPasskey(device dbus.ObjectPath, passkey uint32, entered uint16) *dbus.Error {
	log.Printf("[Agent] DisplayPasskey %06d (entered=%d) for %s\n", passkey, entered, device)
	return nil
}

// RequestConfirmation: numeric comparison confirmation (we accept automatically)
func (a *MyAgent) RequestConfirmation(device dbus.ObjectPath, passkey uint32) *dbus.Error {
	log.Printf("[Agent] RequestConfirmation %06d for %s -> accepting\n", passkey, device)
	return nil
}

// RequestAuthorization: called to authorize a connection
func (a *MyAgent) RequestAuthorization(device dbus.ObjectPath) *dbus.Error {
	log.Printf("[Agent] RequestAuthorization for %s -> allowing\n", device)
	return nil
}

// AuthorizeService: called when a remote requests a particular service
func (a *MyAgent) AuthorizeService(device dbus.ObjectPath, uuid string) *dbus.Error {
	log.Printf("[Agent] AuthorizeService %s for %s -> allowing\n", uuid, device)
	return nil
}

// Cancel: called if pairing is canceled
func (a *MyAgent) Cancel() *dbus.Error {
	log.Println("[Agent] Cancel called")
	return nil
}

// Release: optional cleanup call
func (a *MyAgent) Release() *dbus.Error {
	log.Println("[Agent] Release called")
	return nil
}

// utility to read a line from stdin and trim
func readLineTrim() string {
	reader := bufio.NewReader(os.Stdin)
	s, _ := reader.ReadString('\n')
	return strings.TrimSpace(s)
}
