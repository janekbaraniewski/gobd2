package gobd2

import (
	"errors"
	"log"
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
func (bc *BluetoothConnector) Connect() error {
	log.Println("BluetoothConnector::Connect start")
	log.Println("Enable BT adapter")
	if err := adapter.Enable(); err != nil {
		return err
	}

	// Start scanning
	ch := make(chan bluetooth.ScanResult, 1)
	log.Println("Start scan")
	err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		log.Printf("got scan result - %v", result)
		if strings.EqualFold(strings.ToLower(result.Address.String()), strings.ToLower(bc.deviceAddress)) {
			log.Printf("got matching result - %v", result)
			ch <- result // Send result to channel and handle connection outside the callback
		}
	})
	if err != nil { //nolint
		return err
	}

	// Handle timeout
	select {
	case result := <-ch:
		log.Println("Attempting to connect to device")
		dev, err := adapter.Connect(result.Address, bluetooth.ConnectionParams{})
		if err != nil {
			log.Println("Error connecting to the device:", err)
			return err
		}
		bc.device = &dev
		log.Println("Device connected")
		return nil
	case <-time.After(10 * time.Second):
		log.Println("Failed to find device within 10 seconds")
		adapter.StopScan()
		return errors.New("failed to find device")
	}
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
