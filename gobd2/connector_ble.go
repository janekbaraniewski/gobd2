package gobd2

import (
	"errors"
	"fmt"
	"strings"
	"time"

	dbus "github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
)

// BluetoothConnector handles Bluetooth connections.
type BluetoothConnector struct {
	deviceAddress string
	device        *device.Device1
	adapter       *adapter.Adapter1
}

// NewBluetoothConnector creates a new connector for a Bluetooth device.
func NewBluetoothConnector(deviceAddress string) *BluetoothConnector {
	return &BluetoothConnector{
		deviceAddress: deviceAddress,
	}
}

// Connect initializes the Bluetooth adapter and starts device discovery.
func (bc *BluetoothConnector) Connect() error {
	var err error
	if bc.adapter, err = adapter.GetDefaultAdapter(); err != nil {
		return fmt.Errorf("failed to get default adapter: %w", err)
	}

	if err = bc.adapter.StartDiscovery(); err != nil {
		return fmt.Errorf("failed to start discovery: %w", err)
	}

	// This should be handled better in a real application, with timeout and context handling.
	time.Sleep(10 * time.Second) // Wait for some time to discover devices

	devices, err := bc.adapter.GetDevices()
	if err != nil {
		return fmt.Errorf("failed to get devices: %w", err)
	}

	for _, d := range devices {
		if d.Properties.Address == bc.deviceAddress {
			bc.device = d

			break
		}
	}

	if bc.device == nil {
		return errors.New("device not found")
	}

	err = bc.device.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to device: %w", err)
	}

	return nil
}

// Close terminates the connection to the Bluetooth device.
func (bc *BluetoothConnector) Close() error {
	if bc.device != nil {
		return bc.device.Disconnect()
	}

	return nil
}

// Assuming you have a connected device object.
func (bc *BluetoothConnector) SendCommand(command CommandCode) (string, error) {
	// TODO: find the actual dbus paths for the device
	servicePath := fmt.Sprintf("%s/service0001", bc.device.Path())
	charPath := fmt.Sprintf("%s/char0001", servicePath)

	// Access the D-Bus connection
	conn, err := dbus.SystemBus()
	if err != nil {
		return "", fmt.Errorf("failed to connect to system bus: %w", err)
	}

	// Access the characteristic
	char := conn.Object("org.bluez", dbus.ObjectPath(charPath))
	writeValue := []byte(string(command) + "\r")

	call := char.Call("org.bluez.GattCharacteristic1.WriteValue", 0, writeValue, map[string]interface{}{})
	if call.Err != nil {
		return "", fmt.Errorf("failed to write value: %w", call.Err)
	}

	// Read the response, assuming the characteristic allows reading or notifies after write
	var value []byte

	call = char.Call("org.bluez.GattCharacteristic1.ReadValue", 0, map[string]interface{}{})
	if call.Err != nil {
		return "", fmt.Errorf("failed to read value: %w", call.Err)
	}

	if err := call.Store(&value); err != nil {
		return "", err
	}

	response := strings.Trim(string(value), " \r\n>")

	return response, nil
}
