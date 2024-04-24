package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
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
		var err error

		if useBluetooth {
			if deviceAddress == "" {
				log.Fatal("Bluetooth device address must be provided when using Bluetooth.")
			}
			connector = gobd2.NewBluetoothConnector(deviceAddress)
		} else {
			connector = gobd2.NewSerialConnector(portName, baudRate, &gobd2.RealPortOpener{})
		}

		if err = connector.Connect(); err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		defer connector.Close()

		commander := gobd2.NewCommander(connector)
		if err := runMonitor(commander); err != nil {
			log.Fatalf("Monitor failed: %v", err)
		}
	},
}

// runMonitor initializes the UI and starts the monitoring process.
func runMonitor(commander *gobd2.Commander) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := termui.Init(); err != nil {
		return err
	}
	defer termui.Close()

	pids := []gobd2.CommandCode{
		gobd2.EngineRPMCommand,
		gobd2.VehicleSpeedCommand,
		gobd2.ThrottlePositionCommand,
		gobd2.CoolantTemperatureCommand,
	}
	widgetsList := setupUI(pids)
	startMonitoring(ctx, commander, widgetsList, pids)
	handleUIEvents(ctx)

	return nil
}

// setupUI configures the UI elements for each PID to be monitored.
func setupUI(pids []gobd2.CommandCode) []*widgets.Paragraph {
	termWidth, termHeight := termui.TerminalDimensions()
	widgetsList := make([]*widgets.Paragraph, len(pids))
	grid := termui.NewGrid()
	grid.SetRect(0, 0, termWidth, termHeight)

	for i, pid := range pids {
		widgetsList[i] = widgets.NewParagraph()
		widgetsList[i].Title = "PID: " + string(pid)
		widgetsList[i].Text = "Waiting for data..."
		widgetsList[i].Border = true
		colWidth := float64(12) / float64(len(pids))
		widgetsList[i].SetRect(0, i*termHeight/len(pids), termWidth, (i+1)*termHeight/len(pids))
		grid.Set(termui.NewRow(1.0/float64(len(pids)), termui.NewCol(colWidth, 0, widgetsList[i])))
	}

	termui.Render(grid)

	return widgetsList
}

// startMonitoring begins the data monitoring process for each PID using goroutines.
func startMonitoring(ctx context.Context, commander *gobd2.Commander, widgetsList []*widgets.Paragraph, pids []gobd2.CommandCode) {
	var wg sync.WaitGroup

	for i, pid := range pids {
		wg.Add(1)

		go func(i int, pid gobd2.CommandCode) {
			defer wg.Done()
			ticker := time.NewTicker(2 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					data, err := commander.ExecuteCommand(pid)
					if err != nil {
						widgetsList[i].Text = "Error: " + err.Error()
					} else {
						widgetsList[i].Text = "Data: " + data
					}
					termui.Render(widgetsList[i])
				}
			}
		}(i, pid)
	}

	wg.Wait()
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
