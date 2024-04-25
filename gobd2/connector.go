package gobd2

// Connector defines the interface for connection operations.
type Connector interface {
	Connect() error
	Close() error
	SendCommand(command CommandCode) (string, error)
}
