package gobd2

import (
	"bufio"
	"strings"
	"time"

	"github.com/tarm/serial"
)

// SerialPort defines the interface for interacting with the serial port.
type SerialPort interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Close() error
}

// SerialPortOpener is responsible for opening serial port connections.
type SerialPortOpener interface {
	OpenPort(config *serial.Config) (SerialPort, error)
}

// RealPortOpener implements SerialPortOpener using the actual serial package.
type RealPortOpener struct{}

func (rpo *RealPortOpener) OpenPort(config *serial.Config) (SerialPort, error) {
	return serial.OpenPort(config)
}

// SerialConnector now includes a SerialPortOpener.
type SerialConnector struct {
	portOpener SerialPortOpener
	config     *serial.Config
	connection SerialPort
	reader     *bufio.Reader
	writer     *bufio.Writer
}

func NewSerialConnector(device string, baud int, opener SerialPortOpener) *SerialConnector {
	return &SerialConnector{
		portOpener: opener,
		config:     &serial.Config{Name: device, Baud: baud},
	}
}

func (sc *SerialConnector) Connect() error {
	var err error
	sc.connection, err = sc.portOpener.OpenPort(sc.config)
	if err != nil { //nolint:wsl
		return err
	}

	sc.reader = bufio.NewReader(sc.connection)
	sc.writer = bufio.NewWriter(sc.connection)

	return sc.initializeELM327()
}

func (sc *SerialConnector) Close() error {
	return sc.connection.Close()
}

func (sc *SerialConnector) SendCommand(command CommandCode) (string, error) {
	_, err := sc.writer.WriteString(string(command) + "\r")
	if err != nil {
		return "", err
	}

	sc.writer.Flush()

	// Reading and cleaning up the response to remove command echo and extra characters
	response, err := sc.reader.ReadString('>')
	if err != nil {
		return "", err
	}

	cleanedResponse := strings.Trim(response, " \r\n>")

	return cleanedResponse, nil
}

func (sc *SerialConnector) initializeELM327() error {
	initCommands := []CommandCode{"ATZ", "ATE0", "ATL0", "ATSP0"}
	for _, cmd := range initCommands {
		if _, err := sc.SendCommand(cmd); err != nil {
			return err
		}

		time.Sleep(100 * time.Millisecond) // Delay to allow the ELM327 to reset and apply settings
	}

	return nil
}
