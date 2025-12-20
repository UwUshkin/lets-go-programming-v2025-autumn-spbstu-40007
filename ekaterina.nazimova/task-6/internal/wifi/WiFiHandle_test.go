package wifi_test

import (
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type MockWiFi struct {
	mock.Mock
}

func (m *MockWiFi) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	res := args.Get(0)
	var ifaces []*wifi.Interface
	if res != nil {
		ifaces = res.([]*wifi.Interface)
	}

	return ifaces, args.Error(1)
}
