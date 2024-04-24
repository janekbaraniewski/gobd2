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
	"github.com/janekbaraniewski/gobd2/gobd2" // Update this path as necessary
	"github.com/spf13/cobra"
)

var (
	portName      = "/dev/ttyUSB0" // Set this to your OBD-II device's serial port
	baudRate      = 9600           // Set this to the correct baud rate for your device
	deviceAddress = ""
	useBluetooth  = false
)

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor all available PIDs",
	Long:  `Monitor all available PIDs and display the results in real-time in a full-screen UI.`,
	Run: func(cmd *cobra.Command, args []string) {
		var connector gobd2.Connector
		var err error

		if useBluetooth {
			if deviceAddress == "" {
				log.Fatal("Bluetooth device address must be provided")
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
		gobd2.VehicleSpeedCommand,
	}
	widgetsList := setupUI(pids)
	startMonitoring(ctx, commander, widgetsList, pids)
	handleUIEvents(ctx)

	return nil
}

func setupUI(pids []gobd2.CommandCode) []*widgets.Paragraph {
	termWidth, termHeight := termui.TerminalDimensions()
	widgetsList := make([]*widgets.Paragraph, len(pids))
	grid := termui.NewGrid()
	grid.SetRect(0, 0, termWidth, termHeight)

	for i, pid := range pids {
		widgetsList[i] = widgets.NewParagraph()
		widgetsList[i].Title = "PID: " + string(pid)
		widgetsList[i].Text = "Initializing..."
		widgetsList[i].Border = true
		colWidth := float64(1) / float64(len(pids))
		widgetsList[i].SetRect(0, 0, int(float64(termWidth)*colWidth), termHeight)
		grid.Set(termui.NewRow(colWidth, termui.NewCol(colWidth, 0, widgetsList[i])))
	}

	termui.Render(grid)

	return widgetsList
}

func startMonitoring(ctx context.Context, commander *gobd2.Commander, widgetsList []*widgets.Paragraph, pids []gobd2.CommandCode) { //nolint:lll
	var wg sync.WaitGroup

	for i, pid := range pids {
		wg.Add(1)

		monitorFunc := func(i int, pid gobd2.CommandCode) {
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
		}
		go monitorFunc(i, pid)
	}

	wg.Wait()
}

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

func registerMonitorCommand(rootCmd *cobra.Command) {
	monitorCmd.Flags().StringVarP(&portName, "port", "p", "/dev/ttyUSB0", "Serial port name")
	monitorCmd.Flags().IntVarP(&baudRate, "baud", "b", 9600, "Baud rate for serial connection")
	monitorCmd.Flags().StringVarP(&deviceAddress, "address", "a", "", "Bluetooth device address")
	monitorCmd.Flags().BoolVarP(&useBluetooth, "bluetooth", "l", false, "Use Bluetooth connector instead of serial")

	rootCmd.AddCommand(monitorCmd)
}
