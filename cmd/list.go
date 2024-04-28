//nolint:all
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "Monitors vehicle diagnostics using OBD2 interfaces.",
	Long: `This command supports real-time monitoring of various vehicle diagnostics parameters
from an OBD2 interface via serial or Bluetooth connection. It displays data dynamically in
a full-screen terminal interface powered by termui.`,
	Run: func(cmd *cobra.Command, args []string) {
		devices, err := discoverDevices()
		if err != nil {
			log.Fatalf("Error discovering devices: %v", err)
		}

		fmt.Println("Discovered Devices:") //nolint
		for _, dev := range devices {
			props, _ := dev.GetProperties()
			fmt.Printf("Name: %s, Address: %s\n", props.Name, props.Address) //nolint
		}
	},
}

func discoverDevices() ([]*device.Device1, error) {
	// Get the default Bluetooth adapter.
	a, err := adapter.GetDefaultAdapter()
	if err != nil {
		return nil, fmt.Errorf("failed to get default adapter: %v", err) //nolint
	}

	// Start discovery.
	log.Println("Starting discovery...") //nolint
	err = a.StartDiscovery()
	if err != nil {
		return nil, fmt.Errorf("failed to start discovery: %v", err)
	}
	defer a.StopDiscovery() //nolint:errcheck

	// Allow some time for discovery.
	time.Sleep(10 * time.Second)

	// List discovered devices.
	devices, err := a.GetDevices()
	if err != nil {
		return nil, fmt.Errorf("failed to get devices: %v", err)
	}

	return devices, nil
}

// registerMonitorCommand adds the monitor command to the root command and sets up command line flags.
func registerListCommand(rootCmd *cobra.Command) {
	rootCmd.AddCommand(listCommand)
}
