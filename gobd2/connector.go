package goobd2

import (
	"bufio"
	"github.com/tarm/serial"
)

// Connector defines the interface for connection operations.
type Connector interface {
	Connect() error
	Close() error
	Reader() *bufio.Reader
	Writer() *bufio.Writer
}

// SerialConnector implements Connector interface using a serial port.
type SerialConnector struct {
	config     *serial.Config
	connection *serial.Port
	reader     *bufio.Reader
	writer     *bufio.Writer
}

func NewSerialConnector(device string, baud int) *SerialConnector {
	return &SerialConnector{
		config: &serial.Config{Name: device, Baud: baud},
	}
}

func (sc *SerialConnector) Connect() error {
	var err error
	sc.connection, err = serial.OpenPort(sc.config)
	if err != nil {
		return err
	}
	sc.reader = bufio.NewReader(sc.connection)
	sc.writer = bufio.NewWriter(sc.connection)
	return nil
}

func (sc *SerialConnector) Close() error {
	return sc.connection.Close()
}

func (sc *SerialConnector) Reader() *bufio.Reader {
	return sc.reader
}

func (sc *SerialConnector) Writer() *bufio.Writer {
	return sc.writer
}
