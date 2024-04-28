package gobd2

import (
	"errors"
	"strings"
	"time"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

// BluetoothConnector handles Bluetooth connections.
type BluetoothConnector struct {
	deviceAddress string
	device        *bluetooth.Device
}

// NewBluetoothConnector creates a new connector for a Bluetooth device.
func NewBluetoothConnector(deviceAddress string) *BluetoothConnector {
	return &BluetoothConnector{
		deviceAddress: deviceAddress,
	}
}

// Connect initializes the Bluetooth adapter and starts device discovery.
func (bc *BluetoothConnector) Connect() error {
	if err := adapter.Enable(); err != nil {
		return err
	}

	// Start scanning
	ch := make(chan bluetooth.ScanResult, 1)
	err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if strings.EqualFold(strings.ToLower(result.Address.String()), strings.ToLower(bc.deviceAddress)) {
			dev, err := adapter.Connect(result.Address, bluetooth.ConnectionParams{})
			if err != nil {
				return
			}

			bc.device = &dev

			ch <- result

			if err := adapter.StopScan(); err != nil {
				return
			}
		}
	})
	if err != nil { //nolint
		return err
	}

	select {
	case <-ch:
	case <-time.After(10 * time.Second):
		return errors.New("failed to find device")
	}

	if bc.device == nil {
		return errors.New("device not found")
	}

	// Optionally connect to the device if not automatically handled by the adapter.Connect
	// if !bc.device.Connected() {
	// 	if err := bc.device.Connect(); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

// Close terminates the connection to the Bluetooth device.
func (bc *BluetoothConnector) Close() error {
	if bc.device != nil {
		return bc.device.Disconnect()
	}

	return nil
}

// SendCommand sends a command to a connected Bluetooth device and returns the response.
func (bc *BluetoothConnector) SendCommand(command CommandCode) (string, error) {
	if bc.device == nil {
		return "", errors.New("device not connected")
	}

	// Here you need to use specific GATT profile information, such as service UUID and characteristic UUID
	// Assume we have characteristic UUID for sending command and reading response
	serviceUUID, err := bluetooth.ParseUUID("your-service-uuid")
	if err != nil {
		return "", err
	}

	charUUID, err := bluetooth.ParseUUID("your-char-uuid")
	if err != nil {
		return "", err
	}

	service, err := bc.device.DiscoverServices([]bluetooth.UUID{serviceUUID})
	if err != nil {
		return "", err
	}

	characteristic, err := service[0].DiscoverCharacteristics([]bluetooth.UUID{charUUID})
	if err != nil {
		return "", err
	}

	// Write command to the characteristic
	if _, err := characteristic[0].WriteWithoutResponse([]byte(command + "\r")); err != nil {
		return "", err
	}

	response := []byte{}
	// Assume the device sends a notification with the response
	_, err = characteristic[0].Read(response)
	if err != nil {
		return "", err
	}

	return strings.Trim(string(response), " \r\n>"), nil
}
