package gobd2_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/janekbaraniewski/gobd2/gobd2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tarm/serial"
)

type MockPortOpener struct {
	mock.Mock
}

func (m *MockPortOpener) OpenPort(config *serial.Config) (gobd2.SerialPort, error) {
	args := m.Called(config)
	// Correct
	port, ok := args.Get(0).(gobd2.SerialPort)
	if !ok {
		return nil, errors.New("type assertion failed")
	}

	return port, args.Error(1)
}

type MockSerialPort struct {
	mock.Mock
	buf bytes.Buffer
}

func (m *MockSerialPort) Read(p []byte) (int, error) {
	return m.buf.Read(p)
}

func (m *MockSerialPort) Write(p []byte) (int, error) {
	return m.buf.Write(p)
}

func (m *MockSerialPort) Close() error {
	args := m.Called()

	return args.Error(0)
}

func setupMock() *MockSerialPort {
	mockPort := &MockSerialPort{}
	// Use WriteString to simulate successful responses directly to the buffer
	responses := "ATZ\rOK\r>ATE0\rOK\r>ATL0\rOK\r>ATSP0\rOK\r>"
	mockPort.buf.WriteString(responses)
	mockPort.On("Close").Return(nil)

	return mockPort
}

func TestSerialConnector_Connect_Success(t *testing.T) {
	t.Parallel()

	mockPort := setupMock()
	mockOpener := new(MockPortOpener)

	// Correct setup: Ensure the mock opener returns a mock port and no error
	mockOpener.On("OpenPort", mock.Anything).Return(mockPort, nil)

	connector := gobd2.NewSerialConnector("COM1", 115200, mockOpener)
	err := connector.Connect()
	require.NoError(t, err)

	// Ensure the connection initialization sequence is correctly followed
	mockOpener.AssertCalled(t, "OpenPort", &serial.Config{Name: "COM1", Baud: 115200})
	mockOpener.AssertExpectations(t)
}

func TestSerialConnector_Connect_Failure(t *testing.T) {
	t.Parallel()

	mockOpener := new(MockPortOpener)
	expectedError := errors.New("failed to open port")

	// Simulate the failure scenario: Return `nil` for the SerialPort and `expectedError` for the error.
	mockOpener.On("OpenPort", mock.Anything).Return(&MockSerialPort{}, expectedError)

	connector := gobd2.NewSerialConnector("COM1", 115200, mockOpener)

	// Test the Connect method
	err := connector.Connect()
	require.Error(t, err)
	require.Equal(t, expectedError, err)

	// Verify the expectations were met for the mock
	mockOpener.AssertExpectations(t)
}
