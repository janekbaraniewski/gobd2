package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/janekbaraniewski/gobd2/gobd2"
	"github.com/spf13/cobra"
)

var (
	portName      = "/dev/ttyUSB0" // Default serial port
	baudRate      = 9600           // Default baud rate for serial connections
	deviceAddress = ""             // Bluetooth device address (empty by default)
	useBluetooth  = false          // Flag to toggle Bluetooth connection
)

// monitorCmd defines the command line structure and handling for the monitoring tool.
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitors vehicle diagnostics using OBD2 interfaces.",
	Long: `This command supports real-time monitoring of various vehicle diagnostics parameters
from an OBD2 interface via serial or Bluetooth connection. It displays data dynamically in
a full-screen terminal interface powered by termui.`,
	Run: func(cmd *cobra.Command, args []string) {
		var connector gobd2.Connector

		if useBluetooth {
			if deviceAddress == "" {
				log.Fatal("Bluetooth device address must be provided when using Bluetooth.")
			}
			connector = gobd2.NewBluetoothConnector(deviceAddress)
		} else {
			connector = gobd2.NewSerialConnector(portName, baudRate, &gobd2.RealPortOpener{})
		}
		defer connector.Close()

		if err := connector.Connect(); err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}

		commander := gobd2.NewCommander(connector)
		runMonitor(commander)
	},
}

// createWidgets dynamically creates n widgets.
func createWidgets(n int) []*widgets.Paragraph {
	widgetsList := make([]*widgets.Paragraph, n)
	for i := range widgetsList {
		widgetsList[i] = widgets.NewParagraph()
		widgetsList[i].Text = fmt.Sprintf("Widget %d: Initializing...", i+1)
		widgetsList[i].Border = true
	}

	return widgetsList
}

// setupDynamicGrid arranges widgets into a dynamic grid.
func setupDynamicGrid(widgetsList []*widgets.Paragraph) *termui.Grid {
	grid := termui.NewGrid()
	termWidth, termHeight := termui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	cols := 3                                    // Define the number of columns
	rows := (len(widgetsList) + cols - 1) / cols // Calculate the required rows

	rowHeight := 1.0 / float64(rows) // Calculate the height of each row

	gridRows := []interface{}{}

	// Create slices to hold the row and column configurations
	for rowIndex := range rows {
		rowWidgets := []interface{}{}

		for colIndex := range cols {
			widgetIndex := rowIndex*cols + colIndex
			if widgetIndex >= len(widgetsList) {
				break // No more widgets to place
			}

			colWidth := 1.0 / float64(cols) // Calculate the width of each column
			widget := widgetsList[widgetIndex]
			rowWidgets = append(rowWidgets, termui.NewCol(colWidth, widget))
		}

		// Add a new row to the grid with the widgets for this row
		gridRows = append(gridRows, termui.NewRow(rowHeight, rowWidgets...))
	}

	grid.Set(gridRows...)

	return grid
}

// runMonitor initializes the UI and starts the monitoring process.
func runMonitor(commander *gobd2.Commander) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := termui.Init(); err != nil {
		log.Printf("failed to initialize termui: %v", err)
	}
	defer termui.Close()

	pids := []gobd2.CommandCode{
		gobd2.EngineRPMCommand,
		gobd2.VehicleSpeedCommand,
		gobd2.ThrottlePositionCommand,
		gobd2.CoolantTemperatureCommand,
	}
	widgetsList := createWidgets(len(pids))
	grid := setupDynamicGrid(widgetsList)

	termui.Render(grid) // Render the grid

	for i, widget := range widgetsList {
		go startMonitoring(ctx, widget, commander, pids[i]) // Update each widget at random intervals between 1-5 seconds
	}

	handleUIEvents(ctx)
}

// startMonitoring begins the data monitoring process for each PID using goroutines.
func startMonitoring(ctx context.Context, p *widgets.Paragraph, commander *gobd2.Commander, pid gobd2.CommandCode) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.Title = string(pid)

			data, err := commander.ExecuteCommand(pid)
			if err != nil {
				p.Text = "Error: " + err.Error()
			} else {
				p.Text = "Data: " + data
			}

			termui.Render(p)
		}
	}
}

// handleUIEvents handles user inputs and system signals to gracefully shut down the application.
func handleUIEvents(ctx context.Context) {
	uiEvents := termui.PollEvents()

	for {
		select {
		case e := <-uiEvents:
			if e.Type == termui.KeyboardEvent {
				switch e.ID {
				case "q", "<C-c>":
					return
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

// registerMonitorCommand adds the monitor command to the root command and sets up command line flags.
func registerMonitorCommand(rootCmd *cobra.Command) {
	monitorCmd.Flags().StringVarP(&portName, "port", "p", "/dev/ttyUSB0", "Specify the serial port for connection")
	monitorCmd.Flags().IntVarP(&baudRate, "baud", "b", 9600, "Specify the baud rate for serial connection")
	monitorCmd.Flags().StringVarP(&deviceAddress, "address", "a", "", "Specify the Bluetooth device address")
	monitorCmd.Flags().BoolVarP(&useBluetooth, "bluetooth", "l", false, "Use Bluetooth for connection instead of serial")

	rootCmd.AddCommand(monitorCmd)
}
