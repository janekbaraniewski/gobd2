/*
Package main - OBD2 Diagnostic Tool

This package implements a command-line interface to interact with OBD2 interfaces.
It allows users to monitor vehicle diagnostics using the ELM327 interface or compatible devices.

This tool leverages the Cobra library for creating powerful and flexible CLI applications.

Example usage:

  - Connect to a serial OBD2 device:
    ./gobd2 monitor --port /dev/ttyUSB0 --baud 115200

  - Connect to a Bluetooth OBD2 device:
    ./gobd2 monitor --bluetooth --address "00:1D:A5:68:98:8B"

For more information and updates, visit https://github.com/janekbaraniewski/gobd2.
*/
package main

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gobd2",
	Short: "gobd2 is a command-line tool for interacting with OBD2 interfaces.",
	Long: `gobd2 provides a command-line interface to connect with OBD2 interfaces via serial or Bluetooth connections,
utilizing the ELM327 or compatible devices to monitor vehicle diagnostics, retrieve data, and manage fault codes.

The tool supports a variety of commands to facilitate real-time monitoring and diagnostics operations.`,
}

func main() {
	registerMonitorCommand(rootCmd)
	registerListCommand(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
