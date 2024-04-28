package gobd2_test

// import (
// 	"testing"

// 	"github.com/godbus/dbus/v5"
// 	"github.com/janekbaraniewski/gobd2/gobd2"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// type MockDBusConn struct {
// 	mock.Mock
// }

// func (m *MockDBusConn) Object(dest string, path dbus.ObjectPath) dbus.BusObject {
// 	args := m.Called(dest, path)
// 	return args.Get(0).(dbus.BusObject)
// }

// func TestConnect(t *testing.T) {
// 	mockDBus := new(MockDBusConn)
// 	connector := gobd2.NewBluetoothConnector("00:1D:A5:68:98:8B", mockDBus)

// 	mockDBus.On("Object", "org.bluez", dbus.ObjectPath("/org/bluez/hci0/dev_00_1D_A5_68_98_8B")).Return(&MockBusObject{})
// 	assert.NoError(t, connector.Connect())
// 	mockDBus.AssertExpectations(t)
// }

// func TestSendCommand(t *testing.T) {
// 	mockDBus := new(MockDBusConn)
// 	connector := gobd2.NewBluetoothConnector("00:1D:A5:68:98:8B", mockDBus)
// 	mockObj := &MockBusObject{}
// 	mockDBus.On("Object", "org.bluez", dbus.ObjectPath("/org/bluez/hci0/dev_00_1D_A5_68_98_8B/service001/char001")).Return(mockObj)
// 	mockObj.On("Call", "org.bluez.GattCharacteristic1.WriteValue", 0, []byte("010C\r"), map[string]interface{}{}).Return(nil).Once()
// 	mockObj.On("Call", "org.bluez.GattCharacteristic1.ReadValue", 0, map[string]interface{}{}).Return(dbus.Call{Body: []interface{}{[]byte("OK\r")}}, nil).Once()

// 	data, err := connector.SendCommand("010C")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "OK", data)
// 	mockDBus.AssertExpectations(t)
// 	mockObj.AssertExpectations(t)
// }

// type MockBusObject struct {
// 	mock.Mock
// 	dbus.BusObject
// }

// func (m *MockBusObject) Call(method string, flags dbus.Flags, args ...interface{}) *dbus.Call {
// 	args = m.Called(method, flags, args)
// 	return &dbus.Call{
// 		Err:  args.Error(0),
// 		Body: []interface{}{args.Get(1)},
// 	}
// }
