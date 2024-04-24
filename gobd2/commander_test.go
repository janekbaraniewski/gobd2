package gobd2_test

import (
	"testing"

	"github.com/janekbaraniewski/gobd2/gobd2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockConnector struct {
	mock.Mock
}

func (m *MockConnector) Connect() error {
	args := m.Called()

	return args.Error(0)
}

func (m *MockConnector) Close() error {
	args := m.Called()

	return args.Error(0)
}

func (m *MockConnector) SendCommand(command gobd2.CommandCode) (string, error) {
	args := m.Called(command)

	return args.String(0), args.Error(1)
}

func TestCommander_ExecuteCommand(t *testing.T) {
	t.Parallel()

	mockConnector := new(MockConnector)
	commander := gobd2.NewCommander(mockConnector)

	mockConnector.On("SendCommand", gobd2.EngineRPMCommand).Return("1234 RPM", nil)
	result, err := commander.ExecuteCommand(gobd2.EngineRPMCommand)

	require.NoError(t, err)
	require.Equal(t, "1234 RPM", result)
	mockConnector.AssertExpectations(t)
}
